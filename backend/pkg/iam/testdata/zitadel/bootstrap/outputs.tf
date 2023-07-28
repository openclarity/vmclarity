resource local_file file_vmclarity_data {
  content  = <<-EOT
    {
      "project_id": "${zitadel_project.vmclarity_project.id}",
      "app_client_id": "${zitadel_application_api.vmclarity_app_api.client_id}",
      "app_client_secret": "${zitadel_application_api.vmclarity_app_api.client_secret}",
      "orchestrator_pam": "${zitadel_personal_access_token.vmclarity_orchestrator_pat_key.token}",
      "cli_pam": "${zitadel_personal_access_token.vmclarity_cli_pat_key.token}"
    }
  EOT
  filename = "generated/vmclarity-data.json"
}
