section {
  section {
    variable "person" {
      type = object(person)
    }

    variable "cars" {
      type = list(car)
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
