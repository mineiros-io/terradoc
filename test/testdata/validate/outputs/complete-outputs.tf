output "person" {
  value = var.person
}

output "cars" {
  value = var.cars
}

output "beer" {
  value     = var.beer
  sensitive = true
}

resource "something" {
  name = "name"
}

output "number" {
  value = var.number
}
