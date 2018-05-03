package console

import (
	"github.com/spf13/cobra"
	ucoz "github.com/ichiro18/migrator_UCOZ_2_WP/console/ucoz"
	"github.com/fatih/color"
	"github.com/ichiro18/migrator_UCOZ_2_WP/common/services"
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
	ucoz.Env = env
	// === Child ===
	// Check
	ucozCmd.AddCommand(ucoz.CheckCmd)

}