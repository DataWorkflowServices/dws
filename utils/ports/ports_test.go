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

package ports

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Ports Utilities Test")
}

var _ = Describe("Ports Utilities Test", func() {
	DescribeTable("Validate Ports",
		func(port string, isValid bool) {
			ports := []intstr.IntOrString{intstr.Parse(port)}
			Expect(Validate(ports) == nil).To(Equal(isValid))
		},
		// Integer values
		Entry("Valid smallest integer", "1", true),
		Entry("Valid largest integer", "65535", true),
		Entry("Zero", "0", false),
		Entry("Negative", "-1", false),
		Entry("Funky whitespace", "  1  ", false),
		Entry("Too small", "0", false),
		Entry("Too large", "65536", false),

		// String values
		Entry("Valid port range", "1-65535", true),
		Entry("Invalid port range", "1-", false),
		Entry("Funky whitespace", " 1 -  2", false),
		Entry("Start too small", "0-1", false),
		Entry("Start too large", "65536-65537", false),
		Entry("End too small", "1-0", false),
		Entry("End too large", "1-65536", false),
		Entry("Start equals end", "1-1", false),
		Entry("Start greater than end", "2-1", false),
	)

	It("Port Iterator (Valid)", func() {
		ports := []intstr.IntOrString{
			intstr.FromInt(1),
			intstr.FromInt(2),
			intstr.FromString("3-10"),
			intstr.FromString("11-19"),
			intstr.FromInt(20),
		}

		Expect(Validate(ports)).To(BeNil())

		itr := NewPortIterator(ports)
		expected := uint16(1)
		for port := itr.Next(); port != 0; port = itr.Next() {
			Expect(port).To(Equal(expected))
			expected++
		}

		Expect(itr.Next()).To(Equal(uint16(0)))
	})

	It("Port Iterator (Invalid)", func() {
		ports := []intstr.IntOrString{intstr.FromString("invalid range")}
		itr := NewPortIterator(ports)
		Expect(itr.Next()).To(Equal(InvalidPort), "invalid range")

		ports = []intstr.IntOrString{intstr.FromString("65536-65537")}
		itr = NewPortIterator(ports)
		Expect(itr.Next()).To(Equal(InvalidPort), "start port overflows")

		ports = []intstr.IntOrString{intstr.FromString("1-65536")}
		itr = NewPortIterator(ports)
		Expect(itr.Next()).To(Equal(InvalidPort), "end port overflows")
	})
})
