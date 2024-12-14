// Package config Handles configuration loading
package config

import (
	"fmt"
	"github.com/spf13/viper"
	"it.uniroma2.dicii/goexercise/log"
)

// Host holds a host's address and port number
type Host struct {
	Address string `mapstructure:"address"`
	Port    int64  `mapstructure:"port"`
}

// Config holds the hosts in the configuration
type Config struct {
	Master         Host   `mapstructure:"master"`
	Workers        []Host `mapstructure:"workers"`
	MappersNumber  int    `mapstructure:"mappers_number"`
	ReducersNumber int    `mapstructure:"reducers_number"`
}

const configFile = "config"

var (
	// Current configuration used within the application
	config *Config = nil
)

// GetMaster retrieves the master value from the current configuration
func GetMaster(file ...string) (*Host, error) {
	err := getConfig(file...)
	if err != nil {
		log.Error("Unable to get master from file", err)
		return nil, err
	}
	return &config.Master, nil
}

// GetWorkers retrieves the workers from the current configuration
func GetWorkers(file ...string) (*[]Host, error) {
	err := getConfig(file...)
	if err != nil {
		log.Error("Unable to get workers from file", err)
		return nil, err
	}
	return &config.Workers, nil
}

// GetMappersNumber retrieves the number of mappers from the current configuration
func GetMappersNumber(file ...string) (int, error) {
	err := getConfig(file...)
	if err != nil {
		log.Error("Unable to get mappers number from file", err)
		return -1, err
	}
	return config.MappersNumber, nil
}

// GetReducersNumber retrieves the number of mappers from the current configuration
func GetReducersNumber(file ...string) (int, error) {
	err := getConfig(file...)
	if err != nil {
		log.Error("Unable to get reducers number from file", err)
		return -1, err
	}
	return config.ReducersNumber, nil
}

// getConfig returns the current configuration
// If not loaded yet, loads it from file
func getConfig(file ...string) error {
	if config == nil {
		cfg, err := initConfiguration(file...)
		if err != nil {
			log.Error("Unable to load configuration", err)
			return err
		}
		config = cfg
	}
	return nil
}

// initConfiguration Reads and parses the configuration file
func initConfiguration(file ...string) (*Config, error) {
	var fileName string
	if len(file) == 0 {
		fileName = configFile
	} else {
		fileName = file[0]
		if len(file) > 1 {
			for _, f := range file[1:] {
				viper.AddConfigPath(f)
			}
		}
	}
	viper.SetConfigName(fileName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../config")
	if err := viper.ReadInConfig(); err != nil {
		log.Error("Error reading configuration file", err)
		return nil, err
	}
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Error("Error parsing configuration file", err)
		return nil, err
	}
	log.Info(fmt.Sprintf("conf: %v", config))
	return &config, nil
}
