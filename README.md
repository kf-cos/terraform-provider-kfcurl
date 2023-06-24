# Terraform Provider KFCurl

Available in the [Terraform Registry.](https://registry.terraform.io/providers/kf-cos/kfcurl/latest/docs)

**This is a customised fork of https://github.com/devops-rob/terraform-provider-terracurl. Please use the original if you're 
not part of COS.**

This provider is designed to be a flexible extension of your terraform code to make managed and unmanaged API calls to your 
target endpoint. Platform native providers should be preferred to KFCurl, but for instances where the platform provider does 
not have a resource or data source that you require, KFCurl can be used to make substitute API calls.

## Managed API calls
When using KFCurl, if the API call is creating a change on the target platform, and you would like this change reversed 
upon a destroy, use the `kfcurl_request` resource. This will allow you to enter the API call that should be run 
when `terraform destroy` is run.

### Create ServiceNow CI
```hcl
resource "kfcurl_request" "servicenow" {
  name         = "az-cos-aws-mk"
  url          = "https://insert_the_name_here.service-now.com/api/now/table/u_google_organization_project?sysparm_display_value=true&sysparm_exclude_reference_link=true"
  method       = "POST"

  request_body = <<EOF
        {"u_name":"az-cos-aws-mk",
        "u_status":"5",
        "u_environment":"Non Production",
        "u_ci_owner_group":"Awesome Group",
        "u_date_certified":"2023-06-24",
        "u_certified_by_":"user02",
        "u_project_code_":"Non-Project",
        "u_description":"desc",
        "u_notes_":"NA",
        "u_alert_notes_":"NA",
        "u_special_":"true",
        "u_account_id":"id"}

EOF

  headers = { "Content-type" : "application/json", "Accept" : "application/json", "Authorization": "Basic B4s364ENc0D3dCr3d3nt1al5" }

  response_codes = [
    200,
    201
  ]
}

output "response" {
  value = kfcurl_request.servicenow.response
}
```

### More detailed example of using KFCurl

```hcl
resource "kfcurl_request" "mount" {
  name           = "vault-mount"
  url            = "https://localhost:8200/v1/sys/mounts/aws"
  method         = "POST"
  request_body   = <<EOF
{
  "type": "aws",
  "config": {
    "force_no_cache": true
  }
}

EOF

  headers = {
    X-Vault-Token = "root"
  }

  response_codes = [
    200,
    204
  ]

  cert_file       = "server-vault-0.pem"
  key_file        = "server-vault-0-key.pem"
  ca_cert_file    = "vault-server-ca.pem"
  skip_tls_verify = false


  destroy_url    = "https://localhost:8200/v1/sys/mounts/aws"
  destroy_method = "DELETE"

  destroy_headers = {
    X-Vault-Token = "root"
  }

  destroy_response_codes = [
    204
  ]

  destroy_cert_file       = "server-vault-0.pem"
  destroy_key_file        = "server-vault-0-key.pem"
  destroy_ca_cert_file    = "vault-server-ca.pem"
  destroy_skip_tls_verify = false

}
```
## Unmanaged API calls
For instances where there is no change required on the target platform when the `terraform destroy` command is run, 
use the `kfcurl_request` data source.

```hcl
data "kfcurl_request" "test" {
  name           = "products"
  url            = "https://api.releases.hashicorp.com/v1/products"
  method         = "GET"

  response_codes = [
    200
  ]

  max_retry      = 1
  retry_interval = 10
}
```
## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
-	[Go](https://golang.org/doc/install) >= 1.17

## Automatic Release

GitHub Action will build a new version of this provider every time a new tag is created.

## Building The Provider Manually

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the Go `install` command: 
```sh
$ go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most-up-to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.
