package cmd

import (
	"fmt"
	"github.com/rozdolsky33/ocloud/internal/config"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "ocloud",
	Short: "Interact with Oracle Cloud Infrastructure",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "auth" {
			return config.InitAuth()
		}
		return config.Init()
	},
}

// Execute runs the CLI
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("tenancy", "t", "", "OCI Tenancy OCID (override)")
	rootCmd.PersistentFlags().StringP("compartment", "c", "", "OCI Compartment OCID (required)")
	rootCmd.PersistentFlags().StringP("region", "r", "", "OCI Region (override)")

	viper.BindPFlag("tenancy", rootCmd.PersistentFlags().Lookup("tenancy"))
	viper.BindPFlag("compartment", rootCmd.PersistentFlags().Lookup("compartment"))
	viper.BindPFlag("region", rootCmd.PersistentFlags().Lookup("region"))

	viper.SetEnvPrefix("OCI")
	viper.AutomaticEnv()
}
