<!-- mdtocstart -->

# Table of Contents

- [terradoc](#terradoc)
    - [Installing](#installing)
        - [Go/Git configuration for private Repositories](#gogit-configuration-for-private-repositories)

<!-- mdtocend -->

# terradoc

Terradoc is a lean helper tool that creates human readable documentation from
HCL syntax.

## Installing

To install **terradoc** using Go just run:

```
go install github.com/mineiros-io/terradoc/cmd/terradoc@<version>
```

Where **<version>** is any terradoc [version tag](https://github.com/mineiros-io/terradock/tags),
or if you are feeling adventurous you can just install **latest**:

```
go install github.com/mineiros-io/terradoc/cmd/terradoc@latest
```

We put great effort into keeping the main branch stable, so it should be safe
to use **latest** to play around, but not recommended for long term automation
since you won't get the same build result each time you run the install command.
