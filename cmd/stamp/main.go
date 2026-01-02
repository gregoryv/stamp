// Command stamp generates go source code with build information.
//
//go:generate stamp -go build_stamp.go -clfile ../../changelog.md
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gregoryv/stamp"
)

func main() {
	out := ""
	clfile := "CHANGELOG.md"
	help := false

	s := stamp.NewStamp()
	flag.StringVar(&s.Package, "p", s.Package, "Package for the output source")
	flag.BoolVar(&help, "h", help, "Print this help and exit")
	flag.StringVar(&out, "go", out, "Write Go file, defaults to stdout")
	flag.StringVar(&clfile, "clfile", clfile,
		"Changelog to parse for version, keepachangelog format")
	stamp.InitFlags()
	flag.Parse()
	stamp.AsFlagged()
	if help {
		usage()
	}

	stderr := log.New(os.Stderr, "", 0)
	var err error
	if s.Revision, err = stamp.Revision("."); err != nil {
		stderr.Fatalf("Failed to generate build: %s", err)
	}
	s.ParseChangelog(clfile)
	if err = writeGoSource(s, out); err != nil {
		stderr.Fatalf("Failed to write go source: %s", err)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(0)
}

func writeGoSource(s *stamp.Stamp, file string) (err error) {
	// Write go code
	fh := os.Stdout
	if file != "" {
		if fh, err = os.Create(file); err != nil {
			return
		}
		defer fh.Close()
	}
	if err = stamp.NewGoTemplate().Execute(fh, s); err != nil {
		return
	}
	return
}
