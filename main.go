// Electrostatic: a simple static site generator.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/microcosm-cc/bluemonday"
	// For an example of front matter handling, see:
	//   https://github.com/spf13/hugo/blob/master/parser/frontmatter.go
	_ "github.com/naoina/toml"
	"github.com/russross/blackfriday"
)

type Layout struct {
	Dirs  []string
	Files []string
}

const (
	usage = `
elx: A simple static site generator.
John Gosset 2016 (MIT License)

  elx init [DIR]
  elx build [DIR]
`

	defConfig = `title = "An Elx Static Site"
publisher = "Jerry Q. Hacker"
`
)

var (
	// directory layout to be generated via "init" subcommand
	DefLayout = &Layout{
		Dirs: []string{
			"_site",
			"_includes",
			"_layouts",
			"_posts",
		},
		Files: []string{
			"config.toml",
		},
	}

	src string = "_posts"
	dst string = "_site"
)

type Post struct {
	Title         string
	TimePublished string
	TimeUpdated   string
}

func main() {
	flag.Usage = func() {
		fmt.Println(usage)
		flag.PrintDefaults()
		fmt.Println()
	}
	basePath := flag.String("b", ".", "base path")
	flag.Parse()

	args := flag.Args()

	switch flag.NArg() {
	case 0:
		flag.Usage()
		os.Exit(1)
	case 2:
		basePath = &args[1]
	}

	switch args[0] {
	// init [DIR]
	case "init":
		err := Init(*basePath, DefLayout)
		if err != nil {
			log.Println(err)
		}
	// build [DIR]
	case "build":
		src = filepath.Join(*basePath, src)
		dst = filepath.Join(*basePath, dst)
		err := Build(src, dst)
		if err != nil {
			log.Println(err)
		}
	default:
		flag.Usage()
		os.Exit(1)
	}
}

// Init creates a new directory with the standard layout.
func Init(basepath string, layout *Layout) error {
	defer fmt.Println("Initialized elx directory:", basepath)
	for _, d := range layout.Dirs {
		path := filepath.Join(basepath, d)
		os.MkdirAll(path, 0775)
	}

	for _, f := range layout.Files {
		path := filepath.Join(basepath, f)
		err := ioutil.WriteFile(path, []byte(defConfig), 0644)
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}

// Build renders the markdown provided in srcdir as HTML.
func Build(srcdir, dstdir string) error {
	matches, err := filepath.Glob(srcdir + "/*.md")
	if err != nil {
		return err
	}

	// convert markdown to HTML
	for _, md := range matches {
		dat, err := ioutil.ReadFile(md)
		if err != nil {
			log.Println(err)
		}
		fmt.Fprintf(os.Stdout, "%s\n", toHTML(dat))
	}
	return nil
}

func toHTML(input []byte) []byte {
	unsafe := blackfriday.MarkdownCommon(input)
	return bluemonday.UGCPolicy().SanitizeBytes(unsafe)
}
