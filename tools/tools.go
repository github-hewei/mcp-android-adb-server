package tools

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"mcp-android-adb-server/device"
	"os"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddToolInstallApp adds a tool for installing applications
func AddToolInstallApp(s *server.MCPServer, d *device.AndroidDevice) {
	s.AddTool(mcp.NewTool("install_app",
		mcp.WithDescription("Install an application on the Android device"),
		mcp.WithString("file",
			mcp.Required(),
			mcp.Description("Path to the application package file with .apk extension"),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		file := request.Params.Arguments["file"].(string)

		if err := d.InstallApp(file, true); err != nil {
			return nil, fmt.Errorf("installation failed: %w", err)
		}

		return mcp.NewToolResultText("Installation successful"), nil
	})
}

// AddToolUninstallApp adds a tool for uninstalling applications
func AddToolUninstallApp(s *server.MCPServer, d *device.AndroidDevice) {
	s.AddTool(mcp.NewTool("uninstall_app",
		mcp.WithDescription("Uninstall an application from the Android device"),
		mcp.WithString("package_name",
			mcp.Required(),
			mcp.Description("Android application package name, e.g. com.example.app"),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		packageName := request.Params.Arguments["package_name"].(string)

		if err := d.UninstallApp(packageName); err != nil {
			return nil, fmt.Errorf("uninstallation failed: %w", err)
		}

		return mcp.NewToolResultText("Uninstallation successful"), nil
	})
}

// AddToolTerminateApp adds a tool for terminating running applications
func AddToolTerminateApp(s *server.MCPServer, d *device.AndroidDevice) {
	s.AddTool(mcp.NewTool("terminate_app",
		mcp.WithDescription("Terminate a running application on the Android device"),
		mcp.WithString("package_name",
			mcp.Required(),
			mcp.Description("Android application package name, e.g. com.example.app"),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		packageName := request.Params.Arguments["package_name"].(string)

		if err := d.TerminateApp(packageName); err != nil {
			return nil, fmt.Errorf("failed to terminate application: %w", err)
		}

		return mcp.NewToolResultText("Application terminated"), nil
	})
}

// AddToolLaunchApp adds a tool for launching applications
func AddToolLaunchApp(s *server.MCPServer, d *device.AndroidDevice) {
	s.AddTool(mcp.NewTool("launch_app",
		mcp.WithDescription("Launch an application on the Android device"),
		mcp.WithString("package_name",
			mcp.Required(),
			mcp.Description("Android application package name, e.g. com.example.app"),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		packageName := request.Params.Arguments["package_name"].(string)

		if err := d.LaunchApp(packageName); err != nil {
			return nil, fmt.Errorf("failed to launch application: %w", err)
		}

		return mcp.NewToolResultText("Application launched"), nil
	})
}

// AddToolListApp adds a tool for listing installed applications
func AddToolListApp(s *server.MCPServer, d *device.AndroidDevice) {
	s.AddTool(mcp.NewTool("list_app",
		mcp.WithDescription("List all installed applications on the Android device"),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		apps, err := d.ListApp()
		if err != nil {
			return nil, fmt.Errorf("failed to get application list: %w", err)
		}

		if len(apps) == 0 {
			return mcp.NewToolResultText("No installed applications found on the device"), nil
		}

		// Build application list text
		var result strings.Builder
		result.WriteString("List of installed applications on the device:\n\n")

		for i, app := range apps {
			result.WriteString(fmt.Sprintf("%d. %s\n", i+1, app.PackageName))
		}

		return mcp.NewToolResultText(result.String()), nil
	})
}

// AddToolInstalledApp adds a tool for checking if an application is installed
func AddToolInstalledApp(s *server.MCPServer, d *device.AndroidDevice) {
	s.AddTool(mcp.NewTool("is_app_installed",
		mcp.WithDescription("Check if a specific application is installed on the Android device"),
		mcp.WithString("package_name",
			mcp.Required(),
			mcp.Description("Android application package name, e.g. com.example.app"),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		packageName := request.Params.Arguments["package_name"].(string)

		installed, err := d.InstalledApp(packageName)
		if err != nil {
			return nil, fmt.Errorf("failed to check application installation status: %w", err)
		}

		if installed {
			return mcp.NewToolResultText(fmt.Sprintf("Application %s is installed", packageName)), nil
		} else {
			return mcp.NewToolResultText(fmt.Sprintf("Application %s is not installed", packageName)), nil
		}
	})
}

// AddToolUnlockScreen adds a tool for unlocking the device screen
func AddToolUnlockScreen(s *server.MCPServer, d *device.AndroidDevice) {
	s.AddTool(mcp.NewTool("unlock_screen",
		mcp.WithDescription("Unlock the Android device screen"),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := d.UnlockScreen(); err != nil {
			return nil, fmt.Errorf("failed to unlock screen: %w", err)
		}

		return mcp.NewToolResultText("Screen unlocked"), nil
	})
}

// AddToolLockScreen adds a tool for locking the device screen
func AddToolLockScreen(s *server.MCPServer, d *device.AndroidDevice) {
	s.AddTool(mcp.NewTool("lock_screen",
		mcp.WithDescription("Lock the Android device screen"),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := d.LockScreen(); err != nil {
			return nil, fmt.Errorf("failed to lock screen: %w", err)
		}

		return mcp.NewToolResultText("Screen locked"), nil
	})
}

// AddToolIsScreenLocked adds a tool for checking if the screen is locked
func AddToolIsScreenLocked(s *server.MCPServer, d *device.AndroidDevice) {
	s.AddTool(mcp.NewTool("is_screen_locked",
		mcp.WithDescription("Check if the Android device screen is locked"),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		locked, err := d.IsScreenLocked()
		if err != nil {
			return nil, fmt.Errorf("failed to check screen lock status: %w", err)
		}

		if locked {
			return mcp.NewToolResultText("Screen is currently locked"), nil
		} else {
			return mcp.NewToolResultText("Screen is currently unlocked"), nil
		}
	})
}

// AddToolIsScreenActive adds a tool for checking if the screen is active
func AddToolIsScreenActive(s *server.MCPServer, d *device.AndroidDevice) {
	s.AddTool(mcp.NewTool("is_screen_active",
		mcp.WithDescription("Check if the Android device screen is active"),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		active, err := d.IsScreenActive()
		if err != nil {
			return nil, fmt.Errorf("failed to check screen active status: %w", err)
		}

		if active {
			return mcp.NewToolResultText("Screen is currently active"), nil
		} else {
			return mcp.NewToolResultText("Screen is currently inactive"), nil
		}
	})
}

// AddToolInputText adds a tool for inputting text
func AddToolInputText(s *server.MCPServer, d *device.AndroidDevice) {
	s.AddTool(mcp.NewTool("input_text",
		mcp.WithDescription("Input text on the Android device"),
		mcp.WithString("text",
			mcp.Required(),
			mcp.Description("Text content to input"),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		text := request.Params.Arguments["text"].(string)

		if err := d.InputText(text); err != nil {
			return nil, fmt.Errorf("failed to input text: %w", err)
		}

		return mcp.NewToolResultText("Text input successful"), nil
	})
}

// AddToolInputKey adds a tool for inputting key presses
func AddToolInputKey(s *server.MCPServer, d *device.AndroidDevice) {
	s.AddTool(mcp.NewTool("input_key",
		mcp.WithDescription("Input key press on the Android device"),
		mcp.WithNumber("key_code",
			mcp.Required(),
			mcp.Description("Key code to input, e.g. 3 for Home key, 4 for Back key"),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		keyCode := int(request.Params.Arguments["key_code"].(float64))

		if err := d.InputKey(keyCode); err != nil {
			return nil, fmt.Errorf("failed to input key: %w", err)
		}

		return mcp.NewToolResultText(fmt.Sprintf("Key code %d input successful", keyCode)), nil
	})
}

// AddToolShellCommand adds a tool for executing shell commands
func AddToolShellCommand(s *server.MCPServer, d *device.AndroidDevice) {
	s.AddTool(mcp.NewTool("shell_command",
		mcp.WithDescription("Execute a shell command on the Android device"),
		mcp.WithString("command",
			mcp.Required(),
			mcp.Description("Shell command to execute, input the part after 'shell:', e.g. 'ls -l'"),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		command := request.Params.Arguments["command"].(string)

		output, err := d.RunShellCommand(command)
		if err != nil {
			return nil, fmt.Errorf("failed to execute shell command: %w", err)
		}

		return mcp.NewToolResultText(fmt.Sprintf("Command execution result:\n%s", output)), nil
	})
}

// AddToolSwipeUp adds a tool for swiping up on the screen
func AddToolSwipeUp(s *server.MCPServer, d *device.AndroidDevice) {
	s.AddTool(mcp.NewTool("swipe_up",
		mcp.WithDescription("Perform a swipe up gesture on the Android device screen"),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := d.SwipeUp(); err != nil {
			return nil, fmt.Errorf("failed to swipe up: %w", err)
		}

		return mcp.NewToolResultText("Swipe up gesture performed"), nil
	})
}

// AddToolSwipeDown adds a tool for swiping down on the screen
func AddToolSwipeDown(s *server.MCPServer, d *device.AndroidDevice) {
	s.AddTool(mcp.NewTool("swipe_down",
		mcp.WithDescription("Perform a swipe down gesture on the Android device screen"),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := d.SwipeDown(); err != nil {
			return nil, fmt.Errorf("failed to swipe down: %w", err)
		}

		return mcp.NewToolResultText("Swipe down gesture performed"), nil
	})
}

// AddToolSwipeLeft adds a tool for swiping left on the screen
func AddToolSwipeLeft(s *server.MCPServer, d *device.AndroidDevice) {
	s.AddTool(mcp.NewTool("swipe_left",
		mcp.WithDescription("Perform a swipe left gesture on the Android device screen"),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := d.SwipeLeft(); err != nil {
			return nil, fmt.Errorf("failed to swipe left: %w", err)
		}

		return mcp.NewToolResultText("Swipe left gesture performed"), nil
	})
}

// AddToolSwipeRight adds a tool for swiping right on the screen
func AddToolSwipeRight(s *server.MCPServer, d *device.AndroidDevice) {
	s.AddTool(mcp.NewTool("swipe_right",
		mcp.WithDescription("Perform a swipe right gesture on the Android device screen"),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := d.SwipeRight(); err != nil {
			return nil, fmt.Errorf("failed to swipe right: %w", err)
		}

		return mcp.NewToolResultText("Swipe right gesture performed"), nil
	})
}

// AddToolScreenSize adds a tool for getting screen size information
func AddToolScreenSize(s *server.MCPServer, d *device.AndroidDevice) {
	s.AddTool(mcp.NewTool("screen_size",
		mcp.WithDescription("Get the screen size information of the Android device"),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		width, height, err := d.ScreenSize()
		if err != nil {
			return nil, fmt.Errorf("failed to get screen size: %w", err)
		}

		return mcp.NewToolResultText(fmt.Sprintf("Screen size: %dx%d", width, height)), nil
	})
}

// AddToolScreenDpi adds a tool for getting screen DPI information
func AddToolScreenDpi(s *server.MCPServer, d *device.AndroidDevice) {
	s.AddTool(mcp.NewTool("screen_dpi",
		mcp.WithDescription("Get the screen DPI information of the Android device"),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		dpi, err := d.ScreenDpi()
		if err != nil {
			return nil, fmt.Errorf("failed to get screen DPI: %w", err)
		}

		return mcp.NewToolResultText(fmt.Sprintf("Screen DPI: %d", dpi)), nil
	})
}

// AddToolScreenshot adds a tool for taking screenshots
func AddToolScreenshot(s *server.MCPServer, d *device.AndroidDevice) {
	s.AddTool(mcp.NewTool("screenshot",
		mcp.WithDescription("Take a screenshot of the Android device screen to analyze operations and verify goals"),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		file, err := d.Screenshot()
		if err != nil {
			return nil, fmt.Errorf("failed to take screenshot: %w", err)
		}
		defer file.Close()

		imageBytes, err := os.ReadFile(file.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to read screenshot file: %w", err)
		}

		imageBase64 := base64.StdEncoding.EncodeToString(imageBytes)

		text := "Android device screenshot"
		return mcp.NewToolResultImage(text, imageBase64, "image/png"), nil
	})
}

// AddToolTap adds a tool for tapping on the screen
func AddToolTap(s *server.MCPServer, d *device.AndroidDevice) {
	s.AddTool(mcp.NewTool("tap",
		mcp.WithDescription("Perform a tap operation on the Android device screen"),
		mcp.WithNumber("x",
			mcp.Required(),
			mcp.Description("X coordinate of the tap position"),
		),
		mcp.WithNumber("y",
			mcp.Required(),
			mcp.Description("Y coordinate of the tap position"),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		x := int(request.Params.Arguments["x"].(float64))
		y := int(request.Params.Arguments["y"].(float64))

		if err := d.Tap(x, y); err != nil {
			return nil, fmt.Errorf("failed to perform tap operation: %w", err)
		}

		return mcp.NewToolResultText(fmt.Sprintf("Tap operation performed at coordinates (%d, %d)", x, y)), nil
	})
}

// AddToolLongTap adds a tool for long-pressing on the screen
func AddToolLongTap(s *server.MCPServer, d *device.AndroidDevice) {
	s.AddTool(mcp.NewTool("long_tap",
		mcp.WithDescription("Perform a long press operation on the Android device screen"),
		mcp.WithNumber("x",
			mcp.Required(),
			mcp.Description("X coordinate of the long press position"),
		),
		mcp.WithNumber("y",
			mcp.Required(),
			mcp.Description("Y coordinate of the long press position"),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		x := int(request.Params.Arguments["x"].(float64))
		y := int(request.Params.Arguments["y"].(float64))

		if err := d.LongTap(x, y); err != nil {
			return nil, fmt.Errorf("failed to perform long press operation: %w", err)
		}

		return mcp.NewToolResultText(fmt.Sprintf("Long press operation performed at coordinates (%d, %d)", x, y)), nil
	})
}

// AddToolBack adds a tool for performing back operations
func AddToolBack(s *server.MCPServer, d *device.AndroidDevice) {
	s.AddTool(mcp.NewTool("back",
		mcp.WithDescription("Perform a back operation on the Android device"),
		mcp.WithNumber("steps",
			mcp.DefaultNumber(1),
			mcp.Description("Number of back steps, default is 1 step, e.g. 2 means go back twice"),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		steps := int(request.Params.Arguments["steps"].(float64))

		if err := d.Back(steps); err != nil {
			return nil, fmt.Errorf("failed to perform back operation: %w", err)
		}

		return mcp.NewToolResultText(fmt.Sprintf("Performed %d back operation(s)", steps)), nil
	})
}

// AddToolSystemInfo adds a tool for getting device system information
func AddToolSystemInfo(s *server.MCPServer, d *device.AndroidDevice) {
	s.AddTool(mcp.NewTool("system_info",
		mcp.WithDescription("Get system information of the Android device including model, brand, Android version, etc."),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		info, err := d.SystemInfo()
		if err != nil {
			return nil, fmt.Errorf("failed to get system information: %w", err)
		}

		jsonString, err := json.Marshal(info)
		if err != nil {
			return nil, fmt.Errorf("failed to convert system information to JSON: %w", err)
		}

		return mcp.NewToolResultText(string(jsonString)), nil
	})
}
