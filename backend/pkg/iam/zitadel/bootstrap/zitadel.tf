//----------------------------------------
// Organization configuration
//----------------------------------------
resource zitadel_org openclarity_org {
  name = "openclarity"
}

resource zitadel_domain openclarity_domain {
  org_id      = zitadel_org.openclarity_org.id
  name        = "openclarity.io"
}

resource zitadel_org_idp_github openclarity_org_idp_github {
  org_id              = zitadel_org.openclarity_org.id
  name                = "OpenClarity GitHub"
  client_id           = var.idp_github_client_id
  client_secret       = var.idp_github_client_secret
  scopes              = ["openid", "profile", "email"]
  is_linking_allowed  = true
  is_creation_allowed = true
  is_auto_creation    = true
  is_auto_update      = true
}

//----------------------------------------
// Project configuration
//----------------------------------------
resource zitadel_login_policy openclarity_login_policy {
  org_id                        = zitadel_org.openclarity_org.id
  user_login                    = true
  allow_register                = true
  allow_external_idp            = true
  force_mfa                     = false
  passwordless_type             = "PASSWORDLESS_TYPE_ALLOWED"
  hide_password_reset           = "false"
  password_check_lifetime       = "240h0m0s"
  external_login_check_lifetime = "240h0m0s"
  multi_factor_check_lifetime   = "24h0m0s"
  mfa_init_skip_lifetime        = "720h0m0s"
  second_factor_check_lifetime  = "24h0m0s"
  ignore_unknown_usernames      = true
  default_redirect_uri          = "localhost:8080"
  second_factors                = ["SECOND_FACTOR_TYPE_OTP", "SECOND_FACTOR_TYPE_U2F"]
  multi_factors                 = ["MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION"]
  idps                          = [zitadel_org_idp_github.openclarity_org_idp_github.id]
  allow_domain_discovery        = true
  disable_login_with_email      = false
  disable_login_with_phone      = false
}

resource zitadel_project vmclarity_project {
  name                     = "vmclarity"
  org_id                   = zitadel_org.openclarity_org.id
  project_role_assertion   = true
  project_role_check       = true
  has_project_check        = true
  private_labeling_setting = "PRIVATE_LABELING_SETTING_ENFORCE_PROJECT_RESOURCE_OWNER_POLICY"
}

resource zitadel_project_role vmclarity_project_role {
  for_each = var.project_roles

  org_id       = zitadel_org.openclarity_org.id
  project_id   = zitadel_project.vmclarity_project.id
  group        = split(":", each.value)[0]
  role_key     = each.value
  display_name = each.value
}

resource zitadel_application_api vmclarity_app_api {
  org_id           = zitadel_org.openclarity_org.id
  project_id       = zitadel_project.vmclarity_project.id
  name             = "API"
  auth_method_type = "API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT"
}

resource zitadel_application_key vmclarity_app_api_key {
  org_id          = zitadel_org.openclarity_org.id
  project_id      = zitadel_project.vmclarity_project.id
  app_id          = zitadel_application_api.vmclarity_app_api.id
  key_type        = "KEY_TYPE_JSON"
  expiration_date = ""
}

resource zitadel_application_oidc vmclarity_app_oidc {
  org_id                      = zitadel_org.openclarity_org.id
  project_id                  = zitadel_project.vmclarity_project.id
  name                        = "Web"
  redirect_uris               = ["https://localhost.com"]
  response_types              = ["OIDC_RESPONSE_TYPE_CODE"]
  grant_types                 = ["OIDC_GRANT_TYPE_AUTHORIZATION_CODE"]
  post_logout_redirect_uris   = ["https://localhost.com"]
  app_type                    = "OIDC_APP_TYPE_WEB"
  auth_method_type            = "OIDC_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT"
  version                     = "OIDC_VERSION_1_0"
  dev_mode                    = true
  access_token_type           = "OIDC_TOKEN_TYPE_JWT"
  access_token_role_assertion = true
  id_token_role_assertion     = true
  id_token_userinfo_assertion = true
  additional_origins          = []
}

//----------------------------------------
// Project Service Accounts
//----------------------------------------
resource zitadel_machine_user vmclarity_orchestrator_sa {
  org_id            = zitadel_org.openclarity_org.id
  user_name         = "vmclarity-orchestrator-sa@vmclarity.io"
  name              = "VMClarity Orchestrator Service Account"
  access_token_type = "ACCESS_TOKEN_TYPE_JWT"
}

resource zitadel_machine_key vmclarity_orchestrator_sa_key {
  org_id          = zitadel_org.openclarity_org.id
  user_id         = zitadel_machine_user.vmclarity_orchestrator_sa.id
  key_type        = "KEY_TYPE_JSON"
  expiration_date = ""
}

resource zitadel_user_grant vmclarity_orchestrator_sa_user_grant {
  org_id     = zitadel_org.openclarity_org.id
  project_id = zitadel_project.vmclarity_project.id
  role_keys  = ["api:admin"]
  user_id    = zitadel_machine_user.vmclarity_orchestrator_sa.id
  depends_on = [zitadel_project_role.vmclarity_project_role]
}

resource zitadel_machine_user vmclarity_cli_sa {
  org_id            = zitadel_org.openclarity_org.id
  user_name         = "vmclarity-cli-sa@vmclarity.io"
  name              = "VMClarity CLI Service Account"
  access_token_type = "ACCESS_TOKEN_TYPE_JWT"
}

resource zitadel_machine_key vmclarity_cli_sa_key {
  org_id          = zitadel_org.openclarity_org.id
  user_id         = zitadel_machine_user.vmclarity_cli_sa.id
  key_type        = "KEY_TYPE_JSON"
  expiration_date = ""
}

resource zitadel_user_grant vmclarity_cli_sa_user_grant {
  org_id     = zitadel_org.openclarity_org.id
  project_id = zitadel_project.vmclarity_project.id
  role_keys  = ["api:admin"]
  user_id    = zitadel_machine_user.vmclarity_cli_sa.id
  depends_on = [zitadel_project_role.vmclarity_project_role]
}

//----------------------------------------
// Project role action configuration
//----------------------------------------
resource zitadel_action vmclarity_default_role_action {
  org_id          = zitadel_org.openclarity_org.id
  name            = "vmclarity-default-role"
  script          = <<-EOT
  /**
   * Add a usergrant to a new created/registered user
   *
   * Flow: External Authentication, Trigger: Post creation
   *
   * @param ctx
   * @param api
   */
  function addGrant(ctx, api) {
    api.userGrants.push({
      projectID: '${zitadel_project.vmclarity_project.id}',
      roles: ['${var.default_project_role}']
    });
  }
  EOT
  timeout         = "10s"
  allowed_to_fail = false
}

resource zitadel_trigger_actions vmclarity_default_role_trigger_action {
  org_id       = zitadel_org.openclarity_org.id
  flow_type    = "FLOW_TYPE_EXTERNAL_AUTHENTICATION"
  trigger_type = "TRIGGER_TYPE_POST_AUTHENTICATION"
  action_ids   = [zitadel_action.vmclarity_default_role_action.id]
}

//----------------------------------------
// DEBUG - Management resources
//----------------------------------------
resource zitadel_human_user vmclarity_admin {
  org_id             = zitadel_org.openclarity_org.id
  user_name          = "vmclarity-admin"
  first_name         = "firstname"
  last_name          = "lastname"
  email              = "admin@vmclarity.io"
  is_email_verified  = true
  initial_password   = var.debug_initial_admin_password
}

resource zitadel_org_member org_member {
  org_id     = zitadel_org.openclarity_org.id
  user_id    = zitadel_human_user.vmclarity_admin.id
  roles      = ["ORG_OWNER"]
}

resource zitadel_instance_member instance_member {
  user_id    = zitadel_human_user.vmclarity_admin.id
  roles      = ["IAM_OWNER"]
}
