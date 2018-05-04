package console

import (
	"github.com/fatih/color"
	"github.com/ichiro18/migrator_UCOZ_2_WP/common/services"
	"github.com/ichiro18/migrator_UCOZ_2_WP/console/ucoz"
	"github.com/spf13/cobra"
)

// Интерфейс работы с Ucoz
var ucozCmd = &cobra.Command{
	Use:   "ucoz",
	Short: "This command for work with ucoz",
	Long:  `This command for work with ucoz`,
}

func init() {
	RootCmd.AddCommand(ucozCmd)

	color.Yellow("==== UCOZ interface ====")
	// Load ENV
	env := services.NewEnvService()
	env.Load()

	//	Set config
	if ucoz.Env == nil {
		ucoz.Env = env
	}
	if ucoz.UcozFileStruct == nil {
		str := ucoz.NewUcozStructure()
		ucoz.UcozFileStruct = str
	}
	// === Child ===
	// Check
	ucozCmd.AddCommand(ucoz.CheckCmd)
	// News
	ucozCmd.AddCommand(ucoz.NewsCmd)

}
