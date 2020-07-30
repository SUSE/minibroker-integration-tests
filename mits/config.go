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

// LoadConfig loads a configuration file from where the environment variable points to, into the
// config interface{}.
func LoadConfig(envConfig string, config interface{}) error {
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

// TestsConfig represents the root of the tests configuration document.
type TestsConfig struct {
	MariaDB    TestConfig `yaml:"mariadb"`
	MySQL      TestConfig `yaml:"mysql"`
	PostgreSQL TestConfig `yaml:"postgresql"`
	Redis      TestConfig `yaml:"redis"`
}

// TestConfig represents the configuration for an individual test.
type TestConfig struct {
// TimeoutsConfig represents the root of the timeouts configuration document.
type TimeoutsConfig struct {
	CFPush          Timeout `yaml:"cf_push"`
	CFStart         Timeout `yaml:"cf_start"`
	CFCreateService Timeout `yaml:"cf_create_service"`
}

// Timeout is a wrapper around time.Duration with a custom YAML unmarshal logic.
type Timeout time.Duration

// UnmarshalYAML unmarshals a string into a Timeout.
func (t *Timeout) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}

	d, err := time.ParseDuration(str)
	if err != nil {
		return err
	}

	*t = Timeout(d)

	return nil
}

// Duration returns the Timeout as time.Duration.
func (t *Timeout) Duration() time.Duration {
	return time.Duration(*t)
}
