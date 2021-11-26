[<img src="https://raw.githubusercontent.com/mineiros-io/brand/3bffd30e8bdbbde32c143e2650b2faa55f1df3ea/mineiros-primary-logo.svg" width="400"/>](https://www.mineiros.io)



# root section

This is the root section content.

Section contents support anything markdown and allow us to make references like this one: [mineiros-website]

## sections with variables

### example

- **`name`**: *(Optional `string`)*

  describes the name of the last person who bothered to change this file

  Default is `"nathan"`.

### section of beers

an excuse to mention alcohol

- **`beers`**: *(**Required** `list(beer)`, Forces new resource)*

  a list of beers

  Default is `[]`.

  Example:

  ```hcl
  beers = [
    {
      name = "guinness"
      type = "stout"
      abv  = 4.2
    }
  ]
  ```

  Each object in the list accepts the following attributes:

  - **`name`**: *(Optional `string`)*

    the name of the beer

  - **`type`**: *(Optional `string`, Forces new resource)*

    the type of the beer

  - **`abv`**: *(Optional `number`, Forces new resource)*

    beer's alcohol by volume content


<!-- References -->

[mineiros-website]: https://www.mineiros.io

