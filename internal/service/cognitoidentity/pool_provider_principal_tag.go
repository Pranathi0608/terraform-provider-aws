package cognitoidentity

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentity"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_cognito_identity_pool_provider_principal_tag")
func ResourcePoolProviderPrincipalTag() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourcePoolProviderPrincipalTagCreate,
		ReadWithoutTimeout:   resourcePoolProviderPrincipalTagRead,
		UpdateWithoutTimeout: resourcePoolProviderPrincipalTagUpdate,
		DeleteWithoutTimeout: resourcePoolProviderPrincipalTagDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"identity_pool_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 55),
					validation.StringMatch(regexp.MustCompile(`^[\w-]+:[0-9a-f-]+$`), "see https://docs.aws.amazon.com/cognitoidentity/latest/APIReference/API_SetPrincipalTagAttributeMap.html#API_SetPrincipalTagAttributeMap_ResponseSyntax"),
				),
			},
			"identity_provider_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 128),
				),
			},
			"principal_tags": tftags.TagsSchema(),
			"use_defaults": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourcePoolProviderPrincipalTagCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).CognitoIdentityConn()
	log.Print("[DEBUG] Creating Cognito Identity Provider Principal Tags")

	providerName := d.Get("identity_provider_name").(string)
	poolId := d.Get("identity_pool_id").(string)

	params := &cognitoidentity.SetPrincipalTagAttributeMapInput{
		IdentityPoolId:       aws.String(poolId),
		IdentityProviderName: aws.String(providerName),
	}

	if v, ok := d.GetOk("principal_tags"); ok {
		params.PrincipalTags = flex.ExpandStringMap(v.(map[string]interface{}))
	}

	if v, ok := d.GetOk("use_defaults"); ok {
		params.UseDefaults = aws.Bool(v.(bool))
	}

	_, err := conn.SetPrincipalTagAttributeMapWithContext(ctx, params)
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating Cognito Identity Provider Principal Tags: %s", err)
	}

	d.SetId(fmt.Sprintf("%s:%s", poolId, providerName))

	return append(diags, resourcePoolProviderPrincipalTagRead(ctx, d, meta)...)
}

func resourcePoolProviderPrincipalTagRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).CognitoIdentityConn()
	log.Printf("[DEBUG] Reading Cognito Identity Provider Principal Tags: %s", d.Id())

	poolId, providerName, err := DecodePoolProviderPrincipalTagsID(d.Id())

	if err != nil {
		return create.DiagError(names.CognitoIdentity, create.ErrActionReading, ResNamePoolProviderPrincipalTag, d.Id(), err)
	}

	ret, err := conn.GetPrincipalTagAttributeMapWithContext(ctx, &cognitoidentity.GetPrincipalTagAttributeMapInput{
		IdentityProviderName: aws.String(providerName),
		IdentityPoolId:       aws.String(poolId),
	})

	if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, cognitoidentity.ErrCodeResourceNotFoundException) {
		create.LogNotFoundRemoveState(names.CognitoIdentity, create.ErrActionReading, ResNamePoolProviderPrincipalTag, d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return create.DiagError(names.CognitoIdentity, create.ErrActionReading, ResNamePoolProviderPrincipalTag, d.Id(), err)
	}

	d.Set("identity_pool_id", ret.IdentityPoolId)
	d.Set("identity_provider_name", ret.IdentityProviderName)
	d.Set("use_defaults", ret.UseDefaults)

	if err := d.Set("principal_tags", aws.StringValueMap(ret.PrincipalTags)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting principal_tags: %s", err)
	}

	return diags
}

func resourcePoolProviderPrincipalTagUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).CognitoIdentityConn()
	log.Print("[DEBUG] Updating Cognito Identity Provider Principal Tags")

	poolId, providerName, err := DecodePoolProviderPrincipalTagsID(d.Id())
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "updating Cognito Identity Provider Principal Tags (%s): %s", d.Id(), err)
	}

	params := &cognitoidentity.SetPrincipalTagAttributeMapInput{
		IdentityPoolId:       aws.String(poolId),
		IdentityProviderName: aws.String(providerName),
	}

	if d.HasChanges("principal_tags", "use_defaults") {
		params.PrincipalTags = flex.ExpandStringMap(d.Get("principal_tags").(map[string]interface{}))
		params.UseDefaults = aws.Bool(d.Get("use_defaults").(bool))

		_, err = conn.SetPrincipalTagAttributeMapWithContext(ctx, params)
		if err != nil {
			return sdkdiag.AppendErrorf(diags, "updating Cognito Identity Provider Principal Tags (%s): %s", d.Id(), err)
		}
	}

	return append(diags, resourcePoolProviderPrincipalTagRead(ctx, d, meta)...)
}

func resourcePoolProviderPrincipalTagDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).CognitoIdentityConn()
	log.Printf("[DEBUG] Deleting Cognito Identity Provider Principal Tags: %s", d.Id())

	poolId, providerName, err := DecodePoolProviderPrincipalTagsID(d.Id())
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting Cognito Identity Provider Principal Tags (%s): %s", d.Id(), err)
	}
	emptyList := make(map[string]string)
	params := &cognitoidentity.SetPrincipalTagAttributeMapInput{
		IdentityPoolId:       aws.String(poolId),
		IdentityProviderName: aws.String(providerName),
		UseDefaults:          aws.Bool(true),
		PrincipalTags:        aws.StringMap(emptyList),
	}

	_, err = conn.SetPrincipalTagAttributeMapWithContext(ctx, params)

	if err != nil {
		if tfawserr.ErrCodeEquals(err, cognitoidentity.ErrCodeResourceNotFoundException) {
			return diags
		}
		return sdkdiag.AppendErrorf(diags, "deleting Cognito Identity Provider Principal Tags (%s): %s", d.Id(), err)
	}
	return diags
}

func DecodePoolProviderPrincipalTagsID(id string) (string, string, error) {
	r := regexp.MustCompile(`(?P<ProviderID>[\w-]+:[0-9a-f-]+):(?P<ProviderName>[[:graph:]]+)`)
	idParts := r.FindStringSubmatch(id)
	if len(idParts) <= 2 {
		return "", "", fmt.Errorf("expected ID in format UserPoolID:ProviderName, received: %s", id)
	}
	return idParts[1], idParts[2], nil
}
