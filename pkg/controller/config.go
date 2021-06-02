package controller

import (
	"fmt"
	"os"

	"github.com/suzuki-shunsuke/go-template-unmarshaler/text"
	"github.com/suzuki-shunsuke/tfmigrator/pkg/expr"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Items []Item
}

type Item struct {
	Rule          *expr.Bool
	Exclude       bool
	StateDirname  *text.Template `yaml:"state_dirname"`
	StateBasename *text.Template `yaml:"state_basename"`
	ResourceName  *ResourceName  `yaml:"resource_name"`
	TFBasename    *text.Template `yaml:"tf_basename"`
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
	SourceResourcePath string         `yaml:"source_resource_path"`
	DestResourcePath   string         `yaml:"dest_resource_path"`
	TFBasename         *text.Template `yaml:"tf_basename"`
	StateDirname       *text.Template `yaml:"state_dirname"`
	StateBasename      *text.Template `yaml:"state_basename"`
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
