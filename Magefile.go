//+build mage

package main

import (
    "errors"
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

// Runs Hugo build to generate the website
func Build() error {
	mg.Deps(Clean)

	return sh.RunV("hugo")
}

// Installs theme
func Install() error {
	os.Chdir(THEME_DIR)

	defer os.Chdir("../..")

	return sh.RunV("npm", "install")
}

// Publishes website to GitHub Pages
func Publish() error {
	dirty, err := sh.Output("git", "status", "-s")

	if err != nil {
		return err
	}

	if dirty != "" {
		return errors.New("The working directory is dirty. Please commit any pending changes.")
	}

	fmt.Println("Deleting old publication")

	if err := Clean(); err != nil {
		fmt.Println("Failed to delete old publication")

		return err
	}

	if err := os.Mkdir(OUTPUT_DIR, 0755); err != nil {
		fmt.Println("Failed to create new directory")

		return err
	}

	if err := sh.RunV("git", "worktree", "prune"); err != nil {
		return err
	}

	if err := os.RemoveAll(WORKTREE_DIR); err != nil {
		return err
	}

	fmt.Println("Checking out master branch into public")

	if err := sh.RunV("git", "worktree", "add", "-B", "master", "public", "origin/master"); err != nil {
		return err
	}

	fmt.Println("Removing existing files")

	if err := removeFiles(); err != nil {
		return err
	}

	fmt.Println("Generating site")

	if err := Build(); err != nil {
		fmt.Println("Failed to generate site")

		return err
	}

	fmt.Println("Updating master branch")

	// cd public && git add --all && git commit -m "Publishing to master (publish.sh)" && git push origin master && cd ..

	return nil
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
