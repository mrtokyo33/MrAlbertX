package pc

import (
	"MrAlbertX/server/internal/ports"
	"fmt"
)

type ReindexProgramsUseCase struct {
	indexer ports.IndexerPort
}

func NewReindexProgramsUseCase(indexer ports.IndexerPort) *ReindexProgramsUseCase {
	return &ReindexProgramsUseCase{indexer: indexer}
}

func (uc *ReindexProgramsUseCase) Execute() error {
	fmt.Println("Action: Starting program re-indexing...")
	return uc.indexer.Reindex()
}
