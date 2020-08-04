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

package mits

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

// LoadConfig loads a configuration file from where the environment variable points to.
func LoadConfig(envConfig string, config *Config) error {
	configPath, ok := os.LookupEnv(envConfig)
	if !ok {
		return fmt.Errorf("failed to load config: %q environment variable is not set", envConfig)
	}
	configFile, err := os.Open(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	defer configFile.Close()
	configDecoder := yaml.NewDecoder(configFile)
	configDecoder.SetStrict(true)
	if err := configDecoder.Decode(config); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	return nil
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
