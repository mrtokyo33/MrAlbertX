package ports

import "MrAlbertX/server/internal/core/models"

type LightRepositoryPort interface {
	FindByID(id string) (*models.Light, error)
	GetAll() ([]models.Light, error)
	Save(light models.Light) error
	Update(light models.Light) error
	Delete(id string) error
}