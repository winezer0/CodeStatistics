package statistics

import (
	"time"
)

// FileStats file statistics information
type FileStats struct {
	Extension    string  // File extension
	FileCount    int     // File count
	TotalLines   int     // Total lines
	CodeLines    int     // Code lines (non-empty)
	BlankLines   int     // Blank lines
	CommentLines int     // Comment lines
	FileRatio    float64 // File type ratio
}

// CodeStatistics code statistics analyzer
type CodeStatistics struct {
	RootPath        string
	WhitelistStats  map[string]*FileStats // Whitelist file statistics
	BlacklistStats  map[string]*FileStats // Blacklist file statistics
	TotalFiles      int
	Whitelist       map[string]bool // Code file whitelist
	Blacklist       map[string]bool // Non-code file blacklist
	SkipDirectories []string        // Directories to skip during scanning
	EnableComments  bool            // Whether to enable comment line detection
	OnlyWhite       bool            // Whether to show only whitelist files
	StartTime       time.Time
}

// SummaryData statistics summary data
type SummaryData struct {
	TotalFiles        int
	TotalLines        int
	TotalCodeLines    int
	TotalBlankLines   int
	TotalCommentLines int
	Extensions        []string
}

// WhitelistConfig 后缀白名单配置
type WhitelistConfig struct {
	Add      []string // Extensions to add to whitelist
	Override []string // Extensions to override whitelist
}

// BlacklistConfig 后缀黑名单配置
type BlacklistConfig struct {
	Add      []string // Extensions to add to blacklist
	Override []string // Extensions to override blacklist
}

// BlackDirConfig 目录黑名单配置
type BlackDirConfig struct {
	Add      []string // Directories to add to blacklist
	Override []string // Directories to override blacklist
}

// DefaultWhitelist 默认白名单 - 明确的代码文件扩展名
var DefaultWhitelist = []string{
	// 编程语言
	".go", ".java", ".c", ".cpp", ".cc", ".cxx", ".h", ".hpp", ".hxx",
	".cs", ".vb", ".fs", ".fsx", ".fsi", // .NET languages
	".py", ".pyw", ".pyi", // Python
	".js", ".jsx", ".ts", ".tsx", ".mjs", ".cjs", // JavaScript/TypeScript
	".php", ".php3", ".php4", ".php5", ".phtml", // PHP
	".rb", ".rbw", ".rake", ".gemspec", // Ruby
	".pl", ".pm", ".t", ".pod", // Perl
	".swift", ".kt", ".kts", // Swift, Kotlin
	".rs", ".rlib", // Rust
	".scala", ".sc", // Scala
	".clj", ".cljs", ".cljc", ".edn", // Clojure
	".hs", ".lhs", // Haskell
	".ml", ".mli", ".mll", ".mly", // OCaml
	".elm",     // Elm
	".dart",    // Dart
	".lua",     // Lua
	".r", ".R", // R
	".m", ".mm", // Objective-C
	".f", ".f90", ".f95", ".f03", ".f08", // Fortran
	".pas", ".pp", ".inc", // Pascal
	".asm", ".s", ".S", // Assembly
	".v", ".vh", ".sv", ".svh", // Verilog/SystemVerilog
	".vhd", ".vhdl", // VHDL

	// 模板语言
	".vue", ".svelte", // Frontend frameworks
	".jsx", ".tsx", // React (already included above)
	".hbs", ".handlebars", // Handlebars
	".mustache",                // Mustache
	".twig",                    // Twig
	".jinja", ".jinja2", ".j2", // Jinja
	".erb", ".haml", ".slim", // Ruby templates

	// 前端语言
	".html", ".htm", ".xhtml", ".shtml", // HTML

	// 脚本语言
	".sh", ".bash", ".zsh", ".fish", ".csh", ".tcsh", // Shell scripts
	".bat", ".cmd", ".ps1", ".psm1", ".psd1", // Windows scripts
	".awk", ".sed", // Text processing

}

// DefaultBlacklist 默认黑名单 - 排除非代码文件
var DefaultBlacklist = []string{
	// 数据库
	".sql", ".mysql", ".pgsql", ".sqlite", ".plsql", // SQL

	// 标记和配置语言
	".xml", ".xsl", ".xslt", ".xsd", ".dtd", // XML
	".css", ".scss", ".sass", ".less", ".styl", // CSS
	".json", ".jsonc", ".json5", // JSON
	".yaml", ".yml", // YAML
	".toml",                            // TOML
	".ini", ".cfg", ".conf", ".config", // Config files
	".properties", // Properties files

	// 构建和配置文件
	".mk", ".make", ".makefile", // Makefiles
	".cmake", ".cmake.in", // CMake
	".gradle", ".gradle.kts", // Gradle
	".sbt",           // SBT
	".bazel", ".bzl", // Bazel
	".ninja", // Ninja

	// 文档
	".md", ".markdown", ".mdown", ".mkd", ".mdx", // Markdown
	".rst", ".rest", // reStructuredText
	".tex", ".latex", ".ltx", // LaTeX
	".adoc", ".asciidoc", // AsciiDoc

	// 其他
	".dockerfile", ".containerfile", // Docker
	".gitignore", ".gitattributes", ".gitmodules", // Git
	".editorconfig",            // EditorConfig
	".eslintrc", ".prettierrc", // Linting/Formatting

	// 包管理和依赖文件
	".mod", ".sum", // Go modules
	".lock", ".toml", // Various lock files and configs

	// 图片文件
	".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg", ".ico", ".webp", ".tiff", ".tif",
	// 音视频文件
	".mp3", ".mp4", ".avi", ".mov", ".wmv", ".flv", ".wav", ".ogg", ".mkv", ".webm",
	// 压缩文件
	".zip", ".rar", ".7z", ".tar", ".gz", ".bz2", ".xz", ".lz", ".lzma",
	// 文档文件
	".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".odt", ".ods", ".odp",
	// 字体文件
	".ttf", ".otf", ".woff", ".woff2", ".eot",
	// 其他二进制文件
	".exe", ".dll", ".so", ".dylib", ".bin", ".dat", ".elf", ".o", ".obj", ".lib", ".a",
	// 临时文件
	".tmp", ".temp", ".cache", ".log", ".bak", ".swp", ".swo",
	// 版本控制和配置目录相关
	".git", ".svn", ".hg",
	// 其他非代码文件
	".jar", ".war", ".ear", ".class", // Java compiled
	".pyc", ".pyo", ".pyd", // Python compiled
	".gem",                                         // Ruby gems
	".deb", ".rpm", ".msi", ".dmg", ".pkg", ".pex", // Package files
	"NONE", "", //无后缀文件
}

// DefaultBlackDirs 需要跳过的目录
var DefaultBlackDirs = []string{
	".git", "node_modules", "vendor", ".svn", ".hg", "target", "build", "dist", ".idea", ".vscode",
}
