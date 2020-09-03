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
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

// Load loads a configuration file from configPath.
func Load(configPath string) (*Config, error) {
	configFile, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	defer configFile.Close()
	decoder := yaml.NewDecoder(configFile)
	decoder.SetStrict(true)
	var config Config
	if err := decoder.Decode(&config); err != nil {
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
		Provisioning struct {
			OverrideParams struct {
				Enabled bool `yaml:"enabled"`
			} `yaml:"override_params"`
		} `yaml:"provisioning"`
	} `yaml:"minibroker"`

	Tests struct {
		MariaDB    TestConfig `yaml:"mariadb"`
		MongoDB    TestConfig `yaml:"mongodb"`
		MySQL      TestConfig `yaml:"mysql"`
		PostgreSQL TestConfig `yaml:"postgresql"`
		RabbitMQ   TestConfig `yaml:"rabbitmq"`
		Redis      TestConfig `yaml:"redis"`
	} `yaml:"tests"`

	Timeouts Timeouts `yaml:"timeouts"`
}

// TestConfig represents the configuration for an individual test.
type TestConfig struct {
	Enabled bool   `yaml:"enabled"`
	Class   string `yaml:"class"`
	Plan    string `yaml:"plan"`
}

// Timeouts aggregates the timeouts configuration.
type Timeouts struct {
	CFPush          time.Duration `yaml:"cf_push"`
	CFStart         time.Duration `yaml:"cf_start"`
	CFCreateService time.Duration `yaml:"cf_create_service"`
}
