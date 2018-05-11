package console

import (
	"github.com/ichiro18/migrator_UCOZ_2_WP/common/services"
	"github.com/ichiro18/migrator_UCOZ_2_WP/console/wordpress"
	"github.com/spf13/cobra"
)

// Интерфейс работы с Wordpress
var wordpressCmd = &cobra.Command{
	Use:   "wordpress",
	Short: "This command for work with wordpress",
	Long:  `This command for work with wordpress`,
}

func init() {
	RootCmd.AddCommand(wordpressCmd)
	// Load ENV
	env := services.NewEnvService()
	env.Load()
	if env.Database == nil {
		env.Database = services.NewConnectORM(env.Config.GetStringMapString("wordpress"))
	}
	//	Set config
	if wordpress.Env == nil {
		wordpress.Env = env
	}
	// === Child ===
	// Check
	wordpressCmd.AddCommand(wordpress.CheckCmd)
	wordpressCmd.AddCommand(wordpress.PostCmd)
	wordpressCmd.AddCommand(wordpress.ImageCmd)
}
