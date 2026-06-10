//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const OUTPUT_DIR = "public"
const OUTPUT_FILES = "public/*"
const THEME_DIR = "themes/ferret"
const WORKTREE_DIR = ".git/worktrees/public"
const CONTENT_DIR = "content"
const STDLIB_DOCS_DIR = "content/docs/stdlib"
const STDLIB_AST = "stdlib-docs-rep.yaml"
const STDLIB_TEMPLATE = "templates/docs/stdlib.template"

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

func ensureCleanWorkingTree() error {
	out, err := sh.Output("git", "status", "--porcelain")
	
	if err != nil {
		return err
	}

	if strings.TrimSpace(out) != "" {
		return fmt.Errorf("working directory is dirty; please commit or stash pending changes before publishing")
	}

	return nil
}

func cleanWorktreeContents(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		name := entry.Name()
		
		if name == ".git" {
			continue
		}
		
		if err := os.RemoveAll(filepath.Join(dir, name)); err != nil {
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
	
	return verifySearchIndex()
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

	return sh.RunV("npm", "--prefix", THEME_DIR, "run", "serve:search")
}

// Installs theme
func Install() error {
	os.Chdir(THEME_DIR)

	defer os.Chdir("../..")

	return sh.RunV("npm", "install")
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

// Publishes website to GitHub Pages
func Publish() error {
	if err := ensureCleanWorkingTree(); err != nil {
		return err
	}

	fmt.Println("Deleting old publication")
	if err := Clean(); err != nil {
		return err
	}

	if err := sh.RunV("git", "worktree", "prune"); err != nil {
		return err
	}

	if err := os.RemoveAll(WORKTREE_DIR); err != nil {
		return err
	}

	fmt.Println("Checking out master branch into public")
	if err := sh.RunV("git", "worktree", "add", "-B", "master", OUTPUT_DIR, "origin/master"); err != nil {
		return err
	}

	fmt.Println("Removing existing files")
	if err := cleanWorktreeContents(OUTPUT_DIR); err != nil {
		return err
	}

	fmt.Println("Generating site")
	if err := generateSite(); err != nil {
		return err
	}

	fmt.Println("Updating master branch")
	if err := sh.RunV("git", "-C", OUTPUT_DIR, "add", "--all"); err != nil {
		return err
	}

	status, err := sh.Output("git", "-C", OUTPUT_DIR, "status", "--porcelain")
	if err != nil {
		return err
	}

	if strings.TrimSpace(status) == "" {
		fmt.Println("No changes to publish")
		return nil
	}

	if err := sh.RunV("git", "-C", OUTPUT_DIR, "commit", "-m", "Publishing to master"); err != nil {
		return err
	}

	return sh.RunV("git", "-C", OUTPUT_DIR, "push", "origin", "master")
}
