package helpers

import (
	"errors"
	"testing"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/stretchr/testify/assert"
)

func TestParseClientError(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		action   string
		name     string
		err      error
		expected string
	}{
		"tag_create": {
			action:   "create",
			name:     "sonarr_tag",
			err:      &sonarr.GenericOpenAPIError{},
			expected: "Unable to create sonarr_tag, got error: \nDetails:\n",
		},
		"generic": {
			action:   "create",
			name:     "sonarr_tag",
			err:      errors.New("other error"),
			expected: "Unable to create sonarr_tag, got error: other error",
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expected, ParseClientError(test.action, test.name, test.err))
		})
	}
}

func TestParseNotFoundError(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		kind     string
		field    string
		search   string
		expected string
	}{
		"generic": {
			kind:     "sonarr_tag",
			field:    "label",
			search:   "test",
			expected: "Unable to find sonarr_tag, got error: data source not found: no sonarr_tag with label 'test'",
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expected, ParseNotFoundError(test.kind, test.field, test.search))
		})
	}
}

func TestWrongClient(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		wanted   string
		received interface{}
		expected string
	}{
		"generic": {
			expected: "Expected string, got: int. Please report this issue to the provider developers.",
			wanted:   "string",
			received: 3,
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expected, WrongClient(test.wanted, test.received))
		})
	}
}
