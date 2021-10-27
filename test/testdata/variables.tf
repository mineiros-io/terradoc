variable "module_enabled" {
  type        = bool
  description = "(Optional) Specifies whether resources in the module will be created."
  default     = true
}

variable "local_secondary_indexes" {
  type        = any
  description = "Describe an LSI on the table; these can only be allocated creation so you cannot change this definition after you have created the resource."
  default     = []
}
