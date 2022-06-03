package helpers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type stringMatchValidator struct {
	Slice []string
}

// Description returns a plain text description of the validator's behavior, suitable for a practitioner to understand its impact.
func (v stringMatchValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("string must be one of %#v", v.Slice)
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior, suitable for a practitioner to understand its impact.
func (v stringMatchValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("string must be one of %#v", v.Slice)
}

// Validate runs the main validation logic of the validator, reading configuration data out of `req` and updating `resp` with diagnostics.
func (v stringMatchValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	var str types.String
	diags := tfsdk.ValueAs(ctx, req.AttributeConfig, &str)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
	if str.Unknown || str.Null {
		return
	}

	for _, s := range v.Slice {
		if str.Value == s {
			return
		}
	}

	resp.Diagnostics.AddAttributeError(
		req.AttributePath,
		"Invalid String Content",
		fmt.Sprintf("string must be one of %#v", v.Slice),
	)
}

// StringMatch check that a string is contained in a given slice.
func StringMatch(match []string) stringMatchValidator {
	return stringMatchValidator{
		Slice: match,
	}
}
