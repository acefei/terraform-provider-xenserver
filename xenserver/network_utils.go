package xenserver

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"xenapi"
)

// NetworkResourceModel describes the resource data model.
type networkResourceModel struct {
	NameLabel       types.String `tfsdk:"name_label"`
	NameDescription types.String `tfsdk:"name_description"`
	OtherConfig     types.Map    `tfsdk:"other_config"`
	UUID            types.String `tfsdk:"id"`
}

// Update NetworkResourceModel base on new NetworkRecord
func updateNetworkResourceModel(ctx context.Context, networkRecord xenapi.NetworkRecord, data *networkResourceModel) error {
	data.NameLabel = types.StringValue(networkRecord.NameLabel)

	err := updateNetworkResourceModelComputed(ctx, networkRecord, data)
	if err != nil {
		return err
	}
	return nil
}

// Update NetworkResourceModel computed field base on new NetworkRecord
func updateNetworkResourceModelComputed(ctx context.Context, networkRecord xenapi.NetworkRecord, data *networkResourceModel) error {
	data.UUID = types.StringValue(networkRecord.UUID)
	data.NameDescription = types.StringValue(networkRecord.NameDescription)

	var diags diag.Diagnostics
	data.OtherConfig, diags = types.MapValueFrom(ctx, types.StringType, networkRecord.OtherConfig)
	if diags.HasError() {
		return errors.New("unable to update data for network other_config")
	}
	return nil
}
