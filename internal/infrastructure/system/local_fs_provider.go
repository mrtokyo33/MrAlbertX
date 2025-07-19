package system

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type LocalFSProvider struct{}

func NewLocalFSProvider() *LocalFSProvider {
	return &LocalFSProvider{}
}

func (p *LocalFSProvider) OrganizeFolder(path string) (int, error) {
	files, err := ioutil.ReadDir(path)
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
			continue // Pula arquivos sem extensão
		}

		// Cria o nome da pasta a partir da extensão (ex: ".pdf" -> "pdf")
		destFolderName := strings.TrimPrefix(ext, ".")
		destFolderPath := filepath.Join(path, destFolderName)

		// Cria a pasta de destino se ela não existir
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