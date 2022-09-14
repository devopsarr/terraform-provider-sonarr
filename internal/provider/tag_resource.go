package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr"
	"golift.io/starr/sonarr"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TagResource{}
var _ resource.ResourceWithImportState = &TagResource{}

func NewTagResource() resource.Resource {
	return &TagResource{}
}

// TagResource defines the tag implementation.
type TagResource struct {
	client *sonarr.Sonarr
}

// Tag describes the tag data model.
type Tag struct {
	ID    types.Int64  `tfsdk:"id"`
	Label types.String `tfsdk:"label"`
}

func (r *TagResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tag"
}

func (r *TagResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Tag resource.<br/>For more information refer to [Tags](https://wiki.servarr.com/sonarr/settings#tags) documentation.",
		Attributes: map[string]tfsdk.Attribute{
			"label": {
				MarkdownDescription: "Tag label. It must be lowercase.",
				Required:            true,
				Type:                types.StringType,
				Validators: []tfsdk.AttributeValidator{
					helpers.StringLowercase(),
				},
			},
			"id": {
				MarkdownDescription: "Tag ID.",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
		},
	}, nil
}

func (r *TagResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sonarr.Sonarr)
	if !ok {
		resp.Diagnostics.AddError(
			UnexpectedResourceConfigureType,
			fmt.Sprintf("Expected *sonarr.Sonarr, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *TagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan Tag

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Tag
	request := starr.Tag{
		Label: plan.Label.Value,
	}

	response, err := r.client.AddTagContext(ctx, &request)
	if err != nil {
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to create tag, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "created tag: "+strconv.Itoa(response.ID))
	// Generate resource state struct
	result := writeTag(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func (r *TagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state Tag

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get tag current value
	response, err := r.client.GetTagContext(ctx, int(state.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to read tags, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "read tag: "+strconv.Itoa(response.ID))
	// Map response body to resource schema attribute
	result := writeTag(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func (r *TagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var plan Tag

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update Tag
	request := starr.Tag{
		Label: plan.Label.Value,
		ID:    int(plan.ID.Value),
	}

	response, err := r.client.UpdateTagContext(ctx, &request)
	if err != nil {
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to update tag, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "updated tag: "+strconv.Itoa(response.ID))
	// Generate resource state struct
	result := writeTag(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func (r *TagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Tag

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete tag current value
	err := r.client.DeleteTagContext(ctx, int(state.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to read tags, got error: %s", err))

		return
	}

	tflog.Trace(ctx, "deleted tag: "+strconv.Itoa(int(state.ID.Value)))
	resp.State.RemoveResource(ctx)
}

func (r *TagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported tag: "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func writeTag(tag *starr.Tag) *Tag {
	return &Tag{
		ID:    types.Int64{Value: int64(tag.ID)},
		Label: types.String{Value: tag.Label},
	}
}
