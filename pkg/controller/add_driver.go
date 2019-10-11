package controller

import (
	"stash.us.cray.com/dpm/dws-operator/pkg/controller/driver"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, driver.Add)
}
