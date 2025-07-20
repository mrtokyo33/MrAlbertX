package ports

import "MrAlbertX/server/internal/core/models"

type SystemProviderPort interface {
	OrganizeFolder(path string) (int, error)
	OpenProgram(path string) error
}

// The Search method now only takes one argument (query).
type ProgramFinderPort interface {
	Search(query string) ([]models.Program, error)
}

type IndexerPort interface {
	Reindex() error
}
