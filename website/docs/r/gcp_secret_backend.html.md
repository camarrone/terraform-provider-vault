---
layout: "vault"
page_title: "Vault: vault_gcp_secret_backend resource"
sidebar_current: "docs-vault-resource-gcp-secret-backend"
description: |-
  Creates an GCP secret backend for Vault.
---

# vault\_gcp\_secret\_backend

Creates an GCP Secret Backend for Vault. GCP secret backends can then issue GCP
OAuth token or Service Account keys, once a role has been added to the backend.

~> **Important** All data provided in the resource configuration will be
written in cleartext to state and plan files generated by Terraform, and
will appear in the console output when Terraform runs. Protect these
artifacts accordingly. See
[the main provider documentation](../index.html)
for more details.

## Example Usage

```hcl
resource "vault_gcp_secret_backend" "gcp" {
  credentials = file("credentials.json")
}
```

## Argument Reference

The following arguments are supported:

* `namespace` - (Optional) The namespace to provision the resource in.
  The value should not contain leading or trailing forward slashes.
  The `namespace` is always relative to the provider's configured [namespace](/docs/providers/vault#namespace).
   *Available only for Vault Enterprise*.

* `credentials` - (Optional) The GCP service account credentials in JSON format.

~> **Important** Because Vault does not support reading the configured
credentials back from the API, Terraform cannot detect and correct drift
on `credentials`. Changing the values, however, _will_ overwrite the
previously stored values.

* `path` - (Optional) The unique path this backend should be mounted at. Must
not begin or end with a `/`. Defaults to `gcp`.

* `disable_remount` - (Optional) If set, opts out of mount migration on path updates.
  See here for more info on [Mount Migration](https://www.vaultproject.io/docs/concepts/mount-migration)

* `description` - (Optional) A human-friendly description for this backend.

* `default_lease_ttl_seconds` - (Optional) The default TTL for credentials
issued by this backend. Defaults to '0'.

* `max_lease_ttl_seconds` - (Optional) The maximum TTL that can be requested
for credentials issued by this backend. Defaults to '0'.

* `local` - (Optional) Boolean flag that can be explicitly set to true to enforce local mount in HA environment

## Attributes Reference

No additional attributes are exported by this resource.