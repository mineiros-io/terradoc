
variable "clients" {
  description = "(Optional) A list of objects with the clients definitions."
  type        = any

  # Example:
  #
  # clients = [
  #   {
  #     name                                 = "android-app-client"
  #     allowed_oauth_flows                  = ["client_credentials"]
  #     allowed_oauth_flows_user_pool_client = true
  #     allowed_oauth_scopes                 = ["email", "phone"]
  #     callback_urls                        = ["https://mineiros.io/callback", "https://mineiros.io/callback2"]
  #     default_redirect_uri                 = "https://mineiros.io/callback"
  #     explicit_auth_flows                  = ["USER_PASSWORD_AUTH", "ALLOW_USER_PASSWORD_AUTH"]
  #     generate_secret                      = true
  #     logout_urls                          = "https://mineiros.io/logout"
  #     read_attributes                      = ["phone_number", "birthdate", "preferred_username]
  #     refresh_token_validity               = 15
  #     supported_identity_providers         = []
  #     prevent_user_existence_errors        = "LEGACY"
  #     write_attributes                     = ["phone_number"]
  #   }
  # ]

  default = null
}

variable "default_client_allowed_oauth_flows" {
  description = "(Optional) List of allowed OAuth flows. Possible flows are 'code', 'implicit', and 'client_credentials'."
  type        = list(string)
  default     = null
}

variable "default_client_allowed_oauth_flows_user_pool_client" {
  description = "(Optional) Whether the client is allowed to follow the OAuth protocol when interacting with Cognito user pools."
  type        = bool
  default     = true
}

variable "default_client_allowed_oauth_scopes" {
  description = "(Optional) List of allowed OAuth scopes. Possible values are 'phone', 'email', 'openid', 'profile', and 'aws.cognito.signin.user.admin'."
  type        = list(string)
  default     = null
}

variable "default_client_callback_urls" {
  description = "(Optional) List of allowed callback URLs for the identity providers."
  type        = list(string)
  default     = null
}

variable "default_client_default_redirect_uri" {
  description = "(Optional) The default redirect URI. Must be in the list of callback URLs."
  type        = string
  default     = null
}

variable "default_client_explicit_auth_flows" {
  description = "(Optional) List of authentication flows. Possible values are 'ADMIN_NO_SRP_AUTH', 'CUSTOM_AUTH_FLOW_ONLY', 'USER_PASSWORD_AUTH', 'ALLOW_ADMIN_USER_PASSWORD_AUTH', 'ALLOW_CUSTOM_AUTH', 'ALLOW_USER_PASSWORD_AUTH', 'ALLOW_USER_SRP_AUTH', and 'ALLOW_REFRESH_TOKEN_AUTH'."
  type        = list(string)
  default     = null
}

variable "default_client_generate_secret" {
  description = "(Optional) Should an application secret be generated."
  type        = bool
  default     = false
}

variable "default_client_logout_urls" {
  description = "(Optional) List of allowed logout URLs for the identity providers."
  type        = list(string)
  default     = null
}


variable "default_client_read_attributes" {
  description = "(Optional) List of user pool attributes the application client can read from."
  type        = list(string)
  default     = null
}

variable "default_client_refresh_token_validity" {
  description = "(Optional) The time limit in days refresh tokens are valid for."
  type        = number
  default     = 30
}

variable "default_client_prevent_user_existence_errors" {
  description = "(Optional) Choose which errors and responses are returned by Cognito APIs during authentication, account confirmation, and password recovery when the user does not exist in the user pool. When set to 'ENABLED' and the user does not exist, authentication returns an error indicating either the username or password was incorrect, and account confirmation and password recovery return a response indicating a code was sent to a simulated destination. When set to 'LEGACY', those APIs will return a 'UserNotFoundException' exception if the user does not exist in the user pool."
  type        = string
  default     = "LEGACY"
}

variable "default_client_supported_identity_providers" {
  description = "(Optional) List of provider names for the identity providers that are supported on this client."
  type        = list(string)
  default     = null
}

variable "default_client_write_attributes" {
  description = "(Optional) List of user pool attributes the application client can write to."
  type        = list(string)
  default     = null
}
