package main

import (
	"io"
	"log/slog"
	"mcp-android-adb-server/device"
	"mcp-android-adb-server/tools"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/mark3labs/mcp-go/server"
)

func init() {
	rotateWriter := &lumberjack.Logger{
		Filename:   "mcp-android-adb.log",
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}

	logger := slog.New(slog.NewJSONHandler(io.MultiWriter(os.Stdout, rotateWriter), &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	slog.SetDefault(logger)
}

func main() {
	deviceId := os.Getenv("DEVICE_ID")
	screenLockPassword := os.Getenv("SCREEN_LOCK_PASSWORD")

	d, err := device.NewAndroidDevice(deviceId, device.WithScreenPassword(screenLockPassword))

	if err != nil {
		slog.Error("error connect android device", "error", err)
		return
	}

	s := server.NewMCPServer(
		"mcp-android-adb-server",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)

	// Register all tools
	registerTools(s, d)

	if err := server.ServeStdio(s); err != nil {
		slog.Error("error serving stdio", "error", err)
	}
}

// registerTools registers all Android device tools
func registerTools(s *server.MCPServer, d *device.AndroidDevice) {
	// Define all tool registration functions
	toolList := []func(*server.MCPServer, *device.AndroidDevice){
		tools.AddToolInstallApp,
		tools.AddToolUninstallApp,
		tools.AddToolTerminateApp,
		tools.AddToolLaunchApp,
		tools.AddToolListApp,
		tools.AddToolInstalledApp,
		tools.AddToolUnlockScreen,
		tools.AddToolLockScreen,
		tools.AddToolIsScreenLocked,
		tools.AddToolIsScreenActive,
		tools.AddToolInputText,
		tools.AddToolInputKey,
		tools.AddToolShellCommand,
		tools.AddToolSwipeUp,
		tools.AddToolSwipeDown,
		tools.AddToolSwipeLeft,
		tools.AddToolSwipeRight,
		tools.AddToolScreenSize,
		tools.AddToolScreenDpi,
		tools.AddToolScreenshot,
		tools.AddToolTap,
		tools.AddToolLongTap,
		tools.AddToolBack,
		tools.AddToolSystemInfo,
	}

	// Register all tools
	for _, registerTool := range toolList {
		registerTool(s, d)
	}
}
