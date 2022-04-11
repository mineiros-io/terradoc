section {
  output "person" {
    type = object(person)
  }

  output "number" {
    type = number
  }

  output "cars" {
    type = list(car)
  }
}
