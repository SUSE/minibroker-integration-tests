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

var _ = Describe("PostgreSQL", func() {
	BeforeEach(func() {
		if !mitsConfig.Tests.PostgreSQL.Enabled {
			Skip("Test is disabled")
		}
	})

	It("should deploy and connect", func() {
		mits.SimpleAppAndService(
			testSetup,
			mitsConfig.Tests.PostgreSQL,
			mitsConfig.Timeouts,
			serviceBrokerName,
			"assets/postgresqlapp",
			map[string]interface{}{
				"postgresqlDatabase": generator.PrefixedRandomName(mitsConfig.Tests.PostgreSQL.Class, "db"),
				"postgresqlUsername": generator.PrefixedRandomName(mitsConfig.Tests.PostgreSQL.Class, "user"),
			},
		)
	})
})
