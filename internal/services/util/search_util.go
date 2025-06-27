package util

import (
	"fmt"
	"github.com/blevesearch/bleve/v2"
	bleveQuery "github.com/blevesearch/bleve/v2/search/query"
	"strconv"
	"strings"
)

// BuildIndex creates an in-memory Bleve index from a slice of items using a provided mapping function.
// The mapping function converts each item into an indexable structure for searching.
// Returns the built index or an error if any indexing operation fails.
func BuildIndex[T any](items []T, mapToIndexable func(T) any) (bleve.Index, error) {
	indexMapping := bleve.NewIndexMapping()
	index, err := bleve.NewMemOnly(indexMapping)
	if err != nil {
		return nil, fmt.Errorf("creating index: %w", err)
	}

	for i, item := range items {
		err := index.Index(fmt.Sprintf("%d", i), mapToIndexable(item))
		if err != nil {
			return nil, fmt.Errorf("indexing failed at %d: %w", i, err)
		}
	}
	return index, nil
}

// FuzzySearchIndex performs a fuzzy search on a Bleve index for a given pattern across specified fields.
// It combines fuzzy, prefix, and wildcard queries, limits the results, and returns matched indices or an error.
func FuzzySearchIndex(index bleve.Index, pattern string, fields []string) ([]int, error) {
	var limit = 1000
	var queries []bleveQuery.Query

	for _, field := range fields {
		// Fuzzy match (Levenshtein distance)
		fq := bleve.NewFuzzyQuery(pattern)
		fq.SetField(field)
		fq.SetFuzziness(2)
		queries = append(queries, fq)

		// Prefix match (useful for dev1, splunkdev1, etc.)
		pq := bleve.NewPrefixQuery(pattern)
		pq.SetField(field)
		queries = append(queries, pq)

		// Wildcard match (matches anywhere in token)
		wq := bleve.NewWildcardQuery("*" + pattern + "*")
		wq.SetField(field)
		queries = append(queries, wq)
	}

	// OR all queries together
	combinedQuery := bleve.NewDisjunctionQuery(queries...)

	search := bleve.NewSearchRequestOptions(combinedQuery, limit, 0, false)

	result, err := index.Search(search)
	if err != nil {
		return nil, err
	}

	var hits []int
	for _, hit := range result.Hits {
		idx, err := strconv.Atoi(hit.ID)
		if err != nil {
			continue
		}
		hits = append(hits, idx)
	}

	return hits, nil
}

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
		return "", nil // Return an empty string without error when no tags are found
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
		return "", nil // Return an empty string without error when no tag values are found
	}

	return strings.Join(values, " "), nil
}
