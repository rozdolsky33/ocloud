package objectstorage

import (
	"fmt"
	"strings"

	domain "github.com/rozdolsky33/ocloud/internal/domain/storage/objectstorage"
	"github.com/rozdolsky33/ocloud/internal/services/util"
	"github.com/rozdolsky33/ocloud/internal/tui"
)

// NewBucketListModel builds a TUI list for Buckets.
func NewBucketListModel(b []domain.Bucket) tui.Model {
	return tui.NewModel("Buckets", b, func(b domain.Bucket) tui.ResourceItemData {
		return tui.ResourceItemData{
			ID:          b.OCID,
			Title:       b.Name,
			Description: describeBucket(b),
		}
	})
}

func describeBucket(b domain.Bucket) string {
	size := ""
	if b.ApproximateSize > 0 {
		size = util.HumanizeBytesIEC(b.ApproximateSize)
	}
	count := ""
	if b.ApproximateCount > 0 {
		count = humanCount(b.ApproximateCount) + " objs"
	}
	sizePart := strings.TrimSpace(strings.Join(filterNonEmpty([]string{size, count}), " / "))

	line1 := join(" • ",
		firstNonEmpty(b.Visibility, "Private"),
		firstNonEmpty(b.StorageTier, "Standard"),
		sizePart,
	)

	// Line 2: Protections • Enc • Created
	prot := fmt.Sprintf("Ver:%s Rep:%s RO:%s",
		onOff(b.Versioning == "Enabled"),
		onOff(b.ReplicationEnabled),
		onOff(b.IsReadOnly),
	)

	enc := ""
	if b.Encryption != "" {
		switch strings.ToLower(b.Encryption) {
		case "kms", "customer-managed", "cmk":
			enc = "Enc:KMS"
		case "oracle", "oracle-managed":
			enc = "Enc:Oracle"
		default:
			enc = "Enc:" + b.Encryption
		}
	}

	created := ""
	if !b.TimeCreated.IsZero() {
		created = b.TimeCreated.Format("2006-01-02")
	}

	line2 := join(" • ",
		prot, enc, created,
	)

	return join(" • ", line1, line2)
}

// ---- helpers ----

func onOff(b bool) string {
	if b {
		return "on"
	}
	return "off"
}

func humanCount(n int) string {
	// 0–999 as-is; then K/M/B with one decimal (trim .0)
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}
	val := float64(n)
	suffix := []string{"", "K", "M", "B", "T"}
	i := 0
	for val >= 1000 && i < len(suffix)-1 {
		val /= 1000.0
		i++
	}
	s := fmt.Sprintf("%.1f", val)
	s = strings.TrimSuffix(s, ".0")
	return s + suffix[i]
}

func filterNonEmpty(parts []string) []string {
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if strings.TrimSpace(p) != "" {
			out = append(out, p)
		}
	}
	return out
}

func join(sep string, parts ...string) string {
	return strings.Join(filterNonEmpty(parts), sep)
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
