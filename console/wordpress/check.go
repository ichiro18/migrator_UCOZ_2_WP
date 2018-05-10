package wordpress

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/ichiro18/migrator_UCOZ_2_WP/common/services"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var Env *services.Env

var CheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Validate wordpress files",
	Long:  `Validate wordpress files`,
	Run: func(cmd *cobra.Command, args []string) {
		config := Env.Config.GetStringMapString("wordpress")
		checkFolder(config["path"])
		checkConnect(config)
	},
}

func checkFolder(path string) string {
	isExistFolder, err := afero.DirExists(Env.FileSystem, path)
	if err != nil {
		fmt.Errorf("Can't check config path. ")
	}
	if !isExistFolder {
		color.Yellow("Wordpress folder is exist. Creating...")
		Env.FileSystem.Mkdir(path, 0755)
	}
	return path
}

func checkConnect(cfg map[string]string) {
	if Env.Database == nil {
		config := services.Config{
			Protocol: cfg["database_protocol"],
			Address:  cfg["database_address"],
			Port:     cfg["database_port"],
			Host:     cfg["database_host"],
			Database: cfg["database_name"],
			Login:    cfg["database_login"],
			Password: cfg["database_password"],
		}

		Env.Database = services.ConnectORM(&config)
	}

	color.Green("Connect to Wordpress database is normal")
}
