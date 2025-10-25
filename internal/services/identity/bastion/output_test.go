package bastion

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/rozdolsky33/ocloud/internal/app"
	domain "github.com/rozdolsky33/ocloud/internal/domain/identity"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestPrintBastionInfo_JSON(t *testing.T) {
	tests := []struct {
		name        string
		bastions    []domain.Bastion
		expectEmpty bool
	}{
		{
			name: "single bastion with all fields",
			bastions: []domain.Bastion{
				{
					OCID:                     "ocid1.bastion.oc1..test",
					DisplayName:              "test-bastion",
					BastionType:              "STANDARD",
					LifecycleState:           "ACTIVE",
					TargetVcnName:            "test-vcn",
					TargetSubnetName:         "test-subnet",
					MaxSessionTTL:            10800,
					ClientCidrBlockAllowList: []string{"0.0.0.0/0"},
					PrivateEndpointIpAddress: "10.0.0.1",
				},
			},
			expectEmpty: false,
		},
		{
			name:        "no bastions",
			bastions:    []domain.Bastion{},
			expectEmpty: true,
		},
		{
			name: "multiple bastions",
			bastions: []domain.Bastion{
				{
					OCID:             "ocid1.bastion.oc1..test1",
					DisplayName:      "bastion-1",
					BastionType:      "STANDARD",
					LifecycleState:   "ACTIVE",
					TargetVcnName:    "vcn-1",
					TargetSubnetName: "subnet-1",
					MaxSessionTTL:    10800,
				},
				{
					OCID:             "ocid1.bastion.oc1..test2",
					DisplayName:      "bastion-2",
					BastionType:      "STANDARD",
					LifecycleState:   "CREATING",
					TargetVcnName:    "vcn-2",
					TargetSubnetName: "subnet-2",
					MaxSessionTTL:    7200,
				},
			},
			expectEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			appCtx := &app.ApplicationContext{
				Stdout: buf,
				Logger: logger.NewTestLogger(),
			}

			err := PrintBastionInfo(tt.bastions, appCtx, true)
			assert.NoError(t, err)

			// Verify JSON output
			output := buf.String()
			assert.NotEmpty(t, output)

			if tt.expectEmpty {
				// Empty array or empty object
				var result interface{}
				err = json.Unmarshal([]byte(output), &result)
				assert.NoError(t, err)
			} else {
				// Verify valid JSON array
				var result []domain.Bastion
				err = json.Unmarshal([]byte(output), &result)
				assert.NoError(t, err)
				assert.Equal(t, len(tt.bastions), len(result))
			}
		})
	}
}

func TestPrintBastionInfo_Table(t *testing.T) {
	tests := []struct {
		name           string
		bastions       []domain.Bastion
		expectContains []string
	}{
		{
			name: "bastion with all optional fields",
			bastions: []domain.Bastion{
				{
					OCID:                     "ocid1.bastion.oc1..test",
					DisplayName:              "test-bastion",
					BastionType:              "STANDARD",
					LifecycleState:           "ACTIVE",
					TargetVcnName:            "test-vcn",
					TargetSubnetName:         "test-subnet",
					MaxSessionTTL:            10800,
					ClientCidrBlockAllowList: []string{"0.0.0.0/0", "10.0.0.0/8"},
					PrivateEndpointIpAddress: "10.0.0.1",
				},
			},
			expectContains: []string{
				"test-bastion",
				"STANDARD",
				"ACTIVE",
				"test-vcn",
				"test-subnet",
				"3 hours", // MaxSessionTTL formatted
				"0.0.0.0/0, 10.0.0.0/8",
				"10.0.0.1",
			},
		},
		{
			name: "bastion without optional fields",
			bastions: []domain.Bastion{
				{
					OCID:             "ocid1.bastion.oc1..minimal",
					DisplayName:      "minimal-bastion",
					BastionType:      "STANDARD",
					LifecycleState:   "CREATING",
					TargetVcnName:    "vcn",
					TargetSubnetName: "subnet",
				},
			},
			expectContains: []string{
				"minimal-bastion",
				"STANDARD",
				"CREATING",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			appCtx := &app.ApplicationContext{
				Stdout:          buf,
				Logger:          logger.NewTestLogger(),
				CompartmentName: "TestCompartment",
				TenancyName:     "TestTenancy",
			}

			err := PrintBastionInfo(tt.bastions, appCtx, false)
			assert.NoError(t, err)

			output := buf.String()
			for _, expected := range tt.expectContains {
				assert.Contains(t, output, expected, "Output should contain: %s", expected)
			}
		})
	}
}

func TestPrintBastionInfo_MaxSessionTTLFormatting(t *testing.T) {
	tests := []struct {
		name           string
		maxSessionTTL  int
		expectedFormat string
	}{
		{
			name:           "3 hours",
			maxSessionTTL:  10800,
			expectedFormat: "3 hours (10800 seconds)",
		},
		{
			name:           "2 hours",
			maxSessionTTL:  7200,
			expectedFormat: "2 hours (7200 seconds)",
		},
		{
			name:           "1 hour",
			maxSessionTTL:  3600,
			expectedFormat: "1 hours (3600 seconds)",
		},
		{
			name:          "zero should not display",
			maxSessionTTL: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			appCtx := &app.ApplicationContext{
				Stdout:          buf,
				Logger:          logger.NewTestLogger(),
				CompartmentName: "TestCompartment",
			}

			bastions := []domain.Bastion{
				{
					OCID:             "ocid1.bastion.oc1..test",
					DisplayName:      "test-bastion",
					BastionType:      "STANDARD",
					LifecycleState:   "ACTIVE",
					TargetVcnName:    "vcn",
					TargetSubnetName: "subnet",
					MaxSessionTTL:    tt.maxSessionTTL,
				},
			}

			err := PrintBastionInfo(bastions, appCtx, false)
			assert.NoError(t, err)

			output := buf.String()
			if tt.expectedFormat != "" {
				assert.Contains(t, output, tt.expectedFormat)
			} else {
				assert.NotContains(t, output, "MaxSessionTTL")
			}
		})
	}
}
