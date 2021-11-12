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

### Go/Git configuration for private repositories

While this repository is private, there is some extra work in order to
download and install it using **go install**. There is two main steps.
First you need to configure git to use ssh instead of https:

```
git config --global url.git@github.com:.insteadOf https://github.com/
```

This only needs to be done once and will change the **.gitconfig** on your
host. Then you should always export **GOPRIVATE** for our mineiros-io repos
before running go install:

```
export GOPRIVATE=github.com/mineiros-io
```

More info [here](https://golang.org/ref/mod#private-module-proxy-direct).
