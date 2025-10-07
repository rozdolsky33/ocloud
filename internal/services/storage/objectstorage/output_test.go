package objectstorage

import (
	"bytes"
	"testing"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestPrintBucketsInfo_TableAndJSON(t *testing.T) {
	buf := &bytes.Buffer{}
	appCtx := &app.ApplicationContext{Logger: logger.NewTestLogger(), Stdout: buf}

	b1 := makeBucket(0)
	b2 := makeBucket(1)

	// Table output
	err := PrintBucketsInfo([]Bucket{b1, b2}, appCtx, nil, false)
	assert.NoError(t, err)
	out := buf.String()
	assert.Contains(t, out, b1.Name)
	assert.Contains(t, out, b2.Name)
	buf.Reset()

	// JSON output
	err = PrintBucketsInfo([]Bucket{b1, b2}, appCtx, nil, true)
	assert.NoError(t, err)
	jsonOut := buf.String()
	if assert.NotEmpty(t, jsonOut) {
		first := jsonOut[0]
		if first != '{' && first != '[' {
			t.Fatalf("unexpected JSON start: %q", first)
		}
	}
}

func TestPrintBucketInfo_TableAndJSON(t *testing.T) {
	buf := &bytes.Buffer{}
	appCtx := &app.ApplicationContext{Logger: logger.NewTestLogger(), Stdout: buf}
	b := makeBucket(0)

	// Table output
	err := PrintBucketInfo(&b, appCtx, false)
	assert.NoError(t, err)
	out := buf.String()
	assert.Contains(t, out, b.Name)
	buf.Reset()

	// JSON output
	err = PrintBucketInfo(&b, appCtx, true)
	assert.NoError(t, err)
	jsonOut := buf.String()
	if assert.NotEmpty(t, jsonOut) {
		first := jsonOut[0]
		if first != '{' && first != '[' {
			t.Fatalf("unexpected JSON start: %q", first)
		}
	}
}
