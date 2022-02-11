section {
  section {
    variable "person" {
      type = object(person)
    }
  }

  variable "beer" {
    type = string
  }

  section {
    section {
      variable "number" {
        type = number
      }
    }
  }
}
