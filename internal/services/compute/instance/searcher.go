package instance

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
	bleveQuery "github.com/blevesearch/bleve/v2/search/query"
	"github.com/rozdolsky33/ocloud/internal/domain/compute"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// newInstanceIndexMapping creates a new Bleve index mapping for instances.
func newInstanceIndexMapping() mapping.IndexMapping {
	m := mapping.NewIndexMapping()

	// ngram filter/analyzer for "Google-ish" substring search
	_ = m.AddCustomTokenFilter("ng2_20", map[string]any{
		"type": "ngram", "min": 2, "max": 20,
	})
	_ = m.AddCustomAnalyzer("ngram_lower", map[string]any{
		"type":          "custom",
		"tokenizer":     "unicode",
		"token_filters": []string{"to_lower", "ng2_20"},
	})

	// keyword analyzer (keeps the whole value incl punctuation)
	_ = m.AddCustomAnalyzer("keyword_lower", map[string]any{
		"type":          "custom",
		"tokenizer":     "keyword",
		"token_filters": []string{"to_lower"},
	})

	std := mapping.NewTextFieldMapping() // default wordy search
	raw := mapping.NewTextFieldMapping() // raw/punctuation-preserving
	raw.Analyzer = "keyword_lower"
	ng := mapping.NewTextFieldMapping() // ngram-substring search
	ng.Analyzer = "ngram_lower"

	doc := mapping.NewDocumentMapping()

	// Core fields (tokenized + raw + ngram for the ones you care about)
	for _, f := range []string{
		"Name", "Hostname", "ImageName", "ImageOS", "Shape",
		"PrimaryIP", "OCID", "VcnName", "SubnetName",
		// tags below (logical fields)
		"TagsKV", "TagsVal",
	} {
		doc.AddFieldMappingsAt(f, std)
		doc.AddFieldMappingsAt(f+".raw", raw)
		doc.AddFieldMappingsAt(f+".ng", ng)
	}

	m.DefaultMapping = doc
	return m
}

// BuildIndex builds an in-memory Bleve index for instances.
func BuildIndex(instances []compute.Instance) (bleve.Index, error) {
	idx, err := bleve.NewMemOnly(newInstanceIndexMapping())
	if err != nil {
		return nil, fmt.Errorf("creating index: %w", err)
	}
	for i, inst := range instances {
		doc := mapToIndexableInstance(inst) // see below
		// duplicate into raw/ng variants for the fields you want substring search on
		dup := func(k string) {
			if v, ok := doc[k].(string); ok && v != "" {
				doc[k+".raw"] = v
				doc[k+".ng"] = v
			}
		}
		for _, k := range []string{
			"Name", "Hostname", "ImageName", "ImageOS", "Shape",
			"PrimaryIP", "OCID", "VcnName", "SubnetName",
			"TagsKV", "TagsVal",
		} {
			dup(k)
		}

		if err := idx.Index(strconv.Itoa(i), doc); err != nil {
			return nil, fmt.Errorf("indexing %d: %w", i, err)
		}
	}
	return idx, nil
}

// mapToIndexableInstance converts your domain.Instance into fields apt for search.
func mapToIndexableInstance(inst compute.Instance) map[string]any {
	tagsKV, _ := util.FlattenTags(inst.FreeformTags, inst.DefinedTags)
	tagsVal, _ := util.ExtractTagValues(inst.FreeformTags, inst.DefinedTags)

	return map[string]any{
		"Name":       strings.ToLower(inst.DisplayName),
		"Hostname":   strings.ToLower(inst.Hostname),
		"PrimaryIP":  strings.ToLower(inst.PrimaryIP),
		"ImageName":  strings.ToLower(inst.ImageName),
		"ImageOS":    strings.ToLower(inst.ImageOS),
		"Shape":      strings.ToLower(inst.Shape),
		"OCID":       strings.ToLower(inst.OCID),
		"VcnName":    strings.ToLower(inst.VcnName),
		"SubnetName": strings.ToLower(inst.SubnetName),

		"TagsKV":  strings.ToLower(tagsKV),
		"TagsVal": strings.ToLower(tagsVal),
	}
}

// FuzzySearchInstances searches across core fields + tags.
func FuzzySearchInstances(index bleve.Index, pattern string) ([]int, error) {
	pattern = strings.ToLower(strings.TrimSpace(pattern))
	if pattern == "" {
		return nil, nil
	}

	fields := []string{
		"Name", "Hostname", "PrimaryIP", "ImageName", "ImageOS",
		"Shape", "OCID", "VcnName", "SubnetName",
		"TagsKV", "TagsVal",
	}

	var qs []bleveQuery.Query
	for _, f := range fields {
		// Tokenized fuzzy/prefix
		fq := bleve.NewFuzzyQuery(pattern)
		fq.SetField(f)
		fq.SetFuzziness(2)
		fq.SetBoost(1.2)
		pq := bleve.NewPrefixQuery(pattern)
		pq.SetField(f)
		pq.SetBoost(1.3)
		qs = append(qs, fq, pq)

		// N-gram (the best “Google-like” substring; no wildcard cost)
		mq := bleve.NewMatchQuery(pattern)
		mq.SetField(f + ".ng")
		mq.SetBoost(1.5)
		qs = append(qs, mq)

		// Raw substring (covers punctuation like app-test, foo/bar, ns.key:value)
		wq := bleve.NewWildcardQuery("*" + pattern + "*")
		wq.SetField(f + ".raw")
		wq.SetBoost(1.1)
		qs = append(qs, wq)
	}

	// Slightly favor name/hostname
	for _, boostField := range []string{"Name", "Hostname"} {
		bq := bleve.NewMatchQuery(pattern)
		bq.SetField(boostField + ".ng")
		bq.SetBoost(1.8)
		qs = append(qs, bq)
	}

	req := bleve.NewSearchRequestOptions(bleve.NewDisjunctionQuery(qs...), 1000, 0, false)
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
