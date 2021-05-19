package dwdparse

import (
	"errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"regexp"
	"stash.us.cray.com/dpm/dws-operator/pkg/apis/dws/v1alpha1"
	client "stash.us.cray.com/dpm/dws-operator/pkg/client/clientset/versioned/typed/dws/v1alpha1"
	"strconv"
	"strings"
)

// BuildRulesMap builds a map of the DWDirectives argument parser rules for the specified command
func BuildRulesMap(rules []v1alpha1.DWDirectiveRuleSpec, cmd string) (map[string]v1alpha1.DWDirectiveRuleDef, error) {
	rulesMap := make(map[string]v1alpha1.DWDirectiveRuleDef)

	for _, r := range rules {
		if cmd == r.Command {
			for _, rd := range r.RuleDefs {
				rulesMap[rd.Key] = rd
			}
		}
	}

	if len(rulesMap) == 0 {
		return nil, errors.New("Unsupported #DW command " + cmd)
	}

	return rulesMap, nil
}

// BuildArgsMap builds a map of the DWDirectives arguments args["key"] = value
func BuildArgsMap(dwd string) (map[string]string, error) {
	argsMap := make(map[string]string)
	dwdArgs := strings.Fields(dwd)
	if dwdArgs[0] == "#DW" {
		argsMap["command"] = dwdArgs[1]
		for i := 2; i < len(dwdArgs); i++ {
			keyValue := strings.Split(dwdArgs[i], "=")
			if len(keyValue) == 1 {
				argsMap[keyValue[0]] = "true"
			} else if len(keyValue) == 2 {
				argsMap[keyValue[0]] = keyValue[1]
			} else {
				keyValue := strings.SplitN(dwdArgs[i], "=", 2)
				argsMap[keyValue[0]] = keyValue[1]
			}
		}
	} else {
		return nil, errors.New("missing #DW in directive")
	}
	return argsMap, nil
}

// ValidateArgs validates a map of arguments against the rules
func ValidateArgs(args map[string]string, rules []v1alpha1.DWDirectiveRuleSpec) error {
	command := args["command"]
	rulesMap, err := BuildRulesMap(rules, command)
	if err != nil {
		return err
	}

	// Compile this regex outside the loop for better performance.
	var boolMatcher = regexp.MustCompile("^(true|false|True|False|TRUE|FALSE)$")

	// Iterate over all arguments and validate each based on the associated rule
	for k, v := range args {
		if k != "command" {
			rule, found := rulesMap[k]
			if !found {
				return errors.New("Unsupported argument - " + k)
			}
			if rule.IsValueRequired && len(v) == 0 {
				return errors.New("Malformed keyword[=value]: " + k + "=" + v)
			}
			switch rule.Type {
			case "integer":
				// i,err := strconv.ParseInt(v, 10, 64)
				i, err := strconv.Atoi(v)
				if err != nil {
					return errors.New("Invalid integer argument: " + k + "=" + v)
				}
				if rule.Max != 0 && i > rule.Max {
					return errors.New("Specified integer exceeds maximum " + strconv.Itoa(rule.Max) + ": " + k + "=" + v)
				}
				if rule.Min != 0 && i < rule.Min {
					return errors.New("Specified integer smaller than minimum " + strconv.Itoa(rule.Min) + ": " + k + "=" + v)
				}
			case "bool":
				if rule.Pattern != "" {
					isok := boolMatcher.MatchString(v)
					if !isok {
						return errors.New("Invalid bool argument: " + k + "=" + v)
					}
				}
			case "string":
				if rule.Pattern != "" {
					isok, err := regexp.MatchString(rule.Pattern, v)
					if !isok {
						if err != nil {
							return errors.New("Invalid regexp in rule: " + rule.Pattern)
						}
						return errors.New("Invalid argument: " + k + "=" + v)
					}
				}
			default:
				return errors.New("Unsupported value type: " + rule.Type)
			}
		}
	}
	return nil
}

// GetParserRules is used to get the DWDirective parser rule set
func GetParserRules(ruleSetName string, namespace string) (*v1alpha1.DWDirectiveRule, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	dwsClient, err := client.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	dwdRules, err := dwsClient.DWDirectiveRules(namespace).Get(ruleSetName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return dwdRules, nil
}

// ValidateDWDirectives will validate a set of #DW directives against a specified rule set
func ValidateDWDirectives(directives []string, ruleSetName string, namespace string) error {

	dwdRules, err := GetParserRules(ruleSetName, namespace)
	if err != nil {
		return err
	}

	for _, dwd := range directives {
		// Build a map of the #DW commands and arguments
		argsMap, err := BuildArgsMap(dwd)
		if err != nil {
			return err
		}

		err = ValidateArgs(argsMap, dwdRules.Spec)
		if err != nil {
			return err
		}
	}

	return nil
}
