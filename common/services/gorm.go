package services

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Config struct {
	Protocol string
	Address  string
	Port     string
	Host     string
	Database string
	Login    string
	Password string
}

func NewConnectORM(cfg map[string]string) *gorm.DB {
	config := Config{
		Protocol: cfg["database_protocol"],
		Address:  cfg["database_address"],
		Port:     cfg["database_port"],
		Host:     cfg["database_host"],
		Database: cfg["database_name"],
		Login:    cfg["database_login"],
		Password: cfg["database_password"],
	}
	return ConnectORM(&config)

}

func ConnectORM(cfg *Config) *gorm.DB {
	connectOption := cfg.Login + ":" + cfg.Password + "@" + cfg.Protocol + "(" + cfg.Address + ":" + cfg.Port + ")/" + cfg.Database + "?charset=utf8&parseTime=true"
	db, err := gorm.Open("mysql", connectOption)
	if err != nil {
		fmt.Errorf("Unable connect to GORM database")
	}
	return db
}
