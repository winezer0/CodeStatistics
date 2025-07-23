# CodeStatistics

CodeStatistics 是一个强大的代码统计工具，用于分析项目中的代码文件并生成详细的统计报告。

## 功能特点

- **代码行数统计**：统计项目中的总行数、代码行数、注释行数和空行数
- **文件类型分析**：区分不同类型的文件（基于文件扩展名）
- **灵活的过滤机制**：支持白名单和黑名单机制，可以指定要包含或排除的文件类型和目录
- **报告生成**：可以生成CSV格式的统计报告

## 安装

```bash
# 克隆仓库
git clone https://github.com/WINDOWS/CodeStatistics.git

# 进入项目目录
cd CodeStatistics

# 构建项目
go build
```

## 使用方法

```bash
# 基本用法
./CodeStatistics --path /path/to/your/code

# 生成CSV报告
./CodeStatistics --path /path/to/your/code --output report.csv

# 启用注释行检测
./CodeStatistics --path /path/to/your/code --comments
```

## 命令行选项

| 选项 | 简写 | 描述 |
|------|------|------|
| `--path` | `-p` | 要扫描的代码目录路径 |
| `--output` | `-o` | 输出CSV文件路径 |
| `--comments` | `-c` | 启用注释行检测 |
| `--white-add` | `-w` | 添加扩展名到内置白名单，逗号分隔 (例如: .ext1,.ext2) |
| `--white-cover` | `-W` | 使用扩展名列表覆盖内置白名单，逗号分隔 |
| `--black-add` | `-b` | 添加扩展名到内置黑名单，逗号分隔 |
| `--black-cover` | `-B` | 使用扩展名列表覆盖内置黑名单，逗号分隔 |
| `--only-white` | `-O` | 仅显示白名单文件，跳过黑名单和未知文件分析 |
| `--bdir-add` | `-d` | 添加目录到内置黑名单目录，逗号分隔 (例如: dir1,dir2) |
| `--bdir-cover` | `-D` | 使用目录列表覆盖内置黑名单目录，逗号分隔 |
| `--show-builtin` | `-s` | 显示内置默认白名单/黑名单/黑名单目录 |
| `--lf` |  | 日志文件路径 |
| `--ll` |  | 日志级别 (debug/info/warn/error) |
| `--cf` |  | 控制台日志格式 (T L C M F 组合或 off\|null 禁用) |

## 白名单和黑名单

### 默认白名单（代码文件）

工具内置了常见的代码文件扩展名白名单，包括：
- 编程语言文件：.go, .java, .c, .cpp, .py, .js, .ts 等
- 模板语言文件：.vue, .svelte, .jsx, .tsx 等
- 前端语言文件：.html, .htm 等
- 脚本语言文件：.sh, .bash, .bat, .cmd 等

### 默认黑名单（非代码文件）

工具内置了常见的非代码文件扩展名黑名单，包括：
- 数据库文件：.sql, .mysql 等
- 标记和配置语言文件：.xml, .css, .json, .yaml 等
- 构建和配置文件：.mk, .cmake 等
- 文档文件：.md, .rst, .tex 等
- 图片文件：.jpg, .png, .gif 等
- 音视频文件：.mp3, .mp4, .avi 等
- 压缩文件：.zip, .rar, .7z 等

### 默认跳过目录

工具默认会跳过以下目录：
`.git`, `node_modules`, `vendor`, `.svn`, `.hg`, `target`, `build`, `dist`, `.idea`, `.vscode`

## 示例

### 基本分析

```bash
./CodeStatistics --path ./myproject
```

### 生成CSV报告

```bash
./CodeStatistics --path ./myproject --output stats.csv
```

### 自定义白名单

```bash
./CodeStatistics --path ./myproject --white-add .custom,.special
```

### 自定义黑名单目录

```bash
./CodeStatistics --path ./myproject --bdir-add cache,temp
```

## 输出示例

```
代码统计摘要:
总文件数: 120
总行数: 15000
代码行数: 10000 (66.7%)
空行数: 3000 (20.0%)
注释行数: 2000 (13.3%)

文件类型分布:
.go: 50 文件, 5000 行 (33.3%)
.js: 30 文件, 4000 行 (26.7%)
.html: 20 文件, 3000 行 (20.0%)
.css: 15 文件, 2000 行 (13.3%)
其他: 5 文件, 1000 行 (6.7%)
```
