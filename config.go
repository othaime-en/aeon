package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Zones []ConfigZone `yaml:"zones"`
}

type ConfigZone struct {
	Name     string `yaml:"name"`
	Location string `yaml:"location"`
}

func getConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".aeon.yaml")
}

func loadZonesFromConfig() []Zone {
	configPath := getConfigPath()
	if configPath == "" {
		return nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil
	}

	zones := make([]Zone, 0, len(config.Zones))
	for _, cz := range config.Zones {
		loc, err := time.LoadLocation(cz.Location)
		if err != nil {
			continue
		}

		zones = append(zones, Zone{
			Name:     cz.Name,
			Location: loc,
		})
	}

	return zones
}

func saveZonesToConfig(zones []Zone) error {
	configPath := getConfigPath()
	if configPath == "" {
		return fmt.Errorf("could not determine config path")
	}

	configZones := make([]ConfigZone, len(zones))
	for i, z := range zones {
		configZones[i] = ConfigZone{
			Name:     z.Name,
			Location: z.Location.String(),
		}
	}

	config := Config{Zones: configZones}

	data, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}
