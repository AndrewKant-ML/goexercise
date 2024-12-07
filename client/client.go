// Package client handler client operations
package client

import (
	"fmt"
	"it.uniroma2.dicii/goexercise/config"
	"it.uniroma2.dicii/goexercise/log"
)

var (
	master           *config.Host = nil
	currentWordCount map[string]int
	configFile       string = "" // Default configuration file
)

// initClient initializes the application client
func initClient() {
	masterHost, err := config.GetMasterFromFile(configFile)
	if err != nil {
		log.Error("Unable to get master configuration", err)
		return
	}
	master = masterHost

	currentWordCount = make(map[string]int)
}

func Run() {
	input := ""
	for input != "exit" {
		n, err := fmt.Scanf("%s\n", &input)
		if err != nil {
			log.Error("Unable to read input", err)
			continue
		}
		if n <= 0 {
			log.ErrorMessage(fmt.Sprintf("An error occurred while reading from stardard input. Bytes read: %d", n))
			continue
		}
		if input == "restart" {
			log.Info("Restart word count...")
			currentWordCount = make(map[string]int)
			continue
		}
		// TODO call master gRPC
	}
}
