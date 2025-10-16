package domain

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrNotFound(t *testing.T) {
	assert.NotNil(t, ErrNotFound)
	assert.Equal(t, "not found", ErrNotFound.Error())
	assert.True(t, errors.Is(ErrNotFound, ErrNotFound))
}

func TestNewNotFoundError(t *testing.T) {
	tests := []struct {
		name         string
		resourceType string
		resourceName string
		expected     string
	}{
		{
			name:         "complete error message",
			resourceType: "instance",
			resourceName: "ocid1.instance.oc1..test",
			expected:     "instance 'ocid1.instance.oc1..test': not found",
		},
		{
			name:         "empty resource type",
			resourceType: "",
			resourceName: "ocid1.test",
			expected:     " 'ocid1.test': not found",
		},
		{
			name:         "empty resource name",
			resourceType: "image",
			resourceName: "",
			expected:     "image '': not found",
		},
		{
			name:         "both empty",
			resourceType: "",
			resourceName: "",
			expected:     " '': not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewNotFoundError(tt.resourceType, tt.resourceName)
			assert.Error(t, err)
			assert.Equal(t, tt.expected, err.Error())
			assert.True(t, errors.Is(err, ErrNotFound))
		})
	}
}
