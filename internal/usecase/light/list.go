package light

import (
	"MrAlbertX/server/internal/core/models"
	"MrAlbertX/server/internal/ports"
)

type ListUseCase struct {
	repo ports.LightRepositoryPort
}

func NewListUseCase(repo ports.LightRepositoryPort) *ListUseCase {
	return &ListUseCase{repo: repo}
}

func (uc *ListUseCase) Execute() ([]models.Light, error) {
	return uc.repo.GetAll()
}