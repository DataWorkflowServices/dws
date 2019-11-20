package dwdparse

import (
	"fmt"
	"errors"
	"strings"
)

var validCommands = "jobdw persistentdw stage_in stage-in stage_out stage-out"
var argsWithPath = "source destination"

func BuildArgsMap(dwd string) (map[string]string, error) {
	argsMap := make(map[string]string)
	dwdArgs := strings.Fields(dwd)
	if dwdArgs[0] == "#DW" {
		if strings.Contains(validCommands, dwdArgs[1]) {
			argsMap["command"] = dwdArgs[1]
			//for _, cmd := range dwdArgs {
			for i:= 2; i<len(dwdArgs); i++ {
				keyValue := strings.Split(dwdArgs[i], "=")
				if len(keyValue) == 1 {
					argsMap[keyValue[0]] = "true"
				} else if len(keyValue) == 2 {
					argsMap[keyValue[0]] = keyValue[1]
				} else {
					if strings.Contains(argsWithPath, keyValue[0]) {
						keyValue := strings.SplitN(dwdArgs[i], "=", 2)
						argsMap[keyValue[0]] = keyValue[1]
					} else {
						return nil, errors.New("Invalid: Malformed keyword[=value]: " + dwdArgs[i])
					}
				}
			}
		} else {	
			return nil, errors.New("Invalid: Unsupported #DW command " + dwdArgs[1])
		}
	} else {	
		return nil, errors.New("Invalid: Missing #DW in directive")
	}
	return argsMap, nil
}

func ValidateArgsMap(args map[string]string) error {
	switch args["command"] {
    case "jobdw":
        fmt.Println("jobdw ok")
    case "persistentdw":
        fmt.Println("persistentdw ok")
    case "stage-in":
		fallthrough
    case "stage_in":
        fmt.Println("stage-in ok")
    case "stage-out":
		fallthrough
    case "stage_out":
        fmt.Println("stage-out ok")
    }
	return nil
}
