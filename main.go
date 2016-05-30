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
	usage = `elx: A simple static site generator.
John Gosset 2016 (MIT License)

elx init [<DIR>]
elx build [<DIR>]
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
)

type Post struct {
	Title         string
	TimePublished string
	TimeUpdated   string
}

func main() {
	flag.Usage = func() {
		fmt.Println(usage)
	}
	basePath := flag.String("b", ".", "base path")
	src := flag.String("d", "src", "source directory")
	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	switch args[0] {
	case "init":
		if len(args) == 2 {
			basePath = &args[1]
		}
		err := initDirLayout(*basePath, DefLayout)
		if err != nil {
			log.Println(err)
		}
	}

	// find markdown files
	matches, err := filepath.Glob(*src + "/*.md")
	if err != nil {
		log.Fatal(err)
	}

	// convert markdown to HTML
	for _, m := range matches {
		dat, err := ioutil.ReadFile(m)
		if err != nil {
			log.Println(err)
		}
		fmt.Fprintf(os.Stdout, "%s\n", toHTML(dat))
	}
}

func initDirLayout(basepath string, layout *Layout) error {
	defer fmt.Println("Initialized elx directory:", basepath)
	for _, d := range layout.Dirs {
		path := filepath.Join(basepath, d)
		os.MkdirAll(path, 0775)
	}

	for _, f := range layout.Files {
		path := filepath.Join(basepath, f)
		err := ioutil.WriteFile(path, []byte{}, 0644)
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}

func toHTML(input []byte) []byte {
	unsafe := blackfriday.MarkdownCommon(input)
	return bluemonday.UGCPolicy().SanitizeBytes(unsafe)
}
