// A simple static site generator.
package main // import "github.com/qjcg/elx"

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	// For an example of front matter handling, see:
	//   https://github.com/spf13/hugo/blob/master/parser/frontmatter.go
	"github.com/hashicorp/logutils"
	"github.com/microcosm-cc/bluemonday"
	_ "github.com/naoina/toml"
	"gopkg.in/russross/blackfriday.v2"
)

// Layout is a collection of dirs and files representing an elx site.
type Layout struct {
	Dirs  []string
	Files []string
}

const (
	usage = `
elx: A simple static site generator.
John Gosset 2016 (MIT License)

  elx [OPTS] SUBCOMMAND

  SUBCOMMANDS:
	  init [DIR]
	  build [DIR]
	  version
`

	defConfig = `title = "An Elx Static Site"
publisher = "Jerry Q. Hacker"
`
)

var (
	// DefLayout is the default directory layout to be generated via "init"
	// subcommand.
	DefLayout = &Layout{
		Dirs: []string{
			"_site",
			"_templates",
			"_posts",
		},
		Files: []string{
			"config.toml",
		},
	}

	src = "_posts"
	dst = "_site"
)

// Post contains metadata for a blog post.
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
	debug := flag.Bool("d", false, "print debug messages")
	flag.Parse()

	// Configure levelled logging.
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO"},
		MinLevel: logutils.LogLevel("INFO"),
		Writer:   os.Stdout,
	}

	if *debug {
		filter.MinLevel = logutils.LogLevel("DEBUG")
	}

	log.SetOutput(filter)

	basePath := "." // Default value unless overridden by CLI argument.
	args := flag.Args()
	switch flag.NArg() {
	case 0:
		flag.Usage()
		os.Exit(1)

	case 2:
		basePath = args[1]
	}

	// Run subcommand.
	switch args[0] {

	case "init":
		err := Init(basePath, DefLayout)
		if err != nil {
			log.Fatal(err)
		}

	case "build":
		src = filepath.Join(basePath, src)
		dst = filepath.Join(basePath, dst)
		err := Build(src, dst)
		if err != nil {
			log.Fatal(err)
		}
	case "version":
		fmt.Println(Version)
	default:
		flag.Usage()
		os.Exit(1)
	}
}

// Init creates a new directory with the specified layout.
func Init(basepath string, layout *Layout) error {

	// Create Layout's directories.
	for _, d := range layout.Dirs {
		path := filepath.Join(basepath, d)
		err := os.MkdirAll(path, 0775)
		if err != nil {
			return err
		}
		log.Printf("[DEBUG] Created directory: %s\n", path)
	}

	// Create Layout's files.
	for _, f := range layout.Files {
		path := filepath.Join(basepath, f)
		err := ioutil.WriteFile(path, []byte(defConfig), 0644)
		if err != nil {
			log.Println(err)
		}
		log.Printf("[DEBUG] Created file: %s\n", path)
	}

	return nil
}

// Build renders the markdown provided in srcdir as HTML.
func Build(srcdir, dstdir string) error {

	// Get a slice of markdown files.
	matches, err := filepath.Glob(srcdir + "/*.md")
	if err != nil {
		return err
	}

	// Convert markdown to HTML.
	for _, md := range matches {
		dstFile := strings.Replace(filepath.Base(md), ".md", ".html", 1)
		log.Printf("[DEBUG] markdown file src: %s\n", md)
		log.Printf("[DEBUG] markdown file dst: %s/%s\n", dstdir, dstFile)

		dat, err := ioutil.ReadFile(md)
		if err != nil {
			log.Printf("[INFO] Error reading .md file: %s\n", err)
		}

		// TODO: Write full HTML webpage here, don't just
		err = ioutil.WriteFile(filepath.Join(dstdir, dstFile), toHTML(dat), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func toHTML(input []byte) []byte {
	unsafe := blackfriday.Run(input)
	return bluemonday.UGCPolicy().SanitizeBytes(unsafe)
}
