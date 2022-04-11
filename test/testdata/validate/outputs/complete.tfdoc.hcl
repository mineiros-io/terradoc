section {
  section {
    output "person" {
      type = object(person)
    }

    output "cars" {
      type = list(car)
    }
  }

  output "beer" {
    type = string
  }

  section {
    section {
      output "number" {
        type = number
      }
    }
  }
}
