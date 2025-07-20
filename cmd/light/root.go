package light

import "github.com/spf13/cobra"

var lightCmd = &cobra.Command{
	Use:   "light",
	Short: "Manage your lights",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func InitLightCmd() *cobra.Command {
	return lightCmd
}
