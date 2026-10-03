package helpers

import (
	"errors"
	"net/http"
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

func TestIsNotFound(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		response *http.Response
		expected bool
	}{
		"not found": {
			response: &http.Response{StatusCode: http.StatusNotFound},
			expected: true,
		},
		"unauthorized": {
			response: &http.Response{StatusCode: http.StatusUnauthorized},
			expected: false,
		},
		"ok": {
			response: &http.Response{StatusCode: http.StatusOK},
			expected: false,
		},
		"nil": {
			response: nil,
			expected: false,
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expected, IsNotFound(test.response))
		})
	}
}
