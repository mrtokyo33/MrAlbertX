package pc

import (
	"MrAlbertX/server/cmd/dependencies"

	"github.com/spf13/cobra"
)

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Scans the system for programs and builds a searchable cache",
	RunE: func(cmd *cobra.Command, args []string) error {
		uc := dependencies.GetReindexProgramsUseCase()
		return uc.Execute()
	},
}

func init() {
	pcCmd.AddCommand(indexCmd)
}
