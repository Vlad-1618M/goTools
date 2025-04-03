### `fs_audit.go` Overview:
### Purpose:
The `fs_audit.go` file serves as a **file and directory scanner**. Its main role is to:
- Recursively scan a directory and list files grouped by extensions.
- Provide color-coded terminal outputs based on file extensions, using the `fatih/color` package.
- Help users easily navigate and analyze the contents of directories.
---
### Key Features:
1. **Color-Coded Output**:
      - Files are displayed in different colors according to their extensions:
      - Uses pre-defined color mappings for popular file extensions:

2. **Directory Traversal**:
      - Scans and lists files and directories using `os` and `filepath`

3. **File Extension Grouping**:
      - Tracks file counts by extension:
      - Displays grouped files neatly:

4. **Highly Interactive**:
---
### Included Packages:
### Standard Library:
- **`os`**: For filesystem operations such as reading directories and opening files.
- **`path/filepath`**: To resolve file paths and work with various path formats.
- **`fmt`**, **`bufio`**, and **`strings`**: For formatting output, reading terminal input/output, and string manipulation.
- **`strconv`**: For converting user input to integers.
- **`sort`**: For ordering output data like file or directory lists.

### Third-Party Library:
- **[`fatih/color`](https://github.com/fatih/color)**:
      - For terminal colorization.
      - Allows dynamic styling (bold, brightly colored text) for better readability.
---
### Global Variables
- **Color Definitions**: Multiple colors are predefined for styling output:
```go
var (
      green       = color.New(color.FgGreen).SprintFunc()
      yellow      = color.New(color.FgYellow).SprintFunc()
      cyan        = color.New(color.FgCyan).SprintFunc()
      magenta     = color.New(color.FgMagenta).SprintFunc()
      red         = color.New(color.FgRed).SprintFunc()
      white       = color.New(color.FgWhite, color.Bold).SprintFunc()
      orange      = color.New(color.FgHiRed).SprintFunc()
      darkOrange  = color.New(color.FgMagenta).SprintFunc()
      gray        = color.New(color.FgHiBlack).SprintFunc()
      blue        = color.New(color.FgBlue).SprintFunc()
      lightGreen  = color.New(color.FgHiGreen).SprintFunc()
      lightYellow = color.New(color.FgHiYellow).SprintFunc()
      purple      = color.New(color.FgHiMagenta).SprintFunc()
      brightCyan  = color.New(color.FgHiCyan).SprintFunc()
      brightRed   = color.New(color.FgHiRed, color.Bold).SprintFunc()
      brightBlue  = color.New(color.FgHiBlue).SprintFunc()
)
```
- **File-Extension-to-Color Mapping**: Maps file extensions to their designated colors:
```go
var extensionColors = map[string]func(a ...interface{}) string{
      "sh":   magenta,     // Shell scripts
      "py":   yellow,      // Python scripts
      "yaml": cyan,        // YAML files
      "yml":  cyan,        // YAML files
      "xml":  darkOrange,  // XML files
      "json": green,       // JSON files
      "txt":  white,       // Plain text files
      "go":   lightYellow, // Go source files
      "log":  gray,        // Log files
      "md":   brightCyan,  // Markdown files
      "html": orange,      // HTML files
      "css":  blue,        // CSS files
      "js":   yellow,      // JavaScript files
      "ts":   lightGreen,  // TypeScript files
      "java": brightRed,   // Java files
      "c":    brightBlue,  // C files
      "cpp":  brightBlue,  // C++ files
      "php":  purple,      // PHP files
      "ini":  lightGreen,  // Configuration files (INI)
      "cfg":  lightGreen,  // Configuration files (CFG)
      "env":  gray,        // Environment files
      "toml": cyan,        // TOML files
      "rs":   purple,      // Rust files
      "kt":   brightRed,   // Kotlin files
  }
```
---
### Main Functions:
### **`main()`**
- Script entry point.
- Prompts the user to: <br>
      1. Select a directory to scan: <br>
      2. Resolve and display its absolute path.
- Delegates to `scanFiles()` for processing files within the selected directory.

### **`listDirPaths(base string) []string`**
- Determines immediate subdirectories of a given base directory:
- Filters out non-directory entries:

### **`scanFiles(root string)`**
- Recursively scans the specified directory for files:
- Tracks counts of files by their extensions:
- Categorizes and groups files based on their extensions:
---

### Flow:
1. **Run the script**:
```go
go run fs_audit.go
```
1. **Prompt**:
      - The program lists available directories.
      - Prints the absolute paths for user selection.

2. **Color-Coded Output**:
      - Files and directories are displayed using colors defined in `extensionColors`.

### Notes:
- New file types/extensions can be added by updating the `extensionColors` map.
- Unrecognized extensions default to `white` (defined globally).
---

#### Personal Note :0)
>- ... wrote this for my `.zshrc` to replace existing `shell script` which was too-oo slow for me :) 
```bash
audit_fs() {
    echo -e "$decorator_init"

    local f_count=0 d_count=0
    declare -A ext_count
    declare -A ext_dirs

    while IFS= read -r -d '' item; do
        if [[ -f "$item" ]]; then
            ((f_count++))
            ext="${item##*.}"
            [[ "$ext" == "$item" ]] && ext="(no ext)"
            ((ext_count["$ext"]++))
            parent_dir=$(dirname "$item")
            ext_dirs["$ext"]+="${parent_dir}\n"
        elif [[ -d "$item" ]]; then
            ((d_count++))
        fi
    done < <(find . -print0)

    echo -e "${cyan}Files: ${yellow}${f_count}, ${green}Directories: ${d_count}${_off}\n"

    echo -e "${cyan}File Type Breakdown:${_off}"
    for ext in "${!ext_count[@]}"; do
        echo -e "${yellow}* .$ext${_off} — ${green}${ext_count[$ext]} file(s)${_off}"
        echo -e "  Found in directories:"
        echo -e "${ext_dirs[$ext]}" | sort -u | sed 's/^/    - /'
        echo ""
    done

    echo -e "$decorator_done"
}
```
