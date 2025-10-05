package util

import "fmt"

// HumanizeBytesIEC converts a byte size into a human‑readable string using IEC units (powers of 1024).
// Examples:
//   - 0 -> "0 B"
//   - 1023 -> "1023 B"
//   - 1024 -> "1.00 KiB"
//   - 2.72 * 1024 * 1024 -> "2.72 MiB"
func HumanizeBytesIEC(b int64) string {
	const (
		_        = iota
		KB int64 = 1 << (10 * iota) // 1024
		MB
		GB
		TB
		PB
	)

	switch {
	case b < 1024:
		return fmt.Sprintf("%d B", b)
	case b < MB:
		return fmt.Sprintf("%.2f KiB", float64(b)/float64(KB))
	case b < GB:
		return fmt.Sprintf("%.2f MiB", float64(b)/float64(MB))
	case b < TB:
		return fmt.Sprintf("%.2f GiB", float64(b)/float64(GB))
	case b < PB:
		return fmt.Sprintf("%.2f TiB", float64(b)/float64(TB))
	default:
		return fmt.Sprintf("%.2f PiB", float64(b)/float64(PB))
	}
}
