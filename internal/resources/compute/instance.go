package compute

import (
	"fmt"
	"github.com/rozdolsky33/ocloud/cmd/configuration"
	"github.com/rozdolsky33/ocloud/internal/helpers"
)

func ListInstances(appCtx *configuration.AppContext) error {

	helpers.Logger.V(1).Info("ListInstances()")
	fmt.Println("Inside Instance resources running List Instances" +
		" with tenancyID: " + appCtx.TenancyID)
	return nil
}
