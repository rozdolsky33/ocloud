package oke

import (
	"bytes"
	"testing"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"

	"github.com/stretchr/testify/assert"
)

func TestPrintOKETable(t *testing.T) {
	clusters := []Cluster{
		{
			DisplayName:       "TestCluster1",
			KubernetesVersion: "v1.25.4",
			State:             "ACTIVE",
		},
	}

	var buf bytes.Buffer
	appCtx := &app.ApplicationContext{
		Logger: logger.NewTestLogger(),
		Stdout: &buf,
	}

	err := PrintOKETable(clusters, appCtx, nil, false)
	assert.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "TestCluster1")
}
