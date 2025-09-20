package oke

import (
	"fmt"
	"strings"

	domain "github.com/rozdolsky33/ocloud/internal/domain/compute"
	"github.com/rozdolsky33/ocloud/internal/tui/listx"
)

// NewImageListModel builds a TUI list for images.
func NewImageListModel(cluster []domain.Cluster) listx.Model {
	return listx.NewModel("Oracle Kubernetes Engine", cluster, func(c domain.Cluster) listx.ResourceItemData {
		return listx.ResourceItemData{
			ID:          c.OCID,
			Title:       c.DisplayName,
			Description: description(c),
		}
	})
}

func description(c domain.Cluster) string {
	parts := make([]string, 0, 4)

	if c.State != "" {
		parts = append(parts, c.State)
	}
	if v := strings.TrimSpace(c.KubernetesVersion); v != "" {
		parts = append(parts, v)
	}

	np := len(c.NodePools)
	parts = append(parts, fmt.Sprintf("%d node pool%s", np, plural(np)))

	if !c.TimeCreated.IsZero() {
		parts = append(parts, c.TimeCreated.Format("2006-01-02"))
	}

	return strings.Join(parts, " • ")
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}
