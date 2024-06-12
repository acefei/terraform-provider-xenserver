package xenserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"xenapi"
)

// networkResource defines the resource implementation.
type networkResource struct {
	session *xenapi.Session
}

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &networkResource{}
	_ resource.ResourceWithConfigure   = &networkResource{}
	_ resource.ResourceWithImportState = &networkResource{}
)

// This is a helper function to simplify the provider implementation.
func NewNetworkResource() resource.Resource {
	return &networkResource{}
}

// Set the resource name
func (r *networkResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network"
}

// NetworkSchema returns a map attribute schema for the network field shared by multiple resources.
func NetworkSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The UUID of the virtual network on xenserver",
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		// required
		"name_label": schema.StringAttribute{
			MarkdownDescription: "The name of the virtual network",
			Required:            true,
		},
		"name_description": schema.StringAttribute{
			MarkdownDescription: "The description of the virtual network",
			Optional:            true,
			Computed:            true, // Required to use Default
			Default:             stringdefault.StaticString(""),
		},
		"other_config": schema.MapAttribute{
			MarkdownDescription: "The additional configuration of the virtual network",
			Optional:            true,
			Computed:            true, // Required to use Default
			ElementType:         types.StringType,
			Default:             mapdefault.StaticValue(types.MapValueMust(types.StringType, map[string]attr.Value{})),
		},
	}
}

// Set the defined datamodel of the resource
func (r *networkResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Network Resource",
		Attributes:          NetworkSchema(),
	}
}

// Set the parameter of the resource, pass value from provider
func (r *networkResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
	session, ok := req.ProviderData.(*xenapi.Session)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *xenapi.Session, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.session = session
}

// Read data from Plan, create resource, get data from new source, set to State
// terraform plan/apply
func (r *networkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data networkResourceModel
	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// create new resourcek
	networkRecord := xenapi.NetworkRecord{
		NameLabel:       data.NameLabel.ValueString(),
		NameDescription: data.NameDescription.ValueString(),
		OtherConfig:     make(map[string]string, len(data.OtherConfig.Elements())),
	}

	networkRef, err := xenapi.Network.Create(r.session, networkRecord)
	if err != nil {
		// failed to create network
		resp.Diagnostics.AddError(
			"Unable to create network",
			err.Error(),
		)
		return
	}

	// Overwrite data with refreshed resource state
	record, err := xenapi.Network.GetRecord(r.session, networkRef)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get network record",
			err.Error(),
		)
		return
	}

	err = updateNetworkResourceModelComputed(ctx, record, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update the computed fields of NetworkResourceModel",
			err.Error(),
		)
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read data from State, retrieve the resource's information, update to State
// terraform import
func (r *networkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data networkResourceModel
	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Overwrite data with refreshed resource state
	networkRef, err := xenapi.Network.GetByUUID(r.session, data.UUID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get network ref",
			err.Error(),
		)
		return
	}

	networkRecord, err := xenapi.Network.GetRecord(r.session, networkRef)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get network record",
			err.Error(),
		)
		return
	}
	err = updateNetworkResourceModel(ctx, networkRecord, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update the fields of NetworkResourceModel",
			err.Error(),
		)
		return
	}
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read data from Plan, update resource configuration, Set to State
// terraform plan/apply (+2)
func (r *networkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data networkResourceModel
	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get existing network record
	networkRef, err := xenapi.Network.GetByUUID(r.session, data.UUID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get network ref",
			err.Error(),
		)
		return
	}

	// Update existing network resource with new plan
	err = xenapi.Network.SetNameLabel(r.session, networkRef, data.NameLabel.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set name label of network",
			err.Error(),
		)
		return
	}

	// Overwrite data with refreshed resource state
	record, err := xenapi.Network.GetRecord(r.session, networkRef)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get network record",
			err.Error(),
		)
		return
	}

	err = updateNetworkResourceModelComputed(ctx, record, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update the computed fields of network resource model",
			err.Error(),
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read data from State, delete resource
// terraform destroy
func (r *networkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data networkResourceModel
	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// delete resource
	networkRef, err := xenapi.Network.GetByUUID(r.session, data.UUID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get network ref",
			err.Error(),
		)
		return
	}

	err = xenapi.Network.Destroy(r.session, networkRef)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to destroy network",
			err.Error(),
		)
		return
	}
}

// Import existing resource with id, call Read()
// terraform import
func (r *networkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
