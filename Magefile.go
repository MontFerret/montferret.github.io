//go:build mage
// +build mage

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/MontFerret/montferret.github.io/tools/registryroutes"
	"github.com/MontFerret/montferret.github.io/tools/stdlibdocs"
	"gopkg.in/yaml.v2"

	"github.com/magefile/mage/sh"
)

const OUTPUT_DIR = "public"
const THEME_DIR = "themes/ferret"
const GENERATED_DIR = ".generated"
const STDLIB_DOCS_DIR = ".generated/content/docs/standard-library"
const VERSIONS_DATA = "data/versions.yaml"
const STDLIB_SEARCH_CHECK = "scripts/verify-stdlib-search.mjs"
const REGISTRY_BASE_URL = "https://registry.ferretlang.org/"

type versionsData struct {
	Runtime struct {
		V2 string `yaml:"v2"`
	} `yaml:"runtime"`
}

func generateSite() error {
	if err := sh.RunV("hugo"); err != nil {
		return err
	}

	if err := sh.RunV("npm", "--prefix", THEME_DIR, "run", "build:search"); err != nil {
		return err
	}

	if err := verifySearchIndex(); err != nil {
		return err
	}
	if err := sh.RunV("node", STDLIB_SEARCH_CHECK); err != nil {
		return err
	}

	return generateRegistryRoutes()
}

func generateRegistryRoutes() error {
	client := &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(request *http.Request, via []*http.Request) error {
			if len(via) > 0 && (!strings.EqualFold(request.URL.Scheme, via[0].URL.Scheme) || !strings.EqualFold(request.URL.Host, via[0].URL.Host)) {
				return fmt.Errorf("Registry redirect changed origin")
			}
			return nil
		},
	}

	return registryroutes.Generate(
		context.Background(),
		client,
		REGISTRY_BASE_URL,
		filepath.Join(OUTPUT_DIR, "registry", "index.html"),
		OUTPUT_DIR,
	)
}

func verifySearchIndex() error {
	pagefindDir := filepath.Join(OUTPUT_DIR, "pagefind")

	if _, err := os.Stat(pagefindDir); err != nil {
		return fmt.Errorf("missing Pagefind output: %s", pagefindDir)
	}

	return nil
}

var Default = Serve

// Cleans up build directory
func Clean() error {
	if err := os.RemoveAll(OUTPUT_DIR); err != nil {
		return err
	}

	return os.RemoveAll(GENERATED_DIR)
}

func prepareBuild() error {
	if err := Clean(); err != nil {
		return err
	}

	return GenerateStdlib()
}

// Starts local Hugo server
func Serve() error {
	if err := prepareBuild(); err != nil {
		return err
	}

	return sh.RunV("hugo", "server")
}

// Runs the production Hugo build and generates the search index
func Build() error {
	if err := prepareBuild(); err != nil {
		return err
	}

	return generateSite()
}

// Builds the website search index and serves the generated static site
func ServeSearch() error {
	if err := prepareBuild(); err != nil {
		return err
	}

	if err := sh.RunV("hugo", "--baseURL", "http://localhost:1414/"); err != nil {
		return err
	}

	if err := sh.RunV("npm", "--prefix", THEME_DIR, "run", "build:search"); err != nil {
		return err
	}

	if err := verifySearchIndex(); err != nil {
		return err
	}
	if err := sh.RunV("node", STDLIB_SEARCH_CHECK); err != nil {
		return err
	}

	if err := generateRegistryRoutes(); err != nil {
		return err
	}

	return sh.RunV("npm", "--prefix", THEME_DIR, "run", "serve:search")
}

// Installs theme
func Install() error {
	os.Chdir(THEME_DIR)

	defer os.Chdir("../..")

	return sh.RunV("npm", "ci")
}

// Generates all derived documentation.
func Generate() error {
	return GenerateStdlib()
}

// GenerateStdlib refreshes the unversioned Standard Library reference from the
// exact Ferret runtime version configured in data/versions.yaml.
func GenerateStdlib() error {
	data, err := os.ReadFile(VERSIONS_DATA)
	if err != nil {
		return fmt.Errorf("read Ferret version configuration: %w", err)
	}
	versions := versionsData{}
	if err := yaml.Unmarshal(data, &versions); err != nil {
		return fmt.Errorf("parse Ferret version configuration: %w", err)
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(request *http.Request, via []*http.Request) error {
			if len(via) > 0 && (!strings.EqualFold(request.URL.Scheme, via[0].URL.Scheme) || !strings.EqualFold(request.URL.Host, via[0].URL.Host)) {
				return fmt.Errorf("Ferret API redirect changed origin")
			}

			return nil
		},
	}

	return stdlibdocs.Generate(context.Background(), client, stdlibdocs.Options{
		Version:   versions.Runtime.V2,
		OutputDir: STDLIB_DOCS_DIR,
	})
}
