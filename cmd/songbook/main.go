package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bt108/vaisnava-songbook-generator/internal/generator"
)

func main() {
	dataDir := flag.String("data", "data", "path to the YAML data directory")
	templateDir := flag.String("templates", "templates", "path to the Go template directory")
	outDir := flag.String("out", "out", "path to the generated LaTeX directory")
	flag.Parse()

	if err := generator.Generate(*dataDir, *templateDir, *outDir, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
