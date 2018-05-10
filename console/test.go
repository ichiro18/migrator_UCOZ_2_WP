package console

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "This command for test",
	Long:  `This command for test`,
	Run: func(cmd *cobra.Command, args []string) {
		color.Green("TEST")
	},
}

// --------------------------------------------------------------------

func init() {
	RootCmd.AddCommand(testCmd)
}
