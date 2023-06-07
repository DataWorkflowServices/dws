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

package v1alpha1

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"

	dwsv1alpha2 "github.com/HewlettPackard/dws/api/v1alpha2"
	utilconversion "github.com/HewlettPackard/dws/github/cluster-api/util/conversion"
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
		Hub:   &dwsv1alpha2.SystemConfiguration{},
		Spoke: &SystemConfiguration{},
	}))

	t.Run("for Workflow", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Hub:   &dwsv1alpha2.Workflow{},
		Spoke: &Workflow{},
	}))

}

// Just touch ginkgo, so it's here to interpret any ginkgo args from
// "make test", so that doesn't fail on this test file.
var _ = BeforeSuite(func() {})
