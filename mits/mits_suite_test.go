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

package mits_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"gopkg.in/yaml.v2"

	"github.com/SUSE/minibroker-integration-tests/mits"
)

var testsConfig map[string]mits.TestConfig

func init() {
	if err := loadConfig("CONFIG_TESTS", &testsConfig); err != nil {
		log.Fatal(err)
	}
}

func loadConfig(envConfig string, config interface{}) error {
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
	if err := configDecoder.Decode(config); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	return nil
}

func TestMits(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mits Suite")
}
