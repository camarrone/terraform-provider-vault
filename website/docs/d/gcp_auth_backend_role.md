---
layout: "vault"
page_title: "Vault: vault_gcp_auth_backend_role data source"
sidebar_current: "docs-vault-datasource-gcp-auth-backend-role"
description: |-
  Manages GCP auth backend roles in Vault.
---

# vault\_gcp\_auth\_backend\_role

Reads a GCP auth role from a Vault server.

## Example Usage

```hcl
data "vault_gcp_auth_backend_role" "role" {
  backend   = "my-gcp-backend"
  role_name = "my-role"
}

output "role-id" {
  value = "${data.vault_gcp_auth_backend_role.role.role_id}"
}
```

## Argument Reference

The following arguments are supported:

* `namespace` - (Optional) The namespace of the target resource.
  The value should not contain leading or trailing forward slashes.
  The `namespace` is always relative to the provider's configured [namespace](../index.html#namespace).
  *Available only for Vault Enterprise*.

* `role_name` - (Required) The name of the role to retrieve the Role ID for.

* `backend` - (Optional) The unique name for the GCP backend from which to fetch the role. Defaults to "gcp".

## Attributes Reference

In addition to the above arguments, the following attributes are exported:

* `role_id` - The RoleID of the GCP role.

* `type` - Type of GCP role. Expected values are `iam` or `gce`.

* `bound_service_accounts` - GCP service accounts bound to the role. Returned when `type` is `iam`.

* `bound_projects` - GCP projects bound to the role.

* `bound_zones` - GCP zones bound to the role. Returned when `type` is `gce`.

* `bound_regions` - GCP regions bound to the role. Returned when `type` is `gce`.

* `bound_instance_groups` - GCP regions bound to the role. Returned when `type` is `gce`.

* `bound_labels` - GCP labels bound to the role. Returned when `type` is `gce`.

* `token_policies` - Token policies bound to the role.

### Common Token Attributes

These attributes are common across several Authentication Token resources since Vault 1.2.

* `token_ttl` - The incremental lifetime for generated tokens in number of seconds.
  Its current value will be referenced at renewal time.

* `token_max_ttl` - The maximum lifetime for generated tokens in number of seconds.
  Its current value will be referenced at renewal time.

* `token_period` - (Optional) If set, indicates that the
  token generated using this role should never expire. The token should be renewed within the
  duration specified by this value. At each renewal, the token's TTL will be set to the
  value of this field. Specified in seconds.

* `token_policies` - List of policies to encode onto generated tokens. Depending
  on the auth method, this list may be supplemented by user/group/other values.

* `token_bound_cidrs` - List of CIDR blocks; if set, specifies blocks of IP
  addresses which can authenticate successfully, and ties the resulting token to these blocks
  as well.

* `token_explicit_max_ttl` - If set, will encode an
  [explicit max TTL](https://www.vaultproject.io/docs/concepts/tokens.html#token-time-to-live-periodic-tokens-and-explicit-max-ttls)
  onto the token in number of seconds. This is a hard cap even if `token_ttl` and
  `token_max_ttl` would otherwise allow a renewal.

* `token_no_default_policy` - If set, the default policy will not be set on
  generated tokens; otherwise it will be added to the policies set in token_policies.

* `token_num_uses` - The
  [period](https://www.vaultproject.io/docs/concepts/tokens.html#token-time-to-live-periodic-tokens-and-explicit-max-ttls),
  if any, in number of seconds to set on the token.

* `token_type` - The type of token that should be generated. Can be `service`,
  `batch`, or `default` to use the mount's tuned default (which unless changed will be
  `service` tokens). For token store roles, there are two additional possibilities:
  `default-service` and `default-batch` which specify the type to return unless the client
  requests a different type at generation time.
