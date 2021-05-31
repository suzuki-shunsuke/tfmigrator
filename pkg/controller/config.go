package controller

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Items []Item
}

type Item struct {
	Rule                 string
	Exclude              bool
	StateOut             string                       `yaml:"state_out"`
	ResourceName         string                       `yaml:"resource_name"`
	TFPath               string                       `yaml:"tf_path"`
	CompiledRule         CompiledRule                 `yaml:"-"`
	CompiledResourceName CompiledResourcePathComputer `yaml:"-"`
}

type Param struct {
	ConfigFilePath string
	LogLevel       string
	StatePath      string
	Items          []Item
	SkipState      bool
	DryRun         bool
}

type State struct {
	Values Values `json:"values"`
}

type Values struct {
	RootModule RootModule `json:"root_module"`
}

type RootModule struct {
	Resources []Resource `json:"resources"`
}

type Resource struct {
	Address string                 `json:"address"`
	Type    string                 `json:"type"`
	Name    string                 `json:"name"`
	Values  map[string]interface{} `json:"values"`
}

type DryRunResult struct {
	MigratedResources []MigratedResource `yaml:"migrated_resources"`
	ExcludedResources []string           `yaml:"excluded_resources"`
	NoMatchResources  []string           `yaml:"no_match_resources"`
}

type MigratedResource struct {
	SourceResourcePath string `yaml:"source_resource_path"`
	DestResourcePath   string `yaml:"dest_resource_path"`
	TFPath             string `yaml:"tf_path"`
	StateOut           string `yaml:"state_out"`
}

func (ctrl *Controller) readConfig(param Param, cfg *Config) error {
	cfgFile, err := os.Open(param.ConfigFilePath)
	if err != nil {
		return fmt.Errorf("open a configuration file %s: %w", param.ConfigFilePath, err)
	}
	defer cfgFile.Close()
	if err := yaml.NewDecoder(cfgFile).Decode(&cfg); err != nil {
		return fmt.Errorf("parse a configuration file as YAML %s: %w", param.ConfigFilePath, err)
	}
	return nil
}
