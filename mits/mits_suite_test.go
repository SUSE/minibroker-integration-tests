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
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"

	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	"github.com/cloudfoundry-incubator/cf-test-helpers/config"
	"github.com/cloudfoundry-incubator/cf-test-helpers/generator"
	"github.com/cloudfoundry-incubator/cf-test-helpers/workflowhelpers"

	"github.com/SUSE/minibroker-integration-tests/mits"
)

var (
	tests    mits.TestsConfig
	timeouts mits.TimeoutsConfig

	testSetup         *workflowhelpers.ReproducibleTestSuiteSetup
	serviceBrokerName string
	serviceBrokerURL  string
)

func TestMits(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mits Suite")
}

var _ = BeforeSuite(func() {
	err := mits.LoadConfig("CONFIG_TESTS", &tests)
	Expect(err).NotTo(HaveOccurred())
	err = mits.LoadConfig("CONFIG_TIMEOUTS", &timeouts)
	Expect(err).NotTo(HaveOccurred())

	apiEndpoint := os.Getenv("API_ENDPOINT")
	adminUser := os.Getenv("CF_ADMIN_USERNAME")
	adminPassword := os.Getenv("CF_ADMIN_PASSWORD")
	serviceBrokerName = generator.PrefixedRandomName("mits", "minibroker")
	serviceBrokerURL = os.Getenv("SERVICE_BROKER_URL")

	cfg := config.Config{
		TimeoutScale:      2.0,
		NamePrefix:        "mits",
		ApiEndpoint:       apiEndpoint,
		AdminUser:         adminUser,
		AdminPassword:     adminPassword,
		SkipSSLValidation: true,
	}
	testSetup = workflowhelpers.NewTestSuiteSetup(&cfg)
	testSetup.Setup()

	workflowhelpers.AsUser(testSetup.AdminUserContext(), testSetup.ShortTimeout(), func() {
		Expect(cf.Cf("create-service-broker", serviceBrokerName, "user", "pass", serviceBrokerURL).Wait(testSetup.ShortTimeout())).To(Exit(0))
		if tests.MariaDB.Enabled {
			Expect(cf.Cf("enable-service-access", tests.MariaDB.Class, "-b", serviceBrokerName).Wait(testSetup.ShortTimeout())).
				To(Exit(0))
		}
		if tests.MySQL.Enabled {
			Expect(cf.Cf("enable-service-access", tests.MySQL.Class, "-b", serviceBrokerName).Wait(testSetup.ShortTimeout())).
				To(Exit(0))
		}
		if tests.PostgreSQL.Enabled {
			Expect(cf.Cf("enable-service-access", tests.PostgreSQL.Class, "-b", serviceBrokerName).Wait(testSetup.ShortTimeout())).
				To(Exit(0))
		}
		if tests.Redis.Enabled {
			Expect(cf.Cf("enable-service-access", tests.Redis.Class, "-b", serviceBrokerName).Wait(testSetup.ShortTimeout())).
				To(Exit(0))
		}
	})
})

var _ = AfterSuite(func() {
	workflowhelpers.AsUser(testSetup.AdminUserContext(), testSetup.ShortTimeout(), func() {
		Expect(cf.Cf("delete-service-broker", serviceBrokerName, "-f").Wait(testSetup.ShortTimeout())).To(Exit(0))
	})

	testSetup.Teardown()
})
