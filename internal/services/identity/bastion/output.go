package bastion

import (
	"fmt"
	"strings"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// PrintBastionInfo displays bastion instances in a formatted table or JSON format.
func PrintBastionInfo(bastions []Bastion, appCtx *app.ApplicationContext, useJSON bool) error {

	p := printer.New(appCtx.Stdout)
	if useJSON {
		if len(bastions) == 0 {
			return p.MarshalToJSON(struct{}{})
		}
		return p.MarshalToJSON(bastions)
	}

	for _, b := range bastions {
		bastionInfo := map[string]string{
			"Name":           b.DisplayName,
			"BastionType":    b.BastionType,
			"LifecycleState": b.LifecycleState,
			"TargetVcn":      b.TargetVcnName,
			"TargetSubnet":   b.TargetSubnetName,
		}

		orderedKeys := []string{
			"Name", "BastionType", "LifecycleState", "TargetVcn", "TargetSubnet",
		}

		if b.MaxSessionTTL > 0 {
			hours := b.MaxSessionTTL / 3600
			bastionInfo["MaxSessionTTL"] = fmt.Sprintf("%d hours (%d seconds)", hours, b.MaxSessionTTL)
			orderedKeys = append(orderedKeys, "MaxSessionTTL")
		}

		if len(b.ClientCidrBlockAllowList) > 0 {
			bastionInfo["CIDRAllowList"] = strings.Join(b.ClientCidrBlockAllowList, ", ")
			orderedKeys = append(orderedKeys, "CIDRAllowList")
		}

		if b.PrivateEndpointIpAddress != "" {
			bastionInfo["PrivateEndpointIP"] = b.PrivateEndpointIpAddress
			orderedKeys = append(orderedKeys, "PrivateEndpointIP")
		}

		title := util.FormatColoredTitle(appCtx, b.DisplayName)

		p.PrintKeyValues(title, bastionInfo, orderedKeys)
	}

	return nil
}
