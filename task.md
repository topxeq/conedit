# 任务描述

用 Go 语言编写一个函数，接收一些可选参数（如默认文本内容、指定打开的文件等），启动一个简单命令行文本编辑器（如 Nano 等），用户可以在其中编辑，并选择保存或取消，函数返回相应的结果。

# 说明与要求

- 尽量使用 Go 语言标准库中的相关函数和包，避免依赖外部库。
- 函数的定义如下：
```go
func ConsoleEditText(defaultTextA string, optsA ...string) map[string]interface{}
```
- 函数的固定参数包括默认文本内容（defaultTextA），可选参数包括指定打开的文件路径（例如："-filePath=/mnts/text1.txt"）、是否从 SSH 服务器获取文件（例如："-fromSSH"，这是一个开关参数，存在即表示打开）、指定从 SSH 服务器获取文件时的主机名、端口号、用户名、密码或密钥路径、远端文件路径等信息（例如："-sshHost=192.168.1.100"， "-sshPort=22"， "-sshUser=root"， "-sshPass=abc123"， "-sshKeyPath=/mnts/id_rsa"）、指定临时文件路径（"-tmpPath=/tmpx/buf1.tmp"，如不指定则对于大文件在系统临时文件夹下建立临时文件来处理，对于小文件，例如小于 10MB 的文件，直接在内存中处理，可以通过开关参数"-mem"来强制在内存中处理，不使用临时文件）。
- 函数返回一个 map[string]interface{} 对象，包含当前编辑器中的文本内容（"text"键），退出编辑器时的状态（"status"键，包含"save"、"saveAs" 、"cancel"或"error"这几个值），如果 status 为"save"，表示用户选择了保存，此时"text"键为当前编辑器中的文本内容，如果 status 为"saveAs"则表示用户选择了"另存为"，此时"text"键值为当前编辑器中的文本内容；如果 status 为"cancel"，表示用户取消了保存，此时"text"键为空字符串；如果 status 为"error"，则"text"键为空字符串，同时"error"键是包含错误信息的字符串。如果 status 键的值为 save 或 saveAs，还会有一个 path 键值表示文件路径（包括 SSH 情况下的远端文件路径）
- 用户在编辑器中可以的快捷键操作包括：如保存（Ctrl+S）、另存为（Ctrl+K）、复制（Ctrl+C）、粘贴（Ctrl+V）、撤销（Ctrl+Z）、重做（Ctrl+Y）、查找（Ctrl+F）、替换（Ctrl+H）、跳转到行号（Ctrl+G）、退出（Ctrl+X）、强制快捷退出（Ctrl+Q）、切换自动折行（Ctrl+W）等。
- 查找替换式搜索条件支持正则表达式；
- 编辑器为 UTF-8 编码，要支持中文字符和 ASCII 字符混合显示正常。
- 编辑器内要支持中文输入；
- 编辑器内，要支持光标移动、插入、删除、复制、粘贴等操作，光标移动要考虑中文字符和 ASCII 字符的移动时不能出停在半个字符中间的问题。
- 编辑器内，文本默认自动折行，可以通过热键（Ctrl+W）切换是否自动折行；
- 除必须的状态栏、提示栏外，其他区域都要显示编辑的文本内容；也就是编辑区域要撑满剩余空间；
- 状态栏中，要显示当前光标所在的行号和在行中的位置
- 我准备把该函数做成库传到 github.com/topxeq/conedit 下，让别人可以在 Go 语言程序中，作为库来调用这个函数的功能启动文本编辑器；同时，我们完善测试用的主程序，让它可以做为一个可用的功能相对完整的轻量级命令行编辑器，请为此做调整；

# 项目结构

```
conedit/
├── editor/              # 编辑器库包
│   ├── editor.go        # ConsoleEditText 函数和编辑器主逻辑
│   ├── buffer.go        # 文本缓冲区管理
│   ├── command.go       # 命令定义
│   ├── screen.go        # 屏幕渲染和字符宽度计算
│   ├── input.go         # 参数解析
│   ├── sshclient.go     # SSH 客户端
│   └── *_test.go        # 测试文件
├── cmd/editor/          # 命令行主程序
│   └── main.go
├── go.mod               # Go 模块定义
├── README.md            # 使用文档
└── task.md              # 任务描述
```

# 构建命令

```bash
# 构建库
go build ./editor/...

# 构建命令行程序
go build -o console_editor ./cmd/editor/

# 运行测试
go test ./editor/...

# 交叉编译
GOOS=linux GOARCH=amd64 go build -o console_editor_linux ./cmd/editor/
GOOS=windows GOARCH=amd64 go build -o console_editor.exe ./cmd/editor/
```
