# Conedit - 命令行文本编辑器库

[![Go Reference](https://pkg.go.dev/badge/github.com/topxeq/conedit.svg)](https://pkg.go.dev/github.com/topxeq/conedit)
[![Go Report Card](https://goreportcard.com/badge/github.com/topxeq/conedit)](https://goreportcard.com/report/github.com/topxeq/conedit)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

**[🌏 English Documentation](README.md)**

一个轻量级、可嵌入的 Go 语言终端文本编辑器库，支持 UTF-8/中文和 SSH — 无需 CGO。

## 功能特性

- 功能完整的控制台文本编辑器
- UTF-8 编码，支持中文字符
- SSH 远程文件编辑
- 自动换行（可切换）
- 常用操作快捷键支持
- 查找和替换（支持正则表达式）
- 撤销/重做功能
- 复制/粘贴（支持选区）

## 安装

```bash
go get github.com/topxeq/conedit
```

## 库使用方式

```go
package main

import (
    "fmt"
    "github.com/topxeq/conedit/editor"
)

func main() {
    // 打开编辑器，带默认文本
    result := editor.ConsoleEditText("默认文本内容")
    
    // 或打开文件
    result = editor.ConsoleEditText("", "-filePath=/path/to/file.txt")
    
    // 或通过 SSH 编辑远程文件
    result = editor.ConsoleEditText("", 
        "-fromSSH",
        "-sshHost=192.168.1.100",
        "-sshPort=22",
        "-sshUser=root",
        "-sshPass=password",
        "-filePath=/remote/path/file.txt",
    )
    
    // 检查结果
    if result["status"] == "save" || result["status"] == "saveAs" {
        fmt.Printf("已保存到：%s\n", result["path"])
        fmt.Printf("内容：%s\n", result["text"])
    } else if result["status"] == "cancel" {
        fmt.Println("用户已取消")
    } else if result["status"] == "error" {
        fmt.Printf("错误：%v\n", result["error"])
    }
}
```

## 命令行使用

构建命令行编辑器：

```bash
go build -o conedit ./cmd/editor
```

使用方法：

```bash
# 默认模式 - 文本输入，无文件操作
./conedit

# 打开文件编辑（immediate 模式 - 退出时自动保存）
./conedit file.txt

# 明确选择模式
./conedit -mode=default        # 文本输入模式
./conedit -mode=file file.txt  # 文件模式（保存后返回）
./conedit -mode=immediate file.txt  # immediate 模式（退出时自动保存）

# 通过 SSH 编辑远程文件
./conedit -mode=immediate -fromSSH -sshHost=192.168.1.100 -sshUser=root -filePath=/remote/file.txt

# 显示帮助
./conedit --help
```

### 模式行为

| 模式 | 使用时机 | 行为 | 返回值 |
|------|------|------|------|
| `default` | 无文件参数 | 仅文本输入 | `ok`, `cancel` |
| `file` | 有文件参数 | 编辑，保存时返回 | `save`, `saveAs`, `cancel` |
| `immediate` | 有文件参数 | 编辑，退出时自动保存 | `exit`, `cancel`, `error` |

## 快捷键

| 快捷键 | 功能 |
|--------|------|
| Ctrl+S | 保存 |
| Ctrl+K | 另存为 |
| Ctrl+X | 退出 |
| Ctrl+Q | 强制退出 |
| Ctrl+W | 切换自动换行 |
| Ctrl+C | 复制 |
| Ctrl+V | 粘贴 |
| Ctrl+Z | 撤销 |
| Ctrl+Y | 重做 |
| Ctrl+F | 查找（支持正则） |
| Ctrl+H | 替换（支持正则） |
| Ctrl+G | 跳转到行 |
| Shift+ 方向键 | 选择文本 |

## 选项

| 选项 | 说明 |
|------|------|
| `-filePath=PATH` | 要编辑的文件路径 |
| `-mode=MODE` | 编辑器模式：`default`, `file`, `immediate`（默认：`file`） |
| `-fromSSH` | 在 SSH 服务器上编辑文件 |
| `-sshHost=HOST` | SSH 主机 |
| `-sshPort=PORT` | SSH 端口（默认：22） |
| `-sshUser=USER` | SSH 用户名 |
| `-sshPass=PASS` | SSH 密码 |
| `-sshKeyPath=PATH` | SSH 私钥路径 |
| `-mem` | 强制内存处理（不使用临时文件） |
| `-tmpPath=PATH` | 大文件的自定义临时目录 |

## 返回值

`ConsoleEditText` 函数返回一个 `map[string]interface{}`，包含以下键：

| 键 | 类型 | 说明 |
|-----|------|------|
| `text` | string | 当前编辑器内容（取消或错误时为空） |
| `status` | string | `save`, `saveAs`, `cancel`, `error` 之一 |
| `path` | string | 文件路径（仅当 status 为 `save` 或 `saveAs` 时） |
| `error` | string | 错误信息（仅当 status 为 `error` 时） |

## 系统要求

- Go 1.21 或更高版本
- 支持 UTF-8 的终端

## 依赖项

- [github.com/gdamore/tcell/v2](https://github.com/gdamore/tcell) - 终端屏幕处理
- [golang.org/x/crypto/ssh](https://golang.org/x/crypto) - SSH 支持

## 许可证

[MIT License](LICENSE)
