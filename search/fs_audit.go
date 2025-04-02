package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

// ... color definitions:
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

// ... color grouping fs extensions mapping:
var extensionColors = map[string]func(a ...interface{}) string{
	// Common file types
	"sh":   magenta,     // <-  Shell scripts
	"py":   yellow,      // <-  Python scripts
	"yaml": cyan,        // <-  YAML files
	"yml":  cyan,        // <-  YAML files
	"xml":  darkOrange,  // <-  XML files
	"json": green,       // <-  JSON files
	"txt":  white,       // <-  Plain text files
	"go":   lightYellow, // <-  Go source files
	"log":  gray,        // <-  Log files
	"md":   brightCyan,  // <-  Markdown files
	"html": orange,      // <-  HTML files
	"css":  blue,        // <-  CSS files
	"js":   yellow,      // <-  JavaScript files
	"ts":   lightGreen,  // <-  TypeScript files
	"java": brightRed,   // <-  Java files
	"c":    brightBlue,  // <-  C files
	"cpp":  brightBlue,  // <-  C++ files
	"php":  purple,      // <-  PHP files
	"ini":  lightGreen,  // <-  Configuration files (INI)
	"cfg":  lightGreen,  // <-  Configuration files (CFG)
	"env":  gray,        // <-  Environment files
	"toml": cyan,        // <-  TOML files
	"rs":   purple,      // <-  Rust files
	"kt":   brightRed,   // <-  Kotlin files
}

func main() {
	// ... subdirs fetch list:
	dirs := listDirPaths(".")
	if len(dirs) == 0 {
		fmt.Println(red("No directories found:"))
		return
	}

	// ...  show available dirs selection:
	fmt.Println(yellow("... available directories to scan:"))
	maxWidth := len(fmt.Sprintf("%d", len(dirs))) // ... adjust width to fit indices:
	for i, dir := range dirs {
		absolutePath, err := filepath.Abs(dir)
		if err != nil {
			fmt.Printf("%sError resolving path%s: %v\n", red(""), white(""), err)
			continue
		}
		fmt.Printf("%s[%*d]%s -> %s\n", orange(""), maxWidth, i+1, white(""), darkOrange(absolutePath))
	}

	// ... user promopt | select directory by index value:
	fmt.Println(yellow("... select directory index to scan:"))
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	index, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || index < 1 || index > len(dirs) {
		fmt.Println(red("Invalid index.") + "\n\t" + strings.Repeat(gray("-"), 20))
		return
	}

	// ... resolve selected dir to scan in:
	scanPath := dirs[index-1]
	absoluteScanPath, _ := filepath.Abs(scanPath)

	fmt.Printf("\n%sScanning path:%s %s\n\n", green(""), white(""), darkOrange(absoluteScanPath))
	scanFiles(scanPath)
}

// ... listDirPaths function | lists immediate child dirs only for a given base path:
func listDirPaths(base string) []string {
	var dirs []string
	entries, err := os.ReadDir(base)
	if err != nil {
		fmt.Println(red("Error reading directory:"), err)
		return nil
	}
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, filepath.Join(base, entry.Name()))
		}
	}
	return dirs
}

// ... scanFiles function | scan starts at root and recursive for files & dirs:
func scanFiles(root string) {
	fCount, dCount := 0, 0
	extCount := make(map[string]int)            // ... file count by extension tracking:
	extDirs := make(map[string][]string)        // ... dires for each file extension tracking:
	fileNamesByExt := make(map[string][]string) // ... extensions to corresponding files mapping: 

	// ... dir-tree walk:
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("%sError accessing path:%s %v\n", red(""), white(""), err)
			return nil // ... skip problematic paths:
		}
		if d.IsDir() {
			dCount++ // ... dir count increment:
			return nil
		}

		// ... process files:
		fCount++
		ext := strings.TrimPrefix(filepath.Ext(path), ".")
		if ext == "" {
			ext = "(no ext)"
		}

		// ... count extensions:
		extCount[ext]++

		// ... collect dirs for extensions:
		dir := filepath.Dir(path)
		if !contains(extDirs[ext], dir) {
			extDirs[ext] = append(extDirs[ext], dir)
		}

		// ... collect files by extension:
		fileNamesByExt[ext] = append(fileNamesByExt[ext], path)
		return nil
	})
	if err != nil {
		fmt.Println(red("Error during file scanning:"), err)
		return
	}
	fmt.Printf("\n%sFiles:%s %d, %sDirectories:%s %d\n\n", cyan(""), white(""), fCount, cyan(""), white(""), dCount)
	printBreakdown(extCount, extDirs, fileNamesByExt)
}

// ... checks if a slice contains a specific string:
func contains(slice []string, item string) bool {
	for _, elem := range slice {
		if elem == item {
			return true
		}
	}
	return false
}

func printBreakdown(extCount map[string]int, extDirs map[string][]string, fileNamesByExt map[string][]string) {
	fmt.Println(cyan("File Type Breakdown:"))

	// ... sort extensions alphabetically:
	keys := make([]string, 0, len(extCount))
	for k := range extCount {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// ... print each file type details:
	for _, ext := range keys {
		fmt.Printf("%s* .%s%s — %s%d file(s)%s\n", yellow(""), ext, reset(), green(""), extCount[ext], reset())
		fmt.Println("  Found in directories:")

		// ... print sorted dirs:
		sort.Strings(extDirs[ext])
		for _, dir := range extDirs[ext] {
			fmt.Printf("    - %s\n", gray(dir))
		}

		fmt.Println("  Files:")
		sort.Strings(fileNamesByExt[ext])
		for _, file := range fileNamesByExt[ext] {
			dir := filepath.Dir(file)              // ... get dir part:
			fileName := filepath.Base(file)        // ... get file name:
			colorFunc := getColorForExtension(ext) // ... get extension color:
			fmt.Printf("    - %s/%s\n", gray(dir), colorFunc(fileName))
		}
		fmt.Println()
	}
}

// ... getColorForExtension function returns color function for a given extension:
func getColorForExtension(ext string) func(a ...interface{}) string {
	if colorFunc, exists := extensionColors[ext]; exists {
		return colorFunc
	}
	return white // ... unconfigured extensions default color is white:
}

// ... reset function | resets terminal style:
func reset() string {
	return color.New(color.Reset).Sprint("")

}
