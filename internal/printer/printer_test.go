package printer

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	// Create a buffer to capture output
	var buf bytes.Buffer

	// Create a new printer with the buffer as output
	p := New(&buf)

	// Check that the printer was created correctly
	if p == nil {
		t.Fatal("Expected New to return a non-nil Printer")
	}

	// Check that the output writer was set correctly
	if p.out == nil {
		t.Fatal("Expected Printer.out to be non-nil")
	}
}

func TestMarshalToJSON(t *testing.T) {
	// Create a buffer to capture output
	var buf bytes.Buffer

	// Create a new printer with the buffer as output
	p := New(&buf)

	// Test data
	testData := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	// Marshal the data to JSON
	err := p.MarshalToJSON(testData)

	// Check that there was no error
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check that the output is valid JSON
	var result map[string]string
	err = json.Unmarshal(buf.Bytes(), &result)
	if err != nil {
		t.Fatalf("Expected valid JSON, got error: %v", err)
	}

	// Check that the output contains the expected data
	if result["key1"] != "value1" || result["key2"] != "value2" {
		t.Fatalf("Expected output to contain the test data, got %v", result)
	}
}

func TestPrintKeyValues(t *testing.T) {
	// Create a buffer to capture output
	var buf bytes.Buffer

	// Create a new printer with the buffer as output
	p := New(&buf)

	// Test data
	title := "Test Title"
	data := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	keys := []string{"key1", "key2"}

	// Print the key-values
	p.PrintKeyValues(title, data, keys)

	// Check that the output contains the title
	if !strings.Contains(buf.String(), title) {
		t.Fatalf("Expected output to contain the title '%s', got: %s", title, buf.String())
	}

	// Check that the output contains the keys and values
	for _, key := range keys {
		if !strings.Contains(buf.String(), key) {
			t.Fatalf("Expected output to contain the key '%s', got: %s", key, buf.String())
		}
		if !strings.Contains(buf.String(), data[key]) {
			t.Fatalf("Expected output to contain the value '%s', got: %s", data[key], buf.String())
		}
	}
}

func TestPrintKeyValues_EmptyData(t *testing.T) {
	// Create a buffer to capture output
	var buf bytes.Buffer

	// Create a new printer with the buffer as output
	p := New(&buf)

	// Test with empty data
	title := "Test Title"
	data := map[string]string{}
	keys := []string{}

	// Print the key-values
	p.PrintKeyValues(title, data, keys)

	// Check that the output still contains the title
	if !strings.Contains(buf.String(), title) {
		t.Fatalf("Expected output to contain the title '%s', got: %s", title, buf.String())
	}
}

func TestPrintKeyValues_MissingKey(t *testing.T) {
	// Create a buffer to capture output
	var buf bytes.Buffer

	// Create a new printer with the buffer as output
	p := New(&buf)

	// Test with a key that doesn't exist in the data
	title := "Test Title"
	data := map[string]string{
		"key1": "value1",
	}
	keys := []string{"key1", "key2"}

	// Print the key-values
	p.PrintKeyValues(title, data, keys)

	// Check that the output contains the existing key and value
	if !strings.Contains(buf.String(), "key1") {
		t.Fatalf("Expected output to contain the key 'key1', got: %s", buf.String())
	}
	if !strings.Contains(buf.String(), "value1") {
		t.Fatalf("Expected output to contain the value 'value1', got: %s", buf.String())
	}

	// Check that the output doesn't contain the missing key
	if strings.Contains(buf.String(), "key2") {
		t.Fatalf("Expected output not to contain the key 'key2', got: %s", buf.String())
	}
}
