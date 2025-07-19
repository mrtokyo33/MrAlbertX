package light

import (
	"MrAlbertX/server/cmd/dependencies"
	"fmt"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all available lights and their status",
	RunE: func(cmd *cobra.Command, args []string) error {
		listLightsUC := dependencies.GetListLightUseCase()
		lights, err := listLightsUC.Execute()
		if err != nil {
			return fmt.Errorf("error listing lights: %w", err)
		}
		fmt.Println("--- Available Lights ---")
		if len(lights) == 0 {
			fmt.Println("No lights found. Use 'mr-x light create --name <light_name>' to add one.")
		}
		for _, light := range lights {
			status := "OFF"
			if light.IsOn {
				status = "ON"
			}
			fmt.Printf("- ID: %-25s | Status: %s\n", light.ID, status)
		}
		fmt.Println("------------------------")
		return nil
	},
}

func init() {
	lightCmd.AddCommand(listCmd)
}