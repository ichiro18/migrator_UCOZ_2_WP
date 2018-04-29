package services

import (
	"github.com/spf13/afero"
	"fmt"
	"github.com/spf13/viper"
)

type Env struct {
	FileSystem 		afero.Fs
	Config			*viper.Viper
}

var (
	defaultConfigPath = "./common/config/"
	mainConfigFile = defaultConfigPath + "config.yaml"
	defaultConfigFile = defaultConfigPath + "config.default.yaml"
)

func NewEnvService() *Env{
	env := Env{}
	//env.init()
	return &env
}

func (self *Env) GetConfigFile() (File string, Path string){
	// Find default configPath
	isExistDefaultConfigPath, err := afero.DirExists(self.FileSystem, defaultConfigPath)
	if err != nil{
		fmt.Errorf("Can't check config path. ")
	}
	if !isExistDefaultConfigPath {
		self.FileSystem.Mkdir(defaultConfigPath, 0755)
	}
	// Set config path in Config
	self.Config.Set("configPath", defaultConfigPath)

	// Find config
	isExistConfigFile, err := afero.Exists(self.FileSystem, mainConfigFile)
	if err != nil{
		fmt.Errorf("Can't check config file. ")
	}
	if !isExistConfigFile{
		// Find default config
		isExistDefaultConfigFile, err := afero.Exists(self.FileSystem, defaultConfigFile)
		if err != nil{
			fmt.Errorf("Can't check default config file. ")
		}

		if !isExistDefaultConfigFile{
			fmt.Errorf("Default config file is not exist. ")
		}

		// Copy config file
		defautConfig, err := afero.ReadFile(self.FileSystem, defaultConfigFile)
		err = afero.WriteFile(self.FileSystem, mainConfigFile, defautConfig, 0777)
		if err != nil {
			fmt.Errorf("Can't copy config file fron default. ")
		}
	}

	return "config", defaultConfigPath
}