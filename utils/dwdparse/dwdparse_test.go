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
	"fmt"
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
	deny                         bool = true
	validDWOrAllowUnknownCommand bool = true
	allow                        bool = false
	invalidDW                    bool = false
)

var dwDirectiveTests = []struct {
	directiveList      []string // #DW directive
	failUnknownCommand bool     // deny/allow unknown commands
	validCommand       bool     // expected parse error result compared with nil
}{
	{[]string{"#DW jobdw type=raw    capacity=100GB name=prettyGoodName  "}, deny, validDWOrAllowUnknownCommand},
	{[]string{"#DW jobdw type=xfs    capacity=100GB name=pretty_GoodName "}, deny, validDWOrAllowUnknownCommand},
	{[]string{"#DW jobdw type=gfs2   capacity=100GB name=pretty_GoodName "}, deny, validDWOrAllowUnknownCommand},
	{[]string{"#DW jobdw type=lustre capacity=100GB name=prettyGood-Name "}, deny, validDWOrAllowUnknownCommand},
	{[]string{"#DW jobdw type=raw    capacity=100TB name=__prettyGoodName"}, deny, validDWOrAllowUnknownCommand},
	{[]string{"#DW jobdw type=xfs    capacity=100TB name=-prettyGoodName-"}, deny, validDWOrAllowUnknownCommand},
	{[]string{"#DW jobdw type=lustre capacity=100TB name=0prettyGoodName1"}, deny, validDWOrAllowUnknownCommand},
	{[]string{"#DW jobdw type=raw    capacity=100GB name=prettierGoodName"}, deny, validDWOrAllowUnknownCommand},
	{[]string{"#DW jobdw type=xfs    capacity=100GB name=prettierGoodName"}, deny, validDWOrAllowUnknownCommand},
	{[]string{"#DW jobdw type=lustre capacity=100GB name=prettierGoodName"}, deny, validDWOrAllowUnknownCommand},
	{[]string{"#DW jobdw type=lustre capacity=100GB name=prettierGoodName"}, deny, validDWOrAllowUnknownCommand},

	{[]string{"#DW jobdw type=raw    capacity=100GB name=prettyGoodName  "}, allow, validDWOrAllowUnknownCommand},
	{[]string{"#DW jobdw type=xfs    capacity=100GB name=prettyGoodName  "}, allow, validDWOrAllowUnknownCommand},
	{[]string{"#DW jobdw type=lustre capacity=100GB name=prettyGoodName  "}, allow, validDWOrAllowUnknownCommand},
	{[]string{"#DW jobdw type=raw    capacity=100TB name=prettyGoodName  "}, allow, validDWOrAllowUnknownCommand},
	{[]string{"#DW jobdw type=xfs    capacity=100TB name=prettyGoodName  "}, allow, validDWOrAllowUnknownCommand},
	{[]string{"#DW jobdw type=lustre capacity=100TB name=prettyGoodName  "}, allow, validDWOrAllowUnknownCommand},
	{[]string{"#DW jobdw type=raw    capacity=100GB name=prettierGoodName"}, allow, validDWOrAllowUnknownCommand},
	{[]string{"#DW jobdw type=xfs    capacity=100GB name=prettierGoodName"}, allow, validDWOrAllowUnknownCommand},
	{[]string{"#DW jobdw type=lustre capacity=100GB name=prettierGoodName"}, allow, validDWOrAllowUnknownCommand},

	{[]string{"#DW jobdw type=raw    capacity=100GB name=uniqueName1 ",
		"#DW jobdw type=raw    capacity=100GB name=uniqueName2 "}, allow, validDWOrAllowUnknownCommand},
	{[]string{"#DW jobdw type=raw    capacity=100GB name=conflictName",
		"#DW jobdw type=raw    capacity=100GB name=conflictName"}, deny, invalidDW},

	{[]string{"#DW jobdw type=lustre capacity=100GB name=CoolProfile1 profile=this-TYPE_profile08"}, allow, validDWOrAllowUnknownCommand},
	{[]string{"#DW jobdw type=lustre capacity=100GB name=CoolProfile2 profile=this_TYPE_profile08"}, allow, validDWOrAllowUnknownCommand},
	{[]string{"#DW jobdw type=lustre capacity=100GB name=UncoolProfile1 profile=0this"}, deny, invalidDW},
	{[]string{"#DW jobdw type=lustre capacity=100GB name=UncoolProfile2 profile=this!"}, deny, invalidDW},
	{[]string{"#DW jobdw type=lustre capacity=100GB name=UncoolProfile3 profile=_this"}, deny, invalidDW},
	{[]string{"#DW jobdw type=lustre capacity=100GB name=UncoolProfile4 profile=-this"}, deny, invalidDW},

	{[]string{"#DW jobdw type=lustre capacity=100GB combined_mgtmdt name=prettierGoodName"}, allow, validDWOrAllowUnknownCommand},
	{[]string{"#DW jobdw type=lustre capacity=100GB combined_mgtmdt=true name=prettierGoodName"}, allow, validDWOrAllowUnknownCommand},

	{[]string{"#DW jobdw type=lustre external_mgs=rabbit-01@tcp capacity=100GB name=Extern1"}, allow, validDWOrAllowUnknownCommand},
	{[]string{"#DW jobdw type=lustre external_mgs=rabbit-01@tcp,rabbit-02@tcp:rabbit-03@tcp capacity=100GB name=Extern1"}, allow, validDWOrAllowUnknownCommand},
	{[]string{"#DW jobdw type=lustre external_mgs=rabbit-01@tcp combined_mgtmdt capacity=100GB name=Extern1"}, allow, validDWOrAllowUnknownCommand},
	{[]string{"#DW jobdw type=lustre external_mgs=rabbit-01@tcp0 capacity=100GB name=Extern1"}, allow, validDWOrAllowUnknownCommand},
	{[]string{"#DW jobdw type=lustre external_mgs=10.0.0.1@o2ib capacity=100GB name=Extern2"}, allow, validDWOrAllowUnknownCommand},

	{[]string{"#DW create_persistent type=raw    capacity=100GB name=prettyGoodName  "}, deny, validDWOrAllowUnknownCommand},
	{[]string{"#DW create_persistent type=xfs    capacity=100GB name=prettyGoodName  "}, deny, validDWOrAllowUnknownCommand},
	{[]string{"#DW create_persistent type=gfs2   capacity=100GB name=prettyGoodName  "}, deny, validDWOrAllowUnknownCommand},
	{[]string{"#DW create_persistent type=lustre capacity=100GB name=prettyGoodName  "}, deny, validDWOrAllowUnknownCommand},
	{[]string{"#DW create_persistent type=raw    capacity=100TB name=prettyGoodName  "}, deny, validDWOrAllowUnknownCommand},
	{[]string{"#DW create_persistent type=xfs    capacity=100TB name=prettyGoodName  "}, deny, validDWOrAllowUnknownCommand},
	{[]string{"#DW create_persistent type=lustre capacity=100TB name=prettyGoodName  "}, deny, validDWOrAllowUnknownCommand},
	{[]string{"#DW create_persistent type=raw    capacity=100GB name=prettierGoodName"}, deny, validDWOrAllowUnknownCommand},
	{[]string{"#DW create_persistent type=xfs    capacity=100GB name=prettierGoodName"}, deny, validDWOrAllowUnknownCommand},
	{[]string{"#DW create_persistent type=lustre capacity=100GB name=prettierGoodName"}, deny, validDWOrAllowUnknownCommand},

	{[]string{"#DW create_persistent type=raw    capacity=100GB name=prettyGoodName  "}, allow, validDWOrAllowUnknownCommand},
	{[]string{"#DW create_persistent type=xfs    capacity=100GB name=prettyGoodName  "}, allow, validDWOrAllowUnknownCommand},
	{[]string{"#DW create_persistent type=gfs2   capacity=100GB name=prettyGoodName  "}, allow, validDWOrAllowUnknownCommand},
	{[]string{"#DW create_persistent type=lustre capacity=100GB name=prettyGoodName  "}, allow, validDWOrAllowUnknownCommand},
	{[]string{"#DW create_persistent type=raw    capacity=100TB name=prettyGoodName  "}, allow, validDWOrAllowUnknownCommand},
	{[]string{"#DW create_persistent type=xfs    capacity=100TB name=prettyGoodName  "}, allow, validDWOrAllowUnknownCommand},
	{[]string{"#DW create_persistent type=lustre capacity=100TB name=prettyGoodName  "}, allow, validDWOrAllowUnknownCommand},
	{[]string{"#DW create_persistent type=raw    capacity=100GB name=prettierGoodName"}, allow, validDWOrAllowUnknownCommand},
	{[]string{"#DW create_persistent type=xfs    capacity=100GB name=prettierGoodName"}, allow, validDWOrAllowUnknownCommand},
	{[]string{"#DW create_persistent type=lustre capacity=100GB name=prettierGoodName"}, allow, validDWOrAllowUnknownCommand},

	{[]string{"#DW create_persistent type=raw    capacity=100GB name=uniqueName1 ",
		"#DW create_persistent type=raw    capacity=100GB name=uniqueName2 "}, allow, validDWOrAllowUnknownCommand},
	{[]string{"#DW create_persistent type=raw    capacity=100GB name=conflictName",
		"#DW create_persistent type=raw    capacity=100GB name=conflictName"}, deny, invalidDW},

	{[]string{"#DW create_persistent type=lustre capacity=100GB name=CoolProfile1 profile=this-TYPE_profile08"}, allow, validDWOrAllowUnknownCommand},
	{[]string{"#DW create_persistent type=lustre capacity=100GB name=CoolProfile2 profile=this_TYPE_profile08"}, allow, validDWOrAllowUnknownCommand},
	{[]string{"#DW create_persistent type=lustre capacity=100GB name=UncoolProfile1 profile=0this"}, deny, invalidDW},
	{[]string{"#DW create_persistent type=lustre capacity=100GB name=UncoolProfile2 profile=this!"}, deny, invalidDW},
	{[]string{"#DW create_persistent type=lustre capacity=100GB name=UncoolProfile3 profile=_this"}, deny, invalidDW},
	{[]string{"#DW create_persistent type=lustre capacity=100GB name=UncoolProfile4 profile=-this"}, deny, invalidDW},

	{[]string{"#DW create_persistent type=lustre capacity=100GB combined_mgtmdt name=prettierGoodName"}, allow, validDWOrAllowUnknownCommand},
	{[]string{"#DW create_persistent type=lustre capacity=100GB combined_mgtmdt=true name=prettierGoodName"}, allow, validDWOrAllowUnknownCommand},

	{[]string{"#DW create_persistent type=lustre external_mgs=rabbit-01@tcp capacity=100GB name=Extern1"}, allow, validDWOrAllowUnknownCommand},
	{[]string{"#DW create_persistent type=lustre external_mgs=rabbit-01@tcp combined_mgtmdt capacity=100GB name=Extern1"}, allow, validDWOrAllowUnknownCommand},
	{[]string{"#DW create_persistent type=lustre external_mgs=rabbit-01@tcp0 capacity=100GB name=Extern1"}, allow, validDWOrAllowUnknownCommand},
	{[]string{"#DW create_persistent type=lustre external_mgs=10.0.0.1@o2ib capacity=100GB name=Extern2"}, allow, validDWOrAllowUnknownCommand},

	{[]string{"#DW stage_in  type=file      destination=$DW_JOB_STRIPED source=/pfs/dld-input "}, deny, validDWOrAllowUnknownCommand},
	{[]string{"#DW stage_in  type=directory destination=$DW_JOB_STRIPED source=/pfs/dld-input "}, deny, validDWOrAllowUnknownCommand},
	{[]string{"#DW stage_in  type=list      destination=$DW_JOB_STRIPED source=/pfs/dld-input "}, deny, validDWOrAllowUnknownCommand},
	{[]string{"#DW stage_out type=file      destination=/pfs/dld-output source=$DW_JOB_STRIPED"}, deny, validDWOrAllowUnknownCommand},
	{[]string{"#DW stage_out type=directory destination=/pfs/dld-output source=$DW_JOB_STRIPED"}, deny, validDWOrAllowUnknownCommand},
	{[]string{"#DW stage_out type=list      destination=/pfs/dld-output source=$DW_JOB_STRIPED"}, deny, validDWOrAllowUnknownCommand},

	{[]string{"#DW stage_in  type=file      destination=$DW_JOB_STRIPED source=/pfs/dld-input "}, allow, validDWOrAllowUnknownCommand},
	{[]string{"#DW stage_in  type=directory destination=$DW_JOB_STRIPED source=/pfs/dld-input "}, allow, validDWOrAllowUnknownCommand},
	{[]string{"#DW stage_in  type=list      destination=$DW_JOB_STRIPED source=/pfs/dld-input "}, allow, validDWOrAllowUnknownCommand},
	{[]string{"#DW stage_out type=file      destination=/pfs/dld-output source=$DW_JOB_STRIPED"}, allow, validDWOrAllowUnknownCommand},
	{[]string{"#DW stage_out type=directory destination=/pfs/dld-output source=$DW_JOB_STRIPED"}, allow, validDWOrAllowUnknownCommand},
	{[]string{"#DW stage_out type=list      destination=/pfs/dld-output source=$DW_JOB_STRIPED"}, allow, validDWOrAllowUnknownCommand},

	{[]string{"#DW persistentdw name=evenBetterName"}, deny, validDWOrAllowUnknownCommand},
	{[]string{"#DW persistentdw name=evenBetterName"}, allow, validDWOrAllowUnknownCommand},

	{[]string{"#DW destroy_persistent name=evenBetterName"}, deny, validDWOrAllowUnknownCommand},
	{[]string{"#DW destroy_persistent name=evenBetterName"}, allow, validDWOrAllowUnknownCommand},

	{[]string{"#DW destroy_persistent name=uniqueName1 ",
		"#DW destroy_persistent name=uniqueName2 "}, allow, validDWOrAllowUnknownCommand},
	{[]string{"#DW destroy_persistent name=conflictName",
		"#DW destroy_persistent name=conflictName"}, deny, invalidDW},

	{[]string{"#DW container name=mycontainer spec=some-repo-name                                                                         supervisor=rabbit"}, deny, validDWOrAllowUnknownCommand},
	{[]string{"#DW container name=mycontainer spec=some-repo-name job_storage={stor1}                                                     supervisor=rabbit"}, deny, validDWOrAllowUnknownCommand},
	{[]string{"#DW container name=mycontainer spec=some-repo-name job_storage={stor1,stor2,stor3}                                         supervisor=rabbit"}, deny, validDWOrAllowUnknownCommand},
	{[]string{"#DW container name=mycontainer spec=some-repo-name                                 persistent_storage={perStore,perStore2} supervisor=rabbit"}, deny, validDWOrAllowUnknownCommand},
	{[]string{"#DW container name=mycontainer spec=some-repo-name job_storage={stor1,stor2,stor3} persistent_storage={perStore,perStore2} supervisor=rabbit"}, deny, validDWOrAllowUnknownCommand},

	{[]string{"                                                         "}, deny, invalidDW},
	{[]string{"    jobdw type=raw     capacity=100GB name=prettyGoodName"}, deny, invalidDW},
	{[]string{"#DW       type=xfs     capacity=100TB name=noCommand     "}, deny, invalidDW},
	{[]string{"#DW bogus the_rest_does_not_matter                       "}, deny, invalidDW},
	{[]string{"#DW jobd  type=lustre  capacity=100TB name=badCommand    "}, deny, invalidDW},
	{[]string{"#DW jobdw tye=badtype  capacity=100TB name=badType       "}, deny, invalidDW},
	{[]string{"#DW jobdw type=badtype capacity=100TB name=badType       "}, deny, invalidDW},
	{[]string{"#DW jobdw              capacity=100TB name=missingType   "}, deny, invalidDW},
	{[]string{"#DW jobdw type=file    capacity=100TB name=badType       "}, deny, invalidDW},
	{[]string{"#DW jobdw type=raw     caacity=100TB  name=badCapacity   "}, deny, invalidDW},
	{[]string{"#DW jobdw type=raw     capacity=bad   name=badCapacity   "}, deny, invalidDW},
	{[]string{"#DW jobdw type=xfs     capacity=100TB ame=badName        "}, deny, invalidDW},
	{[]string{"#DW jobdw type=xfs     capacity=100TB name=!!21//\\      "}, deny, invalidDW},
	{[]string{"#DW jobdw                                                "}, deny, invalidDW},

	{[]string{"#DW jobdw type=raw type=raw capacity=100TB name=duplicatedTypes                         "}, deny, invalidDW},
	{[]string{"#DW jobdw type=raw capacity=100TB name=conflictingTypes type=xfs                        "}, deny, invalidDW},
	{[]string{"#DW jobdw type=badtype destination=shouldNotHaveDestination source=shouldNotHaveSource  "}, deny, invalidDW},

	{[]string{"#DW stage_in  type=badtype destination=$DW_JOB_STRIPED source=/pfs/dld-input               "}, deny, invalidDW},
	{[]string{"#DW stage_out type=file    badkey=foobar destination=/pfs/dld-output source=$DW_JOB_STRIPED"}, deny, invalidDW},

	{[]string{"#DW boguscommand the_rest_should_not_matter                                             "}, allow, validDWOrAllowUnknownCommand},
	{[]string{"#DW jobd  type=lustre  capacity=100TB name=badCommand                                   "}, allow, validDWOrAllowUnknownCommand},
	{[]string{"#DW stage_in type=badtype destination=$DW_JOB_STRIPED source=/pfs/dld-input             "}, allow, invalidDW},
	{[]string{"#DW stage_out type=file badkey=foobar destination=/pfs/dld-output source=$DW_JOB_STRIPED"}, allow, invalidDW},
	{[]string{"#DW jobd  type=lustre  capacity=100TB name=badCommand                                   "}, allow, validDWOrAllowUnknownCommand},
	{[]string{"#DW jobdw tye=badtype  capacity=100TB name=badType                                      "}, allow, invalidDW},
	{[]string{"#DW jobdw type=badtype capacity=100TB name=badType                                      "}, allow, invalidDW},
	{[]string{"#DW jobdw              capacity=100TB name=missingType                                  "}, allow, invalidDW},
	{[]string{"#DW jobdw type=file    capacity=100TB name=badType                                      "}, allow, invalidDW},
	{[]string{"#DW jobdw type=raw     caacity=100TB  name=badCapacity                                  "}, allow, invalidDW},
	{[]string{"#DW jobdw type=raw     capacity=bad   name=badCapacity                                  "}, allow, invalidDW},
	{[]string{"#DW jobdw type=xfs     capacity=100TB ame=badName                                       "}, allow, invalidDW},
	{[]string{"#DW jobdw type=xfs     capacity=100TB name=!!21//\\                                     "}, allow, invalidDW},
	{[]string{"#DW jobdw                                                                               "}, allow, invalidDW},

	{[]string{"#DW container name=mycontainer spec=some-repo-name job_storage={stor1,} supervisor=rabbit"}, deny, invalidDW},
	{[]string{"#DW container name=mycontainer spec=some-repo-name job_storage={,stor1} supervisor=rabbit"}, deny, invalidDW},
	{[]string{"#DW container name=mycontainer "}, deny, invalidDW},
}

func parsedw(t *testing.T, directiveList []string, dwRules []DWDirectiveRuleSpec, failUnknownCommand bool) error {

	// Examine each rule. If there is an error with the rule, return that as a failure.
	// Otherwise, continue looking at all the rules to see if you can find a valid rule
	// recording whether we found one. If the DWDirective matches a rule without other errors,
	// return succes.
	uniqueMap := make(map[string]bool)

	for _, directive := range directiveList {
		directiveMatchesARule := false // Anticipate failure
		for i := range dwRules {
			valid, err := ValidateDWDirective(dWDRules[i], directive, uniqueMap, failUnknownCommand)
			if err != nil {
				return err // Errors indicate parsing problems, reject directive
			}

			// The directive matched a rule
			if valid {
				directiveMatchesARule = true
			}
		}

		if !directiveMatchesARule {
			return fmt.Errorf("invalid directive found: %s", directive)
		}
	}

	return nil
}

func TestDWParse(t *testing.T) {
	for index, tt := range dwDirectiveTests {
		validCommand := false
		err := parsedw(t, tt.directiveList, dWDRules, tt.failUnknownCommand)
		if err == nil {
			validCommand = true
		}

		if validCommand != tt.validCommand {
			t.Errorf("TestDWParse(%s)(%d): allowUnsupportedCommand(%v) expect_valid(%v) err(%v)", tt.directiveList, index, tt.failUnknownCommand, tt.validCommand, err)
		}
	}
}
