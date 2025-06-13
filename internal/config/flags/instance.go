// Package flags defines flag types and domain-specific flag collections for the CLI.
package flags

import (
	"github.com/spf13/cobra"
)

// Instance flags
var (
	ListFlag = BoolFlag{
		Name:      FlagNameList,
		Shorthand: FlagShortList,
		Default:   false,
		Usage:     FlagDescList,
	}

	FindFlag = StringFlag{
		Name:      FlagNameFind,
		Shorthand: FlagShortFind,
		Default:   "",
		Usage:     FlagDescFind,
	}

	ImageDetailsFlag = BoolFlag{
		Name:      FlagNameImageDetails,
		Shorthand: FlagShortImageDetails,
		Default:   false,
		Usage:     FlagDescImageDetails,
	}
)

// instanceFlags is a slice of all instance-related flags for batch registration
var instanceFlags = []Flag{
	ListFlag,
	FindFlag,
	ImageDetailsFlag,
}

// AddInstanceFlags adds all instance-related flags to the given command
func AddInstanceFlags(cmd *cobra.Command) {
	for _, f := range instanceFlags {
		f.Add(cmd)
	}
}

// AddListFlags adds flags specific to the instance list command
func AddListFlags(cmd *cobra.Command) {
	ListFlag.Add(cmd)
}

// AddFindFlags adds flags specific to the instance find command
func AddFindFlags(cmd *cobra.Command) {
	FindFlag.Add(cmd)
	ImageDetailsFlag.Add(cmd)
}
