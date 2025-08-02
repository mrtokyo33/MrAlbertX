package ports

import "MrAlbertX/server/internal/core/models"

type FindResult struct {
	Program models.Program
	Score   float64
}

type FindingStrategy interface {
	Find(query string) ([]*FindResult, error)
	Name() string
}
