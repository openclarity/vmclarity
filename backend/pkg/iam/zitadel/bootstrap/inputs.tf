variable "idp_github_client_id" {
  type = string
}

variable "idp_github_client_secret" {
  type = string
}

variable "project_roles" {
  type    = set(string)
  default = ["api:admin", "api:writer", "api:reader"]
}

variable "default_project_role" {
  type = string
  default = "api:admin"
}


//--------------------------------------------------------------------------------------
// DEBUG - Management resources
//--------------------------------------------------------------------------------------
variable "debug_initial_admin_password" {
  type = string
}
