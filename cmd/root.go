package cmd

import (
	"MrAlbertX/server/cmd/light"
	"MrAlbertX/server/cmd/pc"
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mr-x",
	Short: "Mr. Albert X - Your Personal Home Automator",
	Long:  `A command-line tool to manage your home automation projects, built with a robust and clean architecture in Go.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(light.InitLightCmd())
	rootCmd.AddCommand(pc.InitPcCmd())
}