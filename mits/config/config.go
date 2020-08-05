/*
   Copyright 2020 SUSE

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package config

import (
	"fmt"
	"io"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

// ConfigLoader wraps the Load method for loading a configuration file.
type ConfigLoader struct {
	openFile         func(name string) (io.ReadCloser, error)
	newConfigDecoder func(r io.Reader) configDecoder
}

// NewYAMLConfigLoader constructs a new ConfigLoader for loading YAML files.
func NewYAMLConfigLoader() *ConfigLoader {
	return &ConfigLoader{
		openFile: openFile,
		newConfigDecoder: func(r io.Reader) configDecoder {
			decoder := yaml.NewDecoder(r)
			decoder.SetStrict(true)
			return decoder
		},
	}
}

func openFile(name string) (io.ReadCloser, error) {
	return os.Open(name)
}

type configDecoder interface {
	Decode(interface{}) error
}

// Load loads a configuration file from configPath.
func (cl *ConfigLoader) Load(configPath string) (*Config, error) {
	configFile, err := cl.openFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	defer configFile.Close()
	configDecoder := cl.newConfigDecoder(configFile)
	var config Config
	if err := configDecoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return &config, nil
}

// Config is the top-level configuration definition for the test suite.
type Config struct {
	CF struct {
		API struct {
			Endpoint string `yaml:"endpoint"`
		} `yaml:"api"`
		Admin struct {
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		} `yaml:"admin"`
	} `yaml:"cf"`

	Minibroker struct {
		API struct {
			Endpoint string `yaml:"endpoint"`
		} `yaml:"api"`
	} `yaml:"minibroker"`

	Tests struct {
		MariaDB    TestConfig `yaml:"mariadb"`
		MySQL      TestConfig `yaml:"mysql"`
		PostgreSQL TestConfig `yaml:"postgresql"`
		Redis      TestConfig `yaml:"redis"`
	} `yaml:"tests"`

	Timeouts struct {
		CFPush          time.Duration `yaml:"cf_push"`
		CFStart         time.Duration `yaml:"cf_start"`
		CFCreateService time.Duration `yaml:"cf_create_service"`
	} `yaml:"timeouts"`
}

// TestConfig represents the configuration for an individual test.
type TestConfig struct {
	Enabled bool   `yaml:"enabled"`
	Class   string `yaml:"class"`
	Plan    string `yaml:"plan"`
}
