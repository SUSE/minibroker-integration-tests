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
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"

	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	helpersConfig "github.com/cloudfoundry-incubator/cf-test-helpers/config"
	"github.com/cloudfoundry-incubator/cf-test-helpers/generator"
	"github.com/cloudfoundry-incubator/cf-test-helpers/workflowhelpers"

	"github.com/SUSE/minibroker-integration-tests/mits"
)

var (
	config mits.Config

	testSetup         *workflowhelpers.ReproducibleTestSuiteSetup
	serviceBrokerName string
)

func TestMits(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mits Suite")
}

var _ = BeforeSuite(func() {
	err := mits.LoadConfig("CONFIG_PATH", &config)
	Expect(err).NotTo(HaveOccurred())

	serviceBrokerName = generator.PrefixedRandomName("mits", "minibroker")

	cfg := helpersConfig.Config{
		TimeoutScale:      2.0,
		NamePrefix:        "mits",
		ApiEndpoint:       config.CF.API.Endpoint,
		AdminUser:         config.CF.Admin.Username,
		AdminPassword:     config.CF.Admin.Password,
		SkipSSLValidation: true,
	}
	testSetup = workflowhelpers.NewTestSuiteSetup(&cfg)
	testSetup.Setup()

	workflowhelpers.AsUser(testSetup.AdminUserContext(), testSetup.ShortTimeout(), func() {
		Expect(
			cf.Cf("create-service-broker", serviceBrokerName, "user", "pass", config.Minibroker.API.Endpoint).
				Wait(testSetup.ShortTimeout()),
		).To(Exit(0))
		if config.Tests.MariaDB.Enabled {
			Expect(
				cf.Cf("enable-service-access", config.Tests.MariaDB.Class, "-b", serviceBrokerName).
					Wait(testSetup.ShortTimeout()),
			).To(Exit(0))
		}
		if config.Tests.MySQL.Enabled {
			Expect(
				cf.Cf("enable-service-access", config.Tests.MySQL.Class, "-b", serviceBrokerName).
					Wait(testSetup.ShortTimeout()),
			).To(Exit(0))
		}
		if config.Tests.PostgreSQL.Enabled {
			Expect(
				cf.Cf("enable-service-access", config.Tests.PostgreSQL.Class, "-b", serviceBrokerName).
					Wait(testSetup.ShortTimeout()),
			).To(Exit(0))
		}
		if config.Tests.Redis.Enabled {
			Expect(
				cf.Cf("enable-service-access", config.Tests.Redis.Class, "-b", serviceBrokerName).
					Wait(testSetup.ShortTimeout()),
			).To(Exit(0))
		}
	})
})

var _ = AfterSuite(func() {
	workflowhelpers.AsUser(testSetup.AdminUserContext(), testSetup.ShortTimeout(), func() {
		Expect(
			cf.Cf("delete-service-broker", serviceBrokerName, "-f").
				Wait(testSetup.ShortTimeout()),
		).To(Exit(0))
	})

	testSetup.Teardown()
})
