package goexercise

import (
	"it.uniroma2.dicii/goexercise/config"
	"it.uniroma2.dicii/goexercise/log"
)

func main() {
	config, err := config.GetConfig()
	if err != nil {
		log.Error("Exiting...", err)
		return
	}
}
