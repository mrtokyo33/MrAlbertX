package pc

import (
	"MrAlbertX/server/cmd/dependencies"
	"fmt"

	"github.com/spf13/cobra"
)

var organizePath string

var organizeFilesCmd = &cobra.Command{
	Use:   "organize-files",
	Short: "Organizes files in a folder into subdirectories by extension",
	Long: `Scans a folder and moves each file into a subdirectory named after its
extension (e.g., .pdf files go into a 'pdf' folder).
If no path is provided, it defaults to the user's Downloads folder.`,
	Example: `  mr-x pc organize-files
  mr-x pc organize-files --path "C:\My Messy Folder"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		uc := dependencies.GetOrganizeFilesUseCase()
		movedCount, err := uc.Execute(organizePath)
		if err != nil {
			return err
		}
		fmt.Printf("Organization complete. %d file(s) moved.\n", movedCount)
		return nil
	},
}

func init() {
	organizeFilesCmd.Flags().StringVarP(&organizePath, "path", "p", "", "Path to the folder to organize (defaults to user's Downloads folder)")
}
