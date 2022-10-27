package helpers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type IntMatchValidator struct {
	Slice []int64
}

// Description returns a plain text description of the validator's behavior, suitable for a practitioner to understand its impact.
func (v IntMatchValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("number must be one of %#v", v.Slice)
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior, suitable for a practitioner to understand its impact.
func (v IntMatchValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("number must be one of %#v", v.Slice)
}

// Validate runs the main validation logic of the validator, reading configuration data out of `req` and updating `resp` with diagnostics.
func (v IntMatchValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	var str types.Int64
	diags := tfsdk.ValueAs(ctx, req.AttributeConfig, &str)
	resp.Diagnostics.Append(diags...)

	if diags.HasError() {
		return
	}

	if str.IsUnknown() || str.IsNull() {
		return
	}

	for _, s := range v.Slice {
		if str.ValueInt64() == s {
			return
		}
	}

	resp.Diagnostics.AddAttributeError(
		req.AttributePath,
		"Invalid Int64 Content",
		fmt.Sprintf("number must be one of %#v", v.Slice),
	)
}

// IntMatch check that a string is contained in a given slice.
func IntMatch(match []int64) IntMatchValidator {
	return IntMatchValidator{
		Slice: match,
	}
}
