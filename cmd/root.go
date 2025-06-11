package cmd

import (
	"context"
	"fmt"
	"github.com/rozdolsky33/ocloud/cmd/configs"
	"github.com/rozdolsky33/ocloud/internal/helpers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "ocloud",
	Short: "Interact with Oracle Cloud Infrastructure",
	Long:  "",
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&helpers.LogLevel, "log-level", "l", "info", helpers.LogLevelMsg)
	rootCmd.PersistentFlags().BoolVar(&helpers.ColoredOutput, "color", false, helpers.ColoredOutputMsg)
	rootCmd.AddCommand(configs.ConfigCmd)

	// allow ENV overrides, e.g., OCI_CLI_TENANCY, OCI_TENANCY_NAME, OCI_COMPARTMENT
	viper.SetEnvPrefix("OCI")
	viper.AutomaticEnv()
}

func Execute(ctx context.Context) {
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
