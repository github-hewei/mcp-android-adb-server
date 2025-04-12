package main

import (
	"io"
	"log/slog"
	"mcp-android-adb-server/device"
	"mcp-android-adb-server/tools"
	"mcp-android-adb-server/vision"
	"os"
	"path"

	"github.com/mark3labs/mcp-go/mcp"

	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/mark3labs/mcp-go/server"
)

func init() {
	baseDir := getBaseDir()
	rotateWriter := &lumberjack.Logger{
		Filename:   path.Join(baseDir, "mcp-android-adb-server.log"),
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
	wd, _ := os.Getwd()
	slog.Info("server start", "directory", wd)

	deviceId := os.Getenv("DEVICE_ID")
	screenLockPassword := os.Getenv("SCREEN_LOCK_PASSWORD")

	d, err := device.NewAndroidDevice(
		deviceId,
		device.WithScreenPassword(screenLockPassword),
		device.WithScreenshotPath(path.Join(getBaseDir(), "screenshots")))

	if err != nil {
		slog.Error("error connect android device", "error", err)
		return
	}

	s := server.NewMCPServer(
		"mcp-android-adb-server",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
		server.WithHooks(getHooks()),
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
		//tools.AddToolScreenshot,
		tools.AddToolTap,
		tools.AddToolLongTap,
		tools.AddToolBack,
		tools.AddToolSystemInfo,
	}

	// Register all tools
	for _, registerTool := range toolList {
		registerTool(s, d)
	}

	// Register visual tools
	visualModel := os.Getenv("VISUAL_MODEL_ON")

	if visualModel == "true" {
		visualModelApiKey := os.Getenv("VISUAL_MODEL_API_KEY")
		visualModelBaseUrl := os.Getenv("VISUAL_MODEL_BASE_URL")
		visualModelName := os.Getenv("VISUAL_MODEL_NAME")
		m := vision.NewModel(visualModelApiKey, visualModelName, visualModelBaseUrl)
		tools.AddToolScreenshotDescription(s, d, m)
	}
}

// getHooks returns all hooks
func getHooks() *server.Hooks {
	hooks := &server.Hooks{}

	hooks.AddBeforeAny(func(id any, method mcp.MCPMethod, message any) {
		slog.Info("before any hook called",
			"method", method,
			"id", id,
			"message", message)
	})

	hooks.AddOnSuccess(func(id any, method mcp.MCPMethod, message any, result any) {
		slog.Info("operation completed successfully",
			"method", method,
			"id", id,
			"message", message,
			"result", result)
	})

	hooks.AddOnError(func(id any, method mcp.MCPMethod, message any, err error) {
		slog.Error("operation failed",
			"method", method,
			"id", id,
			"message", message,
			"error", err)
	})

	hooks.AddBeforeInitialize(func(id any, message *mcp.InitializeRequest) {
		slog.Info("initializing",
			"id", id,
			"message", message)
	})

	hooks.AddAfterInitialize(func(id any, message *mcp.InitializeRequest, result *mcp.InitializeResult) {
		slog.Info("initialization completed",
			"id", id,
			"message", message,
			"result", result)
	})

	hooks.AddAfterCallTool(func(id any, message *mcp.CallToolRequest, result *mcp.CallToolResult) {
		slog.Info("tool call completed",
			"id", id,
			"message", message,
			"result", result)
	})

	hooks.AddBeforeCallTool(func(id any, message *mcp.CallToolRequest) {
		slog.Info("calling tool",
			"id", id,
			"message", message)
	})

	return hooks
}

// getBaseDir returns the base directory
func getBaseDir() string {
	baseDir, _ := os.UserHomeDir()
	if baseDir == "" {
		baseDir = os.TempDir()
	}

	return path.Join(baseDir, "mcp-android-adb-server")
}
