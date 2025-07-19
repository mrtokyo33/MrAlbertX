package dependencies

import (
	"MrAlbertX/server/internal/infrastructure/mock_hardware"
	"MrAlbertX/server/internal/infrastructure/repository"
	"MrAlbertX/server/internal/infrastructure/system"
	lightuc "MrAlbertX/server/internal/usecase/light"
	pcuc "MrAlbertX/server/internal/usecase/pc"
)

const dbFile = "mr-x-database.db"

func GetCreateLightUseCase() *lightuc.CreateUseCase {
	repo := repository.NewSQLiteLightRepository(dbFile)
	return lightuc.NewCreateUseCase(repo)
}

func GetControlLightUseCase() *lightuc.ControlUseCase {
	repo := repository.NewSQLiteLightRepository(dbFile)
	controller := mock_hardware.NewConsoleLightController()
	return lightuc.NewControlUseCase(repo, controller)
}

func GetDeleteLightUseCase() *lightuc.DeleteUseCase {
	repo := repository.NewSQLiteLightRepository(dbFile)
	return lightuc.NewDeleteUseCase(repo)
}

func GetListLightUseCase() *lightuc.ListUseCase {
	repo := repository.NewSQLiteLightRepository(dbFile)
	return lightuc.NewListUseCase(repo)
}

func GetOrganizeFilesUseCase() *pcuc.OrganizeFilesUseCase {
	provider := system.NewLocalFSProvider()
	return pcuc.NewOrganizeFilesUseCase(provider)
}