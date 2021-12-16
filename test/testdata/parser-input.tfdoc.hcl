header {
  image = "https://raw.githubusercontent.com/mineiros-io/brand/3bffd30e8bdbbde32c143e2650b2faa55f1df3ea/mineiros-primary-logo.svg"
  url = "https://www.mineiros.io"
}

section {
  title = "root section"
  content = <<END
This is the root section content.

Section contents support anything markdown and allow us to make references like this one: [mineiros-website]
END

  toc = true

  section {
    title = "sections with variables"

    section {
      title = "example"

      variable "person" {
        type = object(person)
        description = "describes the last person who bothered to change this file"

        attribute "name" {
          type        = string
          description = "the person's name"
          default     = "nathan"
        }
      }
    }

    section {
      title = "section of beers"
      content = "an excuse to mention alcohol"

      variable "beers" {
        type        = list(beer)

        description = "a list of beers"
        default     = []

        required = true

        forces_recreation = true

        readme_example = <<END
beers = [
  {
    name = "guinness"
    type = "stout"
    abv  = 4.2
    tags = [
      "dark",
      "irish",
    ]
  }
]
END

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

        attribute "tags" {
          type = list(string)

          description = "a list of tags for the beer"

          default = []
        }
      }
    }
  }
}

references {
  ref "mineiros-website" {
    value = "https://www.mineiros.io"
  }
}
