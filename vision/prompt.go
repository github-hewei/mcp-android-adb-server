package vision

// SystemPrompt is the system prompt used by the LLM model to analyze the screenshot.
const SystemPrompt = `
你将收到一张手机屏幕截图，请详细分析截图中的UI内容，并将所有元素按照以下要求整理为JSON格式输出。

**输出JSON数据结构：**

- **整体描述**:
类型：字符串
描述：一段文字，整体概括当前页面的内容和功能，包括所属APP类型（如社交、购物）、主要作用，屏幕分辨率，以及是否有弹窗遮挡等信息。

- **元素列表**:
类型：数组
描述：包含多个元素对象，每个对象代表截图中的一个UI元素，字段如下：

- **类型**:
类型：字符串，非空
描述：元素的类型，例如“文本”、“图标”、“按钮”、“链接”、“列表”、“弹窗”等。

- **坐标**:
类型：字符串，非空
描述：元素的中心点坐标，包含x和y两个数字，逗号隔开例如“100,200”。

- **描述**:
类型：字符串，可为空
描述：推测该元素的功能或作用，例如“按钮用于返回上一页”。

- **方位**:
类型：字符串，非空
描述：元素在屏幕上的大致方位，取值范围：上方、下方、左方、右方、左上、右上、左下、右下。

- **文本**:
类型：字符串，可为空
描述：元素的文本内容（若有）。

- **图标**:
类型：字符串，可为空
描述：推测的图标名称，例如“后退”、“前进”、“菜单”等。

- **交互**:
类型：字符串，可为空
描述：推测该元素是否可交互，例如“可点击”、“可上滑”等。

**注意事项：**
- 忽略截图顶部手机状态栏中的信息（如电量、信号等）。
- 仅输出标准JSON格式文本，不要输出额外内容。内容紧凑不要输出换行和空格等空白字符。

**输出示例：**
{"整体描述":"这是一个社交APP的聊天页面，主要用于用户间发送消息，当前无弹窗遮挡。","元素列表":[{"类型":"按钮","坐标":"100,200","描述":"返回上一页","方位":"左上","文本":"","图标":"后退","交互":"可点击"},{"类型":"文本","坐标":"100,200","描述":"显示聊天内容","方位":"下方","文本":"你好，今天怎么样？","图标":"","交互":""}]}
`

// SystemPromptEnglish is the system prompt used by the LLM model to analyze the screenshot in English.
const SystemPromptEnglish = `
You will receive a screenshot of a mobile phone screen. Please analyze the UI content of the screenshot in detail and organize all elements into JSON format output according to the following requirements.

**Output JSON Data Structure:**

- **Overall Description**:
Type: String
Description: A text that summarizes the content and functionality of the current page, including the APP type (such as social, shopping), main purpose, screen resolution, and whether there are pop-up windows obstructing, etc.

- **Element List**:
Type: Array
Description: Contains multiple element objects, each object represents a UI element in the screenshot, with the following fields:

- **Type**:
  Type: String, Non-empty
  Description: The type of element, such as "text", "icon", "button", "link", "list", "popup", etc.

- **Coordinates**:
  Type: String, Non-empty
  Description: The central point coordinates of the element, including two numbers x and y, separated by a comma, for example, "100,200".

- **Description**:
  Type: String, Nullable
  Description: The inferred function or purpose of the element, such as "button for returning to the previous page".

- **Position**:
  Type: String, Non-empty
  Description: The approximate position of the element on the screen, the value range: top, bottom, left, right, top left, top right, bottom left, bottom right.

- **Text**:
  Type: String, Nullable
  Description: The text content of the element (if any).

- **Icon**:
  Type: String, Nullable
  Description: The inferred icon name, such as "back", "forward", "menu", etc.

- **Interaction**:
  Type: String, Nullable
  Description: Infer whether the element is interactive, such as "clickable", "swipeable", etc.

**Precautions:**

- Ignore the information in the mobile phone status bar at the top of the screenshot (such as battery, signal, etc.).
- Only output standard JSON format text, do not output additional content. The content is compact, do not output line breaks and spaces and other blank characters.

**Output Example:**
{"Overall Description": "This is a chat page of a social APP, mainly used for sending messages between users, there is no pop-up window obstructing at the moment.","Element List": [{"Type": "button","Coordinates": "100,200","Description": "return to previous page","Position": "top left","Text": "","Icon": "back","Interaction": "clickable"},{"Type": "text","Coordinates": "100,200","Description": "show chat content","Position": "bottom","Text": "Hello, how are you today?","Icon": "","Interaction": ""}]}
`
