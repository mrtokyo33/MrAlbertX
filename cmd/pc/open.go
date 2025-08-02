package pc

import (
	"MrAlbertX/server/cmd/dependencies"
	"strings"

	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open [program name...]",
	Short: "Intelligently finds and opens a program",
	Long:  "Finds and launches an application based on a query using a multi-layered search pipeline.",
	Example: `  mr-x pc open visual studio code
  mr-x pc open fl studio`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		programName := strings.Join(args, " ")
		uc := dependencies.GetOpenProgramUseCase()

		if err := uc.Execute(programName); err != nil {
			return err
		}
		return nil
	},
}
