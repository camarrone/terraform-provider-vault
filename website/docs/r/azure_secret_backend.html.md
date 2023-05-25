---
layout: "vault"
page_title: "Vault: vault_azure_secret_backend resource"
sidebar_current: "docs-vault-resource-azure-secret-backend"
description: |-
  Creates an azure secret backend for Vault.
---

# vault\_azure\_secret\_backend

Creates an Azure Secret Backend for Vault.

The Azure secrets engine dynamically generates Azure service principals and role assignments. Vault roles can be mapped to one or more Azure roles, providing a simple, flexible way to manage the permissions granted to generated service principals.

~> **Important** All data provided in the resource configuration will be
written in cleartext to state and plan files generated by Terraform, and
will appear in the console output when Terraform runs. Protect these
artifacts accordingly. See
[the main provider documentation](../index.html)
for more details.

~> It is highly recommended that one transition to the Microsoft Graph API.
See [use_microsoft_graph_api ](https://www.vaultproject.io/api-docs/secret/azure#use_microsoft_graph_api)
for more information. The example below demonstrates how to do this. 

## Example Usage: *vault-1.9 and above*

```hcl
resource "vault_azure_secret_backend" "azure" {
  use_microsoft_graph_api = true
  subscription_id         = "11111111-2222-3333-4444-111111111111"
  tenant_id               = "11111111-2222-3333-4444-222222222222"
  client_id               = "11111111-2222-3333-4444-333333333333"
  client_secret           = "12345678901234567890"
  environment             = "AzurePublicCloud"
}
```

## Example Usage: *vault-1.8 and below*

```hcl
resource "vault_azure_secret_backend" "azure" {
  use_microsoft_graph_api = false
  subscription_id         = "11111111-2222-3333-4444-111111111111"
  tenant_id               = "11111111-2222-3333-4444-222222222222"
  client_id               = "11111111-2222-3333-4444-333333333333"
  client_secret           = "12345678901234567890"
  environment             = "AzurePublicCloud"
}
```

## Argument Reference

The following arguments are supported:

- `namespace` - (Optional) The namespace to provision the resource in.
  The value should not contain leading or trailing forward slashes.
  The `namespace` is always relative to the provider's configured [namespace](/docs/providers/vault#namespace).
   *Available only for Vault Enterprise*.

- `subscription_id` (`string: <required>`) - The subscription id for the Azure Active Directory.

- `use_microsoft_graph_api` (`bool: <optional>`) - Indicates whether the secrets engine should use 
  the Microsoft Graph API. This parameter has been deprecated and will be ignored in `vault-1.12+`. 
  For more information, please refer to the [Vault docs](https://developer.hashicorp.com/vault/api-docs/secret/azure#use_microsoft_graph_api)

- `tenant_id` (`string: <required>`) - The tenant id for the Azure Active Directory.

- `client_id` (`string:""`) - The OAuth2 client id to connect to Azure.

- `client_secret` (`string:""`) - The OAuth2 client secret to connect to Azure.

- `environment` (`string:""`) - The Azure environment.

- `path` (`string: <optional>`) - The unique path this backend should be mounted at. Defaults to `azure`.

- `disable_remount` - (Optional) If set, opts out of mount migration on path updates.
  See here for more info on [Mount Migration](https://www.vaultproject.io/docs/concepts/mount-migration)

## Attributes Reference

No additional attributes are exported by this resource.