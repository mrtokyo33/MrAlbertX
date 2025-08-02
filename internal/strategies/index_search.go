package strategies

import (
	"MrAlbertX/server/internal/core/models"
	"MrAlbertX/server/internal/ports"
	"sort"
	"strings"
)

type IndexSearchStrategy struct {
	finder ports.ProgramFinderPort
}

func NewIndexSearchStrategy(finder ports.ProgramFinderPort) ports.FindingStrategy {
	return &IndexSearchStrategy{finder: finder}
}

func (s *IndexSearchStrategy) Name() string {
	return "Index Search"
}

func (s *IndexSearchStrategy) Find(query string) ([]*ports.FindResult, error) {
	programs, err := s.finder.ListAll()
	if err != nil {
		return nil, nil
	}

	query = strings.ToLower(query)
	var scoredResults []*ports.FindResult

	for _, p := range programs {
		score := s.calculateScore(p, query)
		if score > 0 {
			scoredResults = append(scoredResults, &ports.FindResult{Program: p, Score: score})
		}
	}

	if len(scoredResults) == 0 {
		return nil, nil
	}

	sort.Slice(scoredResults, func(i, j int) bool {
		return scoredResults[i].Score > scoredResults[j].Score
	})

	return scoredResults, nil
}

func (s *IndexSearchStrategy) calculateScore(p models.Program, query string) float64 {
	name := strings.ToLower(p.Name)
	filename := strings.ToLower(p.Filename)
	var score float64

	if name == query || filename == query {
		score += 1000
	}
	if strings.HasPrefix(name, query) || strings.HasPrefix(filename, query) {
		score += 500
	}
	if strings.Contains(name, query) || strings.Contains(filename, query) {
		score += 150
	}
	queryTokens := strings.Fields(query)
	for _, token := range queryTokens {
		for _, keyword := range p.Keywords {
			if strings.Contains(strings.ToLower(keyword), token) {
				score += 50
			}
		}
	}
	return score
}
