section {
  title       = "Module Argument Reference"
  description = "See [variables.tf] and [examples/] for details and use-cases."

  section {
    title = "Module Configuration"

    variable "module_enabled" {
      type        = bool
      description = "Specifies whether resources in the module will be created."
      default     = true
    }
  }

  section {
    title = "Main Resource Configuration"

    variable "local_secondary_indexes" {
      type        = any
      readme_type = "list(local_secondary_index)"

      description = "Describe an LSI on the table; these can only be allocated creation so you cannot change this definition after you have created the resource."
      default     = []

      required = true

      forces_recreation = true

      readme_example = {
        local_secondary_indexes = [
          {
            range_key = "someKey"
          }
        ]
      }

      attribute "range_key" {
        type = string

        description = "The attribute to use as the range (sort) key. Must also be defined as an attribute, see below."

        forces_recreation = true
      }
    }
  }
}
