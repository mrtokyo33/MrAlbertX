package light

import (
	"MrAlbertX/server/cmd/dependencies"
	"fmt"
	"github.com/spf13/cobra"
)

var deleteLightName string

var deleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete a light",
	Example: `  mr-x light delete --name kitchen`,
	RunE: func(cmd *cobra.Command, args []string) error {
		deleteLightUC := dependencies.GetDeleteLightUseCase()
		if err := deleteLightUC.Execute(deleteLightName); err != nil {
			return fmt.Errorf("error deleting light: %w", err)
		}
		fmt.Printf("Light '%s' deleted successfully.\n", deleteLightName)
		return nil
	},
}

func init() {
	lightCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVarP(&deleteLightName, "name", "n", "", "Name of the light to delete (required)")
	deleteCmd.MarkFlagRequired("name")
}