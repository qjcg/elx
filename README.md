# Electrostatic

A simple static site generator.

**NOTE: This is pre-alpha software. Several
features are not yet fully implemented.**


## Features

- Converts directory with markdown content to static HTML
- Templates for:
    - blogs
    - photo galleries
    - presentation slides
- TOML-based configuration
- Optional minification of web content

## Non-Features

- Lots of subcommands
- Lots of options
- Lots of cornercases
- Lots of markup formats


## Directory structure

```
mysite
├── assets
│   ├── css
│   ├── img
│   └── js
├── _includes
├── _layouts
├── _posts
└── _site
```

### References

- https://jekyllrb.com/docs/structure/
- https://gohugo.io/overview/source-directory/
- https://middlemanapp.com/basics/directory-structure/


## Usage

```shell
# Create a new site in ./potatoes.com.
$ elx init potatoes.com

# Create a new site in the current directory.
$ elx init

# Build
$ elx build
```

## Licence

MIT.
