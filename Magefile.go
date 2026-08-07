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

	"github.com/MontFerret/montferret.github.io/internal/registryroutes"
	"gopkg.in/yaml.v2"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const OUTPUT_DIR = "public"
const OUTPUT_FILES = "public/*"
const THEME_DIR = "themes/ferret"
const CONTENT_DIR = "content"
const STDLIB_DOCS_DIR = "content/docs/stdlib"
const STDLIB_AST = "stdlib-docs-rep.yaml"
const STDLIB_TEMPLATE = "templates/docs/stdlib.template"
const REGISTRY_BASE_URL = "https://registry.ferretlang.org/"

type ASTModule struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
}

type AST struct {
	Modules map[string]ASTModule `yaml:"modules"`
}

func removeFiles() error {
	matches, err := filepath.Glob(filepath.Join(OUTPUT_DIR, "*"))

	if err != nil {
		return err
	}

	for _, item := range matches {
		err = os.RemoveAll(item)

		if err != nil {
			return err
		}
	}

	return nil
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
	return os.RemoveAll(OUTPUT_DIR)
}

// Starts local Hugo server
func Serve() error {
	mg.Deps(Clean)

	return sh.RunV("hugo", "server")
}

// Runs the production Hugo build and generates the search index
func Build() error {
	mg.Deps(Clean)

	return generateSite()
}

// Builds the website search index and serves the generated static site
func ServeSearch() error {
	mg.Deps(Clean)

	if err := sh.RunV("hugo", "--baseURL", "http://localhost:1414/"); err != nil {
		return err
	}

	if err := sh.RunV("npm", "--prefix", THEME_DIR, "run", "build:search"); err != nil {
		return err
	}

	if err := verifySearchIndex(); err != nil {
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

// Generates documentation
func Generate() error {
	_, err := os.Stat(STDLIB_AST)

	if err != nil {
		fmt.Println("Missing stdlib data source")

		return err
	}

	content, err := os.ReadFile(STDLIB_AST)

	if err != nil {
		fmt.Println("Failed to read data source")

		return err
	}

	ast := AST{}

	if err := yaml.Unmarshal([]byte(content), &ast); err != nil {
		fmt.Println("Failed to parse data source")

		return err
	}

	for _, module := range ast.Modules {
		name := strings.ReplaceAll(module.Name, "/", "-")

		sh.RunWith(map[string]string{
			"USING_KEY": module.Name,
		}, "frep", "--load", STDLIB_AST, "--overwrite", fmt.Sprintf("%s:%s/%s.md", STDLIB_TEMPLATE, STDLIB_DOCS_DIR, name))
	}

	return nil
}
