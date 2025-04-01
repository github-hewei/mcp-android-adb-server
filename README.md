## 🚀 mcp-android-adb-server

[![smithery badge](https://smithery.ai/badge/@github-hewei/mcp-android-adb-server)](https://smithery.ai/server/@github-hewei/mcp-android-adb-server)

一个MCP服务用于通过adb操作安卓设备。

2025-04-01：增加支持用视觉模型如`qwen2.5-vl`获取屏幕描述内容。

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

### 环境变量

- DEVICE_ID : 必需。Android 设备的 ID，可以通过 adb devices 命令获取。
- SCREEN_LOCK_PASSWORD : 可选。设备的屏幕锁定密码，用于解锁屏幕。
- VISUAL_MODEL_ON : 可选。是否启用视觉模型，默认为 false。
- VISUAL_MODEL_API_KEY : API密钥。
- VISUAL_MODEL_BASE_URL : API BaseURL。
- VISUAL_MODEL_NAME : 模型名称。

### 功能和工具

应用管理
- install_app : 在 Android 设备上安装应用程序
- uninstall_app : 从 Android 设备卸载应用程序
- terminate_app : 终止 Android 设备上运行的应用程序
- launch_app : 启动 Android 设备上的应用程序
- list_app : 列出 Android 设备上安装的所有应用程序
- is_app_installed : 检查特定应用程序是否已安装

屏幕控制
- unlock_screen : 解锁 Android 设备屏幕
- lock_screen : 锁定 Android 设备屏幕
- is_screen_locked : 检查 Android 设备屏幕是否锁定
- is_screen_active : 检查 Android 设备屏幕是否活跃

输入控制

- input_text : 在 Android 设备上输入文本
- input_key : 在 Android 设备上输入按键
- tap : 在屏幕上点击指定位置
- long_tap : 在屏幕上长按指定位置
- back : 执行返回操作

手势控制

- swipe_up : 在 Android 设备屏幕上执行向上滑动手势
- swipe_down : 在 Android 设备屏幕上执行向下滑动手势
- swipe_left : 在 Android 设备屏幕上执行向左滑动手势
- swipe_right : 在 Android 设备屏幕上执行向右滑动手势

设备信息

- screen_size : 获取 Android 设备屏幕尺寸
- screen_dpi : 获取 Android 设备屏幕 DPI
- screenshot_description : 获取 Android 设备屏幕截图描述
- system_info : 获取 Android 设备系统信息

其他功能
- shell_command : 在 Android 设备上执行 shell 命令
