# Ferret Website

Built with Hugo

### Building

1. [Mage](https://magefile.org) is the build tool.
2. [frep](https://github.com/subchen/frep/releases) is the tool used for templating the documentation
pages. You need to install it.
3. Generating `stdlib` documentation requires doc rep YAML. You can then use
`sh generate-docs.sh <REP_FILE>` to generate the documentation.
