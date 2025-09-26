package gateway

import "testing"

// This is a minimal smoke test to ensure the adapter type is constructible in tests.
func TestNewAdapter_Constructible(t *testing.T) {
	var client interface{} // we don't call any methods; just ensure type can be instantiated
	_ = &Adapter{}         // zero-value adapter (no client) is a valid Go value
	if client == nil {     // dummy assertion to avoid unused warnings in some linters
		// pass
	}
}
