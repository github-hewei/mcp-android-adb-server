[![English](https://img.shields.io/badge/Language-English-blue.svg)](./README.md)
[![简体中文](https://img.shields.io/badge/语言-简体中文-red.svg)](./README.zh-CN.md)

## 🚀 mcp-android-adb-server

[![smithery badge](https://smithery.ai/badge/@github-hewei/mcp-android-adb-server)](https://smithery.ai/server/@github-hewei/mcp-android-adb-server)

An MCP service for operating Android devices via ADB.

2025-04-01: Added support for obtaining screen description content using visual models like `qwen2.5-vl`.

### Manual Installation

```sh
# clone the repo and build
git clone https://github.com/github-hewei/mcp-android-adb-server.git
cd mcp-android-adb-server
go build
```

### Configuration

```json
{
  "mcpServers": {
    "mcp-android-adb-server": {
      "command": "D:\\www\\golang\\mcp-android-adb-server\\mcp-android-adb-server.exe",
      "env": {
        "DEVICE_ID": "xxxxx",
        "SCREEN_LOCK_PASSWORD": "123456",
        "VISUAL_MODEL_ON": "true",
        "VISUAL_MODEL_API_KEY": "sk-or-xxxxxxxxxxxxxxxxxxx",
        "VISUAL_MODEL_BASE_URL": "https://openrouter.ai/api/v1/",
        "VISUAL_MODEL_NAME": "qwen/qwen2.5-vl-72b-instruct:free"
      }
    }
  }
}
```

### Environment Variables

- DEVICE_ID : Required. The ID of the Android device, obtainable via the `adb devices` command.
- SCREEN_LOCK_PASSWORD : Optional. The screen lock password of the device, used to unlock the screen.
- VISUAL_MODEL_ON : Optional. Whether to enable the visual model, defaults to false.
- VISUAL_MODEL_API_KEY : API Key.
- VISUAL_MODEL_BASE_URL : API Base URL.
- VISUAL_MODEL_NAME : Model name.

### Features and Tools

Application Management
- install_app : Install an application on the Android device
- uninstall_app : Uninstall an application from the Android device
- terminate_app : Terminate a running application on the Android device
- launch_app : Launch an application on the Android device
- list_app : List all installed applications on the Android device
- is_app_installed : Check if a specific application is installed

Screen Control
- unlock_screen : Unlock the Android device screen
- lock_screen : Lock the Android device screen
- is_screen_locked : Check if the Android device screen is locked
- is_screen_active : Check if the Android device screen is active

Input Control

- input_text : Input text on the Android device
- input_key : Input key press on the Android device
- tap : Perform a tap operation on the screen at a specified position
- long_tap : Perform a long press operation on the screen at a specified position
- back : Perform a back operation

Gesture Control

- swipe_up : Perform a swipe up gesture on the Android device screen
- swipe_down : Perform a swipe down gesture on the Android device screen
- swipe_left : Perform a swipe left gesture on the Android device screen
- swipe_right : Perform a swipe right gesture on the Android device screen

Device Information

- screen_size : Get the screen size of the Android device
- screen_dpi : Get the screen DPI of the Android device
- screenshot_description : Get the Android device screenshot description
- system_info : Get system information of the Android device

Other Functions
- shell_command : Execute a shell command on the Android device

