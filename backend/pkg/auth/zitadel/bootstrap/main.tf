terraform {
  required_providers {
    zitadel = {
      source  = "zitadel/zitadel"
      version = "1.0.0-alpha.16"
    }
  }
}

provider "zitadel" {
  domain           = "localhost"
  insecure         = "true"
  port             = "8080"
  jwt_profile_file = "machinekey/zitadel-admin-sa.json"
}
