package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/MontFerret/montferret.github.io/internal/versions"
)

const defaultVersionFile = "data/versions.yaml"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(arguments []string, stdout, stderr io.Writer) int {
	if len(arguments) == 0 {
		fmt.Fprintln(stderr, "usage: versions update --repository <owner/name> --tag <version> [--file <path>] [--format text|json]")

		return 2
	}

	if arguments[0] != "update" {
		fmt.Fprintf(stderr, "unknown command %q\n", arguments[0])

		return 2
	}

	flags := flag.NewFlagSet("versions update", flag.ContinueOnError)
	flags.SetOutput(stderr)
	repository := flags.String("repository", "", "source repository in owner/name form")
	tag := flags.String("tag", "", "release SemVer tag")
	file := flags.String("file", defaultVersionFile, "website versions YAML file")
	format := flags.String("format", "text", "output format: text or json")

	if err := flags.Parse(arguments[1:]); err != nil {
		return 2
	}

	if flags.NArg() != 0 {
		fmt.Fprintf(stderr, "unexpected positional arguments: %v\n", flags.Args())

		return 2
	}

	if *repository == "" {
		fmt.Fprintln(stderr, "--repository is required")

		return 2
	}

	if *tag == "" {
		fmt.Fprintln(stderr, "--tag is required")

		return 2
	}

	if *format != "text" && *format != "json" {
		fmt.Fprintf(stderr, "unsupported output format %q\n", *format)

		return 2
	}

	result, err := versions.UpdateFile(versions.Request{
		Repository: *repository,
		Tag:        *tag,
		File:       *file,
	})
	if err != nil {
		fmt.Fprintf(stderr, "versions update: %v\n", err)
		return 1
	}

	if result.Status == versions.StatusStale {
		fmt.Fprintf(stderr, "warning: ignored stale release %s for %s; current version is %s\n", result.Version, result.Path, result.Previous)
	}

	if *format == "json" {
		if err := json.NewEncoder(stdout).Encode(result); err != nil {
			fmt.Fprintf(stderr, "write JSON result: %v\n", err)

			return 1
		}

		return 0
	}

	fmt.Fprintf(stdout, "%s %s: %s -> %s\n", result.Status, result.Path, result.Previous, result.Version)

	return 0
}
