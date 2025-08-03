// Project Structure for Zippy - Simple Version Control Tool
package main

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var ZippyVersion = "0.0.1"
var ZippyAuthor = "PixCap Soft"

const (
	ZIPPY_NAME    = "Zippy"
	ZIPPY_DESC    = "Simple Version Control Tool with Zip Storage"
	ZIPPY_REPO    = "https://github.com/pixcapsoft/zippy"
	ZIPPY_STAGE   = "stage.json"
)

// Version represents a version entry
type Version struct {
	Tag         string    `json:"tag"`
	Message     string    `json:"message"`
	Timestamp   time.Time `json:"timestamp"`
	Author      string    `json:"author"`
	FilesCount  int       `json:"files_count"`
	ZipPath     string    `json:"zip_path"`
	Size        int64     `json:"size"`
}

// Repository configuration
type RepoConfig struct {
	Name        string    `json:"name"`
	Author      string    `json:"author"`
	Created     time.Time `json:"created"`
	Description string    `json:"description"`
}

// ZippyIgnore handles .zippyignore file parsing
type ZippyIgnore struct {
	patterns []string
}

// Main CLI structure
type Zippy struct {
	repoPath     string
	zippyPath    string
	configPath   string
	versionsPath string
	storagePath  string
	config       RepoConfig
	stagePath    string
}

func main() {
	if len(os.Args) < 2 {
		showBanner()
		showHelp()
		return
	}

	zippy := &Zippy{}
	command := os.Args[1]

	switch command {
	case "init":
		zippy.initRepo()
	case "add":
		if err := zippy.initPaths(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if len(os.Args) < 3 {
			fmt.Println("Usage: zippy add <files>")
			return
		}
		zippy.addFiles(os.Args[2:])
	case "commit":
		if err := zippy.initPaths(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		zippy.commit()
	case "push":
		if err := zippy.initPaths(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		zippy.push()
	case "list", "ls":
		if err := zippy.initPaths(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		zippy.listVersions()
	case "restore":
		if len(os.Args) < 3 {
			fmt.Println("Usage: zippy restore <version> [path]")
			return
		}
		if err := zippy.initPaths(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		var restorePath string
		if len(os.Args) >= 4 {
			restorePath = os.Args[3]
		}
		zippy.restore(os.Args[2], restorePath)
	case "status":
		if err := zippy.initPaths(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		zippy.status()
	case "diff":
		if len(os.Args) < 4 {
			fmt.Println("Usage: zippy diff <version1> <version2>")
			return
		}
		if err := zippy.initPaths(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		zippy.diff(os.Args[2], os.Args[3])
	case "version", "-v", "--version":
		showVersion()
	case "help", "-h", "--help":
		showHelp()
	case "about", "info":
		showAbout()
	case "patch":
		if len(os.Args) < 4 {
			fmt.Println("Usage: zippy patch <version> <file/folder>")
			return
		}
		if err := zippy.initPaths(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		zippy.patchVersion(os.Args[2], os.Args[3])
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Run 'zippy help' for available commands")
	}
}

// initPaths initializes the Zippy paths and checks if repo exists
func (zippy *Zippy) initPaths() error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}
	
	zippy.repoPath = cwd
	zippy.zippyPath = filepath.Join(cwd, ".zippy")
	zippy.configPath = filepath.Join(zippy.zippyPath, "config.json")
	zippy.versionsPath = filepath.Join(zippy.zippyPath, "versions")
	zippy.storagePath = filepath.Join(zippy.zippyPath, "storage")
	zippy.stagePath = filepath.Join(zippy.zippyPath, ZIPPY_STAGE)

	// Check if repository is initialized
	if _, err := os.Stat(zippy.zippyPath); os.IsNotExist(err) {
		return fmt.Errorf("not a Zippy repository. Run 'zippy init' first")
	}

	// Load config if it exists
	if configData, err := os.ReadFile(zippy.configPath); err == nil {
		json.Unmarshal(configData, &zippy.config)
	}

	return nil
}

func showBanner() {
	fmt.Printf(`
 ███████╗██╗██████╗ ██████╗ ██╗   ██╗
 ╚══███╔╝██║██╔══██╗██╔══██╗╚██╗ ██╔╝
   ███╔╝ ██║██████╔╝██████╔╝ ╚████╔╝ 
  ███╔╝  ██║██╔═══╝ ██╔═══╝   ╚██╔╝  
 ███████╗██║██║     ██║        ██║   
 ╚══════╝╚═╝╚═╝     ╚═╝        ╚═╝   

%s v%s - %s
%s

`, ZIPPY_NAME, ZippyVersion, ZIPPY_DESC, ZIPPY_REPO)
}

func showVersion() {
	fmt.Printf("%s version %s\n", ZIPPY_NAME, ZippyVersion)
	fmt.Printf("Build: %s\n", time.Now().Format("2006-01-02"))
	fmt.Printf("Author: %s\n", ZippyAuthor)
}

func showAbout() {
	showBanner()
	fmt.Printf("About %s:\n", ZIPPY_NAME)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Version: %s\n", ZippyVersion)
	fmt.Printf("Description: %s\n", ZIPPY_DESC)
	fmt.Printf("Author: %s\n", ZippyAuthor)
	fmt.Printf("Repository: %s\n", ZIPPY_REPO)
	fmt.Println()
	fmt.Println("Features:")
	fmt.Println("  • Simple version control with zip storage")
	fmt.Println("  • .zippyignore support (like .gitignore)")
	fmt.Println("  • Human-readable version tags")
	fmt.Println("  • Easy backup and sharing")
	fmt.Println("  • Cross-platform compatibility")
	fmt.Println()
	fmt.Println("For more information, visit:", ZIPPY_REPO)
}

func showHelp() {
	help := `
Zippy - Simple Version Control Tool with Zip Storage

USAGE:
  zippy <command> [arguments]

COMMANDS:
  init
      Initialize a new Zippy repository. Prompts for repo name, author, and description.

  add <files|folders|.>
      Stage files or folders for the next commit. Use '.' to stage all files (except those in .zippyignore).
      Example: zippy add main.go src/ .env

  commit -m "message" -v "tag"
      Create a new version from staged files. Requires a message and a version tag.
      Example: zippy commit -m "Initial commit" -v "v1.0"

  push
      (Placeholder) Save current version to zip file (already done by commit).

  list, ls
      List all saved versions with their tags, dates, authors, and messages.

  restore <version> [path]
      Restore all files from a version, or a specific file/folder if [path] is given.
      Example: zippy restore v1.0 src/main.go

  status
      Show repository status:
        - Files staged for commit
        - Ignored files
        - Changes compared to the latest version

  diff <version1> <version2>
      Show which files were added, removed, or changed between two versions.
      Example: zippy diff v1.0 v2.0

  patch <version> <file/folder>
      Add a file or folder to an existing version (modifies the zip and metadata).
      Example: zippy patch v1.0 README.md

  version, -v, --version
      Show Zippy version information.

  help, -h, --help
      Show this help message.

  about, info
      Show detailed information about Zippy.

FILES:
  .zippyignore
      List files and patterns to ignore (like .gitignore).
  .zippy/
      Zippy metadata directory (do not delete or edit manually).

WORKFLOW EXAMPLES:
  zippy init
  zippy add .
  zippy commit -m "Initial commit" -v "v1.0"
  zippy list
  zippy status
  zippy restore v1.0 src/main.go
  zippy patch v1.0 README.md
  zippy diff v1.0 v2.0

TIPS:
- Use 'zippy add' before every commit to stage files.
- Edit .zippyignore to avoid archiving unwanted files.
- Use 'zippy status' to see what will be committed and what changed.
- Each commit creates a zip in .zippy/storage and a metadata file in .zippy/versions.

For more information: %s
`
	fmt.Printf(help, ZIPPY_REPO)
}

func (zippy *Zippy) initRepo() {
	cwd, _ := os.Getwd()
	zippy.repoPath = cwd
	zippy.zippyPath = filepath.Join(cwd, ".zippy")
	zippy.configPath = filepath.Join(zippy.zippyPath, "config.json")
	zippy.versionsPath = filepath.Join(zippy.zippyPath, "versions")
	zippy.storagePath = filepath.Join(zippy.zippyPath, "storage")

	// Check if already initialized
	if _, err := os.Stat(zippy.zippyPath); !os.IsNotExist(err) {
		fmt.Println("Zippy repository already initialized!")
		return
	}

	// Prompt for repo details
	reader := bufio.NewReader(os.Stdin)
	defaultName := filepath.Base(cwd)
	fmt.Printf("Repository name [%s]: ", defaultName)
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	if name == "" {
		name = defaultName
	}
	fmt.Print("Author [Unknown]: ")
	author, _ := reader.ReadString('\n')
	author = strings.TrimSpace(author)
	if author == "" {
		author = "Unknown"
	}
	fmt.Print("Description [Zippy Repository]: ")
	desc, _ := reader.ReadString('\n')
	desc = strings.TrimSpace(desc)
	if desc == "" {
		desc = "Zippy Repository"
	}

	// Create directories
	dirs := []string{zippy.zippyPath, zippy.versionsPath, zippy.storagePath}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Error creating directory %s: %v\n", dir, err)
			return
		}
	}

	// Create config
	config := RepoConfig{
		Name:        name,
		Author:      author, // Could get from git config or environment
		Created:     time.Now(),
		Description: desc,
	}

	configData, _ := json.MarshalIndent(config, "", "  ")
	if err := os.WriteFile(zippy.configPath, configData, 0644); err != nil {
		fmt.Printf("Error creating config: %v\n", err)
		return
	}

	// Create sample .zippyignore
	zippyignoreContent := `# Zippy ignore file
# Ignore version control directory
.zippy/

# Common files to ignore
*.log
*.tmp
.DS_Store
Thumbs.db
node_modules/
*.exe
*.dll
*.so
.env
.env.local

# Add your patterns here
`
	zippyignorePath := filepath.Join(cwd, ".zippyignore")
	os.WriteFile(zippyignorePath, []byte(zippyignoreContent), 0644)

	fmt.Printf("Initialized Zippy repository in %s\n", cwd)
	fmt.Println("Created .zippyignore file - edit it to specify files to ignore")
}

func (zippy *Zippy) loadZippyIgnore() *ZippyIgnore {
	zippyignore := &ZippyIgnore{}
	zippyignorePath := filepath.Join(zippy.repoPath, ".zippyignore")
	
	file, err := os.Open(zippyignorePath)
	if err != nil {
		return zippyignore // Return empty if file doesn't exist
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if line != "" && !strings.HasPrefix(line, "#") {
			zippyignore.patterns = append(zippyignore.patterns, line)
		}
	}

	return zippyignore
}

func (zippyignore *ZippyIgnore) shouldIgnore(filePath string) bool {
	// Normalize filePath to use forward slashes for matching
	filePath = filepath.ToSlash(filePath)
	for _, pattern := range zippyignore.patterns {
		pattern = strings.TrimSpace(pattern)
		if pattern == "" {
			continue
		}
		// Normalize pattern
		pattern = filepath.ToSlash(pattern)
		// Directory pattern (ends with /)
		if strings.HasSuffix(pattern, "/") {
			if strings.HasPrefix(filePath, pattern) {
				return true
			}
			continue
		}
		// Glob pattern
		matched, err := filepath.Match(pattern, filePath)
		if err == nil && matched {
			return true
		}
		// Also match against just the base name (like .gitignore does)
		if matched, err := filepath.Match(pattern, filepath.Base(filePath)); err == nil && matched {
			return true
		}
	}
	return false
}

// Update addFiles to write staged files to .zippy/stage.json
func (zippy *Zippy) addFiles(paths []string) {
	fmt.Println("Adding files to staging area...")
	zippyignore := zippy.loadZippyIgnore()
	staged := map[string]struct{}{}
	if len(paths) == 1 && paths[0] == "." {
		filepath.Walk(zippy.repoPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			relPath, _ := filepath.Rel(zippy.repoPath, path)
			if relPath == "." {
				return nil
			}
			if zippyignore.shouldIgnore(relPath) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			if !info.IsDir() {
				fmt.Printf("  Added: %s\n", relPath)
				staged[relPath] = struct{}{}
			}
			return nil
		})
	} else {
		for _, p := range paths {
			absPath := filepath.Join(zippy.repoPath, p)
			info, err := os.Stat(absPath)
			if err != nil {
				fmt.Printf("  [Not found or inaccessible]: %s\n", p)
				continue
			}
			if zippyignore.shouldIgnore(p) {
				fmt.Printf("  [Ignored]: %s\n", p)
				continue
			}
			if info.IsDir() {
				filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					relPath, _ := filepath.Rel(zippy.repoPath, path)
					if relPath == "." {
						return nil
					}
					if zippyignore.shouldIgnore(relPath) {
						if info.IsDir() {
							return filepath.SkipDir
						}
						return nil
					}
					if !info.IsDir() {
						fmt.Printf("  Added: %s\n", relPath)
						staged[relPath] = struct{}{}
					}
					return nil
				})
			} else {
				fmt.Printf("  Added: %s\n", p)
				staged[p] = struct{}{}
			}
		}
	}
	if len(staged) == 0 {
		fmt.Println("No files added.")
	} else {
		// Save staged files to .zippy/stage.json
		stageList := make([]string, 0, len(staged))
		for f := range staged {
			stageList = append(stageList, f)
		}
		sort.Strings(stageList)
		data, _ := json.MarshalIndent(stageList, "", "  ")
		os.WriteFile(zippy.stagePath, data, 0644)
	}
}

func (zippy *Zippy) commit() {
	// Parse commit message and version tag from args
	message := "No message"
	version := fmt.Sprintf("v%d", time.Now().Unix())
	for i, arg := range os.Args {
		if arg == "-m" && i+1 < len(os.Args) {
			message = os.Args[i+1]
		}
		if arg == "-v" && i+1 < len(os.Args) {
			version = os.Args[i+1]
		}
	}
	fmt.Printf("Creating version %s: %s\n", version, message)
	// Read staged files
	stageList := []string{}
	if data, err := os.ReadFile(zippy.stagePath); err == nil {
		json.Unmarshal(data, &stageList)
	}
	if len(stageList) == 0 {
		fmt.Println("No files staged. Use 'zippy add <files>' to stage files.")
		return
	}
	// Create zip file
	zipPath := filepath.Join(zippy.storagePath, version+".zip")
	if err := zippy.createZipFileWithList(zipPath, stageList); err != nil {
		fmt.Printf("Error creating zip: %v\n", err)
		return
	}
	// Save version metadata
	versionInfo := Version{
		Tag:       version,
		Message:   message,
		Timestamp: time.Now(),
		Author:    "User", // Could get from config
		ZipPath:   zipPath,
		FilesCount: len(stageList),
	}
	// Get file info
	if stat, err := os.Stat(zipPath); err == nil {
		versionInfo.Size = stat.Size()
	}
	zippy.saveVersionInfo(versionInfo)
	fmt.Printf("Version %s created successfully!\n", version)
	// Clear staging area
	os.Remove(zippy.stagePath)
}

func (zippy *Zippy) createZipFile(zipPath string) error {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	writer := zip.NewWriter(zipFile)
	defer writer.Close()

	zippyignore := zippy.loadZippyIgnore()
	fileCount := 0

	err = filepath.Walk(zippy.repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(zippy.repoPath, path)
		
		// Skip if should be ignored
		if zippyignore.shouldIgnore(relPath) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if info.IsDir() {
			return nil
		}

		// Add file to zip
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		f, err := writer.Create(relPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}

		fileCount++
		return nil
	})

	fmt.Printf("Packed %d files into %s\n", fileCount, filepath.Base(zipPath))
	return err
}

func (zippy *Zippy) createZipFileWithList(zipPath string, files []string) error {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()
	writer := zip.NewWriter(zipFile)
	defer writer.Close()
	for _, relPath := range files {
		absPath := filepath.Join(zippy.repoPath, relPath)
		info, err := os.Stat(absPath)
		if err != nil {
			continue
		}
		if info.IsDir() {
			// Add all files in the directory
			filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				rel, _ := filepath.Rel(zippy.repoPath, path)
				if info.IsDir() {
					return nil
				}
				file, err := os.Open(path)
				if err != nil {
					return err
				}
				defer file.Close()
				f, err := writer.Create(rel)
				if err != nil {
					return err
				}
				_, err = io.Copy(f, file)
				return err
			})
		} else {
			file, err := os.Open(absPath)
			if err != nil {
				continue
			}
			defer file.Close()
			f, err := writer.Create(relPath)
			if err != nil {
				continue
			}
			io.Copy(f, file)
		}
	}
	fmt.Printf("Packed %d files into %s\n", len(files), filepath.Base(zipPath))
	return nil
}

func (zippy *Zippy) push() {
	fmt.Println("Push completed! Version saved as zip file.")
}

func (zippy *Zippy) listVersions() {
	fmt.Println("Available versions:")
	fmt.Println("------------------")
	files, err := os.ReadDir(zippy.versionsPath)
	if err != nil {
		fmt.Printf("Error reading versions: %v\n", err)
		return
	}
	if len(files) == 0 {
		fmt.Println("No versions found.")
		return
	}
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}
		path := filepath.Join(zippy.versionsPath, file.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("  [Error reading %s]\n", file.Name())
			continue
		}
		var v Version
		if err := json.Unmarshal(data, &v); err != nil {
			fmt.Printf("  [Corrupt: %s]\n", file.Name())
			continue
		}
		fmt.Printf("  %s | %s | %s | %s\n", v.Tag, v.Timestamp.Format("2006-01-02 15:04:05"), v.Author, v.Message)
	}
}

func (zippy *Zippy) restore(version string, restorePath string) {
	fmt.Printf("Restoring version %s", version)
	if restorePath != "" {
		fmt.Printf(" (path: %s)", restorePath)
	}
	fmt.Println("...")
	// Load version metadata
	versionFile := filepath.Join(zippy.versionsPath, version+".json")
	data, err := os.ReadFile(versionFile)
	if err != nil {
		fmt.Printf("Version metadata not found: %v\n", err)
		return
	}
	var v Version
	if err := json.Unmarshal(data, &v); err != nil {
		fmt.Printf("Corrupt version metadata: %v\n", err)
		return
	}
	zipPath := v.ZipPath
	zipReader, err := zip.OpenReader(zipPath)
	if err != nil {
		fmt.Printf("Error opening zip: %v\n", err)
		return
	}
	defer zipReader.Close()
	// Normalize restorePath for matching
	restorePath = filepath.ToSlash(restorePath)
	restored := false
	for _, f := range zipReader.File {
		filePath := filepath.Join(zippy.repoPath, f.Name)
		zipEntryPath := filepath.ToSlash(f.Name)
		// If restorePath is set, only restore matching file or folder
		if restorePath != "" {
			if zipEntryPath != restorePath && !strings.HasPrefix(zipEntryPath, restorePath+"/") {
				continue
			}
		}
		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, 0755)
			continue
		}
		os.MkdirAll(filepath.Dir(filePath), 0755)
		outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			fmt.Printf("  [Error restoring %s]\n", f.Name)
			continue
		}
		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			fmt.Printf("  [Error reading %s]\n", f.Name)
			continue
		}
		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			fmt.Printf("  [Error writing %s]\n", f.Name)
			continue
		}
		fmt.Printf("  Restored: %s\n", f.Name)
		restored = true
	}
	if !restored && restorePath != "" {
		fmt.Printf("Error: '%s' not found in version %s.\n", restorePath, version)
	} else if restored {
		fmt.Println("Restore complete.")
	}
}

func (zippy *Zippy) status() {
	fmt.Println("Zippy repository status:")
	zippyignore := zippy.loadZippyIgnore()
	included := []string{}
	ignored := []string{}
	filepath.Walk(zippy.repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, _ := filepath.Rel(zippy.repoPath, path)
		if relPath == "." || strings.HasPrefix(relPath, ".zippy") {
			return nil
		}
		if zippyignore.shouldIgnore(relPath) {
			ignored = append(ignored, relPath)
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if !info.IsDir() {
			included = append(included, relPath)
		}
		return nil
	})
	fmt.Println("\nFiles to be committed:")
	if len(included) == 0 {
		fmt.Println("  (none)")
	} else {
		for _, f := range included {
			fmt.Printf("  %s\n", f)
		}
	}
	fmt.Println("\nIgnored files:")
	if len(ignored) == 0 {
		fmt.Println("  (none)")
	} else {
		for _, f := range ignored {
			fmt.Printf("  %s\n", f)
		}
	}
	// Advanced: Compare with latest version
	files, err := os.ReadDir(zippy.versionsPath)
	if err == nil && len(files) > 0 {
		// Find latest version by timestamp
		latestFile := ""
		latestTime := time.Time{}
		for _, file := range files {
			if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
				continue
			}
			path := filepath.Join(zippy.versionsPath, file.Name())
			data, err := os.ReadFile(path)
			if err != nil {
				continue
			}
			var v Version
			if err := json.Unmarshal(data, &v); err != nil {
				continue
			}
			if v.Timestamp.After(latestTime) {
				latestTime = v.Timestamp
				latestFile = file.Name()
			}
		}
		if latestFile != "" {
			path := filepath.Join(zippy.versionsPath, latestFile)
			data, err := os.ReadFile(path)
			if err == nil {
				var v Version
				if err := json.Unmarshal(data, &v); err == nil {
					zipReader, err := zip.OpenReader(v.ZipPath)
					if err == nil {
						defer zipReader.Close()
						versionFiles := map[string]uint32{}
						for _, f := range zipReader.File {
							if !f.FileInfo().IsDir() {
								versionFiles[f.Name] = f.CRC32
							}
						}
						currentFiles := map[string]uint32{}
						filepath.Walk(zippy.repoPath, func(path string, info os.FileInfo, err error) error {
							relPath, _ := filepath.Rel(zippy.repoPath, path)
							if relPath == "." || strings.HasPrefix(relPath, ".zippy") {
								return nil
							}
							if zippyignore.shouldIgnore(relPath) || info.IsDir() {
								return nil
							}
							crc, _ := fileCRC32(path)
							currentFiles[relPath] = crc
							return nil
						})
						added, removed, changed := []string{}, []string{}, []string{}
						for name, crc := range currentFiles {
							if vcrc, ok := versionFiles[name]; !ok {
								added = append(added, name)
							} else if vcrc != crc {
								changed = append(changed, name)
							}
						}
						for name := range versionFiles {
							if _, ok := currentFiles[name]; !ok {
								removed = append(removed, name)
							}
						}
						fmt.Println("\nCompared to latest version:")
						if len(added) > 0 {
							fmt.Println("  New files:")
							for _, f := range added {
								fmt.Printf("    + %s\n", f)
							}
						}
						if len(removed) > 0 {
							fmt.Println("  Deleted files:")
							for _, f := range removed {
								fmt.Printf("    - %s\n", f)
							}
						}
						if len(changed) > 0 {
							fmt.Println("  Modified files:")
							for _, f := range changed {
								fmt.Printf("    * %s\n", f)
							}
						}
						if len(added) == 0 && len(removed) == 0 && len(changed) == 0 {
							fmt.Println("  No changes since last version.")
						}
					}
				}
			}
		}
	}
}

func (zippy *Zippy) diff(v1, v2 string) {
	fmt.Printf("Comparing %s with %s...\n", v1, v2)
	// Load version metadata
	getZipPath := func(version string) (string, error) {
		versionFile := filepath.Join(zippy.versionsPath, version+".json")
		data, err := os.ReadFile(versionFile)
		if err != nil {
			return "", err
		}
		var v Version
		if err := json.Unmarshal(data, &v); err != nil {
			return "", err
		}
		return v.ZipPath, nil
	}
	zip1, err1 := getZipPath(v1)
	zip2, err2 := getZipPath(v2)
	if err1 != nil || err2 != nil {
		fmt.Printf("Error loading version metadata.\n")
		return
	}
	files1 := make(map[string]uint32)
	files2 := make(map[string]uint32)
	readZip := func(zipPath string, files map[string]uint32) error {
		zr, err := zip.OpenReader(zipPath)
		if err != nil {
			return err
		}
		defer zr.Close()
		for _, f := range zr.File {
			if !f.FileInfo().IsDir() {
				files[f.Name] = f.CRC32
			}
		}
		return nil
	}
	if err := readZip(zip1, files1); err != nil {
		fmt.Printf("Error reading %s: %v\n", v1, err)
		return
	}
	if err := readZip(zip2, files2); err != nil {
		fmt.Printf("Error reading %s: %v\n", v2, err)
		return
	}
	added, removed, changed := []string{}, []string{}, []string{}
	for name, crc := range files2 {
		if crc1, ok := files1[name]; !ok {
			added = append(added, name)
		} else if crc1 != crc {
			changed = append(changed, name)
		}
	}
	for name := range files1 {
		if _, ok := files2[name]; !ok {
			removed = append(removed, name)
		}
	}
	if len(added) > 0 {
		fmt.Println("Added files:")
		for _, f := range added {
			fmt.Printf("  + %s\n", f)
		}
	}
	if len(removed) > 0 {
		fmt.Println("Removed files:")
		for _, f := range removed {
			fmt.Printf("  - %s\n", f)
		}
	}
	if len(changed) > 0 {
		fmt.Println("Changed files:")
		for _, f := range changed {
			fmt.Printf("  * %s\n", f)
		}
	}
	if len(added) == 0 && len(removed) == 0 && len(changed) == 0 {
		fmt.Println("No differences found.")
	}
}

func (zippy *Zippy) saveVersionInfo(version Version) {
	// Save version metadata to JSON file
	versionFile := filepath.Join(zippy.versionsPath, version.Tag+".json")
	data, _ := json.MarshalIndent(version, "", "  ")
	os.WriteFile(versionFile, data, 0644)
}

// Add patchVersion method to Zippy
func (zippy *Zippy) patchVersion(version string, addPath string) {
	fmt.Printf("Patching version %s with %s...\n", version, addPath)
	versionFile := filepath.Join(zippy.versionsPath, version+".json")
	data, err := os.ReadFile(versionFile)
	if err != nil {
		fmt.Printf("Version metadata not found: %v\n", err)
		return
	}
	var v Version
	if err := json.Unmarshal(data, &v); err != nil {
		fmt.Printf("Corrupt version metadata: %v\n", err)
		return
	}
	zipPath := v.ZipPath
	// Create temp dir
	tempDir, err := os.MkdirTemp("", "zippy_patch_*")
	if err != nil {
		fmt.Printf("Error creating temp dir: %v\n", err)
		return
	}
	defer os.RemoveAll(tempDir)
	// Extract zip to temp dir
	zipReader, err := zip.OpenReader(zipPath)
	if err != nil {
		fmt.Printf("Error opening zip: %v\n", err)
		return
	}
	for _, f := range zipReader.File {
		outPath := filepath.Join(tempDir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(outPath, 0755)
			continue
		}
		os.MkdirAll(filepath.Dir(outPath), 0755)
		outFile, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			continue
		}
		io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
	}
	zipReader.Close()
	// Copy new file/folder into temp dir
	srcPath := filepath.Join(zippy.repoPath, addPath)
	info, err := os.Stat(srcPath)
	if err != nil {
		fmt.Printf("File/folder to add not found: %v\n", err)
		return
	}
	var copyErr error
	if info.IsDir() {
		copyErr = copyDir(srcPath, filepath.Join(tempDir, addPath))
	} else {
		os.MkdirAll(filepath.Dir(filepath.Join(tempDir, addPath)), 0755)
		copyErr = copyFile(srcPath, filepath.Join(tempDir, addPath))
	}
	if copyErr != nil {
		fmt.Printf("Error copying file/folder: %v\n", copyErr)
		return
	}
	// Re-zip everything in temp dir
	newZip, err := os.Create(zipPath)
	if err != nil {
		fmt.Printf("Error creating zip: %v\n", err)
		return
	}
	writer := zip.NewWriter(newZip)
	fileCount := 0
	filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(tempDir, path)
		if rel == "." {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		f, err := writer.Create(rel)
		if err != nil {
			return err
		}
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		_, err = io.Copy(f, in)
		in.Close()
		if err != nil {
			return err
		}
		fileCount++
		return nil
	})
	writer.Close()
	newZip.Close()
	// Update version metadata
	if stat, err := os.Stat(zipPath); err == nil {
		v.Size = stat.Size()
	}
	v.FilesCount = fileCount
	zippy.saveVersionInfo(v)
	fmt.Println("Patch complete.")
}
// Helper functions for copying files and directories
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		return copyFile(path, target)
	})
}
// Helper to get CRC32 of a file
func fileCRC32(path string) (uint32, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	h := crc32.NewIEEE()
	_, err = io.Copy(h, f)
	if err != nil {
		return 0, err
	}
	return h.Sum32(), nil
}