package dependencies

import (
	"MrAlbertX/server/internal/indexer"
	"MrAlbertX/server/internal/infrastructure/mock_hardware"
	"MrAlbertX/server/internal/infrastructure/repository"
	"MrAlbertX/server/internal/infrastructure/system"
	"MrAlbertX/server/internal/ports"
	"MrAlbertX/server/internal/strategies"
	lightuc "MrAlbertX/server/internal/usecase/light"
	pcuc "MrAlbertX/server/internal/usecase/pc"
	"os"
	"path/filepath"
)

const dbFile = "mr-x-database.db"

var (
	sqliteRepo    ports.LightRepositoryPort
	programFinder ports.ProgramFinderPort
	osProvider    ports.SystemProviderPort
	hwMock        ports.LightControllerPort
	progIndexer   ports.IndexerPort
)

func init() {
	configDir, _ := os.UserConfigDir()
	cacheDir := filepath.Join(configDir, "MrAlbertX")
	os.MkdirAll(cacheDir, 0755)
	cachePath := filepath.Join(cacheDir, "program_index.json")

	sqliteRepo = repository.NewSQLiteLightRepository(dbFile)
	programFinder = system.NewProgramFinder(cachePath)
	osProvider = system.NewLocalOSProvider()
	hwMock = mock_hardware.NewConsoleLightController()
	progIndexer = indexer.NewProgramIndexer(cachePath)
}

func GetProgramFinder() ports.ProgramFinderPort {
	return programFinder
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
	indexSearcher := strategies.NewIndexSearchStrategy(programFinder)
	tagSearcher := strategies.NewTagSearchStrategy(programFinder)
	heuristicSearcher := strategies.NewHeuristicSearchStrategy(indexSearcher)

	pipeline := []ports.FindingStrategy{
		indexSearcher,
		tagSearcher,
		heuristicSearcher,
	}

	return pcuc.NewOpenProgramUseCase(pipeline, osProvider)
}

func GetReindexProgramsUseCase() *pcuc.ReindexProgramsUseCase {
	return pcuc.NewReindexProgramsUseCase(progIndexer)
}
