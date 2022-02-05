package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Tag -
type Tag struct {
	ID    types.Int64  `tfsdk:"id"`
	Label types.String `tfsdk:"label"`
}

//TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
// Tags -
type Tags struct {
	ID   types.String `tfsdk:"id"`
	Tags []Tag        `tfsdk:"tags"`
}
