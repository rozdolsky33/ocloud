package instance

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/blevesearch/bleve/v2"
	_ "github.com/blevesearch/bleve/v2/analysis/analyzer/custom"
	_ "github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
	_ "github.com/blevesearch/bleve/v2/analysis/token/lowercase"
	_ "github.com/blevesearch/bleve/v2/analysis/tokenizer/unicode"
	"github.com/blevesearch/bleve/v2/mapping"
	bleveQuery "github.com/blevesearch/bleve/v2/search/query"

	"github.com/rozdolsky33/ocloud/internal/domain/compute"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// newInstanceIndexMapping creates a new Bleve index mapping for instances.
func newInstanceIndexMapping() mapping.IndexMapping {
	m := mapping.NewIndexMapping()

	_ = m.AddCustomAnalyzer("ngram_lower", map[string]any{
		"type":          "custom",
		"tokenizer":     "unicode",
		"token_filters": []string{"to_lower"},
	})

	// Use built-in keyword analyzer for raw fields; values are already lowercased at ingestion
	std := mapping.NewTextFieldMapping() // default wordy search
	raw := mapping.NewTextFieldMapping() // raw/punctuation-preserving
	raw.Analyzer = "keyword"
	ng := mapping.NewTextFieldMapping() // ngram-substring search
	ng.Analyzer = "ngram_lower"

	doc := mapping.NewDocumentMapping()

	// Core fields (tokenized + raw + ngram for the ones you care about)
	for _, f := range []string{
		"Name", "Hostname", "ImageName", "ImageOS", "Shape",
		"PrimaryIP", "OCID", "VcnName", "SubnetName",
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
		"FD":         strings.ToLower(inst.FaultDomain),
		"AD":         strings.ToLower(inst.AvailabilityDomain),
		"VcnName":    strings.ToLower(inst.VcnName),
		"SubnetName": strings.ToLower(inst.SubnetName),

		"TagsKV":  strings.ToLower(tagsKV),
		"TagsVal": strings.ToLower(tagsVal),
	}
}

// FuzzySearchInstances searches across core fields + tags.
// It applies heuristics: if the pattern looks very specific (e.g., IP, long name with punctuation),
// it performs precise matching on raw fields first (exact equals, then substring). Falls back to
// broader fuzzy/prefix/ngram search only if no precise hits are found.
func FuzzySearchInstances(index bleve.Index, pattern string) ([]int, error) {
	pattern = strings.ToLower(strings.TrimSpace(pattern))
	if pattern == "" {
		return nil, nil
	}

	fields := []string{
		"Name", "Hostname", "PrimaryIP", "ImageName", "ImageOS",
		"Shape", "OCID", "VcnName", "SubnetName", "AD", "FD",
		"TagsKV", "TagsVal",
	}

	looksSpecific := func(p string) bool {
		if len(p) >= 15 { // long tokens are likely specific
			return true
		}
		// contains punctuation typical for exact IDs/hosts/IPs
		if strings.ContainsAny(p, ".:-_/[]@") {
			return true
		}
		// naive IPv4 check
		if strings.Count(p, ".") == 3 {
			return true
		}
		return false
	}

	// helper to run a bleve search and convert to []int
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

	// 1) If specific, try exact equals on raw fields
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
			return hits, nil // return only exact matches
		}

		// 2) Next, substring on raw fields only
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

	// 3) Broad fallback: fuzzy + prefix + ngram + raw substring across all fields
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

	return collect(bleve.NewDisjunctionQuery(qs...), 1000)
}
