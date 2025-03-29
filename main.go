package main

import (
	"log/slog"
	"mcp-android-adb-server/device"
	"mcp-android-adb-server/tools"
	"os"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	deviceId := os.Getenv("DEVICE_ID")
	if deviceId == "" {
		slog.Error("DEVICE_ID environment variable not set")
		return
	}

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
