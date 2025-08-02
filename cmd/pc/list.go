package pc

import (
	"MrAlbertX/server/cmd/dependencies"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all indexed programs in a scrollable view",
	RunE: func(cmd *cobra.Command, args []string) error {
		finder := dependencies.GetProgramFinder()
		programs, err := finder.ListAll()
		if err != nil {
			return err
		}

		if len(programs) == 0 {
			fmt.Println("No indexed programs found.")
			return nil
		}

		var buffer bytes.Buffer

		for _, p := range programs {
			buffer.WriteString(fmt.Sprintf("Name: %s\n", p.Name))
			buffer.WriteString(fmt.Sprintf("Path: %s\n", p.Path))
			if len(p.Tags) > 0 {
				buffer.WriteString(fmt.Sprintf("Tags: %s\n", strings.Join(p.Tags, ", ")))
			}
			buffer.WriteString("----------------------------------------\n")
		}

		buffer.WriteString(fmt.Sprintf("\nTotal programs found: %d\n", len(programs)))

		return runPager(buffer.String())
	},
}

func runPager(content string) error {
	pagerCmd := os.Getenv("PAGER")
	if pagerCmd == "" {
		pagerCmd = "less"
	}

	tempFile, err := os.CreateTemp("", "mr-x-list-*.tmp")
	if err != nil {
		return fmt.Errorf("could not create temporary file: %w", err)
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.WriteString(content); err != nil {
		return fmt.Errorf("could not write to temporary file: %w", err)
	}
	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("could not close temporary file: %w", err)
	}

	cmd := exec.Command(pagerCmd, "-R", tempFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
