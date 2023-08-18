/*
 * Copyright 2023 Hewlett Packard Enterprise Development LP
 * Other additional copyright holders may be indicated within.
 *
 * The entirety of this work is licensed under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 *
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v1alpha2

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Resource error", func() {
	var status1 string
	var status2 string
	var err error

	DescribeTable("translates severity values to status error values",
		func(severityStr string, severity ResourceErrorSeverity, expectedStatus string) {
			status1, err = severity.ToStatus()
			Expect(err).To(BeNil())
			Expect(status1).To(Equal(expectedStatus))
			status2, err = SeverityStringToStatus(severityStr)
			Expect(err).To(BeNil())
			Expect(status2).To(Equal(expectedStatus))
		},
		Entry("empty severity", "", SeverityMinor, StatusRunning),
		Entry("minor severity", string(SeverityMinor), SeverityMinor, StatusRunning),
		Entry("major severity", string(SeverityMajor), SeverityMajor, StatusTransientCondition),
		Entry("fatal severity", string(SeverityFatal), SeverityFatal, StatusError),
	)

	It("detects unknown severity", func() {
		status2, err = SeverityStringToStatus("squishy")
		Expect(err).ToNot(BeNil())
	})
})
