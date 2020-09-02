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
	. "github.com/onsi/ginkgo"

	"github.com/cloudfoundry-incubator/cf-test-helpers/generator"

	"github.com/SUSE/minibroker-integration-tests/mits"
)

var _ = Describe("MariaDB", func() {
	BeforeEach(func() {
		if !mitsConfig.Tests.MariaDB.Enabled {
			Skip("All MariaDB tests are disabled")
		}
	})

	Context("Without overrideParams set", func() {
		BeforeEach(func() {
			if mitsConfig.Minibroker.Provisioning.OverrideParams.Enabled {
				Skip("overrideParams are set")
			}
		})

		It("should deploy and connect WITH extra provisioning parameters", func() {
			mits.SimpleAppAndService(
				testSetup,
				mitsConfig.Tests.MariaDB,
				mitsConfig.Timeouts,
				serviceBrokerName,
				"assets/mysqlapp",
				map[string]interface{}{
					"db": map[string]interface{}{
						"name": generator.PrefixedRandomName(mitsConfig.Tests.MariaDB.Class, "db"),
						"user": generator.PrefixedRandomName(mitsConfig.Tests.MariaDB.Class, "user"),
					},
					"replication": map[string]interface{}{
						"enabled": false,
					},
				},
			)
		})
	})

	Context("With overrideParams set", func() {
		BeforeEach(func() {
			if !mitsConfig.Minibroker.Provisioning.OverrideParams.Enabled {
				Skip("overrideParams are not set")
			}
		})

		It("should deploy and connect WITHOUT extra provisioning parameters", func() {
			mits.SimpleAppAndService(
				testSetup,
				mitsConfig.Tests.MariaDB,
				mitsConfig.Timeouts,
				serviceBrokerName,
				"assets/mysqlapp",
				nil,
			)
		})
	})
})
