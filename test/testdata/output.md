# root section

i am the root section

## sub section with no description

- **`name`**: *(Optional `string`)*

  describes the name of the last person who bothered to change this file

  Default is `"nathan"`.

## section of beers

an excuse to mention alcohol

- **`beers`**: *(**Required** `list(beer)`, Forces new resource)*

  a list of beers

  Default is `[]`.

  Example:

  ```terraform
  beers = [{
      abv  = 4.2
      name = "guinness"
      type = "stout"
    }]
  ```

  `list(beer)` is a `list` of `any` with the following attributes:

  - **`name`**: *(Optional `string`)*

    the name of the beer

  - **`type`**: *(Optional `string`, Forces new resource)*

    the type of the beer

  - **`abv`**: *(Optional `number`, Forces new resource)*

    beer's alcohol by volume content

