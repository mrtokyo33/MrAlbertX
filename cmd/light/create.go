package light

import (
	"MrAlbertX/server/cmd/dependencies"
	"fmt"
	"github.com/spf13/cobra"
)

var createLightName string

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a new light",
	Example: `  mr-x light create --name kitchen`,
	RunE: func(cmd *cobra.Command, args []string) error {
		createLightUC := dependencies.GetCreateLightUseCase()
		if err := createLightUC.Execute(createLightName); err != nil {
			return fmt.Errorf("error creating light: %w", err)
		}
		fmt.Printf("Light '%s' created successfully.\n", createLightName)
		return nil
	},
}

func init() {
	lightCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&createLightName, "name", "n", "", "Name of the light to create (required)")
	createCmd.MarkFlagRequired("name")
}