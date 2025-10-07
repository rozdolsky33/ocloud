package vcn

import (
	"bytes"
	"testing"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestPrintVCNInfo_JSONAndTable(t *testing.T) {
	buf := &bytes.Buffer{}
	appCtx := &app.ApplicationContext{Logger: logger.NewTestLogger(), Stdout: buf}

	v := makeVCN(0)

	// Table output
	err := PrintVCNInfo(v, appCtx, false, false, false, false, false, false)
	assert.NoError(t, err)
	out := buf.String()
	assert.Contains(t, out, v.DisplayName)
	assert.Contains(t, out, "OCID")
	buf.Reset()

	// JSON output
	err = PrintVCNInfo(v, appCtx, true, true, true, true, true, true)
	assert.NoError(t, err)
	jsonOut := buf.String()
	if assert.NotEmpty(t, jsonOut) {
		first := jsonOut[0]
		if first != '{' && first != '[' {
			t.Fatalf("unexpected JSON start: %q", first)
		}
	}
}

func TestPrintVCNsInfo_List_JSONAndTable(t *testing.T) {
	buf := &bytes.Buffer{}
	appCtx := &app.ApplicationContext{Logger: logger.NewTestLogger(), Stdout: buf}

	vcns := []VCN{makeVCN(0), makeVCN(1)}

	// Table
	err := PrintVCNsInfo(vcns, appCtx, nil, false, false, false, false, false, false)
	assert.NoError(t, err)
	out := buf.String()
	assert.Contains(t, out, vcns[0].DisplayName)
	assert.Contains(t, out, vcns[1].DisplayName)
	buf.Reset()

	// JSON
	err = PrintVCNsInfo(vcns, appCtx, nil, true, true, true, true, true, true)
	assert.NoError(t, err)
	jsonOut := buf.String()
	if assert.NotEmpty(t, jsonOut) {
		first := jsonOut[0]
		if first != '{' && first != '[' {
			t.Fatalf("unexpected JSON start: %q", first)
		}
	}
}
