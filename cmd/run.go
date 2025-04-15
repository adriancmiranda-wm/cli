package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the WM API",
	Long:  "Run the WM API",
	Example: `
	# Generate the schema for the WM API
	wm schema

	# Generate the schema for the WM API and save to a file
	wm schema -o /path/to/schema.json
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Running the WM API")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
