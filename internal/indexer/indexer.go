package indexer

import (
	"MrAlbertX/server/internal/core/models"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	lnk "github.com/parsiya/golnk"
)

type ProgramIndexer struct {
	cachePath string
}

func NewProgramIndexer(cachePath string) *ProgramIndexer {
	return &ProgramIndexer{cachePath: cachePath}
}

func GetDefaultScanPaths() []string {
	windowsUsername := os.Getenv("USERNAME")
	if windowsUsername == "" {
		windowsUsername = "User"
	}

	return []string{
		"/mnt/c/Program Files",
		"/mnt/c/Program Files (x86)",
		filepath.Join("/mnt/c/Users", windowsUsername, "AppData/Roaming/Microsoft/Windows/Start Menu/Programs"),
		"/mnt/c/ProgramData/Microsoft/Windows/Start Menu/Programs",
	}
}

func (i *ProgramIndexer) Reindex() error {
	fmt.Println("Starting program indexing... (This may take a moment)")
	programs := make(map[string]models.Program)
	searchPaths := GetDefaultScanPaths()

	for _, path := range searchPaths {
		if path == "" {
			continue
		}
		filepath.WalkDir(path, func(s string, d os.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			ext := strings.ToLower(filepath.Ext(s))
			if ext == ".exe" || ext == ".lnk" {
				name := strings.TrimSuffix(d.Name(), ext)
				targetPath := s
				if ext == ".lnk" {
					lnkInfo, err := lnk.File(s)
					if err == nil && lnkInfo.LinkInfo.LocalBasePath != "" {
						targetPath = lnkInfo.LinkInfo.LocalBasePath
						if lnkInfo.StringData.NameString != "" {
							name = lnkInfo.StringData.NameString
						}
					}
				}
				key := strings.ToLower(name)
				if _, exists := programs[key]; !exists {
					programs[key] = models.Program{Name: name, Path: targetPath}
				}
			}
			return nil
		})
	}

	var programList []models.Program
	for _, p := range programs {
		programList = append(programList, p)
	}

	bytes, err := json.MarshalIndent(programList, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal program list: %w", err)
	}

	tempFile, err := os.CreateTemp(filepath.Dir(i.cachePath), "index-*.json")
	if err != nil {
		return fmt.Errorf("failed to create temp file for index: %w", err)
	}

	if _, err := tempFile.Write(bytes); err != nil {
		tempFile.Close()
		os.Remove(tempFile.Name())
		return fmt.Errorf("failed to write to temp index file: %w", err)
	}
	tempFile.Close()

	if err := os.Rename(tempFile.Name(), i.cachePath); err != nil {
		return fmt.Errorf("failed to replace old index with new one: %w", err)
	}

	fmt.Printf("Indexing complete. Found %d programs.\n", len(programList))
	return nil
}
