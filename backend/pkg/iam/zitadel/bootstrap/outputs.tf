resource local_file file_vmclarity_cli_sa_key {
  content  = zitadel_machine_key.vmclarity_cli_sa_key.key_details
  filename = "machinekey/vmclarity-cli-sa-admin.json"
}

resource local_file file_vmclarity_orchestrator_sa_key {
  content  = zitadel_machine_key.vmclarity_orchestrator_sa_key.key_details
  filename = "machinekey/vmclarity-orchestrator-sa-admin.json"
}

resource local_file file_vmclarity_app_api_key {
  content  = zitadel_application_key.vmclarity_app_api_key.key_details
  filename = "machinekey/vmclarity-api-key.json"
}

output vmclarity_project_id {
  value = zitadel_project.vmclarity_project.id
}
