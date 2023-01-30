/*
 * Copyright 2021, 2022 Hewlett Packard Enterprise Development LP
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

package dwdparse

import (
	"testing"
)

var dWDRules = []DWDirectiveRuleSpec{
	{
		Command: "jobdw",
		RuleDefs: []DWDirectiveRuleDef{
			{
				Key:             "type",
				Type:            "string",
				Pattern:         "^(raw|xfs|gfs2|lustre)$",
				IsRequired:      true,
				IsValueRequired: true,
			},
			{
				Key:             "capacity",
				Type:            "string",
				Pattern:         "^\\d+(KiB|KB|MiB|MB|GiB|GB|TiB|TB)$",
				IsRequired:      true,
				IsValueRequired: true,
			},
			{
				Key:             "name",
				Type:            "string",
				Pattern:         "^([A-Za-z0-9_-]+)$",
				IsRequired:      true,
				IsValueRequired: true,
				UniqueWithin:    "jobdw_name",
			},
			{
				Key:             "profile",
				Type:            "string",
				Pattern:         "^[A-Za-z][A-Za-z0-9_-]+$",
				IsRequired:      false,
				IsValueRequired: true,
			},
			{
				Key:             "combined_mgtmdt",
				Type:            "bool",
				IsRequired:      false,
				IsValueRequired: false,
			},
			{
				Key:             "external_mgs",
				Type:            "string",
				Pattern:         "^[A-Za-z0-9\\-_\\.@,:]+$",
				IsRequired:      false,
				IsValueRequired: true,
			},
		},
	},
	{
		Command: "create_persistent",
		RuleDefs: []DWDirectiveRuleDef{
			{
				Key:             "type",
				Type:            "string",
				Pattern:         "^(raw|xfs|gfs2|lustre)$",
				IsRequired:      true,
				IsValueRequired: true,
			},
			{
				Key:             "capacity",
				Type:            "string",
				Pattern:         "^\\d+(KiB|KB|MiB|MB|GiB|GB|TiB|TB)$",
				IsRequired:      true,
				IsValueRequired: true,
			},
			{
				Key:             "name",
				Type:            "string",
				Pattern:         "^([A-Za-z0-9_-]+)$",
				IsRequired:      true,
				IsValueRequired: true,
				UniqueWithin:    "create_persistent_name",
			},
			{
				Key:             "profile",
				Type:            "string",
				Pattern:         "^[A-Za-z][A-Za-z0-9_-]+$",
				IsRequired:      false,
				IsValueRequired: true,
			},
			{
				Key:             "combined_mgtmdt",
				Type:            "bool",
				IsRequired:      false,
				IsValueRequired: false,
			},
			{
				Key:             "external_mgs",
				Type:            "string",
				Pattern:         "^[A-Za-z0-9\\-_\\.@,:]+$",
				IsRequired:      false,
				IsValueRequired: true,
			},
		},
	},
	{
		Command: "stage_in",
		RuleDefs: []DWDirectiveRuleDef{
			{
				Key:             "type",
				Type:            "string",
				Pattern:         "^(directory|file|list)$",
				IsRequired:      true,
				IsValueRequired: true,
			},
			{
				Key:             "source",
				Type:            "string",
				IsRequired:      true,
				IsValueRequired: true,
			},
			{
				Key:             "destination",
				Type:            "string",
				IsRequired:      true,
				IsValueRequired: true,
			},
		},
	},
	{
		Command: "stage_out",
		RuleDefs: []DWDirectiveRuleDef{
			{
				Key:             "type",
				Type:            "string",
				Pattern:         "^(directory|file|list)$",
				IsRequired:      true,
				IsValueRequired: true,
			},
			{
				Key:             "source",
				Type:            "string",
				IsRequired:      true,
				IsValueRequired: true,
			},
			{
				Key:             "destination",
				Type:            "string",
				IsRequired:      true,
				IsValueRequired: true,
			},
		},
	},
	{
		Command: "persistentdw",
		RuleDefs: []DWDirectiveRuleDef{
			{
				Key:             "name",
				Type:            "string",
				Pattern:         "^([A-Za-z0-9_-]+)$",
				IsRequired:      true,
				IsValueRequired: true,
			},
		},
	},
	{
		Command: "destroy_persistent",
		RuleDefs: []DWDirectiveRuleDef{
			{
				Key:             "name",
				Type:            "string",
				Pattern:         "^([A-Za-z0-9_-]+)$",
				IsRequired:      true,
				IsValueRequired: true,
				UniqueWithin:    "destroy_persistent_name",
			},
		},
	},
	{
		Command: "container",
		RuleDefs: []DWDirectiveRuleDef{
			{
				Key:             "name",
				Type:            "string",
				Pattern:         "^([A-Za-z0-9_-]+)$",
				IsRequired:      true,
				IsValueRequired: true,
			},
			{
				Key:             "spec",
				Type:            "string",
				IsRequired:      true,
				IsValueRequired: true,
			},
			{
				Key:             "job_storage",
				Type:            "string",
				Pattern:         "^{([A-Za-z0-9_-]+)(,([A-Za-z0-9_-]+))*}$",
				IsRequired:      false,
				IsValueRequired: true,
			},
			{
				Key:             "persistent_storage",
				Type:            "string",
				Pattern:         "^{([A-Za-z0-9_-]+)(,([A-Za-z0-9_-]+))*}$",
				IsRequired:      false,
				IsValueRequired: true,
			},
			{
				Key:             "supervisor",
				Type:            "string",
				Pattern:         "^(rabbit|compute)$",
				IsRequired:      false,
				IsValueRequired: true,
			},
		},
	},
}

const (
	pass = true
	fail = false
)

var dwDirectiveTests = []struct {
	directiveList []string // #DW directive
	result        bool
}{
	{[]string{"#DW jobdw type=raw    capacity=100GB name=prettyGoodName  "}, pass},
	{[]string{"#DW jobdw type=xfs    capacity=100GB name=pretty_GoodName "}, pass},
	{[]string{"#DW jobdw type=gfs2   capacity=100GB name=pretty_GoodName "}, pass},
	{[]string{"#DW jobdw type=lustre capacity=100GB name=prettyGood-Name "}, pass},
	{[]string{"#DW jobdw type=raw    capacity=100TB name=__prettyGoodName"}, pass},
	{[]string{"#DW jobdw type=xfs    capacity=100TB name=-prettyGoodName-"}, pass},
	{[]string{"#DW jobdw type=lustre capacity=100TB name=0prettyGoodName1"}, pass},
	{[]string{"#DW jobdw type=raw    capacity=100GB name=prettierGoodName"}, pass},
	{[]string{"#DW jobdw type=xfs    capacity=100GB name=prettierGoodName"}, pass},
	{[]string{"#DW jobdw type=lustre capacity=100GB name=prettierGoodName"}, pass},
	{[]string{"#DW jobdw type=lustre capacity=100GB name=prettierGoodName"}, pass},

	{[]string{"#DW jobdw type=lustre capacity=100GB name=CoolProfile1 profile=this-TYPE_profile08"}, pass},
	{[]string{"#DW jobdw type=lustre capacity=100GB name=CoolProfile2 profile=this_TYPE_profile08"}, pass},
	{[]string{"#DW jobdw type=lustre capacity=100GB name=UncoolProfile1 profile=0this"}, fail},
	{[]string{"#DW jobdw type=lustre capacity=100GB name=UncoolProfile2 profile=this!"}, fail},
	{[]string{"#DW jobdw type=lustre capacity=100GB name=UncoolProfile3 profile=_this"}, fail},
	{[]string{"#DW jobdw type=lustre capacity=100GB name=UncoolProfile4 profile=-this"}, fail},

	{[]string{"#DW jobdw type=lustre capacity=100GB combined_mgtmdt name=prettierGoodName"}, pass},
	{[]string{"#DW jobdw type=lustre capacity=100GB combined_mgtmdt=true name=prettierGoodName"}, pass},

	{[]string{"#DW jobdw type=lustre external_mgs=rabbit-01@tcp capacity=100GB name=Extern1"}, pass},
	{[]string{"#DW jobdw type=lustre external_mgs=rabbit-01@tcp,rabbit-02@tcp:rabbit-03@tcp capacity=100GB name=Extern1"}, pass},
	{[]string{"#DW jobdw type=lustre external_mgs=rabbit-01@tcp combined_mgtmdt capacity=100GB name=Extern1"}, pass},
	{[]string{"#DW jobdw type=lustre external_mgs=rabbit-01@tcp0 capacity=100GB name=Extern1"}, pass},
	{[]string{"#DW jobdw type=lustre external_mgs=10.0.0.1@o2ib capacity=100GB name=Extern2"}, pass},

	{[]string{"#DW create_persistent type=raw    capacity=100GB name=prettyGoodName  "}, pass},
	{[]string{"#DW create_persistent type=xfs    capacity=100GB name=prettyGoodName  "}, pass},
	{[]string{"#DW create_persistent type=gfs2   capacity=100GB name=prettyGoodName  "}, pass},
	{[]string{"#DW create_persistent type=lustre capacity=100GB name=prettyGoodName  "}, pass},
	{[]string{"#DW create_persistent type=raw    capacity=100TB name=prettyGoodName  "}, pass},
	{[]string{"#DW create_persistent type=xfs    capacity=100TB name=prettyGoodName  "}, pass},
	{[]string{"#DW create_persistent type=lustre capacity=100TB name=prettyGoodName  "}, pass},
	{[]string{"#DW create_persistent type=raw    capacity=100GB name=prettierGoodName"}, pass},
	{[]string{"#DW create_persistent type=xfs    capacity=100GB name=prettierGoodName"}, pass},
	{[]string{"#DW create_persistent type=lustre capacity=100GB name=prettierGoodName"}, pass},

	{[]string{"#DW create_persistent type=lustre capacity=100GB name=CoolProfile1 profile=this-TYPE_profile08"}, pass},
	{[]string{"#DW create_persistent type=lustre capacity=100GB name=CoolProfile2 profile=this_TYPE_profile08"}, pass},
	{[]string{"#DW create_persistent type=lustre capacity=100GB name=UncoolProfile1 profile=0this"}, fail},
	{[]string{"#DW create_persistent type=lustre capacity=100GB name=UncoolProfile2 profile=this!"}, fail},
	{[]string{"#DW create_persistent type=lustre capacity=100GB name=UncoolProfile3 profile=_this"}, fail},
	{[]string{"#DW create_persistent type=lustre capacity=100GB name=UncoolProfile4 profile=-this"}, fail},

	{[]string{"#DW create_persistent type=lustre capacity=100GB combined_mgtmdt name=prettierGoodName"}, pass},
	{[]string{"#DW create_persistent type=lustre capacity=100GB combined_mgtmdt=true name=prettierGoodName"}, pass},

	{[]string{"#DW create_persistent type=lustre external_mgs=rabbit-01@tcp capacity=100GB name=Extern1"}, pass},
	{[]string{"#DW create_persistent type=lustre external_mgs=rabbit-01@tcp combined_mgtmdt capacity=100GB name=Extern1"}, pass},
	{[]string{"#DW create_persistent type=lustre external_mgs=rabbit-01@tcp0 capacity=100GB name=Extern1"}, pass},
	{[]string{"#DW create_persistent type=lustre external_mgs=10.0.0.1@o2ib capacity=100GB name=Extern2"}, pass},

	{[]string{"#DW stage_in  type=file      destination=$DW_JOB_STRIPED source=/pfs/dld-input "}, pass},
	{[]string{"#DW stage_in  type=directory destination=$DW_JOB_STRIPED source=/pfs/dld-input "}, pass},
	{[]string{"#DW stage_in  type=list      destination=$DW_JOB_STRIPED source=/pfs/dld-input "}, pass},
	{[]string{"#DW stage_out type=file      destination=/pfs/dld-output source=$DW_JOB_STRIPED"}, pass},
	{[]string{"#DW stage_out type=directory destination=/pfs/dld-output source=$DW_JOB_STRIPED"}, pass},
	{[]string{"#DW stage_out type=list      destination=/pfs/dld-output source=$DW_JOB_STRIPED"}, pass},

	{[]string{"#DW stage_in  type=file      destination=$DW_JOB_STRIPED source=/pfs/dld-input "}, pass},
	{[]string{"#DW stage_in  type=directory destination=$DW_JOB_STRIPED source=/pfs/dld-input "}, pass},
	{[]string{"#DW stage_in  type=list      destination=$DW_JOB_STRIPED source=/pfs/dld-input "}, pass},
	{[]string{"#DW stage_out type=file      destination=/pfs/dld-output source=$DW_JOB_STRIPED"}, pass},
	{[]string{"#DW stage_out type=directory destination=/pfs/dld-output source=$DW_JOB_STRIPED"}, pass},
	{[]string{"#DW stage_out type=list      destination=/pfs/dld-output source=$DW_JOB_STRIPED"}, pass},

	{[]string{"#DW persistentdw name=evenBetterName"}, pass},
	{[]string{"#DW persistentdw name=evenBetterName"}, pass},

	{[]string{"#DW destroy_persistent name=evenBetterName"}, pass},
	{[]string{"#DW destroy_persistent name=evenBetterName"}, pass},

	{[]string{"#DW destroy_persistent name=uniqueName1 ",
		"#DW destroy_persistent name=uniqueName2 "}, pass},
	{[]string{"#DW destroy_persistent name=conflictName",
		"#DW destroy_persistent name=conflictName"}, fail},

	{[]string{"                                                         "}, fail},
	{[]string{"    jobdw type=raw     capacity=100GB name=prettyGoodName"}, fail},
	{[]string{"#DW       type=xfs     capacity=100TB name=noCommand     "}, fail},
	{[]string{"#DW bogus the_rest_does_not_matter                       "}, fail},
	{[]string{"#DW jobd  type=lustre  capacity=100TB name=badCommand    "}, fail},
	{[]string{"#DW jobdw tye=badtype  capacity=100TB name=badType       "}, fail},
	{[]string{"#DW jobdw type=badtype capacity=100TB name=badType       "}, fail},
	{[]string{"#DW jobdw              capacity=100TB name=missingType   "}, fail},
	{[]string{"#DW jobdw type=file    capacity=100TB name=badType       "}, fail},
	{[]string{"#DW jobdw type=raw     caacity=100TB  name=badCapacity   "}, fail},
	{[]string{"#DW jobdw type=raw     capacity=bad   name=badCapacity   "}, fail},
	{[]string{"#DW jobdw type=xfs     capacity=100TB ame=badName        "}, fail},
	{[]string{"#DW jobdw type=xfs     capacity=100TB name=!!21//\\      "}, fail},
	{[]string{"#DW jobdw                                                "}, fail},

	{[]string{"#DW jobdw type=raw type=raw capacity=100TB name=duplicatedTypes                         "}, fail},
	{[]string{"#DW jobdw type=raw capacity=100TB name=conflictingTypes type=xfs                        "}, fail},
	{[]string{"#DW jobdw type=badtype destination=shouldNotHaveDestination source=shouldNotHaveSource  "}, fail},

	{[]string{"#DW stage_in  type=badtype destination=$DW_JOB_STRIPED source=/pfs/dld-input               "}, fail},
	{[]string{"#DW stage_out type=file    badkey=foobar destination=/pfs/dld-output source=$DW_JOB_STRIPED"}, fail},

	{[]string{"#DW boguscommand the_rest_should_not_matter                                             "}, fail},
	{[]string{"#DW jobd  type=lustre  capacity=100TB name=badCommand                                   "}, fail},
	{[]string{"#DW stage_in type=badtype destination=$DW_JOB_STRIPED source=/pfs/dld-input             "}, fail},
	{[]string{"#DW stage_out type=file badkey=foobar destination=/pfs/dld-output source=$DW_JOB_STRIPED"}, fail},
	{[]string{"#DW jobd  type=lustre  capacity=100TB name=badCommand                                   "}, fail},
	{[]string{"#DW jobdw tye=badtype  capacity=100TB name=badType                                      "}, fail},
	{[]string{"#DW jobdw type=badtype capacity=100TB name=badType                                      "}, fail},
	{[]string{"#DW jobdw              capacity=100TB name=missingType                                  "}, fail},
	{[]string{"#DW jobdw type=file    capacity=100TB name=badType                                      "}, fail},
	{[]string{"#DW jobdw type=raw     caacity=100TB  name=badCapacity                                  "}, fail},
	{[]string{"#DW jobdw type=raw     capacity=bad   name=badCapacity                                  "}, fail},
	{[]string{"#DW jobdw type=xfs     capacity=100TB ame=badName                                       "}, fail},
	{[]string{"#DW jobdw type=xfs     capacity=100TB name=!!21//\\                                     "}, fail},
	{[]string{"#DW jobdw                                                                               "}, fail},

	// NOTE: Please do not add new test cases here. Instead use the new test framework below that
	//       contains individualized test cases.
}

func _TestDWParse(t *testing.T) {
	for index, tt := range dwDirectiveTests {

		err := Validate(dWDRules, tt.directiveList, func(int, DWDirectiveRuleSpec) {})

		if (tt.result == pass && err != nil) || (tt.result == fail && err == nil) {
			t.Errorf("TestDWParse(%s)(%d): expect_valid(%v) err(%v)", tt.directiveList, index, tt.result, err)
		}
	}
}

// testCase defines an individual test case
type testCase struct {
	rules      []DWDirectiveRuleSpec
	directives []string
	result     bool
}

// test provides a common method to test a series of test cases against a set of rules
func test(t *testing.T, rules []DWDirectiveRuleSpec, tests []testCase) {
	for index, tc := range tests {
		err := Validate(rules, tc.directives, func(int, DWDirectiveRuleSpec) {})

		if (tc.result == pass && err != nil) || (tc.result == fail && err == nil) {
			t.Errorf("test(%s)(%d): expected(%v) err(%v)", tc.directives, index, tc.result, err)
		}
	}
}

func _TestUniqueWithin(t *testing.T) {
	rules := []DWDirectiveRuleSpec{{
		Command: "unique",
		RuleDefs: []DWDirectiveRuleDef{{
			Key:          "name",
			Type:         "string",
			IsRequired:   true,
			UniqueWithin: "jobdw_name",
		}},
	}}

	tests := []testCase{{
		directives: []string{
			"#DW unique name=uniqueName1 ",
			"#DW unique name=uniqueName2 "},
		result: pass,
	}, {
		directives: []string{
			"#DW unique name=sameName",
			"#DW unique name=sameName"},
		result: fail,
	}}

	test(t, rules, tests)
}

func _TestKeyRegexp(t *testing.T) {
	rules := []DWDirectiveRuleSpec{{
		Command: "regexp",
		RuleDefs: []DWDirectiveRuleDef{{
			Key:  `^(PREFIX_1_|PREFIX_2_)\w+`, // test a simple prefix regexp
			Type: "string",
		}},
	}}

	tests := []testCase{
		{directives: []string{"#DW regexp PREFIX_1_good=value"}, result: pass},
		{directives: []string{"#DW regexp PREFIX_2_good=value"}, result: pass},
		{directives: []string{"#DW regexp PREFIX_2_=value"}, result: fail}, // missing word after prefix strings
		{directives: []string{"#DW regexp INVALID=value"}, result: fail},
	}

	test(t, rules, tests)
}

func TestKeyMalformedRegexp(t *testing.T) {
	rules := []DWDirectiveRuleSpec{{
		Command: "regexp",
		RuleDefs: []DWDirectiveRuleDef{{
			Key:  `*`, // missing argument to repetition operator '*'
			Type: "string",
		}},
	}}

	tests := []testCase{
		{directives: []string{"#DW regexp KEY=VALUE"}, result: fail},
	}

	test(t, rules, tests)
}
