package dependencies

import (
	"MrAlbertX/server/internal/indexer"
	"MrAlbertX/server/internal/infrastructure/mock_hardware"
	"MrAlbertX/server/internal/infrastructure/repository"
	"MrAlbertX/server/internal/infrastructure/system"
	lightuc "MrAlbertX/server/internal/usecase/light"
	pcuc "MrAlbertX/server/internal/usecase/pc"
	"os"
	"path/filepath"
)

const dbFile = "mr-x-database.db"

var (
	sqliteRepo    = repository.NewSQLiteLightRepository(dbFile)
	programFinder = system.NewProgramFinder(getCachePath())
	osProvider    = system.NewLocalOSProvider()
	hwMock        = mock_hardware.NewConsoleLightController()
	progIndexer   = indexer.NewProgramIndexer(getCachePath())
)

func getCachePath() string {
	configDir, _ := os.UserConfigDir()
	cacheDir := filepath.Join(configDir, "MrAlbertX")
	os.MkdirAll(cacheDir, 0755)
	return filepath.Join(cacheDir, "program_index.json")
}

func GetCreateLightUseCase() *lightuc.CreateUseCase {
	return lightuc.NewCreateUseCase(sqliteRepo)
}

func GetControlLightUseCase() *lightuc.ControlUseCase {
	return lightuc.NewControlUseCase(sqliteRepo, hwMock)
}

func GetDeleteLightUseCase() *lightuc.DeleteUseCase {
	return lightuc.NewDeleteUseCase(sqliteRepo)
}

func GetListLightUseCase() *lightuc.ListUseCase {
	return lightuc.NewListUseCase(sqliteRepo)
}

func GetOrganizeFilesUseCase() *pcuc.OrganizeFilesUseCase {
	return pcuc.NewOrganizeFilesUseCase(osProvider)
}

func GetOpenProgramUseCase() *pcuc.OpenProgramUseCase {
	return pcuc.NewOpenProgramUseCase(programFinder, osProvider, progIndexer, getCachePath())
}

func GetReindexProgramsUseCase() *pcuc.ReindexProgramsUseCase {
	return pcuc.NewReindexProgramsUseCase(progIndexer)
}
