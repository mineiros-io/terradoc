section {
  title       = "root section"
  description = "i am the root section"

  section {
    title = "sub section with no description"

    variable "name" {
      type        = string
      description = "describes the name of the last person who bothered to change this file"
      default     = "nathan"
    }
  }

  section {
    title = "section of beers"
    description = "an excuse to mention alcohol"

    variable "beers" {
      type        = any
      readme_type = "list(beer)"

      description = "a list of beers"
      default     = []

      required = true

      forces_recreation = true

      readme_example = {
        beers = [
          {
            name = "guinness"
            type = "stout"
            abv = 4.2
          }
        ]
      }

      attribute "name" {
        type = string

        description = "the name of the beer"
      }

      attribute "type" {
        type = string

        description = "the type of the beer"

        forces_recreation = true
      }

      attribute "abv" {
        type = number

        description = "beer's alcohol by volume content"

        forces_recreation = true
      }
    }
  }
}
