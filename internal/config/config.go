package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (config *Config) SetUser(username string) error {
	// Update the current username
	config.CurrentUserName = username

	// Write the updated config file data
	return writeConfigFile(*config)
}

func ReadConfig() (Config, error) {
	var config Config

	// Get the path to the config file
	fileName, err := getConfigFilePath()
	if err != nil {
		return config, fmt.Errorf("Failed to get config file path: %w\n", err)
	}

	// open the file
	configFile, err := os.Open(fileName)
	if err != nil {
		return config, fmt.Errorf("Failed to open config file: %w\n", err)
	}
	defer configFile.Close()

	// Read all the file’s content
	bytes, err := io.ReadAll(configFile)
	if err != nil {
		return config, fmt.Errorf("Failed to read data from config file: %w\n", err)
	}

	// Unmarshal JSON into the struct
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return config, fmt.Errorf("Error converting JSON to Config struct: %w\n", err)
	}

	return config, nil
}

func getConfigFilePath() (string, error) {
	// Get the path to the config file in the users home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configFileFullPath := filepath.Join(homeDir, configFileName)
	return configFileFullPath, nil
}

func writeConfigFile(cfg Config) error {
	// Serialize struct to JSON with indentation
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("Failed to convert config to JSON: %w\n", err)
	}

	// Get the path to the config file
	configFile, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("Failed to get config file path: %w\n", err)
	}

	// Write JSON data to file
	err = os.WriteFile(configFile, data, 0644)
	if err != nil {
		return fmt.Errorf("Failed to write data to config file: %w\n", err)
	}

	return nil
}
