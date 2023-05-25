// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-vault/internal/provider"
	"github.com/hashicorp/terraform-provider-vault/util"
)

var (
	approleAuthBackendRoleBackendFromPathRegex = regexp.MustCompile("^auth/(.+)/role/.+$")
	approleAuthBackendRoleNameFromPathRegex    = regexp.MustCompile("^auth/.+/role/(.+)$")
)

func approleAuthBackendRoleResource() *schema.Resource {
	fields := map[string]*schema.Schema{
		"role_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name of the role.",
			ForceNew:    true,
		},
		"role_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The RoleID of the role. Autogenerated if not set.",
		},
		"bind_secret_id": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Whether or not to require secret_id to be present when logging in using this AppRole.",
		},
		"secret_id_bound_cidrs": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "List of CIDR blocks that can log in using the AppRole.",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"secret_id_num_uses": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Number of times which a particular SecretID can be used to fetch a token from this AppRole, after which the SecretID will expire. Leaving this unset or setting it to 0 will allow unlimited uses.",
		},
		"secret_id_ttl": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Number of seconds a SecretID remains valid for.",
		},
		"backend": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Unique name of the auth backend to configure.",
			ForceNew:    true,
			Default:     "approle",
			// standardise on no beginning or trailing slashes
			StateFunc: func(v interface{}) string {
				return strings.Trim(v.(string), "/")
			},
		},
	}

	addTokenFields(fields, &addTokenFieldsConfig{})

	return &schema.Resource{
		CreateContext: approleAuthBackendRoleCreate,
		ReadContext:   ReadContextWrapper(approleAuthBackendRoleRead),
		UpdateContext: approleAuthBackendRoleUpdate,
		DeleteContext: approleAuthBackendRoleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: fields,
	}
}

func approleAuthBackendRoleUpdateFields(d *schema.ResourceData, data map[string]interface{}, create bool) {
	updateTokenFields(d, data, create)

	if create {
		if v, ok := d.GetOkExists("bind_secret_id"); ok {
			data["bind_secret_id"] = v.(bool)
		}

		if v, ok := d.GetOk("secret_id_num_uses"); ok {
			data["secret_id_num_uses"] = v.(int)
		}

		if v, ok := d.GetOk("secret_id_ttl"); ok {
			data["secret_id_ttl"] = v.(int)
		}

		if v, ok := d.GetOk("secret_id_bound_cidrs"); ok {
			data["secret_id_bound_cidrs"] = v.(*schema.Set).List()
		}
	} else {
		if d.HasChange("bind_secret_id") {
			data["bind_secret_id"] = d.Get("bind_secret_id").(bool)
		}

		if d.HasChange("secret_id_num_uses") {
			data["secret_id_num_uses"] = d.Get("secret_id_num_uses").(int)
		}

		if d.HasChange("secret_id_ttl") {
			data["secret_id_ttl"] = d.Get("secret_id_ttl").(int)
		}

		if d.HasChange("secret_id_bound_cidrs") {
			data["secret_id_bound_cidrs"] = d.Get("secret_id_bound_cidrs").(*schema.Set).List()
		}
	}
}

func approleAuthBackendRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, e := provider.GetClient(d, meta)
	if e != nil {
		return diag.FromErr(e)
	}

	backend := d.Get("backend").(string)
	role := d.Get("role_name").(string)

	path := approleAuthBackendRolePath(backend, role)

	log.Printf("[DEBUG] Writing AppRole auth backend role %q", path)

	diags := diag.Diagnostics{}

	data := map[string]interface{}{}
	approleAuthBackendRoleUpdateFields(d, data, true)

	_, err := client.Logical().Write(path, data)
	if err != nil {
		diags = append(diags,
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("error writing AppRole auth backend role %q: %s", path, err),
			},
		)
		return diags
	}

	d.SetId(path)
	log.Printf("[DEBUG] Wrote AppRole auth backend role %q", path)

	if v, ok := d.GetOk("role_id"); ok {
		log.Printf("[DEBUG] Writing AppRole auth backend role %q RoleID", path)
		_, err := client.Logical().Write(path+"/role-id", map[string]interface{}{
			"role_id": v.(string),
		})
		if err != nil {
			diags = append(diags,
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("error writing AppRole auth backend role %q's RoleID: %s", path, err),
				},
			)
			return diags
		}

		log.Printf("[DEBUG] Wrote AppRole auth backend role %q RoleID", path)
	}

	return append(diags, approleAuthBackendRoleRead(ctx, d, meta)...)
}

func approleAuthBackendRoleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, e := provider.GetClient(d, meta)
	if e != nil {
		return diag.FromErr(e)
	}

	path := d.Id()

	backend, err := approleAuthBackendRoleBackendFromPath(path)
	if err != nil {
		return diag.Errorf("invalid path %q for AppRole auth backend role: %s", path, err)
	}

	role, err := approleAuthBackendRoleNameFromPath(path)
	if err != nil {
		return diag.Errorf("invalid path %q for AppRole auth backend role: %s", path, err)
	}

	log.Printf("[DEBUG] Reading AppRole auth backend role %q", path)
	resp, err := client.Logical().Read(path)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Read AppRole auth backend role %q", path)
	if resp == nil {
		log.Printf("[WARN] AppRole auth backend role %q not found, removing from state", path)
		d.SetId("")
		return nil
	}

	if err := d.Set("backend", backend); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("role_name", role); err != nil {
		return diag.FromErr(err)
	}

	if err := readTokenFields(d, resp); err != nil {
		return diag.FromErr(err)
	}

	if v, ok := resp.Data["secret_id_bound_cidrs"]; ok {
		if err := d.Set("secret_id_bound_cidrs", v); err != nil {
			return diag.FromErr(err)
		}
	}

	for _, k := range []string{"bind_secret_id", "secret_id_num_uses", "secret_id_ttl"} {
		if err := d.Set(k, resp.Data[k]); err != nil {
			return diag.FromErr(err)
		}
	}

	log.Printf("[DEBUG] Reading AppRole auth backend role %q RoleID", path)
	resp, err = client.Logical().Read(path + "/role-id")
	if err != nil {
		return diag.Errorf("error reading AppRole auth backend role %q RoleID: %s", path, err)
	}
	log.Printf("[DEBUG] Read AppRole auth backend role %q RoleID", path)
	if resp != nil {
		if err := d.Set("role_id", resp.Data["role_id"]); err != nil {
			return diag.FromErr(err)
		}
	}

	diags := checkCIDRs(d, TokenFieldBoundCIDRs)

	return diags
}

func approleAuthBackendRoleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, e := provider.GetClient(d, meta)
	if e != nil {
		return diag.FromErr(e)
	}

	path := d.Id()

	log.Printf("[DEBUG] Updating AppRole auth backend role %q", path)

	data := map[string]interface{}{}
	approleAuthBackendRoleUpdateFields(d, data, false)

	_, err := client.Logical().Write(path, data)

	d.SetId(path)

	diags := diag.Diagnostics{}
	if err != nil {
		return diag.Errorf("error updating AppRole auth backend role %q: %s", path, err)
	}
	log.Printf("[DEBUG] Updated AppRole auth backend role %q", path)

	if d.HasChange("role_id") {
		log.Printf("[DEBUG] Updating AppRole auth backend role %q RoleID", path)
		_, err := client.Logical().Write(path+"/role-id", map[string]interface{}{
			"role_id": d.Get("role_id").(string),
		})
		if err != nil {
			return diag.Errorf("error updating AppRole auth backend role %q's RoleID: %s", path, err)
		}
		log.Printf("[DEBUG] Updated AppRole auth backend role %q RoleID", path)
	}

	return append(diags, approleAuthBackendRoleRead(ctx, d, meta)...)
}

func approleAuthBackendRoleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, e := provider.GetClient(d, meta)
	if e != nil {
		return diag.FromErr(e)
	}

	path := d.Id()

	log.Printf("[DEBUG] Deleting AppRole auth backend role %q", path)
	_, err := client.Logical().Delete(path)
	if err != nil {
		if util.Is404(err) {
			log.Printf("[DEBUG] AppRole auth backend role %q not found, removing from state", path)
			d.SetId("")
			return nil
		} else {
			return diag.Errorf("error deleting AppRole auth backend role %q, err=%s", path, err)
		}
	}

	log.Printf("[DEBUG] Deleted AppRole auth backend role %q", path)

	return nil
}

func approleAuthBackendRolePath(backend, role string) string {
	return "auth/" + strings.Trim(backend, "/") + "/role/" + strings.Trim(role, "/")
}

func approleAuthBackendRoleNameFromPath(path string) (string, error) {
	if !approleAuthBackendRoleNameFromPathRegex.MatchString(path) {
		return "", fmt.Errorf("no role found")
	}
	res := approleAuthBackendRoleNameFromPathRegex.FindStringSubmatch(path)
	if len(res) != 2 {
		return "", fmt.Errorf("unexpected number of matches (%d) for role", len(res))
	}
	return res[1], nil
}

func approleAuthBackendRoleBackendFromPath(path string) (string, error) {
	if !approleAuthBackendRoleBackendFromPathRegex.MatchString(path) {
		return "", fmt.Errorf("no backend found")
	}
	res := approleAuthBackendRoleBackendFromPathRegex.FindStringSubmatch(path)
	if len(res) != 2 {
		return "", fmt.Errorf("unexpected number of matches (%d) for backend", len(res))
	}
	return res[1], nil
}