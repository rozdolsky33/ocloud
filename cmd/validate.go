package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var validateCmd = &cobra.Command{
	Use:   "validate-session",
	Short: "Validate OCI CLI session token",
	RunE: func(cmd *cobra.Command, args []string) error {
		profile := viper.GetString("profile")
		validate := exec.Command("oci", "session", "validate", "--profile", profile)
		validate.Stdout = os.Stdout
		validate.Stderr = os.Stderr
		if err := validate.Run(); err != nil {
			return errors.Wrap(err, "OCI session validation failed")
		}
		fmt.Println("✅ Session is valid.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
	validateCmd.Flags().StringP("profile", "p", "DEFAULT", "OCI config profile to validate")
	viper.BindPFlag("profile", validateCmd.Flags().Lookup("profile"))
}
