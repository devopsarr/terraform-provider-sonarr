package helpers

import (
	"context"
	"testing"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/stretchr/testify/assert"
)

func TestDataSourceConfigure(t *testing.T) {
	t.Parallel()

	var diags diag.Diagnostics

	diags.AddError("Unexpected DataSource Configure Type", "Expected *sonarr.APIClient, got: string. Please report this issue to the provider developers.")

	tests := map[string]struct {
		expected    any
		errorString diag.Diagnostics
	}{
		"working": {
			expected: sonarr.NewAPIClient(sonarr.NewConfiguration()),
		},
		"nil": {
			expected: (*sonarr.APIClient)(nil),
		},
		"error": {
			expected:    "abc",
			errorString: diags,
		},
	}
	for name, test := range tests {
		test := test
		req := datasource.ConfigureRequest{ProviderData: test.expected}
		resp := datasource.ConfigureResponse{}

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := DataSourceConfigure(context.TODO(), req, &resp)
			if !resp.Diagnostics.HasError() {
				assert.Equal(t, test.expected, client)
			}
			assert.Equal(t, test.errorString, resp.Diagnostics)
		})
	}
}

func TestResourceConfigure(t *testing.T) {
	t.Parallel()

	var diags diag.Diagnostics

	diags.AddError("Unexpected Resource Configure Type", "Expected *sonarr.APIClient, got: string. Please report this issue to the provider developers.")

	tests := map[string]struct {
		expected    any
		errorString diag.Diagnostics
	}{
		"working": {
			expected: sonarr.NewAPIClient(sonarr.NewConfiguration()),
		},
		"nil": {
			expected: (*sonarr.APIClient)(nil),
		},
		"error": {
			expected:    "abc",
			errorString: diags,
		},
	}
	for name, test := range tests {
		test := test
		req := resource.ConfigureRequest{ProviderData: test.expected}
		resp := resource.ConfigureResponse{}

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := ResourceConfigure(context.TODO(), req, &resp)
			if !resp.Diagnostics.HasError() {
				assert.Equal(t, test.expected, client)
			}
			assert.Equal(t, test.errorString, resp.Diagnostics)
		})
	}
}
