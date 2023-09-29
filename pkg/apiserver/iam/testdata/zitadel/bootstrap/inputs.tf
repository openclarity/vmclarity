variable "idp_github_client_id" {
  type    = string
  default = "123"
}

variable "idp_github_client_secret" {
  type    = string
  default = "123"
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
  type    = string
  default = "Password1!"
}

variable "debug_initial_admin_username" {
  type = string
  default = "admin@vmclarity.io"
}