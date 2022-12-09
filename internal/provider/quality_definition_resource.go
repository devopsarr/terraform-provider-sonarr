package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/tools"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr"
	"golift.io/starr/sonarr"
)

const qualityDefinitionResourceName = "quality_definition"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &QualityDefinitionResource{}
	_ resource.ResourceWithImportState = &QualityDefinitionResource{}
)

func NewQualityDefinitionResource() resource.Resource {
	return &QualityDefinitionResource{}
}

// QualityDefinitionResource defines the quality definition implementation.
type QualityDefinitionResource struct {
	client *sonarr.Sonarr
}

// QualityDefinition describes the quality definition data model.
type QualityDefinition struct {
	Title       types.String  `tfsdk:"title"`
	QualityName types.String  `tfsdk:"quality_name"`
	Source      types.String  `tfsdk:"source"`
	MinSize     types.Float64 `tfsdk:"min_size"`
	MaxSize     types.Float64 `tfsdk:"max_size"`
	ID          types.Int64   `tfsdk:"id"`
	QualityID   types.Int64   `tfsdk:"quality_id"`
	Resolution  types.Int64   `tfsdk:"resolution"`
}

func (r *QualityDefinitionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + qualityDefinitionResourceName
}

func (r *QualityDefinitionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Definitions -->Quality Definition resource.\nFor more information refer to [Quality Definition](https://wiki.servarr.com/sonarr/settings#quality-1) documentation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Quality Definition ID.",
				Required:            true,
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "Quality Definition Title.",
				Required:            true,
			},
			"min_size": schema.Float64Attribute{
				MarkdownDescription: "Minimum size MB/min.",
				Optional:            true,
				Computed:            true,
			},
			"max_size": schema.Float64Attribute{
				MarkdownDescription: "Maximum size MB/min.",
				Optional:            true,
				Computed:            true,
			},
			"quality_id": schema.Int64Attribute{
				MarkdownDescription: "Quality ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"resolution": schema.Int64Attribute{
				MarkdownDescription: "Quality Resolution.",
				Computed:            true,
			},
			"quality_name": schema.StringAttribute{
				MarkdownDescription: "Quality Name.",
				Computed:            true,
			},
			"source": schema.StringAttribute{
				MarkdownDescription: "Quality source.",
				Computed:            true,
			},
		},
	}
}

func (r *QualityDefinitionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sonarr.Sonarr)
	if !ok {
		resp.Diagnostics.AddError(
			tools.UnexpectedResourceConfigureType,
			fmt.Sprintf("Expected *sonarr.Sonarr, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *QualityDefinitionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var definition *QualityDefinition

	resp.Diagnostics.Append(req.Plan.Get(ctx, &definition)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Create resource
	data := definition.read()

	// Read to get the quality ID
	read, err := r.client.GetQualityDefinitionContext(ctx, data.ID)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", qualityDefinitionResourceName, err))

		return
	}

	data.Quality.ID = read.Quality.ID

	// Create new QualityDefinition
	response, err := r.client.UpdateQualityDefinitionContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", qualityDefinitionResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+qualityDefinitionResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	definition.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &definition)...)
}

func (r *QualityDefinitionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var definition *QualityDefinition

	resp.Diagnostics.Append(req.State.Get(ctx, &definition)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get qualitydefinition current value
	response, err := r.client.GetQualityDefinitionContext(ctx, definition.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", qualityDefinitionResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+qualityDefinitionResourceName+": "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	definition.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &definition)...)
}

func (r *QualityDefinitionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var definition *QualityDefinition

	resp.Diagnostics.Append(req.Plan.Get(ctx, &definition)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Update resource
	data := definition.read()

	// Update QualityDefinition
	response, err := r.client.UpdateQualityDefinitionContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", qualityDefinitionResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+qualityDefinitionResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	definition.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &definition)...)
}

func (r *QualityDefinitionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// QualityDefinitions cannot be really deleted just removing configuration
	tflog.Trace(ctx, "decoupled "+qualityDefinitionResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *QualityDefinitionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			tools.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+qualityDefinitionResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (p *QualityDefinition) write(definition *sonarr.QualityDefinition) {
	p.ID = types.Int64Value(definition.ID)
	p.MinSize = types.Float64Value(definition.MinSize)
	p.MaxSize = types.Float64Value(definition.MaxSize)
	p.Title = types.StringValue(definition.Title)
	p.QualityName = types.StringValue(definition.Quality.Name)
	p.QualityID = types.Int64Value(definition.Quality.ID)
	p.Source = types.StringValue(definition.Quality.Source)
	p.Resolution = types.Int64Value(int64(definition.Quality.Resolution))
}

func (p *QualityDefinition) read() *sonarr.QualityDefinition {
	return &sonarr.QualityDefinition{
		ID:      p.ID.ValueInt64(),
		MinSize: p.MinSize.ValueFloat64(),
		MaxSize: p.MaxSize.ValueFloat64(),
		Title:   p.Title.ValueString(),
		Quality: &starr.BaseQuality{
			ID:         p.QualityID.ValueInt64(),
			Name:       p.QualityName.ValueString(),
			Source:     p.Source.ValueString(),
			Resolution: int(p.Resolution.ValueInt64()),
		},
	}
}
