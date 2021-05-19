package dwdparse

import (
	"stash.us.cray.com/dpm/dws-operator/pkg/apis/dws/v1alpha1"
	"testing"
)

var dWDRules = []v1alpha1.DWDirectiveRuleSpec{
	{
		Command: "stage_in",
		RuleDefs: []v1alpha1.DWDirectiveRuleDef{
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
		RuleDefs: []v1alpha1.DWDirectiveRuleDef{
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
}

// Directives with proper syntax
var DWDirectives0 = []string{
	"#DW stage_in type=file destination=$DW_JOB_STRIPED source=/pfs/dld-input",
	"#DW stage_out type=file destination=/pfs/dld-output source=$DW_JOB_STRIPED",
	"#DW stage_in type=directory destination=$DW_JOB_STRIPED source=/pfs/dld-input",
	"#DW stage_out type=directory destination=/pfs/dld-output source=$DW_JOB_STRIPED",
	"#DW stage_in type=list destination=$DW_JOB_STRIPED source=/pfs/dld-input",
	"#DW stage_out type=list destination=/pfs/dld-output source=$DW_JOB_STRIPED",
}

// Directives with bad/unsupported syntax
var DWDirectives1 = []string{
	"#DW boguscommand the_rest_should_not_matter",
	"#DW jobdw type=scratch capacity=10GiB access_mode=striped max_mds=yes",
	"#DW stage_in type=badtype destination=$DW_JOB_STRIPED source=/pfs/dld-input",
	"#DW stage_out type=file badkey=foobar destination=/pfs/dld-output source=$DW_JOB_STRIPED",
}

// Test parser with valid #DW syntax
func TestDWDParseValid(t *testing.T) {
	for _, dwd := range DWDirectives0 {
		// Build a map of the #DW commands and arguments
		argsMap, err := BuildArgsMap(dwd)
		if err != nil {
			t.Errorf("#DW parsing error: %v", err)
		}

		err = ValidateArgs(argsMap, dWDRules)
		if err != nil {
			t.Errorf("#DW parsing error: %v", err)
		}
	}
}

// Test parser with invalid #DW syntax
func TestDWDParseInvalid(t *testing.T) {
	for _, dwd := range DWDirectives1 {
		// Build a map of the #DW commands and arguments
		argsMap, err := BuildArgsMap(dwd)
		if err != nil {
			t.Errorf("#DW syntax error: %v", err)
		}

		err = ValidateArgs(argsMap, dWDRules)
		if err == nil {
			t.Errorf("#DW parsing error not detected:\n%v\n", dwd)
		}
	}
}
