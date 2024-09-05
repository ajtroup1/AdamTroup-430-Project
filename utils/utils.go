package utils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/ajtroup1/GoDoc/internal/models"
)

type SettingManager struct {
	Settings models.Settings
}

const settingsFile = "settings.json"

// NewSettings initializes a SettingManager with retrieved or default settings
func NewSettings() (*SettingManager, error) {
	settings, err := retrieveOrCreateSettings()
	if err != nil {
		return nil, err
	}

	return &SettingManager{
		Settings: settings,
	}, nil
}

// retrieveOrCreateSettings checks if the settings file exists, or creates a default one
func retrieveOrCreateSettings() (models.Settings, error) {
	var settings models.Settings

	// Check if settings.json exists
	if _, err := os.Stat(settingsFile); errors.Is(err, os.ErrNotExist) {
		// File doesn't exist, create default settings
		defaultSettings := models.Settings{
			ProjectName: "",
			ProjectDesc: "",
			ProjectPath: "./",
			DocGenPath:  "./",
		}

		// Save default settings to a new JSON file
		err := saveSettings(defaultSettings)
		if err != nil {
			return models.Settings{}, err
		}

		return defaultSettings, nil
	}

	// If file exists, read the settings
	data, err := ioutil.ReadFile(settingsFile)
	if err != nil {
		return models.Settings{}, err
	}

	err = json.Unmarshal(data, &settings)
	if err != nil {
		return models.Settings{}, err
	}

	return settings, nil
}

// saveSettings writes the settings to a JSON file
func saveSettings(settings models.Settings) error {
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(settingsFile, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
