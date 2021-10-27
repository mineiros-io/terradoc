# Module Argument Reference

See [variables.tf] and [examples/] for details and use-cases.

## Module Configuration

- **`module_enabled`**: *(Optional `bool`)*

  Specifies whether resources in the module will be created.
  Default is `true`.

## Main Resource Configuration

- **`local_secondary_indexes`**: **(Required `list(local_secondary_index)`, Forces new resource)**

  Describe an LSI on the table; these can only be allocated at creation so you cannot change this definition after you have created the resource.
  Default is `[]`.

  `` `terraform
  local_secondary_indexes = [
    {
      range_key = "someKey"
    }
  ]
  `` `

  Each element in the list of `local_secondary_indexes` is an object with the following attributes:

  - *`range_key`*: *(Optional, Forces new resource)*

    The attribute to use as the range (sort) key. Must also be defined as an attribute, see below.
