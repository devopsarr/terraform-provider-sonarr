package tools

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type StringLowercaseValidator struct{}

// Description returns a plain text description of the validator's behavior, suitable for a practitioner to understand its impact.
func (v StringLowercaseValidator) Description(ctx context.Context) string {
	return "string must be lowercase"
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior, suitable for a practitioner to understand its impact.
func (v StringLowercaseValidator) MarkdownDescription(ctx context.Context) string {
	return "string must be lowercase"
}

// Validate runs the main validation logic of the validator, reading configuration data out of `req` and updating `resp` with diagnostics.
func (v StringLowercaseValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	var str types.String
	diags := tfsdk.ValueAs(ctx, req.AttributeConfig, &str)
	resp.Diagnostics.Append(diags...)

	if diags.HasError() {
		return
	}

	if str.IsUnknown() || str.IsNull() {
		return
	}

	upper, _ := regexp.Match(`^.*[A-Z]+.*$`, []byte(str.ValueString()))
	if upper {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid String Content",
			"String cannot contains uppercase values",
		)

		return
	}
}

// StringLowercase checks that a string is lowercase.
func StringLowercase() StringLowercaseValidator {
	return StringLowercaseValidator{}
}
