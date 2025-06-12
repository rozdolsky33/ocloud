package compute

import (
	"fmt"
	"github.com/rozdolsky33/ocloud/internal/helpers"
)

func ListInstances() {
	helpers.Logger.V(1).Info("ListInstances()")
	fmt.Println("Inside Instance resources running List Instances")

}
