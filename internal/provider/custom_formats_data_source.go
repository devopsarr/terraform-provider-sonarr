package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const customFormatsDataSourceName = "custom_formats"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &CustomFormatsDataSource{}

func NewCustomFormatsDataSource() datasource.DataSource {
	return &CustomFormatsDataSource{}
}

// CustomFormatsDataSource defines the download clients implementation.
type CustomFormatsDataSource struct {
	client *sonarr.APIClient
}

// CustomFormats describes the download clients data model.
type CustomFormats struct {
	CustomFormats types.Set    `tfsdk:"custom_formats"`
	ID            types.String `tfsdk:"id"`
}

func (d *CustomFormatsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + customFormatsDataSourceName
}

func (d *CustomFormatsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Profiles -->List all available [Custom Formats](../resources/custom_format).",
		Attributes: map[string]schema.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.StringAttribute{
				Computed: true,
			},
			"custom_formats": schema.SetNestedAttribute{
				MarkdownDescription: "Custom Format list.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"include_custom_format_when_renaming": schema.BoolAttribute{
							MarkdownDescription: "Include custom format when renaming flag.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Custom Format name.",
							Computed:            true,
						},
						"id": schema.Int64Attribute{
							MarkdownDescription: "Custom Format ID.",
							Computed:            true,
						},
						"specifications": schema.SetNestedAttribute{
							MarkdownDescription: "Specifications.",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"negate": schema.BoolAttribute{
										MarkdownDescription: "Negate flag.",
										Computed:            true,
									},
									"required": schema.BoolAttribute{
										MarkdownDescription: "Computed flag.",
										Computed:            true,
									},
									"name": schema.StringAttribute{
										MarkdownDescription: "Specification name.",
										Computed:            true,
									},
									"implementation": schema.StringAttribute{
										MarkdownDescription: "Implementation.",
										Computed:            true,
									},
									// Field values
									"value": schema.StringAttribute{
										MarkdownDescription: "Value.",
										Computed:            true,
									},
									"min": schema.Int64Attribute{
										MarkdownDescription: "Min.",
										Computed:            true,
									},
									"max": schema.Int64Attribute{
										MarkdownDescription: "Max.",
										Computed:            true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *CustomFormatsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *CustomFormatsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *CustomFormats

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get download clients current value
	response, _, err := d.client.CustomFormatApi.ListCustomFormat(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", customFormatsDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+customFormatsDataSourceName)
	// Map response body to resource schema attribute
	profiles := make([]CustomFormat, len(response))
	for i, p := range response {
		profiles[i].write(ctx, p)
	}

	tfsdk.ValueFrom(ctx, profiles, data.CustomFormats.Type(context.Background()), &data.CustomFormats)
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	data.ID = types.StringValue(strconv.Itoa(len(response)))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
