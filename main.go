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

var (
	// directory layout to be generated via "init" subcommand
	DirLayout = []string{
		"aaa",
		"bbb",
		"ccc",
	}
)

type Post struct {
	Title         string
	TimePublished string
	TimeUpdated   string
}

func main() {
	src := flag.String("d", "src", "source directory")
	flag.Parse()

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

func toHTML(input []byte) []byte {
	unsafe := blackfriday.MarkdownCommon(input)
	return bluemonday.UGCPolicy().SanitizeBytes(unsafe)
}
