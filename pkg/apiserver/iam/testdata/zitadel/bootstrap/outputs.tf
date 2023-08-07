resource local_file file_vmclarity_data {
  content  = <<-EOT
    IAM_ENABLED=true
    FAKE_DATA=true

    AUTH_OIDC_ISSUER=http://localhost:8080
    AUTH_OIDC_CLIENT_ID=${zitadel_application_api.vmclarity_app_api.client_id}
    AUTH_OIDC_CLIENT_SECRET=${zitadel_application_api.vmclarity_app_api.client_secret}

    ROLESYNCER_JWT_ROLE_CLAIM=urn:zitadel:iam:org:project:${zitadel_project.vmclarity_project.id}:roles

    AUTHZ_LOCAL_RBAC_RULE_FILEPATH=path-to-folder/rbac_rule_policy_example.csv

    APISERVER_BEARER_TOKEN_ENV_VAR=APISERVER_BEARER_TOKEN
    APISERVER_BEARER_TOKEN=${zitadel_personal_access_token.vmclarity_orchestrator_pat_key.token}
  EOT
  filename = "generated/vmclarity-data.env"
}
