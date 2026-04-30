package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/envoy-diff/internal/differ"
	"github.com/envoy-diff/internal/loader"
)

func main() {
	var (
		baseFile   = flag.String("base", "", "path to base snapshot JSON file (required)")
		headFile   = flag.String("head", "", "path to head snapshot JSON file (required)")
		format     = flag.String("format", "text", "output format: text or json")
		showSame   = flag.Bool("show-unchanged", false, "include unchanged resources in output")
	)
	flag.Parse()

	if *baseFile == "" || *headFile == "" {
		fmt.Fprintln(os.Stderr, "error: --base and --head are required")
		flag.Usage()
		os.Exit(1)
	}

	base, err := loader.FromFile(*baseFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading base: %v\n", err)
		os.Exit(1)
	}

	head, err := loader.FromFile(*headFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading head: %v\n", err)
		os.Exit(1)
	}

	results := differ.Compare(base, head)

	out, err := differ.Render(results, *format, *showSame)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error rendering output: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(out)
}
