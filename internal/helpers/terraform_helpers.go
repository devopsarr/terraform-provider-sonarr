package helpers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ImportStatePassthroughIntID is a helper function to set the import
// identifier to a given state attribute path. The attribute must accept a
// int value.
// extends https://github.com/hashicorp/terraform-plugin-framework/blob/main/resource/import_state.go.
func ImportStatePassthroughIntID(ctx context.Context, attrPath path.Path, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %s", req.ID),
		)
	}

	if attrPath.Equal(path.Empty()) {
		resp.Diagnostics.AddError(
			"Resource Import Passthrough Missing Attribute Path",
			"This is always an error in the provider. Please report the following to the provider developer:\n\n"+
				"Resource ImportState method call to ImportStatePassthroughIntID path must be set to a valid attribute path that can accept a int value.",
		)
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, attrPath, id)...)
}
