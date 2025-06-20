package util

import (
	"fmt"
	"strings"
)

// FlattenTags flattens freeform and defined tags into a single string with a specific format suitable for indexing.
// Freeform tags are processed as key:value pairs, while defined tags include namespace, key, and value.
// Returns the flattened string or an empty string if no valid tags are found.
func FlattenTags(freeform map[string]string, defined map[string]map[string]interface{}) (string, error) {
	var parts []string
	// Handle Freeform Tags
	if freeform != nil {
		for k, v := range freeform {
			if k == "" || v == "" {
				continue // skip empty keys/values
			}
			parts = append(parts, fmt.Sprintf("%s:%s", strings.ToLower(k), strings.ToLower(v)))
		}
	}
	// Handle Defined Tags
	if defined != nil {
		for ns, kv := range defined {
			if ns == "" || kv == nil {
				continue
			}
			for k, v := range kv {
				if k == "" || v == nil {
					continue
				}

				// You can restrict this to specific types if desired (e.g., string only)
				// Convert to string safely
				var valueStr string
				switch val := v.(type) {
				case string:
					valueStr = val
				default:
					// fallback to fmt.Sprintf
					valueStr = fmt.Sprintf("%v", val)
				}

				parts = append(parts, fmt.Sprintf("%s.%s:%s", strings.ToLower(ns), strings.ToLower(k), valueStr))
			}
		}
	}

	if len(parts) == 0 {
		return "", nil // Return empty string without error when no tags are found
	}

	return strings.Join(parts, " "), nil
}

// ExtractTagValues extracts only the values from freeform and defined tags into a single space-separated string.
// This is useful for making tag values directly searchable without requiring the key prefix.
// Returns the extracted values string or an empty string if no valid tag values are found.
func ExtractTagValues(freeform map[string]string, defined map[string]map[string]interface{}) (string, error) {
	var values []string

	// Extract values from Freeform Tags
	if freeform != nil {
		for _, v := range freeform {
			if v == "" {
				continue // skip empty values
			}
			values = append(values, strings.ToLower(v))
		}
	}

	// Extract values from Defined Tags
	if defined != nil {
		for _, kv := range defined {
			if kv == nil {
				continue
			}
			for _, v := range kv {
				if v == nil {
					continue
				}

				// Convert to string safely
				var valueStr string
				switch val := v.(type) {
				case string:
					valueStr = val
				default:
					// fallback to fmt.Sprintf
					valueStr = fmt.Sprintf("%v", val)
				}

				if valueStr != "" {
					values = append(values, strings.ToLower(valueStr))
				}
			}
		}
	}

	if len(values) == 0 {
		return "", nil // Return empty string without error when no tag values are found
	}

	return strings.Join(values, " "), nil
}
