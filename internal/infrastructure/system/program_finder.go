package system

import (
	"MrAlbertX/server/internal/core/models"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

type ProgramFinder struct {
	cachePath string
}

type scoredResult struct {
	program models.Program
	score   int
}

func NewProgramFinder(cachePath string) *ProgramFinder {
	return &ProgramFinder{cachePath: cachePath}
}

func (f *ProgramFinder) Search(query string) ([]models.Program, error) {
	programs, err := f.ListAll()
	if err != nil {
		return nil, err
	}

	queryLower := strings.ToLower(query)
	var scoredResults []scoredResult

	for _, p := range programs {
		currentScore := 0
		nameLower := strings.ToLower(p.Name)
		filenameLower := strings.ToLower(p.Filename)

		if nameLower == queryLower || filenameLower == queryLower {
			currentScore += 1000
		}
		if strings.HasPrefix(nameLower, queryLower) || strings.HasPrefix(filenameLower, queryLower) {
			currentScore += 500
		}
		for _, tag := range p.Tags {
			if tag == queryLower {
				currentScore += 300
				break
			}
		}
		if strings.Contains(nameLower, queryLower) || strings.Contains(filenameLower, queryLower) {
			currentScore += 150
		}
		for _, keyword := range p.Keywords {
			if strings.Contains(strings.ToLower(keyword), queryLower) {
				currentScore += 50
				break
			}
		}

		if currentScore > 0 {
			scoredResults = append(scoredResults, scoredResult{program: p, score: currentScore})
		}
	}

	if len(scoredResults) == 0 {
		return nil, nil
	}

	sort.Slice(scoredResults, func(i, j int) bool {
		if scoredResults[i].score == scoredResults[j].score {
			return len(scoredResults[i].program.Name) < len(scoredResults[j].program.Name)
		}
		return scoredResults[i].score > scoredResults[j].score
	})

	var finalResults []models.Program
	for _, sr := range scoredResults {
		finalResults = append(finalResults, sr.program)
	}

	const maxResults = 10
	if len(finalResults) > maxResults {
		return finalResults[:maxResults], nil
	}
	return finalResults, nil
}

func (f *ProgramFinder) ListAll() ([]models.Program, error) {
	bytes, err := os.ReadFile(f.cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("the program index was not found. Please run 'mr-x pc index' first")
		}
		return nil, err
	}

	var programs []models.Program
	if err := json.Unmarshal(bytes, &programs); err != nil {
		return nil, fmt.Errorf("failed to parse the program index: %w", err)
	}
	return programs, nil
}
