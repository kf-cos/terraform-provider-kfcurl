terraform {
  required_providers {
    kfcurl = {
      source  = "kf-cos/kfcurl"
      version = "1.2.0"
    }
  }
}

provider "kfcurl" {}

data "kfcurl_request" "test" {
  name   = "products"
  url    = "https://api.releases.hashicorp.com/v1/products"
  method = "GET"

  response_codes = [
    200
  ]

  max_retry      = 1
  retry_interval = 10
}

output "response" {
  value = jsondecode(data.kfcurl_request.test.response)
}