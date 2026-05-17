package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// This file handles reading and writing to the json file

// "db_url": "connection_string_goes_here",
// "current_user_name": "username_goes_here"

const configFileName = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c *Config) SetUser(newUserName string) error {
	// 1. Update current user name
	c.CurrentUserName = newUserName

	// 2. Write the config struct to the json file
	err := write(*c)
	if err != nil {
		return fmt.Errorf("failed to write config struct to json file: %v\n", err)
	}

	return nil
}

// Helper function to get full filepath

func getConfigFilePath() (string, error) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve home directory: %v\n", err)
	}

	configPath := filepath.Join(userHome, configFileName)

	return configPath, nil
}

// Helper function that abstracts out the write function in SetUser

func write(cfg Config) error {
	// Step 1: Get the file path
	configPath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("could not retrieve config path: %v\n", err)
	}

	// Step 2: Marshal the struct into JSON
	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config to JSON: %v\n", err)
	}

	// Step 3: Write JSON to the file
	err = os.WriteFile(configPath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write JSON data to disk: %v\n", err)
	}

	return nil

}

// Reads the json file and returns a Config struct
func Read() (Config, error) {
	// Gets the location of the user's home directory

	// Appends the config file to the full path so we can read it
	configPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, fmt.Errorf("could not retrieve config path: %\n", err)
	}

	// Decode the JSON string into a new Config struct
	// Open the file, defer its close until the end of the function
	jsonFile, err := os.Open(configPath)
	if err != nil {
		return Config{}, fmt.Errorf("could not open config path: %v\n", err)
	}
	defer jsonFile.Close()

	// Read the file into a byte slice
	byteSliceJson, err := io.ReadAll(jsonFile)
	if err != nil {
		return Config{}, fmt.Errorf("could not read from json file data: %v\n", err)
	}

	// Declare new config struct
	var newConfig Config

	// Unmarshal the json data into the struct
	if err := json.Unmarshal(byteSliceJson, &newConfig); err != nil {
		return Config{}, fmt.Errorf("Error unmarshalling JSON: %v\n", err)
	}

	return newConfig, nil
}
