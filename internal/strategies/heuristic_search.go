package strategies

import (
	"MrAlbertX/server/internal/ports"
	"strings"
)

type HeuristicSearchStrategy struct {
	indexStrategy ports.FindingStrategy
	synonyms      map[string][]string
}

func NewHeuristicSearchStrategy(indexStrategy ports.FindingStrategy) ports.FindingStrategy {
	return &HeuristicSearchStrategy{
		indexStrategy: indexStrategy,
		synonyms: map[string][]string{
			"code":      {"vscode", "visual studio code"},
			"fl":        {"flstudio", "fruity loops"},
			"fl studio": {"flstudio", "fruity loops", "fl64"},
			"word":      {"winword"},
			"excel":     {"msexcel"},
		},
	}
}

func (s *HeuristicSearchStrategy) Name() string {
	return "Heuristic Search"
}

func (s *HeuristicSearchStrategy) Find(query string) ([]*ports.FindResult, error) {
	queries := s.generateQueries(query)

	for _, q := range queries {
		results, err := s.indexStrategy.Find(q)
		if err != nil {
			continue // Try next heuristic
		}
		if len(results) > 0 {
			return results, nil
		}
	}

	return nil, nil
}

func (s *HeuristicSearchStrategy) generateQueries(original string) []string {
	queries := make(map[string]struct{})
	original = strings.ToLower(original)
	queries[original] = struct{}{}

	if synonyms, ok := s.synonyms[original]; ok {
		for _, syn := range synonyms {
			queries[syn] = struct{}{}
		}
	}

	parts := strings.Fields(original)
	if len(parts) > 1 {
		joined := strings.Join(parts, "")
		queries[joined] = struct{}{}
	}

	var uniqueQueries []string
	for q := range queries {
		uniqueQueries = append(uniqueQueries, q)
	}

	return uniqueQueries
}
