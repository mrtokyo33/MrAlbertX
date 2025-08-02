package indexer

import (
	"MrAlbertX/server/internal/core/models"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	lnk "github.com/parsiya/golnk"
	"github.com/saferwall/pe"
	"github.com/schollz/progressbar/v3"
)

type ProgramIndexer struct {
	cachePath string
}

func NewProgramIndexer(cachePath string) *ProgramIndexer {
	return &ProgramIndexer{cachePath: cachePath}
}

func askForDisksToScan() []string {
	var disks []string
	entries, err := os.ReadDir("/mnt/")
	if err != nil {
		fmt.Println("Warning: Could not list disks in /mnt/. Defaulting to 'C'.")
		return []string{"/mnt/c"}
	}

	for _, e := range entries {
		if e.IsDir() && len(e.Name()) == 1 {
			disks = append(disks, strings.ToUpper(e.Name()))
		}
	}

	if len(disks) == 0 {
		fmt.Println("Warning: No disks found in /mnt/. Defaulting to 'C'.")
		return []string{"/mnt/c"}
	}

	disksToScan := []string{}
	prompt := &survey.MultiSelect{
		Message: "Which disks do you want to scan?",
		Options: disks,
	}
	survey.AskOne(prompt, &disksToScan)

	var pathsToScan []string
	for _, disk := range disksToScan {
		pathsToScan = append(pathsToScan, filepath.Join("/mnt", strings.ToLower(disk)))
	}
	return pathsToScan
}

func buildExclusionList(rootPath string) []string {
	return []string{
		filepath.Join(rootPath, "Windows"),
		filepath.Join(rootPath, "$Recycle.Bin"),
		filepath.Join(rootPath, "System Volume Information"),
		filepath.Join(rootPath, "ProgramData"),
		filepath.Join(rootPath, "Users", "All Users"),
		filepath.Join(rootPath, "Recovery"),
	}
}

func (i *ProgramIndexer) Reindex() error {
	searchPaths := askForDisksToScan()
	if len(searchPaths) == 0 {
		fmt.Println("No disks selected. Aborting.")
		return nil
	}

	fmt.Println("Mapping files for indexing... (This may take a while)")
	var filesToIndex []string

	for _, path := range searchPaths {
		exclusions := buildExclusionList(path)
		filepath.WalkDir(path, func(s string, d os.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			for _, excludedPath := range exclusions {
				if strings.HasPrefix(s, excludedPath) {
					return filepath.SkipDir
				}
			}
			if !d.IsDir() {
				ext := strings.ToLower(filepath.Ext(s))
				if ext == ".exe" || ext == ".lnk" {
					filesToIndex = append(filesToIndex, s)
				}
			}
			return nil
		})
	}

	fmt.Printf("   -> Found %d files. Starting metadata extraction.\n", len(filesToIndex))
	programs := make(map[string]models.Program)
	bar := progressbar.Default(int64(len(filesToIndex)), "Extracting metadata...")

	for _, path := range filesToIndex {
		bar.Add(1)
		ext := strings.ToLower(filepath.Ext(path))
		name := strings.TrimSuffix(filepath.Base(path), ext)
		targetPath := path
		var keywords []string

		if ext == ".lnk" {
			if lnkInfo, err := lnk.File(path); err == nil && lnkInfo.LinkInfo.LocalBasePath != "" {
				resolvedPath := lnkInfo.LinkInfo.LocalBasePath
				if !strings.HasSuffix(strings.ToLower(resolvedPath), ".exe") {
					continue
				}
				targetPath = resolvedPath
				if lnkInfo.StringData.NameString != "" {
					name = lnkInfo.StringData.NameString
				}
			} else {
				continue
			}
		}

		if strings.HasSuffix(strings.ToLower(targetPath), ".exe") {
			if peFile, err := pe.New(targetPath, &pe.Options{}); err == nil {
				defer peFile.Close()
				if res, err := peFile.ParseVersionResources(); err == nil {
					if val, ok := res["ProductName"]; ok && val != "" {
						name = val
					}
					if val, ok := res["FileDescription"]; ok {
						keywords = append(keywords, strings.Fields(val)...)
					}
					if val, ok := res["CompanyName"]; ok {
						keywords = append(keywords, strings.Fields(val)...)
					}
				}
			}
		} else {
			continue
		}

		key := strings.ToLower(targetPath)
		if _, exists := programs[key]; !exists {
			programs[key] = models.Program{
				Name:     name,
				Path:     targetPath,
				Filename: filepath.Base(targetPath),
				Keywords: append(keywords, strings.Fields(name)...),
				Tags:     inferTags(name, targetPath),
			}
		}
	}

	var programList []models.Program
	for _, p := range programs {
		programList = append(programList, p)
	}
	bytes, err := json.MarshalIndent(programList, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize program list: %w", err)
	}
	if err := os.WriteFile(i.cachePath, bytes, 0644); err != nil {
		return fmt.Errorf("failed to write index file: %w", err)
	}
	fmt.Printf("\nIndexing complete. Found %d unique programs.\n", len(programList))
	return nil
}

func GetDefaultWatcherPaths() []string {
	windowsUser := os.Getenv("USERNAME")
	if windowsUser == "" {
		windowsUser = "User"
	}
	winUserPath := "/mnt/c/Users/" + windowsUser

	return []string{
		"/mnt/c/Program Files",
		"/mnt/c/Program Files (x86)",
		filepath.Join(winUserPath, "Desktop"),
		filepath.Join(winUserPath, "AppData/Roaming/Microsoft/Windows/Start Menu/Programs"),
		"/mnt/c/ProgramData/Microsoft/Windows/Start Menu/Programs",
	}
}

func inferTags(name, path string) []string {
	lowerName := strings.ToLower(name)
	tags := []string{}
	tagMap := map[string]string{
		"steam": "game", "epic games": "game", "riot games": "game",
		"code": "development", "visual studio": "development", "pycharm": "development", "goland": "development",
		"firefox": "browser", "chrome": "browser", "edge": "browser", "brave": "browser",
		"photoshop": "design", "illustrator": "design", "figma": "design", "blender": "design",
		"office": "office", "word": "office", "excel": "office", "powerpoint": "office",
	}
	for key, tag := range tagMap {
		if strings.Contains(lowerName, key) {
			found := false
			for _, existingTag := range tags {
				if existingTag == tag {
					found = true
					break
				}
			}
			if !found {
				tags = append(tags, tag)
			}
		}
	}
	return tags
}
