package console

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/ichiro18/migrator_UCOZ_2_WP/common/services"
	"github.com/fatih/color"
)

var (
	Env *services.Env
	// Флаги
	cfgFile string
)

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
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $PROJECT/common/config/.config.yaml)")
}

// initConfig читает конфиги и ENV переменные.
func initConfig() {
	Env = services.NewEnvService()
	if cfgFile != "" {
		// Используем файл конфига, установленного в флаге
		Env.Config.SetConfigFile(cfgFile)
	} else {
		config, configPath := Env.GetConfigFile()

		Env.Config.AddConfigPath(configPath)
		Env.Config.SetConfigName(config)
	}

	Env.Config.AutomaticEnv() // собираем переменные окружения ENV

	// Если конфиг найден, читаем его
	if err := Env.Config.ReadInConfig(); err == nil {
		color.White("Using config file: %v", Env.Config.ConfigFileUsed())
	}
}