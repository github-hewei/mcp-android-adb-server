package device

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/electricbubble/gadb"
)

// TempPath is the path to the temporary directory on the device.
var TempPath = "/data/local/tmp"

// App represents an Android application.
type App struct {
	PackageName string `json:"package_name"`
}

// SystemInfo represents the system information of the device.
type SystemInfo struct {
	Model            string `json:"model"`             // Device model
	Brand            string `json:"brand"`             // Device brand
	Manufacturer     string `json:"manufacturer"`      // Manufacturer
	AndroidVersion   string `json:"android_version"`   // Android version
	SDK              int    `json:"sdk"`               // SDK version
	Battery          int    `json:"battery"`           // Battery percentage
	BatteryStatus    string `json:"battery_status"`    // Battery status (charging/discharging)
	ScreenWidth      int    `json:"screen_width"`      // Screen width
	ScreenHeight     int    `json:"screen_height"`     // Screen height
	ScreenDensity    int    `json:"screen_density"`    // Screen density
	NetworkType      string `json:"network_type"`      // Network type (WiFi/Mobile data)
	WifiSSID         string `json:"wifi_ssid"`         // WiFi SSID
	WifiSignal       int    `json:"wifi_signal"`       // WiFi signal strength
	LocationEnabled  bool   `json:"location_enabled"`  // Location service enabled
	IMEI             string `json:"imei"`              // IMEI
	SerialNumber     string `json:"serial_number"`     // Serial number
	TotalRAM         int64  `json:"total_ram"`         // Total RAM (bytes)
	AvailableRAM     int64  `json:"available_ram"`     // Available RAM (bytes)
	TotalStorage     int64  `json:"total_storage"`     // Total storage (bytes)
	AvailableStorage int64  `json:"available_storage"` // Available storage (bytes)
}

// Option is a function that can be used to configure an AndroidDevice instance.
type Option func(*AndroidDevice)

// AndroidDevice represents an Android device.
type AndroidDevice struct {
	id              string
	adb             gadb.Device
	swipeDuration   time.Duration
	longTapDuration time.Duration
	sleepDuration   time.Duration
	screenshotPath  string
	screenPassword  string
}

// NewAndroidDevice creates a new AndroidDevice instance.
func NewAndroidDevice(id string, opts ...Option) (*AndroidDevice, error) {
	adb, err := getAdb(id)
	if err != nil {
		return nil, err
	}

	wd, _ := os.Getwd()
	d := &AndroidDevice{
		id:              id,
		adb:             adb,
		swipeDuration:   time.Millisecond * 500,
		longTapDuration: time.Second * 2,
		sleepDuration:   time.Second,
		screenshotPath:  path.Join(wd, "screenshot"),
	}

	for _, opt := range opts {
		opt(d)
	}

	return d, nil
}

// getAdb returns a new gadb.Device instance.
func getAdb(id string) (gadb.Device, error) {
	adb, err := gadb.NewClient()
	if err != nil {
		return gadb.Device{}, err
	}

	deviceList, err := adb.DeviceList()
	if err != nil {
		return gadb.Device{}, err
	}

	for _, device := range deviceList {
		if device.Serial() == id {
			return device, nil
		}
	}

	return gadb.Device{}, nil
}

// WithSwipeDuration sets the swipe duration for the device.
func WithSwipeDuration(duration time.Duration) Option {
	return func(d *AndroidDevice) {
		d.swipeDuration = duration
	}
}

// WithLongTapDuration sets the long tap duration for the device.
func WithLongTapDuration(duration time.Duration) Option {
	return func(d *AndroidDevice) {
		d.longTapDuration = duration
	}
}

// WithSleepDuration sets the sleep duration for the device.
func WithSleepDuration(duration time.Duration) Option {
	return func(d *AndroidDevice) {
		d.sleepDuration = duration
	}
}

// WithScreenshotPath sets the screenshot path for the device.
func WithScreenshotPath(path string) Option {
	return func(d *AndroidDevice) {
		d.screenshotPath = path
	}
}

// WithScreenPassword sets the screen password for the device.
func WithScreenPassword(password string) Option {
	return func(d *AndroidDevice) {
		d.screenPassword = password
	}
}

// InstallApp installs an app on the device.
func (d *AndroidDevice) InstallApp(filePath string, reinstall ...bool) (err error) {
	apkName := filepath.Base(filePath)
	if !strings.HasSuffix(strings.ToLower(apkName), ".apk") {
		return fmt.Errorf("apk file must have an extension of .apk: %s", filePath)
	}

	var appFile *os.File
	if appFile, err = os.Open(filePath); err != nil {
		return fmt.Errorf("apk open: %w", err)
	}

	remotePath := path.Join(TempPath, apkName)
	if err := d.adb.PushFile(appFile, remotePath); err != nil {
		return fmt.Errorf("apk push: %w", err)
	}

	var shellOutput string
	if len(reinstall) != 0 && reinstall[0] {
		shellOutput, err = d.RunShellCommand("pm install", "-r", remotePath)
	} else {
		shellOutput, err = d.RunShellCommand("pm install", remotePath)
	}

	if err != nil {
		return fmt.Errorf("apk install: %w", err)
	}

	if !strings.Contains(shellOutput, "Success") {
		return fmt.Errorf("apk installed: %s", shellOutput)
	}

	return
}

// UninstallApp uninstalls an app from the device.
func (d *AndroidDevice) UninstallApp(packageName string, keepData ...bool) (err error) {
	var shellOutput string
	if len(keepData) != 0 && keepData[0] {
		shellOutput, err = d.RunShellCommand("pm uninstall", "-k", packageName)
	} else {
		shellOutput, err = d.RunShellCommand("pm uninstall", packageName)
	}

	if err != nil {
		return fmt.Errorf("apk uninstall: %w", err)
	}

	if !strings.Contains(shellOutput, "Success") {
		return fmt.Errorf("apk uninstall: %s", shellOutput)
	}

	return
}

// TerminateApp terminates an app on the device.
func (d *AndroidDevice) TerminateApp(packageName string) (err error) {
	_, err = d.RunShellCommand("am force-stop", packageName)
	return
}

// LaunchApp launches an app on the device.
func (d *AndroidDevice) LaunchApp(packageName string) (err error) {
	var shellOutput string
	if shellOutput, err = d.RunShellCommand("monkey -p", packageName, "-c android.intent.category.LAUNCHER 1"); err != nil {
		return err
	}

	if strings.Contains(shellOutput, "monkey aborted") {
		return fmt.Errorf("app launch: %s", strings.TrimSpace(shellOutput))
	}

	return
}

// ListApp lists installed apps on the device.
func (d *AndroidDevice) ListApp() ([]App, error) {
	out, err := d.RunShellCommand("pm list packages -f")
	if err != nil {
		return nil, fmt.Errorf("list packages: %w", err)
	}

	lines := strings.Split(out, "\n")
	apps := make([]App, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 格式: package:/data/app/~~XXXXX/com.example.app-XXXXX/base.apk=com.example.app
		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			continue
		}

		packageName := parts[1]
		apps = append(apps, App{
			PackageName: packageName,
		})
	}

	return apps, nil
}

// InstalledApp checks if an app is installed on the device.
func (d *AndroidDevice) InstalledApp(packageName string) (bool, error) {
	out, err := d.RunShellCommand("pm list packages|grep ", packageName)

	if err != nil {
		return false, err
	}

	return len(out) != 0, nil
}

// UnlockScreen unlocks the screen of the device.
func (d *AndroidDevice) UnlockScreen() error {
	locked, err := d.IsScreenLocked()
	if err != nil {
		return err
	}

	if !locked {
		return nil
	}

	active, err := d.IsScreenActive()

	if err != nil {
		return err
	}

	if !active {
		if err = d.InputKey(KeycodePower); err != nil {
			return err
		}
		d.Sleep()
	}

	if err = d.InputKey(KeycodeMenu); err != nil {
		return err
	}
	d.Sleep()

	if d.screenPassword != "" {
		if err = d.InputText(d.screenPassword); err != nil {
			return err
		}
		d.Sleep()
	}

	return nil
}

// LockScreen locks the screen of the device.
func (d *AndroidDevice) LockScreen() error {
	locked, err := d.IsScreenLocked()
	if err != nil {
		return err
	}

	if locked {
		return nil
	}

	// Press power button to lock screen
	if err = d.InputKey(KeycodePower); err != nil {
		return err
	}

	return nil
}

// IsScreenLocked checks if the screen of the device is locked.
func (d *AndroidDevice) IsScreenLocked() (bool, error) {
	out, err := d.RunShellCommand("dumpsys window | grep mDreamingLockscreen=")
	if err != nil {
		return false, err
	}

	return strings.Contains(out, "mDreamingLockscreen=true"), nil
}

// IsScreenActive checks if the screen of the device is active.
func (d *AndroidDevice) IsScreenActive() (bool, error) {
	out, err := d.RunShellCommand("dumpsys power | grep mWakefulness=")
	if err != nil {
		return false, err
	}

	return strings.Contains(out, "mWakefulness=Awake"), nil
}

// InputText inputs text on the device.
func (d *AndroidDevice) InputText(text string) (err error) {
	_, err = d.RunShellCommand("input", "text", text)
	return
}

// InputKey inputs a key on the device.
func (d *AndroidDevice) InputKey(keyCode int) (err error) {
	_, err = d.RunShellCommand("input", "keyevent", strconv.Itoa(keyCode))
	return
}

// RunShellCommand runs a shell command on the device.
func (d *AndroidDevice) RunShellCommand(cmd string, args ...string) (string, error) {
	return d.adb.RunShellCommand(cmd, args...)
}

// Swipe swipes on the device from one point to another point.
func (d *AndroidDevice) Swipe(x, y, x2, y2 int, duration ...time.Duration) (err error) {
	if len(duration) == 0 {
		duration = []time.Duration{d.swipeDuration}
	}

	dur := int(duration[0] / time.Millisecond)
	_, err = d.RunShellCommand("input", "swipe", strconv.Itoa(x), strconv.Itoa(y), strconv.Itoa(x2),
		strconv.Itoa(y2), strconv.Itoa(dur))
	return
}

// SwipeUp swipes up on the device.
// e.g. SwipeUp(0.7, 0.3, 1000)
func (d *AndroidDevice) SwipeUp(args ...any) error {
	width, height, err := d.ScreenSize()
	if err != nil {
		return err
	}

	r1, r2, duration := d.getSwipeArgs(0.7, 0.3, args...)
	return d.Swipe(width/2, int(float64(height)*r1), width/2, int(float64(height)*r2), duration)
}

// SwipeDown swipes down on the device.
// e.g. SwipeDown(0.3, 0.7, 1000)
func (d *AndroidDevice) SwipeDown(args ...any) error {
	width, height, err := d.ScreenSize()
	if err != nil {
		return err
	}

	r1, r2, duration := d.getSwipeArgs(0.3, 0.7, args...)
	return d.Swipe(width/2, int(float64(height)*r1), width/2, int(float64(height)*r2), duration)
}

// SwipeLeft swipes left on the device.
// e.g. SwipeLeft(0.8, 0.2, 1000)
func (d *AndroidDevice) SwipeLeft(args ...any) error {
	width, height, err := d.ScreenSize()
	if err != nil {
		return err
	}

	r1, r2, duration := d.getSwipeArgs(0.8, 0.2, args...)
	return d.Swipe(int(float64(width)*r1), height/2, int(float64(width)*r2), height/2, duration)
}

// SwipeRight swipes right on the device.
// e.g. SwipeRight(0.2, 0.8, 1000)
func (d *AndroidDevice) SwipeRight(args ...any) error {
	width, height, err := d.ScreenSize()
	if err != nil {
		return err
	}

	r1, r2, duration := d.getSwipeArgs(0.2, 0.8, args...)
	return d.Swipe(int(float64(width)*r1), height/2, int(float64(width)*r2), height/2, duration)
}

// getSwipeArgs returns the swipe arguments.
func (d *AndroidDevice) getSwipeArgs(r1, r2 float64, args ...any) (float64, float64, time.Duration) {
	duration := d.swipeDuration

	if len(args) > 2 {
		duration = args[2].(time.Duration)
		r1, r2 = args[0].(float64), args[1].(float64)
	} else if len(args) > 1 {
		r1, r2 = args[0].(float64), args[1].(float64)
	} else if len(args) > 0 {
		r1, r2 = args[0].(float64), 1-args[0].(float64)
	}

	return r1, r2, duration
}

// ScreenSize returns the screen size of the device.
func (d *AndroidDevice) ScreenSize() (int, int, error) {
	out, err := d.RunShellCommand("wm size")
	if err != nil {
		return 0, 0, err
	}

	// Regular expression to match width and height values in "Physical size: 1080x2400"
	reg := regexp.MustCompile(`Physical size: (\d+)x(\d+)`)
	match := reg.FindStringSubmatch(out)
	if len(match) != 3 {
		return 0, 0, fmt.Errorf("failed to parse screen size: %s", out)
	}

	width, _ := strconv.Atoi(match[1])
	height, _ := strconv.Atoi(match[2])
	return width, height, nil
}

// ScreenDpi returns the screen DPI of the device.
func (d *AndroidDevice) ScreenDpi() (int, error) {
	out, err := d.RunShellCommand("wm density")
	if err != nil {
		return 0, err
	}

	// Regular expression to match DPI value in "Physical density: 440"
	reg := regexp.MustCompile(`Physical density: (\d+)`)
	match := reg.FindStringSubmatch(out)
	if len(match) != 2 {
		return 0, fmt.Errorf("failed to parse screen DPI: %s", out)
	}

	dpi, _ := strconv.Atoi(match[1])
	return dpi, nil
}

// Screenshot takes a screenshot of the device and saves it to the screenshotPath.
func (d *AndroidDevice) Screenshot() (*os.File, error) {
	if d.screenshotPath == "" {
		return nil, fmt.Errorf("screenshot path is not set")
	}

	if err := os.MkdirAll(d.screenshotPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create screenshot directory: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")

	rep := strings.NewReplacer(" ", "_", ":", "_", ".", "_")
	filename := fmt.Sprintf("screenshot_%s_%s.png", rep.Replace(d.id), timestamp)
	localPath := filepath.Join(d.screenshotPath, filename)

	remotePath := path.Join(TempPath, filename)

	_, err := d.RunShellCommand("screencap", "-p", remotePath)
	if err != nil {
		return nil, fmt.Errorf("failed to capture screenshot: %w", err)
	}

	file, err := os.Create(localPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create local file: %w", err)
	}

	if err := d.adb.Pull(remotePath, file); err != nil {
		file.Close()
		os.Remove(localPath)
		return nil, fmt.Errorf("failed to pull screenshot: %w", err)
	}

	_, _ = d.RunShellCommand("rm", remotePath)

	if _, err := file.Seek(0, 0); err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to reset file pointer: %w", err)
	}

	return file, nil
}

// Tap taps on the device at the specified coordinates.
func (d *AndroidDevice) Tap(x, y int) (err error) {
	_, err = d.RunShellCommand("input", "tap", strconv.Itoa(x), strconv.Itoa(y))
	return
}

// LongTap long taps on the device at the specified coordinates.
func (d *AndroidDevice) LongTap(x, y int, duration ...time.Duration) error {
	if len(duration) == 0 {
		duration = []time.Duration{d.longTapDuration}
	}

	return d.Swipe(x, y, x, y, duration[0])
}

// Back presses the back button on the device.
func (d *AndroidDevice) Back(num ...int) error {
	if len(num) == 0 {
		num = []int{1}
	}

	for i := 0; i < num[0]; i++ {
		err := d.InputKey(KeycodeBack)
		if err != nil {
			return err
		}
		d.Sleep()
	}

	return nil
}

// Sleep sleeps for a specified duration.
func (d *AndroidDevice) Sleep(delay ...time.Duration) {
	if len(delay) != 0 {
		time.Sleep(delay[0])
	} else {
		time.Sleep(d.sleepDuration)
	}
}

// SystemInfo returns the system information of the device.
func (d *AndroidDevice) SystemInfo() (info *SystemInfo, err error) {
	info = &SystemInfo{}

	// Get device basic information
	propOutput, err := d.RunShellCommand("getprop")
	if err != nil {
		return nil, fmt.Errorf("failed to get device properties: %w", err)
	}

	// Parse device model
	if match := regexp.MustCompile(`\[ro.product.model\]:\s*\[(.*?)\]`).FindStringSubmatch(propOutput); len(match) > 1 {
		info.Model = match[1]
	}

	// Parse device brand
	if match := regexp.MustCompile(`\[ro.product.brand\]:\s*\[(.*?)\]`).FindStringSubmatch(propOutput); len(match) > 1 {
		info.Brand = match[1]
	}

	// Parse manufacturer
	if match := regexp.MustCompile(`\[ro.product.manufacturer\]:\s*\[(.*?)\]`).FindStringSubmatch(propOutput); len(match) > 1 {
		info.Manufacturer = match[1]
	}

	// Parse Android version
	if match := regexp.MustCompile(`\[ro.build.version.release\]:\s*\[(.*?)\]`).FindStringSubmatch(propOutput); len(match) > 1 {
		info.AndroidVersion = match[1]
	}

	// Parse SDK version
	if match := regexp.MustCompile(`\[ro.build.version.sdk\]:\s*\[(.*?)\]`).FindStringSubmatch(propOutput); len(match) > 1 {
		info.SDK, _ = strconv.Atoi(match[1])
	}

	// Parse serial number
	if match := regexp.MustCompile(`\[ro.serialno\]:\s*\[(.*?)\]`).FindStringSubmatch(propOutput); len(match) > 1 {
		info.SerialNumber = match[1]
	}

	// Get battery information
	batteryOutput, err := d.RunShellCommand("dumpsys battery")
	if err == nil {
		// Parse battery level
		if match := regexp.MustCompile(`level: (\d+)`).FindStringSubmatch(batteryOutput); len(match) > 1 {
			info.Battery, _ = strconv.Atoi(match[1])
		}

		// Parse battery status
		if strings.Contains(batteryOutput, "status: 2") {
			info.BatteryStatus = "Charging"
		} else if strings.Contains(batteryOutput, "status: 5") {
			info.BatteryStatus = "Full"
		} else {
			info.BatteryStatus = "Discharging"
		}
	}

	// Get screen information
	info.ScreenWidth, info.ScreenHeight, _ = d.ScreenSize()
	info.ScreenDensity, _ = d.ScreenDpi()

	// Get network information
	wifiOutput, err := d.RunShellCommand("dumpsys wifi")
	if err == nil {
		// Check if WiFi is connected
		if strings.Contains(wifiOutput, "mNetworkInfo") && strings.Contains(wifiOutput, "state: CONNECTED") {
			info.NetworkType = "WiFi"

			// Get WiFi SSID
			if match := regexp.MustCompile(`"(.*?)"`).FindStringSubmatch(wifiOutput); len(match) > 1 {
				info.WifiSSID = match[1]
			}

			// Get WiFi signal strength
			if match := regexp.MustCompile(`RSSI: (-?\d+)`).FindStringSubmatch(wifiOutput); len(match) > 1 {
				info.WifiSignal, _ = strconv.Atoi(match[1])
			}
		} else {
			// Check mobile data
			dataOutput, err := d.RunShellCommand("dumpsys telephony.registry")
			if err == nil && strings.Contains(dataOutput, "mDataConnectionState=2") {
				info.NetworkType = "Mobile Data"
			} else {
				info.NetworkType = "No Network"
			}
		}
	}

	// Get location service status
	locationOutput, err := d.RunShellCommand("settings get secure location_providers_allowed")
	if err == nil {
		info.LocationEnabled = locationOutput != ""
	}

	// Get IMEI (requires READ_PHONE_STATE permission)
	imeiOutput, err := d.RunShellCommand("service call iphonesubinfo 1")
	if err == nil && !strings.Contains(imeiOutput, "Exception") {
		// Parse IMEI (complex format, needs special handling)
		var imei strings.Builder
		for _, match := range regexp.MustCompile(`'(.*?)'`).FindAllStringSubmatch(imeiOutput, -1) {
			if len(match) > 1 {
				imei.WriteString(match[1])
			}
		}
		info.IMEI = strings.Replace(imei.String(), ".", "", -1)
	}

	// Get memory information
	memOutput, err := d.RunShellCommand("cat /proc/meminfo")
	if err == nil {
		// Parse total memory
		if match := regexp.MustCompile(`MemTotal:\s+(\d+) kB`).FindStringSubmatch(memOutput); len(match) > 1 {
			totalKB, _ := strconv.ParseInt(match[1], 10, 64)
			info.TotalRAM = totalKB * 1024 // Convert to bytes
		}

		// Parse available memory
		if match := regexp.MustCompile(`MemAvailable:\s+(\d+) kB`).FindStringSubmatch(memOutput); len(match) > 1 {
			availableKB, _ := strconv.ParseInt(match[1], 10, 64)
			info.AvailableRAM = availableKB * 1024 // Convert to bytes
		}
	}

	// Get storage information
	storageOutput, err := d.RunShellCommand("df /data")
	if err == nil {
		lines := strings.Split(storageOutput, "\n")
		if len(lines) > 1 {
			fields := strings.Fields(lines[1])
			if len(fields) >= 4 {
				// Total storage
				totalBlocks, _ := strconv.ParseInt(fields[1], 10, 64)
				info.TotalStorage = totalBlocks * 1024 // Convert to bytes

				// Available storage
				availableBlocks, _ := strconv.ParseInt(fields[3], 10, 64)
				info.AvailableStorage = availableBlocks * 1024 // Convert to bytes
			}
		}
	}

	return info, nil
}
