package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// BuildConfig is the build configuration for a project.
type BuildConfig struct {
	Remote    string `yaml:"remote"`
	Image     string `yaml:"image"`
	Command   string `yaml:"command"`
	MountPath string `yaml:"mountPath"`
}

// ReadConfig reads the build config and returns a BuildConfig.
func ReadConfig() (*BuildConfig, error) {
	var config BuildConfig

	yamlFile, err := ioutil.ReadFile("build.yaml")
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
