package dwdparse

import (
	"testing"
    "stash.us.cray.com/dpm/dws-operator/pkg/apis/dws/v1alpha1"
)

var DWDRules = []v1alpha1.DWDirectiveRuleSpec {
	v1alpha1.DWDirectiveRuleSpec {
		Command: "stage_in",
		RuleDefs: []v1alpha1.DWDirectiveRuleDef {
			v1alpha1.DWDirectiveRuleDef {
				Key: "type",
				Type: "string",
				Pattern: "^(directory|file|list)$",
				IsRequired: true,
				IsValueRequired: true,
			},
			v1alpha1.DWDirectiveRuleDef {
				Key: "source",
				Type: "string",
				IsRequired: true,
				IsValueRequired: true,
			},
			v1alpha1.DWDirectiveRuleDef {
				Key: "destination",
				Type: "string",
				IsRequired: true,
				IsValueRequired: true,
			},
		},
	},

	v1alpha1.DWDirectiveRuleSpec {
		Command: "stage_out",
		RuleDefs: []v1alpha1.DWDirectiveRuleDef {
			v1alpha1.DWDirectiveRuleDef {
				Key: "type",
				Type: "string",
				Pattern: "^(directory|file|list)$",
				IsRequired: true,
				IsValueRequired: true,
			},
			v1alpha1.DWDirectiveRuleDef {
				Key: "source",
				Type: "string",
				IsRequired: true,
				IsValueRequired: true,
			},
			v1alpha1.DWDirectiveRuleDef {
				Key: "destination",
				Type: "string",
				IsRequired: true,
				IsValueRequired: true,
			},
		},
	},
}


// Directives with proper syntax
var DWDirectives0 = []string {
	"#DW stage_in type=file destination=$DW_JOB_STRIPED source=/pfs/dld-input",
	"#DW stage_out type=file destination=/pfs/dld-output source=$DW_JOB_STRIPED",
	"#DW stage_in type=directory destination=$DW_JOB_STRIPED source=/pfs/dld-input",
	"#DW stage_out type=directory destination=/pfs/dld-output source=$DW_JOB_STRIPED",
	"#DW stage_in type=list destination=$DW_JOB_STRIPED source=/pfs/dld-input",
	"#DW stage_out type=list destination=/pfs/dld-output source=$DW_JOB_STRIPED",
}

// Directives with bad/unsupported syntax
var DWDirectives1 = []string {
	"#DW jobdw type=scratch capacity=10GiB access_mode=striped max_mds=yes",
	"#DW stage_in type=badtype destination=$DW_JOB_STRIPED source=/pfs/dld-input",
	"#DW stage_out type=file badkey=foobar destination=/pfs/dld-output source=$DW_JOB_STRIPED",
}

// Test parser with valid #DW syntax
func TestDWDParse0(t *testing.T) {
    for _, dwd := range DWDirectives0 {
        // Build a map of the #DW commands and arguments
        argsMap,err := BuildArgsMap(dwd)
        if err != nil {
			t.Errorf("#DW parsing error: %v", err)
        }

        err = ValidateArgs(argsMap, DWDRules)
        if err != nil {
			t.Errorf("#DW parsing error: %v", err)
        }
    }
}

// Test parser with invalid #DW syntax
func TestDWDParse1(t *testing.T) {
    for _, dwd := range DWDirectives1 {
        // Build a map of the #DW commands and arguments
        argsMap,err := BuildArgsMap(dwd)
        if err != nil {
			t.Errorf("#DW syntax error: %v", err)
        }

        err = ValidateArgs(argsMap, DWDRules)
        if err == nil {
			t.Errorf("#DW parsing error not detected:\n%v\n", dwd)
        }
    }
}
