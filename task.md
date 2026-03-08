# 任务描述

用Go语言编写一个函数，接收一些可选参数（如默认文本内容、指定打开的文件等），启动一个简单命令行文本编辑器（如Nano等），用户可以在其中编辑，并选择保存或取消，函数返回相应的结果。

# 说明与要求

- 尽量使用Go语言标准库中的相关函数和包，避免依赖外部库。
- 函数的定义如下：
```go
func ConsoleEditText(defaultTextA string, optsA ...string) map[string]interface{}
```
- 函数的固定参数包括默认文本内容（defaultTextA），可选参数包括指定打开的文件路径（例如："-filePath=/mnts/text1.txt"）、是否从SSH服务器获取文件（例如："-fromSSH"，这是一个开关参数，存在即表示打开）、指定从SSH服务器获取文件时的主机名、端口号、用户名、密码或密钥路径、远端文件路径等信息（例如："-sshHost=192.168.1.100"， "-sshPort=22"， "-sshUser=root"， "-sshPass=abc123"， "-sshKeyPath=/mnts/id_rsa"）、指定临时文件路径（"-tmpPath=/tmpx/buf1.tmp"，如不指定则对于大文件在系统临时文件夹下建立临时文件来处理，对于小文件，例如小于10MB的文件，直接在内存中处理，可以通过开关参数"-mem"来强制在内存中处理，不使用临时文件）。
- 函数返回一个map[string]interface{}对象，包含当前编辑器中的文本内容（"text"键），退出编辑器时的状态（"status"键，包含"save"、"saveAs" 、"cancel"或"error"这几个值），如果status为"save"，表示用户选择了保存，此时"text"键为当前编辑器中的文本内容，如果status为"saveAs"则表示用户选择了“另存为”，此时"text"键值为当前编辑器中的文本内容；如果status为"cancel"，表示用户是强制快捷退出的（按Ctrl+Q），此时"text"键为空字符串；如果status为"error"，则"text"键为空字符串，同时"error"键是包含错误信息的字符串。如果status键的值为save或saveAs，还会有一个path键值表示文件路径（包括SSH情况下的远端文件路径）
- 用户在编辑器中可以的快捷键操作包括：如保存（Ctrl+S）、另存为（Ctrl+K）、复制（Ctrl+C）、粘贴（Ctrl+V）、撤销（Ctrl+Z）、重做（Ctrl+Y）、查找（Ctrl+F）、替换（Ctrl+H）、跳转到行号（Ctrl+G）、正常退出（Ctrl+X，如果文本有变动需要提示保存）、不保存强制退出（Ctrl+Q，即使文本做了改动也不保存）、切换自动折行（Ctrl+W）等。
- 查找替换式搜索条件支持正则表达式；
- 编辑器为UTF-8编码，要支持中文字符和ASCII字符混合显示正常。
- 编辑器内要支持中文输入；
- 编辑器内，要支持光标移动、插入、删除、复制、粘贴等操作，光标移动要考虑中文字符和ASCII字符的移动时不能出停在半个字符中间的问题。
- 编辑器内，文本默认自动折行，可以通过热键（Ctrl+W）切换是否自动折行；
- 除必须的状态栏、提示栏外，其他区域都要显示编辑的文本内容；也就是编辑区域要撑满剩余空间；
- 状态栏中，要显示当前光标所在的行号和在行中的位置
- 我准备把该函数做成库传到github.com/topxeq/conedit下，让别人可以在Go语言程序中，作为库来调用这个函数的功能启动文本编辑器；同时，我们完善测试用的主程序，让它可以做为一个可用的功能相对完整的轻量级命令行编辑器，请为此做调整；
- 主程序的名字为conedit（Windows下为conedit.exe）
- ConsoleEditText增加参数“-mode=immediate”，模式包括default、file、immediate这几种，编辑器默认是default模式（应用场景主要是获取用户的一段输入，并返回给调用这个函数的程序进行后续处理，因此不涉及文件操作，所有文件、SSH有关的参数也都无效，状态栏中热键Ctrl-X的提示变为“确认”，Ctrl-Q的提示变为“取消”，其他有关保存和另存的热键提示不应出现，函数主要通过返回status为“ok”和“cancel”表示用户是否确认）。file模式则是可以打开文件或SSH远端文件编辑，但结果还是以函数调用方去处理，因此一旦用户按了Ctrl-S进行保存或者按了Ctrl-K进行另存操作，都是立即执行结束返回结果，status为对应的save或saveAs，path中为文件路径。如果指定了immediate模式，则文件保存、另存、ssh的保存、另存都在函数执行时用户按键后直接进行，只有用户按Ctrl-X退出或用Ctrl-Q强制退出编辑器时，函数才会返回，那时候status只有cancel（用户按Ctrl-Q退出）、error、exit（表示用户正常退出，即按Ctrl-X键退出）

