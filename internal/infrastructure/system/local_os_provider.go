package system

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type LocalOSProvider struct{}

func NewLocalOSProvider() *LocalOSProvider {
	return &LocalOSProvider{}
}

func (p *LocalOSProvider) OrganizeFolder(path string) (int, error) {
	files, err := os.ReadDir(path) // Troca de ioutil.ReadDir para os.ReadDir
	if err != nil {
		return 0, fmt.Errorf("could not read directory %s: %w", path, err)
	}

	movedFilesCount := 0
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		ext := filepath.Ext(file.Name())
		if ext == "" {
			continue
		}

		destFolderName := strings.TrimPrefix(ext, ".")
		destFolderPath := filepath.Join(path, destFolderName)

		if err := os.MkdirAll(destFolderPath, 0755); err != nil {
			fmt.Printf("Warning: could not create destination directory %s: %v\n", destFolderPath, err)
			continue
		}

		sourcePath := filepath.Join(path, file.Name())
		destPath := filepath.Join(destFolderPath, file.Name())

		fmt.Printf("Moving %s to %s...\n", file.Name(), destFolderPath)
		if err := os.Rename(sourcePath, destPath); err != nil {
			fmt.Printf("Warning: could not move file %s: %v\n", file.Name(), err)
			continue
		}
		movedFilesCount++
	}
	return movedFilesCount, nil
}

func (p *LocalOSProvider) OpenProgram(path string) error {
	finalPath := path

	if os.Getenv("WSL_DISTRO_NAME") != "" {
		out, err := exec.Command("wslpath", "-w", path).Output()
		if err == nil {
			finalPath = strings.TrimSpace(string(out))
		}
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" || os.Getenv("WSL_DISTRO_NAME") != "" {
		cmd = exec.Command("cmd.exe", "/C", "start", `""`, finalPath)
	} else if runtime.GOOS == "darwin" {
		cmd = exec.Command("open", finalPath)
	} else {
		cmd = exec.Command("xdg-open", finalPath)
	}

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to open path '%s' (tried '%s'): %w", path, finalPath, err)
	}

	return nil
}
