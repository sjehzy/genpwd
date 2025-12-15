# genpwd
用 Go 语言编写的随机密码生成工具，兼具命令行交互（CLI）和可编译为 Windows 可执行文件（.exe）的特性，支持自定义密码长度、字符类型、排除易混淆字符，且生成效率高、体积小（编译后仅几 MB）

参数简化规则：


1、 环境准备
安装 Go 环境：Go 官方下载（Windows 版，安装时勾选 Add Go to PATH）；
验证安装：打开 CMD/PowerShell，执行 go version，输出版本号即成功。

2、 编译命令
将代码保存为 gpwd.go，在脚本所在目录执行：
```
$ ls
gpwd.go

```
3. 编译为 Windows 64 位可执行文件
```
go build -o gpwd.exe gpwd.go
```
编译完成后，目录下会生成 gpwd.exe，可直接双击运行，或通过 CMD/PowerShell 调用。

# 使用方法（Windows 下）
1. 基础用法（默认配置）
双击 gpwd.exe，或在 CMD 执行：
powershell
.\gpwd.exe
默认生成 12 位 密码（包含大小写字母、数字、特殊符号，排除易混淆字符）。
2. 自定义参数示例
```
# 测试1：生成16位无特殊符号、排除易混淆字符的密码
.\gpwd.exe -l 16 -s false
# 输出示例：=== 生成的随机密码 ===
# 1. 7R9k8S7p2L8m7R8t5（无特殊符号，无0/O/l/I等）

# 测试2：生成8位纯数字（包含易混淆字符）
.\gpwd.exe -l 8 -u false -w false -s false -n false
# 输出示例：=== 生成的随机密码 ===
# 1. 08917689（包含0/8/9等易混淆数字）

# 测试3：批量生成3个密码，复制第一个到剪贴板
.\gpwd.exe -b 3 -c true
# 输出示例：
# === 生成的随机密码 ===
# 1. 9s8K7$6p5L4m3R2t
# 2. 8k7J6%5o4K3n2M1s
# 3. 7j6H5&4n3J2b1N0r
# ✅ 第一个密码已复制到剪贴板！
```

# 对应参数：
| 原参数     | 缩写参数 | 说明           |
|:---------: |:-------: |:-------------: |
| -len       | -l       | 密码长度       |
| -upper     | -u       | 包含大写字母   |
| -lower     | -w       | 包含小写字母   |
| -digit     | -d       | 包含数字       |
| -symbol    | -s       | 包含特殊符号   |
| -no-similar| -n       | 排除易混淆字符 |
| -batch     | -b       | 批量生成数量   |
| -copy      | -c       | 复制到剪贴板   |

> 注：lower 首字母 l 与 len 冲突，改用 w（小写=lower=小写w），更易记。

