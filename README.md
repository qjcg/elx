# Electrostatic: a simple static site generator


## Features

- Converts directory with markdown content to static HTML
- Uses Go Templates
- Default template handles photo galleries and text articles
- Configure site with TOML
- Optional minification of web content

## Anti-Features

- Lots of subcommands
- Lots of options
- Lots of cornercases
- Lots of markup formats


## Directory structure

(WIP, try and create something simple, inspired by below refs)

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
$ elx init

$ elx build

$ elx watch
```
