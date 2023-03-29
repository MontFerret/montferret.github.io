# Ferret Website

Built with Hugo

## Prerequisites
- [Go](https://golang.org/doc/install)
- [Hugo](https://gohugo.io/getting-started/installing/)
- [Node.js](https://nodejs.org/en/download/)
- [Mage](https://magefile.org) for building the website.
- [frep](https://github.com/subchen/frep/releases)(optional) is the tool used for templating the documentation

## Getting Started
### Installing dependencies
```bash
go mod tidy && mage install
```

### Dev server
```bash
mage serve
```

### Generating API docs
Generating `stdlib` documentation requires doc rep YAML.

```bash
mage generate
```

### Publishing
```bash
mage publish
```