package pc

import (
	"MrAlbertX/server/cmd/dependencies"

	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open [program_name]",
	Short: "Opens an application or file",
	Long: `Attempts to open a program using the system's default handler.
It includes aliases for common programs (e.g., 'vscode' for 'code').`,
	Example: `  mr-x pc open vscode
  mr-x pc open firefox
  mr-x pc open notepad`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		programName := args[0]
		uc := dependencies.GetOpenProgramUseCase()

		if err := uc.Execute(programName); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	pcCmd.AddCommand(openCmd)
}
