# 代码行数统计工具

一个用Go语言编写的代码行数统计工具，可以扫描指定目录下的所有代码文件，统计各种文件类型的行数信息。

## 功能特性

1. **指定代码文件夹路径** - 支持扫描任意指定的目录
2. **自动统计文件类型和数量** - 自动识别并统计所有文件类型
3. **代码行数统计** - 统计总行数、代码行数、空行数、注释行数
4. **文件类型占比** - 计算每种文件类型在项目中的占比
5. **黑名单模式** - 自动排除非代码文件（图片、音视频、压缩包等）
6. **CSV报告输出** - 生成详细的CSV格式统计报告

## 项目结构

- `main.go` - 主程序入口和命令行处理
- `types.go` - 数据结构定义
- `config.go` - 配置信息和黑名单设置
- `analyzer.go` - 文件分析和目录扫描逻辑
- `reporter.go` - 报告生成和输出格式化

## 使用方法

### 编译程序
```bash
go build .
```

### 运行程序

#### 基本用法
```bash
# 扫描当前目录
./CodeStatistics

# 扫描指定目录
./CodeStatistics -p /path/to/your/project

# 指定输出文件
./CodeStatistics -p /path/to/your/project -o report.csv
```

#### 命令行参数
- `-p, --path` : 要扫描的代码目录路径 (默认: 当前目录)
- `-o, --output` : 输出CSV文件路径 (默认: code_statistics.csv)
- `-h, --help` : 显示帮助信息

#### 示例
```bash
# Windows
./CodeStatistics.exe -p "C:\Users\YourName\Projects\MyProject" -o "C:\Reports\stats.csv"

# Linux/Mac
./CodeStatistics -p ~/projects/myproject -o ~/reports/stats.csv
```

## 输出格式

### 控制台输出
程序会在控制台显示格式化的统计表格，包含：
- 文件类型
- 文件数量和占比
- 总行数、代码行数、空行数、注释行数
- 代码占比

### CSV报告
生成的CSV文件包含以下列：
- 文件类型
- 文件数量
- 文件占比(%)
- 总行数
- 代码行数
- 空行数
- 注释行数
- 代码占比(%)

## 支持的文件类型

工具支持识别多种编程语言的注释格式：
- Go, JavaScript, TypeScript, Java, C/C++, C#, PHP, Swift, Kotlin
- Python, Shell, Ruby, Perl, YAML
- HTML, XML, Vue
- CSS, SCSS, SASS
- SQL

## 黑名单文件类型

自动排除以下类型的文件：
- 图片文件: .jpg, .png, .gif, .svg 等
- 音视频文件: .mp3, .mp4, .avi 等
- 压缩文件: .zip, .rar, .7z 等
- 文档文件: .pdf, .doc, .xls 等
- 二进制文件和临时文件

## 跳过的目录

自动跳过以下目录：
- .git, .svn, .hg (版本控制)
- node_modules, vendor (依赖包)
- target, build, dist (构建输出)
- .idea, .vscode (IDE配置)

## 技术实现

- 使用结构化设计模式，每个文件不超过300行
- 使用go-flags库处理命令行参数
- 支持跨平台运行 (Windows, Linux, macOS)
- 高效的文件遍历和内容分析
