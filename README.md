# Ferret Website

Built with Hugo

### Building

1. [frep](https://github.com/subchen/frep/releases) is the tool used for templating the documentation
pages. You need to install it.
2. [yq](https://github.com/kislyuk/yq) is a YAML and XML processor - jq wrapper for YAML/XML documents.
3. Generating `stdlib` documentation requires doc rep YAML. You can then use
`sh generate-docs.sh <REP_FILE>` to generate the documentation.
