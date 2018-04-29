package console

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ichiro18/migrator_UCOZ_2_WP/common/services"
)

var cfgFile string

// RootCmd представляет базовую команду для вызова любых команд
var RootCmd = &cobra.Command{
	Use:   "migrator",
	Short: "Migrator Ucoz to Wordpress",
	Long: `This application is modern CLI applications 
for migrate content from UCOZ to Wordpress`,
}

// Execute добавляет все дочерние команды в корневую вместе с флагами
// Вызывается в main.main(). Данная функция необходима только для rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Устанавливаем глобальные флаги
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is common/config/.config.yaml)")
}

// initConfig читает конфиги и ENV переменные.
func initConfig() {
	if cfgFile != "" {
		// Используем файл конфига, установленного в флаге
		viper.SetConfigFile(cfgFile)
	} else {

		env := services.NewEnvService()
		config, configPath := env.GetConfigFile()

		// Search config in home directory with name ".cobra-example" (without extension).
		viper.AddConfigPath(configPath)
		viper.SetConfigName(config)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}