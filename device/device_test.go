package device_test

import (
	"mcp-android-adb-server/device"
	"os"
	"testing"
	"time"
)

// Default test device ID, used if not specified in environment variables
const defaultDeviceID = "E6EDU20723063683"

// TestNewAndroidDevice tests creating a new AndroidDevice instance
func TestNewAndroidDevice(t *testing.T) {
	// Get device ID from environment variable
	deviceID := os.Getenv("TEST_DEVICE_ID")
	if deviceID == "" {
		deviceID = defaultDeviceID // Use default test ID
	}

	d, err := device.NewAndroidDevice(deviceID)
	if err != nil {
		t.Fatalf("Failed to create AndroidDevice: %v", err)
	}

	if d == nil {
		t.Fatal("Created AndroidDevice instance is nil")
	}
}

// TestWithOptions tests device option functionality
func TestWithOptions(t *testing.T) {
	deviceID := os.Getenv("TEST_DEVICE_ID")
	if deviceID == "" {
		deviceID = defaultDeviceID
	}

	customSwipeDuration := time.Millisecond * 800
	customLongTapDuration := time.Second * 3
	customSleepDuration := time.Millisecond * 1500
	customScreenshotPath := "custom_screenshots"
	customPassword := "1234"

	d, err := device.NewAndroidDevice(deviceID,
		device.WithSwipeDuration(customSwipeDuration),
		device.WithLongTapDuration(customLongTapDuration),
		device.WithSleepDuration(customSleepDuration),
		device.WithScreenshotPath(customScreenshotPath),
		device.WithScreenPassword(customPassword),
	)

	if err != nil {
		t.Fatalf("Failed to create AndroidDevice with options: %v", err)
	}

	if d == nil {
		t.Fatal("Created AndroidDevice instance is nil")
	}

	// Note: Since AndroidDevice fields are private, we cannot directly verify if options took effect
	// We can only verify that the creation process was successful
}

// TestScreenSize tests getting screen size
func TestScreenSize(t *testing.T) {
	deviceID := os.Getenv("TEST_DEVICE_ID")
	if deviceID == "" {
		deviceID = defaultDeviceID
	}

	d, err := device.NewAndroidDevice(deviceID)
	if err != nil {
		t.Fatalf("Failed to create AndroidDevice: %v", err)
	}

	width, height, err := d.ScreenSize()
	if err != nil {
		t.Fatalf("Failed to get screen size: %v", err)
	}

	if width <= 0 || height <= 0 {
		t.Errorf("Invalid screen size: %dx%d", width, height)
	}

	t.Logf("Device screen size: %dx%d", width, height)
}

// TestScreenDpi tests getting screen DPI
func TestScreenDpi(t *testing.T) {
	deviceID := os.Getenv("TEST_DEVICE_ID")
	if deviceID == "" {
		deviceID = defaultDeviceID
	}

	d, err := device.NewAndroidDevice(deviceID)
	if err != nil {
		t.Fatalf("Failed to create AndroidDevice: %v", err)
	}

	dpi, err := d.ScreenDpi()
	if err != nil {
		t.Fatalf("Failed to get screen DPI: %v", err)
	}

	if dpi <= 0 {
		t.Errorf("Invalid screen DPI: %d", dpi)
	}

	t.Logf("Device screen DPI: %d", dpi)
}

// TestListApp tests getting installed application list
func TestListApp(t *testing.T) {
	deviceID := os.Getenv("TEST_DEVICE_ID")
	if deviceID == "" {
		deviceID = defaultDeviceID
	}

	d, err := device.NewAndroidDevice(deviceID)
	if err != nil {
		t.Fatalf("Failed to create AndroidDevice: %v", err)
	}

	apps, err := d.ListApp()
	if err != nil {
		t.Fatalf("Failed to get application list: %v", err)
	}

	if len(apps) == 0 {
		t.Error("Application list is empty, this might be an error")
	}

	t.Logf("Found %d installed applications", len(apps))
	if len(apps) > 0 {
		t.Logf("First application package name: %s", apps[0].PackageName)
	}
}

// TestIsScreenLocked tests checking if screen is locked
func TestIsScreenLocked(t *testing.T) {
	deviceID := os.Getenv("TEST_DEVICE_ID")
	if deviceID == "" {
		deviceID = defaultDeviceID
	}

	d, err := device.NewAndroidDevice(deviceID)
	if err != nil {
		t.Fatalf("Failed to create AndroidDevice: %v", err)
	}

	locked, err := d.IsScreenLocked()
	if err != nil {
		t.Fatalf("Failed to check screen lock status: %v", err)
	}

	t.Logf("Screen lock status: %v", locked)
}

// TestSystemInfo tests getting system information
func TestSystemInfo(t *testing.T) {
	deviceID := os.Getenv("TEST_DEVICE_ID")
	if deviceID == "" {
		deviceID = defaultDeviceID
	}

	d, err := device.NewAndroidDevice(deviceID)
	if err != nil {
		t.Fatalf("Failed to create AndroidDevice: %v", err)
	}

	info, err := d.SystemInfo()
	if err != nil {
		t.Fatalf("Failed to get system information: %v", err)
	}

	if info == nil {
		t.Fatal("System information is nil")
	}

	t.Logf("Device model: %s", info.Model)
	t.Logf("Brand: %s", info.Brand)
	t.Logf("Android version: %s", info.AndroidVersion)
	t.Logf("SDK version: %d", info.SDK)
	t.Logf("Screen resolution: %dx%d", info.ScreenWidth, info.ScreenHeight)
}

// TestRunShellCommand tests running shell commands
func TestRunShellCommand(t *testing.T) {
	deviceID := os.Getenv("TEST_DEVICE_ID")
	if deviceID == "" {
		deviceID = defaultDeviceID
	}

	d, err := device.NewAndroidDevice(deviceID)
	if err != nil {
		t.Fatalf("Failed to create AndroidDevice: %v", err)
	}

	output, err := d.RunShellCommand("echo", "hello world")
	if err != nil {
		t.Fatalf("Failed to run shell command: %v", err)
	}

	expected := "hello world\n"
	if output != expected {
		t.Errorf("Shell command output does not match expected: expected %q, got %q", expected, output)
	}
}
