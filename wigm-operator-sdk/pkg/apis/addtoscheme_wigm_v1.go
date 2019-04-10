package apis

import (
	"github.com/carsonoid/talk-kubernetes-operators/wigm-operator-sdk/pkg/apis/wigm/v1"
)

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	AddToSchemes = append(AddToSchemes, v1.SchemeBuilder.AddToScheme)
}
