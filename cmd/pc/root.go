package pc

import "github.com/spf13/cobra"

var pcCmd = &cobra.Command{
	Use:   "pc",
	Short: "Manage your personal computer",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func InitPcCmd() *cobra.Command {
	return pcCmd
}