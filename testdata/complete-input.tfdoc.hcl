section {
  title = "Section #1"
  description = "this is the first section"

  section {
    title = "Section #1.1"
    description = "first subsection"
  }

  section {
    title = "Section #1.2"
    description = "second subsection"
  }
}

section {
  title = "Section #2"
  description = "just to test multiple sections (2)"
}

section {
  title = "Variable types"
  description = "Test all variable types!"

  section {
    title = "Basic"
    description = "Variable definition with required attributes"

    variable "string-var" {
      type = string
      description = "a string"
    }

    variable "number-var" {
      type = number
      description = "a number"
    }

    variable "bool-var" {
      type = bool
      description = "a bool"
    }

    variable "list-or-tuple-var" {
      type = list(string)
      description = "a list (or tuple)"
    }

    variable "tuple-or-list-var" {
      type = tuple([number])
      description = "a tuple (or list)"
    }

    variable "set-var" {
      type = set(number)
      description = "a set"
    }

    variable "object-or-map-var" {
      type = object({age=number, name=string})
      description = "an object (or map)"
    }

    variable "map-or-object-var" {
      type = map(string)
      description = "a map (or object)"
    }
  }

  section {
    title = "Required"
    description = "Required variable definitions"

    variable "string-var" {
      type = string
      description = "a string"
      required = true
    }

    variable "number-var" {
      type = number
      description = "a number"
      required = true
    }

    variable "bool-var" {
      type = bool
      description = "a bool"
      required = true
    }

    variable "list-or-tuple-var" {
      type = list(string)
      description = "a list (or tuple)"
      required = true
    }

    variable "tuple-or-list-var" {
      type = tuple([number])
      description = "a tuple (or list)"
      required = true
    }

    variable "set-var" {
      type = set(number)
      description = "a set"
      required = true
    }

    variable "object-or-map-var" {
      type = object({age=number, name=string})
      description = "an object (or map)"
      required = true
    }

    variable "map-or-object-var" {
      type = map(string)
      description = "a map (or object)"
      required = true
    }
  }

  section {
    title = "With defaults"
    description = "Required variable definitions with default"

    variable "string-var" {
      type = string
      description = "a string"
      required = true
      default = "i am a string"
    }

    variable "number-var" {
      type = number
      description = "a number"
      required = true
      default = 30
    }

    variable "bool-var" {
      type = bool
      description = "a bool"
      required = true
      default = true
    }

    variable "list-or-tuple-var" {
      type = list(string)
      description = "a list (or tuple)"
      required = true
      default = ["a", "b", "c"]
    }

    variable "tuple-or-list-var" {
      type = tuple([number])
      description = "a tuple (or list)"
      required = true
      default = [1,2,3,4]
    }

    variable "set-var" {
      type = set(number)
      description = "a set"
      required = true
      default = [1,2]
    }

    variable "object-or-map-var" {
      type = object({age=number, name=string})
      description = "an object (or map)"
      required = true
      default = {"age": 30, name: "Nathan"}
    }

    variable "map-or-object-var" {
      type = map(string)
      description = "a map (or object)"
      required = true
      default = {"age": 30, name: "Nathan"}
    }
  }

  section {
    title = "Forces recreation"

    variable "string-var" {
      type = string
      description = "a string"
      required = true
      default = "i am a string"
      forces_recreation = true
    }

    variable "number-var" {
      type = number
      description = "a number"
      required = true
      default = 30
      forces_recreation = true
    }

    variable "bool-var" {
      type = bool
      description = "a bool"
      required = true
      default = true
      forces_recreation = true
    }

    variable "list-or-tuple-var" {
      type = list(string)
      description = "a list (or tuple)"
      required = true
      default = ["a", "b", "c"]
      forces_recreation = true
    }

    variable "tuple-or-list-var" {
      type = tuple([number])
      description = "a tuple (or list)"
      required = true
      default = [1,2,3,4]
      forces_recreation = true
    }

    variable "set-var" {
      type = set(number)
      description = "a set"
      required = true
      default = [1,2]
      forces_recreation = true
    }

    variable "object-or-map-var" {
      type = object({age=number, name=string})
      description = "an object (or map)"
      required = true
      default = {"age": 30, name: "Nathan"}
      forces_recreation = true
    }

    variable "map-or-object-var" {
      type = map(string)
      description = "a map (or object)"
      required = true
      default = {"age": 30, name: "Nathan"}
      forces_recreation = true
    }
  }

  section {
    title = "Readme Example and Readme Type"

    variable "complex_variable_list" {
      type = any
      readme_type = "list(complex_variable)"
      default = []
      description = "A complex variable"
      required = true
      readme_example = {
        complex_variable_list = [
          {
            foo = "bar",
            bar = 5
          }
        ]
      }

      attribute "foo" {
        type = string
        description = "a string"
        forces_recreation = false
      }

      attribute "bar" {
        type = number
        description = "a number"
        forces_recreation = true
      }
    }
  }
}
