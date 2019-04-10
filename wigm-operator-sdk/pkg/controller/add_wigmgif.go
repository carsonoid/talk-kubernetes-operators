package controller

import (
	"github.com/carsonoid/talk-kubernetes-operators/wigm-operator-sdk/pkg/controller/wigmgif"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, wigmgif.Add)
}
