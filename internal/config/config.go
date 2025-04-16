package config

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mathly/internal/repository"
	"mathly/internal/service"
	"os"
	"path"
)

type Configuration struct {
	Databases repository.DatabasesConfig `json:"databases"`
	OAuth     AuthOAuth                  `json:"authOAuth"`
	Services  service.ServiceConfig      `json:"services"`
}

type AuthOAuth struct {
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
	CallbackURL  string `json:"callbackURL"`
}

func new() *Configuration {
	configuration, err := readConfigurationFromJson()
	if err != nil {
		log.Fatalf("failed to read configuration from json. Details: %v", err)
	}

	return configuration
}

func readConfigurationFromJson() (*Configuration, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get config path. Details: %v", err)
	}

	confFile, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open configuration file. Details: %v", err)
	}
	defer confFile.Close()

	conf, err := io.ReadAll(confFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file. Details: %v", err)
	}

	configuration := Configuration{}
	err = json.Unmarshal(conf, &configuration)
	if err != nil {
		return nil, fmt.Errorf("failed to parse configuration file. Details: %v", err)
	}

	return &configuration, nil
}

func getConfigPath() (string, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./configuration/config.json"
	}

	if !path.IsAbs(configPath) {
		workingDirectory, err := os.Getwd()
		if err != nil {
			return configPath, err
		}
		configPath = path.Join(workingDirectory, configPath)
	}

	return configPath, nil
}

var AppConfig = new()
