package helpers

import (
	"errors"
	"testing"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/stretchr/testify/assert"
)

func TestErrDataNotFoundError(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		kind, field, search, expected string
	}{
		"tag": {"sonarr_tag", "label", "test", "data source not found: no sonarr_tag with label 'test'"},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expected, ErrDataNotFoundError(test.kind, test.field, test.search).Error())
		})
	}
}

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
			name:     "radarr_tag",
			err:      errors.New("other error"),
			expected: "Unable to create radarr_tag, got error: other error",
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
