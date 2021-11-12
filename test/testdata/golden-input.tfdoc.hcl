section {
  title = "Module Argument Reference"
  description = "See [variables.tf] and [examples/] for details and use-cases."

  section {
    title = "Top-level Arguments"

    section {
      title = "Module Configuration"

      variable "module_enabled" {
        type = bool
        description = "Specifies whether resources in the module will be created."
        default = true
      }

      variable "module_depends_on" {
        type = any
        readme_type = "list(dependencies)"
        description = "A list of dependencies. Any object can be _assigned_ to this list to define a hidden external dependency."
        readme_example = {
          module_depends_on = [
            "google_network.network"
          ]
        }
      }
    }

    section {
      title = "Main Resource Configuration"

      variable "secret_id" {
        required = true
        type = string
        description = "The id of the secret."
      }

      variable "members" {
        type = set(string)
        default = []
        description = <<END
Identities that will be granted the privilege in role. Each entry can have one of the following values:
  - `allUsers`: A special identifier that represents anyone who is on the internet; with or without a Google account.
  - `allAuthenticatedUsers`: A special identifier that represents anyone who is authenticated with a Google account or a service account.
  - `user:{emailid}`: An email address that represents a specific Google account. For example, alice@gmail.com or joe@example.com.
  - `serviceAccount:{emailid}`: An email address that represents a service account. For example, my-other-app@appspot.gserviceaccount.com.
  - `group:{emailid}`: An email address that represents a Google group. For example, admins@example.com.
  - `domain:{domain}`: A G Suite domain (primary, instead of alias) name that represents all the users of that domain. For example, google.com or example.com.
END
      }

      variable "role" {
        type = string
        description = "The role that should be applied. Note that custom roles must be of the format `[projects|organizations]/{parent-name}/roles/{role-name}`."
      }

      variable "project" {
        type = string
        description = "The resource name of the project the policy is attached to. Its format is `projects/{project_id}`."
      }

      variable "authoritative" {
        description = "Whether to exclusively set (authoritative mode) or add (non-authoritative/additive mode) members to the role."
        type = bool
        default = true
      }

      variable "policy_bindings" {
        type = list(any)
        readme_type = "list(policy_bindings)"
        description = "A list of IAM policy bindings."
        readme_example = {
          policy_bindings = [{
            role    = "roles/secretmanager.secretAccessor"
            members = ["user:member@example.com"]
          }]
        }


        attribute "role" {
          description = "The role that should be applied."
          required = true
          type = string
        }

        attribute "members" {
          required = true
          type = string
          default = "var.members"
          description = "Identities that will be granted the privilege in `role`."
        }

        attribute "condition" {
          type = any
          readme_type = "object(condition)"
          description = "An IAM Condition for a given binding."
          readme_example = {
            condition = {
              expression = "request.time < timestamp(\"2022-01-01T00:00:00Z\")"
              title      = "expires_after_2021_12_31"
            }
          }

          attribute "expression" {
            type = string
            required = true
            description = "Textual representation of an expression in Common Expression Language syntax."
          }

          attribute "title" {
            type = string
            required = true
            description = "A title for the expression, i.e. a short string describing its purpose."
          }

          attribute "description" {
            type = string
            description = "An optional description of the expression. This is a longer text which describes the expression, e.g. when hovered over it in a UI."
          }
        }
      }
    }
  }
}
