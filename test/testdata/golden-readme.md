# Module Argument Reference

See [variables.tf] and [examples/] for details and use-cases.

## Top-level Arguments

### Module Configuration

- **`module_enabled`**: *(Optional `bool`)*

  Specifies whether resources in the module will be created.

  Default is `true`.

- **`module_depends_on`**: *(Optional `list(dependencies)`)*

  A list of dependencies. Any object can be _assigned_ to this list to define a hidden external dependency.

  Example:

  ```terraform
  module_depends_on = ["google_network.network"]
  ```

### Main Resource Configuration

- **`secret_id`**: *(**Required** `string`)*

  The id of the secret.

- **`members`**: *(Optional `set`)*

  Identities that will be granted the privilege in role. Each entry can have one of the following values:
  - `allUsers`: A special identifier that represents anyone who is on the internet; with or without a Google account.
  - `allAuthenticatedUsers`: A special identifier that represents anyone who is authenticated with a Google account or a service account.
  - `user:{emailid}`: An email address that represents a specific Google account. For example, alice@gmail.com or joe@example.com.
  - `serviceAccount:{emailid}`: An email address that represents a service account. For example, my-other-app@appspot.gserviceaccount.com.
  - `group:{emailid}`: An email address that represents a Google group. For example, admins@example.com.
  - `domain:{domain}`: A G Suite domain (primary, instead of alias) name that represents all the users of that domain. For example, google.com or example.com.


  Default is `[]`.

- **`role`**: *(Optional `string`)*

  The role that should be applied. Note that custom roles must be of the format `[projects|organizations]/{parent-name}/roles/{role-name}`.

- **`project`**: *(Optional `string`)*

  The resource name of the project the policy is attached to. Its format is `projects/{project_id}`.

- **`authoritative`**: *(Optional `bool`)*

  Whether to exclusively set (authoritative mode) or add (non-authoritative/additive mode) members to the role.

  Default is `true`.

- **`policy_bindings`**: *(Optional `list(policy_bindings)`)*

  A list of IAM policy bindings.

  Example:

  ```terraform
  policy_bindings = [{
      members = ["user:member@example.com"]
      role    = "roles/secretmanager.secretAccessor"
    }]
  ```

  `list(policy_bindings)` is a `list` of `any` with the following attributes:

  - **`role`**: *(**Required** `string`)*

    The role that should be applied.

  - **`members`**: *(**Required** `string`)*

    Identities that will be granted the privilege in `role`.

    Default is `"var.members"`.

  - **`condition`**: *(Optional `object(condition)`)*

    An IAM Condition for a given binding.

    Example:

    ```terraform
    condition = {
        expression = "request.time < timestamp(\"2022-01-01T00:00:00Z\")"
        title      = "expires_after_2021_12_31"
      }
    ```

    `object(condition)` is a `any` with the following attributes:

    - **`expression`**: *(**Required** `string`)*

      Textual representation of an expression in Common Expression Language syntax.

    - **`title`**: *(**Required** `string`)*

      A title for the expression, i.e. a short string describing its purpose.

    - **`description`**: *(Optional `string`)*

      An optional description of the expression. This is a longer text which describes the expression, e.g. when hovered over it in a UI.

