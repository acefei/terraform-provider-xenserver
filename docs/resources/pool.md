---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "xenserver_pool Resource - xenserver"
subcategory: ""
description: |-
  This provides a pool resource.
  -> Note: During the execution of terraform destroy for this particular resource, all of the hosts that are part of the pool will be separated and converted into standalone hosts.
---

# xenserver_pool (Resource)

This provides a pool resource.

-> **Note:** During the execution of `terraform destroy` for this particular resource, all of the hosts that are part of the pool will be separated and converted into standalone hosts.

## Example Usage

```terraform
resource "xenserver_sr_nfs" "nfs" {
  name_label       = "NFS shared storage"
  name_description = "A test NFS storage repository"
  version          = "3"
  storage_location = format("%s:%s", local.env_vars["NFS_SERVER"], local.env_vars["NFS_SERVER_PATH"])
}

data "xenserver_pif" "pif" {
  device = "eth0"
}

data "xenserver_pif" "pif1" {
  device = "eth3"
}

locals {
  pif1_data = tomap({for element in data.xenserver_pif.pif1.data_items: element.uuid => element})
}

resource "xenserver_pif_configure" "pif_update" {
  for_each = local.pif1_data
  uuid     = each.key
  interface = {
    mode = "DHCP"
  }
}

# Configure default SR and Management Network of the pool
resource "xenserver_pool" "pool" {
  name_label   = "pool"
  default_sr = xenserver_sr_nfs.nfs.uuid
  management_network = data.xenserver_pif.pif.data_items[0].network
}

# Join supporter into the pool
resource "xenserver_pool" "pool" {
  name_label   = "pool"
  join_supporters = [
    {
      host = local.env_vars["SUPPORTER_HOST"]
      username = local.env_vars["SUPPORTER_USERNAME"]
      password = local.env_vars["SUPPORTER_PASSWORD"]
    }
  ]
}

# Eject supporter from the pool
data "xenserver_host" "supporter" {
  is_coordinator = false
}

resource "xenserver_pool" "pool" {
  name_label   = "pool"
  eject_supporters = [ data.xenserver_host.supporter.data_items[1].uuid ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name_label` (String) The name of the pool.

### Optional

- `default_sr` (String) The default SR UUID of the pool. this SR should be shared SR.
- `eject_supporters` (Set of String) The set of pool supporters which will be ejected from the pool.
- `join_supporters` (Attributes Set) The set of pool supporters which will join the pool.

-> **Note:** 1. It would raise error if a supporter is in both join_supporters and eject_supporters.<br>2. The join operation would be performed only when the host, username, and password are provided.<br> (see [below for nested schema](#nestedatt--join_supporters))
- `management_network` (String) The management network UUID of the pool.

-> **Note:** 1. The management network would be reconfigured only when the management network UUID is provided.<br>2. All of the hosts in the pool should have the same management network with network configuration, and you can set network configuration by resource `pif_configure`.<br>3. It is not recommended to set the `management_network` with the `join_supporters` and `eject_supporters` attributes together.<br>
- `name_description` (String) The description of the pool, default to be `""`.

### Read-Only

- `id` (String) The test ID of the pool.
- `uuid` (String) The UUID of the pool.

<a id="nestedatt--join_supporters"></a>
### Nested Schema for `join_supporters`

Optional:

- `host` (String) The address of the host.
- `password` (String, Sensitive) The password of the host.
- `username` (String) The user name of the host.

## Import

Import is supported using the following syntax:

```shell
terraform import xenserver_pool.pool 00000000-0000-0000-0000-000000000000
```
