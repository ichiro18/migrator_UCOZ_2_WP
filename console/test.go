package console

import (
	"github.com/spf13/cobra"
	"strings"
	"github.com/fatih/color"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "This command for test",
	Long:  `This command for test`,
	Run: func(cmd *cobra.Command, args []string) {
		//color.Yellow("==== Test interface ====")
		//color.Red("VAR: %v", conf)
		str := "\t Hello, World\n "
		color.Yellow("Before Trim Length: %d String:%v\n", len(str), str)
		trim := strings.TrimSpace(str)
		color.Yellow("After Trim Length: %d String:%v\n", len(trim), trim)
	},
}

func init() {
	RootCmd.AddCommand(testCmd)
}
