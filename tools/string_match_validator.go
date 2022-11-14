package tools

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type StringMatchValidator struct {
	Slice []string
}

// Description returns a plain text description of the validator's behavior, suitable for a practitioner to understand its impact.
func (v StringMatchValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("string must be one of %#v", v.Slice)
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior, suitable for a practitioner to understand its impact.
func (v StringMatchValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("string must be one of %#v", v.Slice)
}

// Validate runs the main validation logic of the validator, reading configuration data out of `req` and updating `resp` with diagnostics.
func (v StringMatchValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	var str types.String
	diags := tfsdk.ValueAs(ctx, req.AttributeConfig, &str)
	resp.Diagnostics.Append(diags...)

	if diags.HasError() {
		return
	}

	if str.IsUnknown() || str.IsNull() {
		return
	}

	for _, s := range v.Slice {
		if str.ValueString() == s {
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
func StringMatch(match []string) StringMatchValidator {
	return StringMatchValidator{
		Slice: match,
	}
}
