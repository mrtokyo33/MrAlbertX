package strategies

import (
	"MrAlbertX/server/internal/ports"
	"strings"
)

type TagSearchStrategy struct {
	finder ports.ProgramFinderPort
}

func NewTagSearchStrategy(finder ports.ProgramFinderPort) ports.FindingStrategy {
	return &TagSearchStrategy{finder: finder}
}

func (s *TagSearchStrategy) Name() string {
	return "Tag Search"
}

func (s *TagSearchStrategy) Find(query string) ([]*ports.FindResult, error) {
	programs, err := s.finder.ListAll()
	if err != nil {
		return nil, nil
	}

	query = strings.ToLower(query)
	var results []*ports.FindResult

	for _, p := range programs {
		for _, tag := range p.Tags {
			if tag == query {
				results = append(results, &ports.FindResult{
					Program: p,
					Score:   2000, // High score for a direct tag match
				})
				break // Move to next program once a match is found
			}
		}
	}

	if len(results) == 0 {
		return nil, nil
	}

	return results, nil
}
