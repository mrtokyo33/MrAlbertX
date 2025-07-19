package light

import (
	"MrAlbertX/server/internal/ports"
	"fmt"
)

type DeleteUseCase struct {
	repo ports.LightRepositoryPort
}

func NewDeleteUseCase(repo ports.LightRepositoryPort) *DeleteUseCase {
	return &DeleteUseCase{repo: repo}
}

func (uc *DeleteUseCase) Execute(lightID string) error {
	fmt.Printf("[BUSINESS LOGIC] Attempting to delete light '%s'.\n", lightID)
	return uc.repo.Delete(lightID)
}