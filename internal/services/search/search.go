package search

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
	bleveQuery "github.com/blevesearch/bleve/v2/search/query"
	// Register analyzers used in field mappings
	_ "github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
	_ "github.com/blevesearch/bleve/v2/analysis/analyzer/simple"
)

// Indexable represents an object that can be converted into a searchable document.
type Indexable interface {
	ToIndexable() map[string]any
}

// NewIndexMapping creates a new Bleve index mapping with analyzers.
func NewIndexMapping(fields []string) mapping.IndexMapping {
	m := mapping.NewIndexMapping()

	std := mapping.NewTextFieldMapping()
	raw := mapping.NewTextFieldMapping()
	raw.Analyzer = "keyword"
	ng := mapping.NewTextFieldMapping()
	// use simple analyzer (unicode + lowercase)
	ng.Analyzer = "simple"

	doc := mapping.NewDocumentMapping()

	for _, f := range fields {
		doc.AddFieldMappingsAt(f, std)
		doc.AddFieldMappingsAt(f+".raw", raw)
		doc.AddFieldMappingsAt(f+".ng", ng)
	}

	m.DefaultMapping = doc
	return m
}

// BuildIndex builds an in-memory Bleve index for a slice of Indexable items.
func BuildIndex[T Indexable](items []T, indexMapping mapping.IndexMapping) (bleve.Index, error) {
	idx, err := bleve.NewMemOnly(indexMapping)
	if err != nil {
		return nil, fmt.Errorf("creating index: %w", err)
	}

	for i, item := range items {
		doc := item.ToIndexable()
		dup := func(k string) {
			if v, ok := doc[k].(string); ok && v != "" {
				doc[k+".raw"] = v
				doc[k+".ng"] = v
			}
		}
		for k := range doc {
			dup(k)
		}

		if err := idx.Index(strconv.Itoa(i), doc); err != nil {
			return nil, fmt.Errorf("indexing %d: %w", i, err)
		}
	}
	return idx, nil
}

// FuzzySearch performs a fuzzy search on the given index.
func FuzzySearch(index bleve.Index, pattern string, fields, boostedFields []string) ([]int, error) {
	pattern = strings.ToLower(strings.TrimSpace(pattern))
	if pattern == "" {
		return nil, nil
	}

	looksSpecific := func(p string) bool {
		if len(p) >= 15 {
			return true
		}
		if strings.ContainsAny(p, ".:-_/[]@") {
			return true
		}
		if strings.Count(p, ".") == 3 {
			return true
		}
		return false
	}

	collect := func(q bleveQuery.Query, size int) ([]int, error) {
		req := bleve.NewSearchRequestOptions(q, size, 0, false)
		res, err := index.Search(req)
		if err != nil {
			return nil, err
		}
		out := make([]int, 0, len(res.Hits))
		for _, h := range res.Hits {
			if n, err := strconv.Atoi(h.ID); err == nil {
				out = append(out, n)
			}
		}
		return out, nil
	}

	if looksSpecific(pattern) {
		var eqQs []bleveQuery.Query
		for _, f := range fields {
			tq := bleve.NewTermQuery(pattern)
			tq.SetField(f + ".raw")
			eqQs = append(eqQs, tq)
		}
		if hits, err := collect(bleve.NewDisjunctionQuery(eqQs...), 200); err != nil {
			return nil, err
		} else if len(hits) > 0 {
			return hits, nil
		}

		var subQs []bleveQuery.Query
		for _, f := range fields {
			wq := bleve.NewWildcardQuery("*" + pattern + "*")
			wq.SetField(f + ".raw")
			wq.SetBoost(1.0)
			subQs = append(subQs, wq)
		}
		if hits, err := collect(bleve.NewDisjunctionQuery(subQs...), 500); err != nil {
			return nil, err
		} else if len(hits) > 0 {
			return hits, nil
		}
	}

	var qs []bleveQuery.Query
	for _, f := range fields {
		fq := bleve.NewFuzzyQuery(pattern)
		fq.SetField(f)
		fq.SetFuzziness(2)
		fq.SetBoost(1.2)
		pq := bleve.NewPrefixQuery(pattern)
		pq.SetField(f)
		pq.SetBoost(1.3)
		qs = append(qs, fq, pq)

		mq := bleve.NewMatchQuery(pattern)
		mq.SetField(f + ".ng")
		mq.SetBoost(1.5)
		qs = append(qs, mq)

		wq := bleve.NewWildcardQuery("*" + pattern + "*")
		wq.SetField(f + ".raw")
		wq.SetBoost(1.1)
		qs = append(qs, wq)
	}

	for _, boostField := range boostedFields {
		bq := bleve.NewMatchQuery(pattern)
		bq.SetField(boostField + ".ng")
		bq.SetBoost(1.8)
		qs = append(qs, bq)
	}

	return collect(bleve.NewDisjunctionQuery(qs...), 1000)
}
