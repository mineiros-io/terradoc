section {
  variable "person" {
    type = object(person)
  }

  variable "number" {
    type = number
  }

  variable "cars" {
    type = list(car)
  }
}
