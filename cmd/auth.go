// cmd/auth.go
package cmd

import (
	"bufio"
	"fmt"
	"github.com/rozdolsky33/ocloud/internal/config"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with OCI and refresh session tokens",
	Long: `Runs the OCI CLI's session authenticate under the hood:

    oci session authenticate --profile-name <PROFILE> --region <REGION>

Interactively lets you pick your desired region.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		profile := viper.GetString("profile")
		// 1) Init SDK provider (tenancy, profile, region if present)
		if err := config.InitAuth(); err != nil {
			return errors.Wrap(err, "initializing auth config")
		}

		// 2) Static list of OCI regions
		regions := []string{
			"af-johannesburg-1", "ap-batam-1", "ap-chiyoda-1", "ap-chuncheon-1",
			"ap-chuncheon-2", "ap-dcc-canberra-1", "ap-dcc-gazipur-1", "ap-hyderabad-1",
			"ap-ibaraki-1", "ap-melbourne-1", "ap-mumbai-1", "ap-osaka-1", "ap-seoul-1",
			"ap-seoul-2", "ap-singapore-1", "ap-singapore-2", "ap-suwon-1", "ap-sydney-1",
			"ap-tokyo-1", "ca-montreal-1", "ca-toronto-1", "eu-amsterdam-1", "eu-crissier-1",
			"eu-dcc-dublin-1", "eu-dcc-dublin-2", "eu-dcc-milan-1", "eu-dcc-milan-2",
			"eu-dcc-rating-1", "eu-dcc-rating-2", "eu-dcc-zurich-1", "eu-frankfurt-1",
			"eu-frankfurt-2", "eu-jovanovac-1", "eu-madrid-1", "eu-madrid-2",
			"eu-marseille-1", "eu-milan-1", "eu-paris-1", "eu-stockholm-1",
			"eu-zurich-1", "il-jerusalem-1", "me-abudhabi-1", "me-abudhabi-2",
			"me-abudhabi-3", "me-abudhabi-4", "me-alain-1", "me-dcc-doha-1",
			"me-dcc-muscat-1", "me-dubai-1", "me-jeddah-1", "me-riyadh-1",
			"mx-monterrey-1", "mx-queretaro-1", "sa-bogota-1", "sa-santiago-1",
			"sa-saopaulo-1", "sa-valparaiso-1", "sa-vinhedo-1", "uk-cardiff-1",
			"uk-gov-cardiff-1", "uk-gov-london-1", "uk-london-1", "us-abilene-1",
			"us-ashburn-1", "us-chicago-1", "us-dallas-1", "us-gov-ashburn-1",
			"us-gov-chicago-1", "us-gov-phoenix-1", "us-langley-1", "us-luke-1",
			"us-phoenix-1", "us-saltlake-2", "us-sanjose-1", "us-somerset-1",
			"us-thames-1",
		}

		// 3) Prompt region selection
		fmt.Println("Select a region:")
		for i, r := range regions {
			fmt.Printf("%d: %s\n", i+1, r)
		}
		fmt.Print("Enter region number or name: ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return errors.Wrap(err, "reading region input")
		}
		input = strings.TrimSpace(input)
		var chosen string
		if idx, err := strconv.Atoi(input); err == nil && idx >= 1 && idx <= len(regions) {
			chosen = regions[idx-1]
		} else {
			chosen = input
		}
		fmt.Printf("Using region: %s\n", chosen)
		// 4) Authenticate via OCI CLI
		ociCmd := exec.Command("oci", "session", "authenticate", "--profile-name", profile, "--region", chosen)
		ociCmd.Stdout = os.Stdout
		ociCmd.Stderr = os.Stderr
		if err := ociCmd.Run(); err != nil {
			return errors.Wrap(err, "failed to run `oci session authenticate`")
		}
		// 5) Reload provider with chosen profile/region
		os.Setenv("OCI_PROFILE", profile)
		os.Setenv("OCI_REGION", chosen)
		if err := config.InitAuth(); err != nil {
			return errors.Wrap(err, "reloading config after auth")
		}
		// 6) Fetch root compartment (tenancy) OCID
		tenancyOCID, err := config.ConfigProvider.TenancyOCID()
		if err != nil {
			return errors.Wrap(err, "fetching tenancy OCID")
		}
		// 7) Print export for compartment
		fmt.Printf("export OCI_COMPARTMENT=%s\n", tenancyOCID)
		fmt.Println("✅ Authentication complete. Run `eval $(ocloud auth)` to set your env.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.Flags().StringP("profile", "p", "DEFAULT", "OCI config profile to authenticate")
	viper.BindPFlag("profile", authCmd.Flags().Lookup("profile"))
}
