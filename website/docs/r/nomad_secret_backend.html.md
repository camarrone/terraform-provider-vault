---
layout: "vault"
page_title: "Vault: vault_nomad_secret_backend resource"
sidebar_current: "docs-vault-resource-nomad-secret-backend"
description: |-
  Creates a Nomad secret backend for Vault.
---

# vault\_nomad\_secret\_backend

Creates a Nomad Secret Backend for Vault. The Nomad secret backend for Vault
generates Nomad ACL tokens dynamically based on pre-existing Nomad ACL policies.

~> **Important** All data provided in the resource configuration will be
written in cleartext to state and plan files generated by Terraform, and
will appear in the console output when Terraform runs. Protect these
artifacts accordingly. See
[the main provider documentation](../index.html)
for more details.

## Example Usage

```hcl
resource "vault_nomad_secret_backend" "config" {
	backend                   = "nomad"
	description               = "test description"
	default_lease_ttl_seconds = "3600"
	max_lease_ttl_seconds     = "7200"
	max_ttl                   = "240"
	address                   = "https://127.0.0.1:4646"
	token                     = "ae20ceaa-..."
	ttl                       = "120"
}
```

## Argument Reference

The following arguments are supported:

* `namespace` - (Optional) The namespace to provision the resource in.
  The value should not contain leading or trailing forward slashes.
  The `namespace` is always relative to the provider's configured [namespace](/docs/providers/vault#namespace).
   *Available only for Vault Enterprise*.

* `backend` - (Optional) The unique path this backend should be mounted at. Must
not begin or end with a `/`. Defaults to `nomad`.

* `disable_remount` - (Optional) If set, opts out of mount migration on path updates.
  See here for more info on [Mount Migration](https://www.vaultproject.io/docs/concepts/mount-migration)

* `address` - (Optional) Specifies the address of the Nomad instance, provided
as "protocol://host:port" like "http://127.0.0.1:4646".

* `ca_cert` - (Optional) CA certificate to use when verifying the Nomad server certificate, must be
x509 PEM encoded.

* `client_cert` - (Optional) Client certificate to provide to the Nomad server, must be x509 PEM encoded.

* `client_key` - (Optional) Client certificate key to provide to the Nomad server, must be x509 PEM encoded.

* `default_lease_ttl_seconds` - (Optional) Default lease duration for secrets in seconds.

* `description` - (Optional) Human-friendly description of the mount for the Active Directory backend.

* `local` - (Optional) Mark the secrets engine as local-only. Local engines are not replicated or removed by
replication.Tolerance duration to use when checking the last rotation time.

* `max_token_name_length` - (Optional) Specifies the maximum length to use for the name of the Nomad token
generated with Generate Credential. If omitted, 0 is used and ignored, defaulting to the max value allowed
by the Nomad version.

* `max_ttl` - (Optional) Maximum possible lease duration for secrets in seconds.

* `token` - (Optional) Specifies the Nomad Management token to use.

* `ttl` - (Optional) Specifies the ttl of the lease for the generated token.



## Attributes Reference

No additional attributes are exported by this resource.

## Import

Nomad secret backend can be imported using the `backend`, e.g.

```
$ terraform import vault_nomad_secret_backend.nomad nomad
```