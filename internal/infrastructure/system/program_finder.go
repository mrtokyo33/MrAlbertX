package system

import (
	"MrAlbertX/server/internal/core/models"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"unicode"
)

type ProgramFinder struct {
	cachePath string
}

type scoredResult struct {
	program models.Program
	score   int
}

func tokenize(s string) []string {
	var tokens []string
	for _, word := range strings.Fields(s) {
		if word == "" {
			continue
		}
		var currentToken strings.Builder
		runes := []rune(word)
		for i, r := range runes {
			isUpper := unicode.IsUpper(r)
			isLower := unicode.IsLower(r)
			if i > 0 {
				prevIsLower := unicode.IsLower(runes[i-1])
				prevIsUpper := unicode.IsUpper(runes[i-1])
				if isUpper && prevIsLower {
					if currentToken.Len() > 0 {
						tokens = append(tokens, currentToken.String())
					}
					currentToken.Reset()
				} else if isLower && prevIsUpper && currentToken.Len() > 1 {
					lastChar := runes[i-1]
					currentToken.Reset()
					currentToken.WriteRune(lastChar)
				}
			}
			currentToken.WriteRune(r)
		}
		if currentToken.Len() > 0 {
			tokens = append(tokens, currentToken.String())
		}
	}
	return tokens
}

func NewProgramFinder(cachePath string) *ProgramFinder {
	return &ProgramFinder{cachePath: cachePath}
}

func (f *ProgramFinder) Search(query string) ([]models.Program, error) {
	verbose := os.Getenv("MRX_VERBOSE_SEARCH") == "1"

	bytes, err := os.ReadFile(f.cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("program index not found. Please run 'mr-x pc index' first")
		}
		return nil, err
	}

	var programs []models.Program
	if err := json.Unmarshal(bytes, &programs); err != nil {
		return nil, fmt.Errorf("failed to parse program index: %w", err)
	}

	queryLower := strings.ToLower(query)
	queryTokens := tokenize(query)
	var scoredResults []scoredResult

	if verbose {
		fmt.Printf("\n--- VERBOSE SEARCH FOR: '%s' ---\n", query)
	}

	for _, p := range programs {
		currentScore := 0
		nameLower := strings.ToLower(p.Name)
		nameTokens := tokenize(p.Name)

		if verbose {
			fmt.Printf("\n[Evaluating]: %s\n", p.Name)
		}

		// 1. Exact Match
		if nameLower == queryLower {
			currentScore += 1000
			if verbose {
				fmt.Printf("  -> Exact Match: +1000\n")
			}
		}

		// 2. Prefix Match
		if strings.HasPrefix(nameLower, queryLower) {
			currentScore += 500
			if verbose {
				fmt.Printf("  -> Prefix Match: +500\n")
			}
		}

		// 3. Substring Match
		if strings.Contains(nameLower, queryLower) {
			currentScore += 100
			if verbose {
				fmt.Printf("  -> Substring Match: +100\n")
			}
		}

		// 4. Token-based matching
		if len(queryTokens) > 0 {
			matchedTokens := 0
			for _, queryToken := range queryTokens {
				for _, nameToken := range nameTokens {
					if strings.HasPrefix(strings.ToLower(nameToken), strings.ToLower(queryToken)) {
						matchedTokens++
						currentScore += 200
						if verbose {
							fmt.Printf("  -> Token Match ('%s' in '%s'): +200\n", queryToken, nameToken)
						}
						break
					}
				}
			}
			if matchedTokens == len(queryTokens) {
				currentScore += 400
				if verbose {
					fmt.Printf("  -> All Tokens Matched Bonus: +400\n")
				}
			}
		}

		if currentScore > 0 {
			if verbose {
				fmt.Printf("  ==> FINAL SCORE: %d\n", currentScore)
			}
			scoredResults = append(scoredResults, scoredResult{program: p, score: currentScore})
		}
	}

	if len(scoredResults) == 0 {
		return nil, nil
	}

	sort.Slice(scoredResults, func(i, j int) bool {
		return scoredResults[i].score > scoredResults[j].score
	})

	var finalResults []models.Program
	for _, sr := range scoredResults {
		finalResults = append(finalResults, sr.program)
	}

	if verbose {
		fmt.Println("\n--- TOP 5 RESULTS ---")
		for i, r := range scoredResults {
			if i >= 5 {
				break
			}
			fmt.Printf("[%d] %s (Score: %d)\n", i+1, r.program.Name, r.score)
		}
		fmt.Println("--------------------")
	}

	const maxResults = 5
	if len(finalResults) > maxResults {
		return finalResults[:maxResults], nil
	}
	return finalResults, nil
}
