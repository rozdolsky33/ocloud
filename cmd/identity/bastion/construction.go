package bastion

import (
	"github.com/common-nighthawk/go-figure"
)

// ShowConstructionAnimation displays a simple "Under Construction" banner
func ShowConstructionAnimation() {
	// Print ASCII banner in yellow
	figure.NewColorFigure("Under Construction", "", "yellow", true).Print()
}
