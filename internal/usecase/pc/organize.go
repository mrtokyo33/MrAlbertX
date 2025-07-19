package pc

import (
	"MrAlbertX/server/internal/ports"
	"fmt"
	"os/user"
	"path/filepath"
)

type OrganizeFilesUseCase struct {
	sysProvider ports.SystemProviderPort
}

func NewOrganizeFilesUseCase(sysProvider ports.SystemProviderPort) *OrganizeFilesUseCase {
	return &OrganizeFilesUseCase{sysProvider: sysProvider}
}

func (uc *OrganizeFilesUseCase) Execute(targetPath string) (int, error) {
	if targetPath == "" {
		usr, err := user.Current()
		if err != nil {
			return 0, fmt.Errorf("could not get current user: %w", err)
		}
		targetPath = filepath.Join(usr.HomeDir, "Downloads")
	}

	fmt.Printf("Scanning folder: %s\n", targetPath)
	return uc.sysProvider.OrganizeFolder(targetPath)
}