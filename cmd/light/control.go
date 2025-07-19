package light

import (
	"MrAlbertX/server/cmd/dependencies"
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

var (
	controlLightName  string
	controlLightState string
)

var controlCmd = &cobra.Command{
	Use:     "control",
	Short:   "Control the state of a light",
	Example: `  mr-x light control --name kitchen --state on`,
	RunE: func(cmd *cobra.Command, args []string) error {
		state := strings.ToLower(controlLightState)
		if state != "on" && state != "off" {
			return fmt.Errorf("invalid state: '%s'. Please use 'on' or 'off'", controlLightState)
		}
		isOn := (state == "on")
		controlLightUC := dependencies.GetControlLightUseCase()
		if err := controlLightUC.Execute(controlLightName, isOn); err != nil {
			return fmt.Errorf("error controlling light: %w", err)
		}
		fmt.Printf("Light '%s' was successfully turned %s.\n", controlLightName, strings.ToUpper(state))
		return nil
	},
}

func init() {
	lightCmd.AddCommand(controlCmd)
	controlCmd.Flags().StringVarP(&controlLightName, "name", "n", "", "Name of the light to control (required)")
	controlCmd.Flags().StringVarP(&controlLightState, "state", "s", "", "Desired state: 'on' or 'off' (required)")
	controlCmd.MarkFlagRequired("name")
	controlCmd.MarkFlagRequired("state")
}