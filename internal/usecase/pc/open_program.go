package pc

import (
	"MrAlbertX/server/internal/core/models"
	"MrAlbertX/server/internal/ports"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

type OpenProgramUseCase struct {
	pipeline    []ports.FindingStrategy
	sysProvider ports.SystemProviderPort
}

func NewOpenProgramUseCase(pipeline []ports.FindingStrategy, sysProvider ports.SystemProviderPort) *OpenProgramUseCase {
	return &OpenProgramUseCase{
		pipeline:    pipeline,
		sysProvider: sysProvider,
	}
}

func (uc *OpenProgramUseCase) Execute(query string) error {
	var finalResults []*ports.FindResult

	for _, strategy := range uc.pipeline {
		fmt.Printf("Trying Layer: %s...\n", strategy.Name())
		results, err := strategy.Find(query)
		if err != nil {
			fmt.Printf("  -> Layer %s failed: %v\n", strategy.Name(), err)
			continue
		}

		if len(results) > 0 {
			fmt.Printf("  -> Layer %s found %d potential match(es).\n", strategy.Name(), len(results))
			finalResults = results
			break
		}
		fmt.Printf("  -> No results from this layer.\n")
	}

	if len(finalResults) == 0 {
		return fmt.Errorf("could not find any program matching '%s' after trying all strategies", query)
	}

	var programToLaunch models.Program
	isConfident := len(finalResults) == 1 || (finalResults[0].Score > finalResults[1].Score*1.5)

	if isConfident {
		programToLaunch = finalResults[0].Program
	} else {
		var options []string
		programMap := make(map[string]models.Program)
		for i, res := range finalResults {
			if i >= 5 {
				break
			} // Limit options
			optionStr := fmt.Sprintf("%s (%s)", res.Program.Name, res.Program.Filename)
			options = append(options, optionStr)
			programMap[optionStr] = res.Program
		}

		selection := ""
		prompt := &survey.Select{
			Message: "Found multiple close matches, please choose one:",
			Options: options,
		}
		survey.AskOne(prompt, &selection)

		if selection == "" {
			fmt.Println("Aborted.")
			return nil
		}
		programToLaunch = programMap[selection]
	}

	fmt.Printf("Launching '%s'...\n", programToLaunch.Name)
	return uc.sysProvider.OpenProgram(programToLaunch.Path)
}
