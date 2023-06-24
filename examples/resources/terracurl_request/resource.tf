terraform {
  required_providers {
    kfcurl = {
      source  = "kf-cos/kfcurl"
      version = "1.2.1"
    }
  }
}

provider "kfcurl" {}

resource "kfcurl_request" "mount" {
  name         = "vault-mount"
  url          = "https://localhost:8200/v1/sys/mounts/aws"
  method       = "POST"
  request_body = <<EOF
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

output "response" {
  value = kfcurl_request.mount.response
}