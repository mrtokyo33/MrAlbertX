package light

import (
	"MrAlbertX/server/internal/core/models"
	"MrAlbertX/server/internal/ports"
	"fmt"
)

type CreateUseCase struct {
	repo ports.LightRepositoryPort
}

func NewCreateUseCase(repo ports.LightRepositoryPort) *CreateUseCase {
	return &CreateUseCase{repo: repo}
}

func (uc *CreateUseCase) Execute(lightID string) error {
	fmt.Printf("[BUSINESS LOGIC] Attempting to create light '%s'.\n", lightID)
	newLight := models.Light{
		ID:   lightID,
		IsOn: false,
	}
	return uc.repo.Save(newLight)
}