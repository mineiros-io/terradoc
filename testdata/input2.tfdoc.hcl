section {
  title       = "Module Argument Reference"
  description = "See [variables.tf] and [examples/] for details and use-cases. Section level 1"

  section {
    title = "Module Configuration"
    description = "Section level 2"

    variable "module_enabled" {
      type        = bool
      description = "Specifies whether resources in the module will be created."
      default     = true
    }

    section {
      title = "Sub section to test this shit up!"
      description = "It must have descriptions as well. Section level 3"

      section {
        title = "will it work with even another section?"
        description = "Section level 4"

        section {
          title = "let's see!"
          description = "this is the 5th level"

          variable "foo" {
            type = string
            description = "just a regular string"
            default = "i am the default value!"
            forces_recreation = true
            required = true
          }
        }
      }
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
            range_key = {
              a = "foo"
              b = "bar"
            }
          }
        ]
      }

      attribute "range_key" {
        type = object({a=string, b=string})

        description = "The attribute to use as the range (sort) key. Must also be defined as an attribute, see below."
        forces_recreation = true

        attribute "a" {
          type = string
          description = "a string"
        }

        attribute "b" {
          type = string
          description = "another string"
        }
      }
    }
  }
}
