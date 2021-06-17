package dwdparse

import (
	"testing"

	"stash.us.cray.com/dpm/dws-operator/pkg/apis/dws/v1alpha1"
)

var dWDRules = []v1alpha1.DWDirectiveRuleSpec{
	{
		Command: "jobdw",
		RuleDefs: []v1alpha1.DWDirectiveRuleDef{
			{
				Key:             "type",
				Type:            "string",
				Pattern:         "^(raw|xfs|lustre)$",
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
			},
			{
				Key:             "modifiers",
				Type:            "string",
				IsRequired:      false,
				IsValueRequired: true,
			},
		},
	},
	{
		Command: "create_persistent",
		RuleDefs: []v1alpha1.DWDirectiveRuleDef{
			{
				Key:             "type",
				Type:            "string",
				Pattern:         "^(raw|xfs|lustre)$",
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
			},
		},
	},
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
	{
		Command: "persistentdw",
		RuleDefs: []v1alpha1.DWDirectiveRuleDef{
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
		RuleDefs: []v1alpha1.DWDirectiveRuleDef{
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
		Command: "container",
		RuleDefs: []v1alpha1.DWDirectiveRuleDef{
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
	directive          string // #DW directive
	failUnknownCommand bool   // deny/allow unknown commands
	errEqualNil        bool   // expected parse error result compared with nil
}{
	{"#DW jobdw type=raw    capacity=100GB name=prettyGoodName  ", deny, validDWOrAllowUnknownCommand},
	{"#DW jobdw type=xfs    capacity=100GB name=pretty_GoodName ", deny, validDWOrAllowUnknownCommand},
	{"#DW jobdw type=lustre capacity=100GB name=prettyGood-Name ", deny, validDWOrAllowUnknownCommand},
	{"#DW jobdw type=raw    capacity=100TB name=__prettyGoodName", deny, validDWOrAllowUnknownCommand},
	{"#DW jobdw type=xfs    capacity=100TB name=-prettyGoodName-", deny, validDWOrAllowUnknownCommand},
	{"#DW jobdw type=lustre capacity=100TB name=0prettyGoodName1", deny, validDWOrAllowUnknownCommand},
	{"#DW jobdw type=raw    capacity=100GB name=prettierGoodName", deny, validDWOrAllowUnknownCommand},
	{"#DW jobdw type=xfs    capacity=100GB name=prettierGoodName", deny, validDWOrAllowUnknownCommand},
	{"#DW jobdw type=lustre capacity=100GB name=prettierGoodName", deny, validDWOrAllowUnknownCommand},
	{"#DW jobdw type=lustre capacity=100GB name=prettierGoodName", deny, validDWOrAllowUnknownCommand},

	{"#DW jobdw type=raw    capacity=100GB name=prettyGoodName  ", allow, validDWOrAllowUnknownCommand},
	{"#DW jobdw type=xfs    capacity=100GB name=prettyGoodName  ", allow, validDWOrAllowUnknownCommand},
	{"#DW jobdw type=lustre capacity=100GB name=prettyGoodName  ", allow, validDWOrAllowUnknownCommand},
	{"#DW jobdw type=raw    capacity=100TB name=prettyGoodName  ", allow, validDWOrAllowUnknownCommand},
	{"#DW jobdw type=xfs    capacity=100TB name=prettyGoodName  ", allow, validDWOrAllowUnknownCommand},
	{"#DW jobdw type=lustre capacity=100TB name=prettyGoodName  ", allow, validDWOrAllowUnknownCommand},
	{"#DW jobdw type=raw    capacity=100GB name=prettierGoodName", allow, validDWOrAllowUnknownCommand},
	{"#DW jobdw type=xfs    capacity=100GB name=prettierGoodName", allow, validDWOrAllowUnknownCommand},
	{"#DW jobdw type=lustre capacity=100GB name=prettierGoodName", allow, validDWOrAllowUnknownCommand},

	{"#DW create_persistent type=raw    capacity=100GB name=prettyGoodName  ", deny, validDWOrAllowUnknownCommand},
	{"#DW create_persistent type=xfs    capacity=100GB name=prettyGoodName  ", deny, validDWOrAllowUnknownCommand},
	{"#DW create_persistent type=lustre capacity=100GB name=prettyGoodName  ", deny, validDWOrAllowUnknownCommand},
	{"#DW create_persistent type=raw    capacity=100TB name=prettyGoodName  ", deny, validDWOrAllowUnknownCommand},
	{"#DW create_persistent type=xfs    capacity=100TB name=prettyGoodName  ", deny, validDWOrAllowUnknownCommand},
	{"#DW create_persistent type=lustre capacity=100TB name=prettyGoodName  ", deny, validDWOrAllowUnknownCommand},
	{"#DW create_persistent type=raw    capacity=100GB name=prettierGoodName", deny, validDWOrAllowUnknownCommand},
	{"#DW create_persistent type=xfs    capacity=100GB name=prettierGoodName", deny, validDWOrAllowUnknownCommand},
	{"#DW create_persistent type=lustre capacity=100GB name=prettierGoodName", deny, validDWOrAllowUnknownCommand},

	{"#DW create_persistent type=raw    capacity=100GB name=prettyGoodName  ", allow, validDWOrAllowUnknownCommand},
	{"#DW create_persistent type=xfs    capacity=100GB name=prettyGoodName  ", allow, validDWOrAllowUnknownCommand},
	{"#DW create_persistent type=lustre capacity=100GB name=prettyGoodName  ", allow, validDWOrAllowUnknownCommand},
	{"#DW create_persistent type=raw    capacity=100TB name=prettyGoodName  ", allow, validDWOrAllowUnknownCommand},
	{"#DW create_persistent type=xfs    capacity=100TB name=prettyGoodName  ", allow, validDWOrAllowUnknownCommand},
	{"#DW create_persistent type=lustre capacity=100TB name=prettyGoodName  ", allow, validDWOrAllowUnknownCommand},
	{"#DW create_persistent type=raw    capacity=100GB name=prettierGoodName", allow, validDWOrAllowUnknownCommand},
	{"#DW create_persistent type=xfs    capacity=100GB name=prettierGoodName", allow, validDWOrAllowUnknownCommand},
	{"#DW create_persistent type=lustre capacity=100GB name=prettierGoodName", allow, validDWOrAllowUnknownCommand},

	{"#DW stage_in  type=file      destination=$DW_JOB_STRIPED source=/pfs/dld-input ", deny, validDWOrAllowUnknownCommand},
	{"#DW stage_in  type=directory destination=$DW_JOB_STRIPED source=/pfs/dld-input ", deny, validDWOrAllowUnknownCommand},
	{"#DW stage_in  type=list      destination=$DW_JOB_STRIPED source=/pfs/dld-input ", deny, validDWOrAllowUnknownCommand},
	{"#DW stage_out type=file      destination=/pfs/dld-output source=$DW_JOB_STRIPED", deny, validDWOrAllowUnknownCommand},
	{"#DW stage_out type=directory destination=/pfs/dld-output source=$DW_JOB_STRIPED", deny, validDWOrAllowUnknownCommand},
	{"#DW stage_out type=list      destination=/pfs/dld-output source=$DW_JOB_STRIPED", deny, validDWOrAllowUnknownCommand},

	{"#DW stage_in  type=file      destination=$DW_JOB_STRIPED source=/pfs/dld-input ", allow, validDWOrAllowUnknownCommand},
	{"#DW stage_in  type=directory destination=$DW_JOB_STRIPED source=/pfs/dld-input ", allow, validDWOrAllowUnknownCommand},
	{"#DW stage_in  type=list      destination=$DW_JOB_STRIPED source=/pfs/dld-input ", allow, validDWOrAllowUnknownCommand},
	{"#DW stage_out type=file      destination=/pfs/dld-output source=$DW_JOB_STRIPED", allow, validDWOrAllowUnknownCommand},
	{"#DW stage_out type=directory destination=/pfs/dld-output source=$DW_JOB_STRIPED", allow, validDWOrAllowUnknownCommand},
	{"#DW stage_out type=list      destination=/pfs/dld-output source=$DW_JOB_STRIPED", allow, validDWOrAllowUnknownCommand},

	{"#DW persistentdw name=evenBetterName", deny, validDWOrAllowUnknownCommand},
	{"#DW persistentdw name=evenBetterName", allow, validDWOrAllowUnknownCommand},

	{"#DW destroy_persistent name=evenBetterName", deny, validDWOrAllowUnknownCommand},
	{"#DW destroy_persistent name=evenBetterName", allow, validDWOrAllowUnknownCommand},

	{"#DW container name=mycontainer spec=some-repo-name                                                                         supervisor=rabbit", deny, validDWOrAllowUnknownCommand},
	{"#DW container name=mycontainer spec=some-repo-name job_storage={stor1}                                                     supervisor=rabbit", deny, validDWOrAllowUnknownCommand},
	{"#DW container name=mycontainer spec=some-repo-name job_storage={stor1,stor2,stor3}                                         supervisor=rabbit", deny, validDWOrAllowUnknownCommand},
	{"#DW container name=mycontainer spec=some-repo-name                                 persistent_storage={perStore,perStore2} supervisor=rabbit", deny, validDWOrAllowUnknownCommand},
	{"#DW container name=mycontainer spec=some-repo-name job_storage={stor1,stor2,stor3} persistent_storage={perStore,perStore2} supervisor=rabbit", deny, validDWOrAllowUnknownCommand},

	{"#DW       type=xfs     capacity=100TB name=noCommand    ", deny, invalidDW},
	{"#DW bogus the_rest_does_not_matter               ", deny, invalidDW},
	{"#DW jobd  type=lustre  capacity=100TB name=badCommand   ", deny, invalidDW},
	{"#DW jobdw tye=badtype  capacity=100TB name=badType      ", deny, invalidDW},
	{"#DW jobdw type=badtype capacity=100TB name=badType      ", deny, invalidDW},
	{"#DW jobdw              capacity=100TB name=missingType  ", deny, invalidDW},
	{"#DW jobdw type=file    capacity=100TB name=badType      ", deny, invalidDW},
	{"#DW jobdw type=raw     caacity=100TB  name=badCapacity  ", deny, invalidDW},
	{"#DW jobdw type=raw     capacity=bad   name=badCapacity  ", deny, invalidDW},
	{"#DW jobdw type=xfs     capacity=100TB ame=badName       ", deny, invalidDW},
	{"#DW jobdw type=xfs     capacity=100TB name=!!21//\\     ", deny, invalidDW},
	{"#DW jobdw                                               ", deny, invalidDW},

	{"#DW jobdw type=raw type=raw capacity=100TB name=duplicatedTypes                         ", deny, invalidDW},
	{"#DW jobdw type=raw capacity=100TB name=conflictingTypes type=xfs                        ", deny, invalidDW},
	{"#DW jobdw type=badtype destination=shouldNotHaveDestination source=shouldNotHaveSource  ", deny, invalidDW},

	{"#DW stage_in  type=badtype destination=$DW_JOB_STRIPED source=/pfs/dld-input               ", deny, invalidDW},
	{"#DW stage_out type=file    badkey=foobar destination=/pfs/dld-output source=$DW_JOB_STRIPED", deny, invalidDW},

	{"#DW boguscommand the_rest_should_not_matter                                             ", allow, validDWOrAllowUnknownCommand},
	{"#DW jobd  type=lustre  capacity=100TB name=badCommand                                   ", allow, validDWOrAllowUnknownCommand},
	{"#DW stage_in type=badtype destination=$DW_JOB_STRIPED source=/pfs/dld-input             ", allow, invalidDW},
	{"#DW stage_out type=file badkey=foobar destination=/pfs/dld-output source=$DW_JOB_STRIPED", allow, invalidDW},
	{"#DW jobd  type=lustre  capacity=100TB name=badCommand                                   ", allow, validDWOrAllowUnknownCommand},
	{"#DW jobdw tye=badtype  capacity=100TB name=badType                                      ", allow, invalidDW},
	{"#DW jobdw type=badtype capacity=100TB name=badType                                      ", allow, invalidDW},
	{"#DW jobdw              capacity=100TB name=missingType                                  ", allow, invalidDW},
	{"#DW jobdw type=file    capacity=100TB name=badType                                      ", allow, invalidDW},
	{"#DW jobdw type=raw     caacity=100TB  name=badCapacity                                  ", allow, invalidDW},
	{"#DW jobdw type=raw     capacity=bad   name=badCapacity                                  ", allow, invalidDW},
	{"#DW jobdw type=xfs     capacity=100TB ame=badName                                       ", allow, invalidDW},
	{"#DW jobdw type=xfs     capacity=100TB name=!!21//\\                                     ", allow, invalidDW},
	{"#DW jobdw                                                                               ", allow, invalidDW},

	{"#DW container name=mycontainer spec=some-repo-name job_storage={stor1,} supervisor=rabbit", deny, invalidDW},
	{"#DW container name=mycontainer spec=some-repo-name job_storage={,stor1} supervisor=rabbit", deny, invalidDW},
	{"#DW container name=mycontainer ", deny, invalidDW},
}

func parsedw(t *testing.T, dwd string, dwRules []v1alpha1.DWDirectiveRuleSpec, failUnknownCommand bool) bool {
	// Build a map of the next #DW command and arguments
	// The command is not checked at this point!
	argsMap, err := BuildArgsMap(dwd)
	if err != nil {
		return err == nil
	}

	err = ValidateArgs(argsMap, dWDRules, failUnknownCommand)
	return err == nil
}

func TestDWParse(t *testing.T) {
	for index, tt := range dwDirectiveTests {
		actual := parsedw(t, tt.directive, dWDRules, tt.failUnknownCommand)

		if actual != tt.errEqualNil {
			t.Errorf("TestDWParse(%s)(%d): allowUnsupportedCommand(%v) expect_errEqualsNil(%v) actualFailure(%v)", tt.directive, index, tt.failUnknownCommand, tt.errEqualNil, actual)
		}
	}
}
