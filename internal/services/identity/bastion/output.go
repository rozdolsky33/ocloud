package bastion

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

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
			"Name":           b.Name,
			"BastionType":    string(b.BastionType),
			"LifecycleState": string(b.LifecycleState),
			"TargetVcn":      b.TargetVcnName,
			"TargetSubnet":   b.TargetSubnetName,
		}

		orderedKeys := []string{
			"Name", "BastionType", "LifecycleState", "TargetVcn", "TargetSubnet",
		}

		title := util.FormatColoredTitle(appCtx, b.Name)

		p.PrintKeyValues(title, bastionInfo, orderedKeys)
	}

	return nil
}
