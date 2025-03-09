package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/lkebin/repository/generator"
)

var (
	typeNames = flag.String("type", "", "type names, comma-separated")
	output    = flag.String("output", "", "output file name; default srcdir/<type>_impl.go")
	buildTags = flag.String("tags", "", "comma-separated list of build tags to apply")
)

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of repository:\n")
	fmt.Fprintf(os.Stderr, "\trepository [flags] -type T [directory]\n")
	fmt.Fprintf(os.Stderr, "\trepository [flags] -type T files... # Must be a single package\n")
	fmt.Fprintf(os.Stderr, "For more information, see:\n")
	fmt.Fprintf(os.Stderr, "\thttps://github.com/lkebin/repository\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

// isDirectory reports whether the named file is a directory.
func isDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}

// baseName that will put the generated code together with pkg.
func baseName(typename string) string {
	suffix := "impl.go"
	return fmt.Sprintf("%s_%s", generator.ToSnakeCase(typename), suffix)
}

func main() {
	log.SetFlags(0)
	flag.Usage = Usage
	flag.Parse()
	if len(*typeNames) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	types := strings.Split(*typeNames, ",")
	var tags []string
	if len(*buildTags) > 0 {
		tags = strings.Split(*buildTags, ",")
	}

	args := flag.Args()
	if len(args) == 0 {
		// Default: process whole package in current directory.
		args = []string{"."}
	}

	// Parse the package once.
	var dir string
	// TODO(suzmue): accept other patterns for packages (directories, list of files, import paths, etc).
	if len(args) == 1 && isDirectory(args[0]) {
		dir = args[0]
	} else {
		if len(tags) != 0 {
			log.Fatal("-tags option applies only to directories, not when files are specified")
		}
		dir = filepath.Dir(args[0])
	}

	specs := generator.ParseRepository(types, args, tags)
	if len(specs) == 0 {
		log.Fatalf("no types found in %s", args)
	}

	for _, spec := range specs {
		src, err := generator.GenerateRepositoryImplements(&spec)
		if err != nil {
			log.Fatalf("generating implements: %v", err)
		}

		// Write to file.
		outputName := *output
		if outputName == "" {
			outputName = filepath.Join(dir, baseName(spec.Name))
		}

		if err = os.WriteFile(outputName, src, 0644); err != nil {
			log.Fatalf("writing output: %s", err)
		}
	}
}
