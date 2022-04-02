variable "person" {
  type = any
}

variable "cars" {
  type = any
}

variable "beer" {
  type      = string
  sensitive = true
}

resource "something" {
  name = "name"
}

variable "number" {
  type = number
}
