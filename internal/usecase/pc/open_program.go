package pc

import (
	"MrAlbertX/server/internal/core/models"
	"MrAlbertX/server/internal/ports"
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type OpenProgramUseCase struct {
	finder      ports.ProgramFinderPort
	sysProvider ports.SystemProviderPort
	indexer     ports.IndexerPort
	aliases     map[string]string
	cachePath   string
}

func NewOpenProgramUseCase(finder ports.ProgramFinderPort, sysProvider ports.SystemProviderPort, indexer ports.IndexerPort, cachePath string) *OpenProgramUseCase {
	return &OpenProgramUseCase{
		finder:      finder,
		sysProvider: sysProvider,
		indexer:     indexer,
		cachePath:   cachePath,
		aliases: map[string]string{
			"vscode":    "code",
			"browser":   "firefox",
			"edge":      "msedge",
			"flstudio":  "FL64",
			"fl studio": "FL64",
		},
	}
}

func (uc *OpenProgramUseCase) Execute(query string) error {
	info, err := os.Stat(uc.cachePath)
	if os.IsNotExist(err) || time.Since(info.ModTime()) > 24*time.Hour {
		log.Println("Program index is missing or stale. Forcing re-index...")
		if err := uc.indexer.Reindex(); err != nil {
			return fmt.Errorf("failed to perform fallback re-index: %w", err)
		}
	}

	searchQuery := query
	if aliasTarget, isAlias := uc.aliases[strings.ToLower(query)]; isAlias {
		fmt.Printf("Alias found: '%s' -> '%s'. Searching...\n", query, aliasTarget)
		searchQuery = aliasTarget
	}

	results, err := uc.finder.Search(searchQuery)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		return fmt.Errorf("no program found matching '%s'", query)
	}

	var programToLaunch models.Program
	if len(results) == 1 {
		programToLaunch = results[0]
		fmt.Printf("Found one match: '%s'. Launching...\n", programToLaunch.Name)
	} else {
		fmt.Println("Found multiple matches. Please choose one:")
		for i, p := range results {
			fmt.Printf("  [%d] %s\n", i+1, p.Name)
		}

		fmt.Print("Enter number: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		choice, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil || choice < 1 || choice > len(results) {
			return fmt.Errorf("invalid choice")
		}
		programToLaunch = results[choice-1]
	}

	fmt.Printf("Action: Attempting to open '%s'...\n", programToLaunch.Name)
	return uc.sysProvider.OpenProgram(programToLaunch.Path)
}
