package main

import (
	"./common/services"
)

func main() {
	env := services.NewEnvService()
	env.LoadConfig()
}
