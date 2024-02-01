/*
 * Copyright 2023-2024 Hewlett Packard Enterprise Development LP
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

package v1alpha1

import (
	"testing"

	fuzz "github.com/google/gofuzz"
	. "github.com/onsi/ginkgo/v2"
	"k8s.io/apimachinery/pkg/api/apitesting/fuzzer"
	runtimeserializer "k8s.io/apimachinery/pkg/runtime/serializer"

	dwsv1alpha2 "github.com/DataWorkflowServices/dws/api/v1alpha2"
	utilconversion "github.com/DataWorkflowServices/dws/github/cluster-api/util/conversion"
)

func TestFuzzyConversion(t *testing.T) {

	t.Run("for ClientMount", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Hub:   &dwsv1alpha2.ClientMount{},
		Spoke: &ClientMount{},
	}))

	t.Run("for Computes", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Hub:   &dwsv1alpha2.Computes{},
		Spoke: &Computes{},
	}))

	t.Run("for DWDirectiveRule", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Hub:   &dwsv1alpha2.DWDirectiveRule{},
		Spoke: &DWDirectiveRule{},
	}))

	t.Run("for DirectiveBreakdown", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Hub:   &dwsv1alpha2.DirectiveBreakdown{},
		Spoke: &DirectiveBreakdown{},
	}))

	t.Run("for PersistentStorageInstance", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Hub:   &dwsv1alpha2.PersistentStorageInstance{},
		Spoke: &PersistentStorageInstance{},
	}))

	t.Run("for Servers", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Hub:   &dwsv1alpha2.Servers{},
		Spoke: &Servers{},
	}))

	t.Run("for Storage", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Hub:   &dwsv1alpha2.Storage{},
		Spoke: &Storage{},
	}))

	t.Run("for SystemConfiguration", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Hub:         &dwsv1alpha2.SystemConfiguration{},
		Spoke:       &SystemConfiguration{},
		FuzzerFuncs: []fuzzer.FuzzerFuncs{SystemConfigurationFuzzFunc},
	}))

	t.Run("for Workflow", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Hub:   &dwsv1alpha2.Workflow{},
		Spoke: &Workflow{},
	}))

}

func SystemConfigurationFuzzFunc(_ runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		SystemConfigurationComputesFuzzer,
	}
}

// Use the same compute names in both spec.ComputeNodes and spec.StorageNodes.
func SystemConfigurationComputesFuzzer(in *SystemConfigurationSpec, c fuzz.Continue) {
	// Tell the fuzzer to begin by fuzzing everything in the object.
	c.FuzzNoCustom(in)

	// Now pull any fuzzed compute names out of the spec.StorageNodes list and
	// use them to build a new spec.ComputeNodes list.
	if in.StorageNodes != nil {
		computes := make([]SystemConfigurationComputeNode, 0)
		for sidx := range in.StorageNodes {
			for cidx := range in.StorageNodes[sidx].ComputesAccess {
				name := c.RandString()
				in.StorageNodes[sidx].ComputesAccess[cidx].Name = name
				computes = append(computes, SystemConfigurationComputeNode{Name: name})
			}
		}
		in.ComputeNodes = computes
	} else {
		in.ComputeNodes = nil
	}

}

// Just touch ginkgo, so it's here to interpret any ginkgo args from
// "make test", so that doesn't fail on this test file.
var _ = BeforeSuite(func() {})
