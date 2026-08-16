package stdlibdocs

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/MontFerret/specs/pkg/api"
	apicatalog "github.com/MontFerret/specs/pkg/api/catalog"
)

type (
	namespaceNode struct {
		Name      string
		Segment   string
		Real      bool
		Functions []api.Function
		Children  map[string]*namespaceNode
	}

	functionPage struct {
		QualifiedName string
		FunctionID    string
		Breadcrumb    *breadcrumbView
		Signatures    []signatureView
	}

	breadcrumbView struct {
		CategoryTitle string
		CategoryURL   string
		FunctionName  string
	}

	functionCategory struct {
		ID    string
		Title string
	}

	signatureView struct {
		ID                    string
		Label                 string
		Code                  string
		Variadic              bool
		DescriptionParagraphs [][]string
		Parameters            []api.Parameter
		Return                *api.Return
		Throws                []api.Throw
		DeprecatedParagraphs  [][]string
	}

	categoryPage struct {
		Title       string
		Description string
		Functions   []categoryFunctionView
	}

	categoryFunctionView struct {
		Name    string
		URL     string
		Summary string
	}

	sectionOptions struct {
		SidebarTitle      string
		Title             string
		Description       string
		Kind              string
		Aliases           []string
		Weight            int
		SidebarHidden     bool
		StdlibLandingHide bool
		Body              []byte
	}
)

var (
	legacyFunctionAliases = []string{
		"/docs/standard-library/arrays/",
		"/docs/standard-library/collections/",
		"/docs/standard-library/datetime/",
		"/docs/standard-library/math/",
		"/docs/standard-library/objects/",
		"/docs/standard-library/path/",
		"/docs/standard-library/strings/",
		"/docs/standard-library/types/",
		"/docs/standard-library/utils/",
		"/docs/stdlib/arrays/",
		"/docs/stdlib/collections/",
		"/docs/stdlib/datetime/",
		"/docs/stdlib/math/",
		"/docs/stdlib/objects/",
		"/docs/stdlib/path/",
		"/docs/stdlib/strings/",
		"/docs/stdlib/types/",
		"/docs/stdlib/utils/",
	}

	functionTemplate = template.Must(template.New("function").Parse(`
{{ with .Breadcrumb }}<nav class="stdlib-api-breadcrumbs" aria-label="Breadcrumb">
  <a href="/docs/standard-library/">Standard Library</a><span aria-hidden="true">/</span>
  <a href="{{ .CategoryURL }}">{{ .CategoryTitle }}</a><span aria-hidden="true">/</span>
  <span aria-current="page">{{ .FunctionName }}</span>
</nav>
{{ end }}<h1 id="{{ .FunctionID }}">
  <a class="stdlib-api-entity-link" href="#{{ .FunctionID }}" aria-label="{{ .QualifiedName }}" title="Link to {{ .QualifiedName }}">
    <code>{{ .QualifiedName }}</code><span aria-hidden="true">#</span>
  </a>
</h1>
{{ range .Signatures }}
<section class="stdlib-api-signature" aria-labelledby="{{ .ID }}">
  <h2 id="{{ .ID }}">
    <a class="stdlib-api-entity-link" href="#{{ .ID }}">{{ .Label }}<span aria-hidden="true">#</span></a>
  </h2>
  <p class="stdlib-api-signature-code"><code>{{ .Code }}</code></p>
  {{ if .DeprecatedParagraphs }}
  <div class="stdlib-api-deprecated"><strong>Deprecated.</strong>{{ range .DeprecatedParagraphs }}<p>{{ range $index, $line := . }}{{ if $index }}<br>{{ end }}{{ $line }}{{ end }}</p>{{ end }}</div>
  {{ end }}
  {{ range .DescriptionParagraphs }}<p>{{ range $index, $line := . }}{{ if $index }}<br>{{ end }}{{ $line }}{{ end }}</p>{{ end }}
  <dl class="stdlib-api-details">
    <div>
      <dt>Parameters</dt>
      <dd>{{ if .Parameters }}<ul>{{ range .Parameters }}<li><span class="stdlib-api-value-heading"><code>{{ .Name }}</code>{{ if .Type }}<code class="stdlib-api-value-type">{{ .Type }}</code>{{ end }}</span>{{ if .Description }}<span>{{ .Description }}</span>{{ end }}</li>{{ end }}</ul>{{ else }}<span>None</span>{{ end }}</dd>
    </div>
    {{- if .Variadic }}<div><dt>Signature</dt><dd><span class="stdlib-api-parameter-kind">Variadic</span></dd></div>{{ end -}}
    {{- with .Return }}<div><dt>Returns</dt><dd class="stdlib-api-value"><span class="stdlib-api-value-heading"><code>{{ .Type }}</code></span><span>{{ .Description }}</span></dd></div>{{ end -}}
    {{- if .Throws }}<div><dt>Throws</dt><dd><ul>{{ range .Throws }}<li class="stdlib-api-value"><span class="stdlib-api-value-heading"><code>{{ .Error }}</code></span><span>{{ .Description }}</span></li>{{ end }}</ul></dd></div>{{ end }}
  </dl>
</section>
{{ end }}`))

	categoryTemplate = template.Must(template.New("category").Parse(`<h1>{{ .Title }}</h1>

<p>{{ .Description }}</p>

<ul class="stdlib-category-functions">
{{ range .Functions }}  <li>
    <a href="{{ .URL }}"><code>{{ .Name }}</code></a>{{ with .Summary }}
    <p>{{ . }}</p>{{ end }}
  </li>
{{ end }}</ul>
`))
)

func renderReference(root string, reference *api.Reference, catalog *apicatalog.Catalog) error {
	rootNode := buildNamespaceTree(reference)
	if catalog == nil {
		return renderFlatReference(root, rootNode)
	}

	return renderCatalogReference(root, rootNode, catalog)
}

func buildNamespaceTree(reference *api.Reference) *namespaceNode {
	namespaces := append([]api.Namespace(nil), reference.Namespaces...)
	sort.Slice(namespaces, func(i, j int) bool { return namespaces[i].Name < namespaces[j].Name })

	rootNode := &namespaceNode{Children: make(map[string]*namespaceNode)}
	for _, namespace := range namespaces {
		functions := append([]api.Function(nil), namespace.Functions...)
		sort.Slice(functions, func(i, j int) bool { return functions[i].Name < functions[j].Name })

		if namespace.Name == "" {
			rootNode.Functions = functions

			continue
		}

		node := rootNode
		segments := strings.Split(namespace.Name, "::")
		qualified := make([]string, 0, len(segments))
		for _, segment := range segments {
			qualified = append(qualified, segment)
			child := node.Children[segment]
			if child == nil {
				child = &namespaceNode{
					Name:     strings.Join(qualified, "::"),
					Segment:  segment,
					Children: make(map[string]*namespaceNode),
				}
				node.Children[segment] = child
			}

			node = child
		}

		node.Real = true
		node.Functions = functions
	}

	return rootNode
}

func renderFlatReference(root string, rootNode *namespaceNode) error {
	paths := make(map[string]string)
	if len(rootNode.Functions) > 0 {
		functionsRoot := filepath.Join(root, "functions")
		if err := registerRoute(paths, "functions", "global functions"); err != nil {
			return err
		}

		if err := writeSection(functionsRoot, sectionOptions{
			SidebarTitle: "Functions",
			Title:        "Global functions",
			Description:  "Functions in the global Ferret namespace.",
			Aliases:      legacyFunctionAliases,
		}); err != nil {
			return err
		}

		if err := writeFunctions(functionsRoot, "", rootNode.Functions, paths, "functions", nil); err != nil {
			return err
		}
	}

	for _, segment := range sortedNodeKeys(rootNode.Children) {
		if err := writeNamespace(root, rootNode.Children[segment], paths, nil, 0, ""); err != nil {
			return err
		}
	}

	return nil
}

func renderCatalogReference(root string, rootNode *namespaceNode, catalog *apicatalog.Catalog) error {
	paths := make(map[string]string)
	functionsRoot := filepath.Join(root, "functions")
	if err := registerRoute(paths, "functions", "global functions"); err != nil {
		return err
	}

	compatibilityBody := []byte("# Global functions\n\nGlobal functions are organized by category on the [Standard Library page](/docs/standard-library/).\n")
	if err := writeSection(functionsRoot, sectionOptions{
		SidebarTitle:      "Functions",
		Title:             "Global functions",
		Description:       "Global functions organized by documentation category.",
		SidebarHidden:     true,
		StdlibLandingHide: true,
		Body:              compatibilityBody,
	}); err != nil {
		return err
	}

	globalFunctions := make(map[string]api.Function, len(rootNode.Functions))
	for _, function := range rootNode.Functions {
		globalFunctions[function.Name] = function
	}

	functionCategories := make(map[string]functionCategory, len(rootNode.Functions))
	for index, category := range catalog.Categories {
		if err := registerRoute(paths, category.ID, "category "+category.ID); err != nil {
			return err
		}

		functions := make([]api.Function, 0, len(category.Functions))
		for _, name := range category.Functions {
			function, exists := globalFunctions[name]
			if !exists {
				return fmt.Errorf("cannot render category %q: global function %q is missing", category.ID, name)
			}

			functions = append(functions, function)
			functionCategories[name] = functionCategory{ID: category.ID, Title: category.Title}
		}

		if err := writeCategory(filepath.Join(root, category.ID), category, functions, (index+1)*10); err != nil {
			return err
		}
	}

	if err := writeFunctions(functionsRoot, "", rootNode.Functions, paths, "functions", functionCategories); err != nil {
		return err
	}

	for index, rootName := range catalog.NamespaceRoots {
		node := rootNode.Children[rootName]
		if node == nil {
			return fmt.Errorf("cannot render namespace root %q: API namespace tree is missing", rootName)
		}

		weight := (len(catalog.Categories) + index + 1) * 10
		if err := writeNamespace(root, node, paths, nil, weight, "Namespace"); err != nil {
			return err
		}
	}

	return nil
}

func writeCategory(root string, category apicatalog.Category, functions []api.Function, weight int) error {
	page := categoryPage{
		Title:       category.Title,
		Description: category.Description,
		Functions:   make([]categoryFunctionView, 0, len(functions)),
	}

	for _, function := range functions {
		page.Functions = append(page.Functions, categoryFunctionView{
			Name:    function.Name,
			URL:     "/docs/standard-library/functions/" + function.Name + "/",
			Summary: functionSummary(function),
		})
	}

	body := bytes.NewBuffer(nil)
	if err := categoryTemplate.Execute(body, page); err != nil {
		return fmt.Errorf("render category %s: %w", category.ID, err)
	}

	return writeSection(root, sectionOptions{
		SidebarTitle: category.Title,
		Title:        category.Title,
		Description:  category.Description,
		Kind:         "Global functions",
		Aliases:      []string{"/docs/stdlib/" + category.ID + "/"},
		Weight:       weight,
		Body:         body.Bytes(),
	})
}

func writeNamespace(root string, node *namespaceNode, paths map[string]string, parentSegments []string, weight int, kind string) error {
	segments := append(append([]string(nil), parentSegments...), node.Segment)
	relative := filepath.Join(segments...)
	if err := registerRoute(paths, relative, "namespace "+node.Name); err != nil {
		return err
	}

	description := "Namespaces beneath the " + node.Name + " Ferret namespace prefix."
	if node.Real {
		description = "Functions in the " + node.Name + " Ferret namespace."
	}

	if err := writeSection(filepath.Join(root, relative), sectionOptions{
		SidebarTitle: node.Segment,
		Title:        node.Name,
		Description:  description,
		Kind:         kind,
		Aliases:      namespaceAliases(node.Name),
		Weight:       weight,
	}); err != nil {
		return err
	}

	if err := writeFunctions(filepath.Join(root, relative), node.Name, node.Functions, paths, relative, nil); err != nil {
		return err
	}

	for index, segment := range sortedNodeKeys(node.Children) {
		childWeight := 0
		if weight != 0 {
			childWeight = (index + 1) * 10
		}

		if err := writeNamespace(root, node.Children[segment], paths, segments, childWeight, kind); err != nil {
			return err
		}
	}

	return nil
}

func writeFunctions(root, namespace string, functions []api.Function, paths map[string]string, parentRoute string, categories map[string]functionCategory) error {
	for _, function := range functions {
		route := filepath.Join(parentRoute, function.Name)
		qualified := qualifiedName(namespace, function.Name)
		if err := registerRoute(paths, route, "function "+qualified); err != nil {
			return err
		}

		var category *functionCategory
		if categories != nil {
			value, exists := categories[function.Name]
			if !exists {
				return fmt.Errorf("cannot render global function %q without a category", function.Name)
			}

			category = &value
		}

		if err := writeFunction(filepath.Join(root, function.Name+".md"), namespace, function, category); err != nil {
			return err
		}
	}

	return nil
}

func writeSection(root string, options sectionOptions) error {
	frontMatter := map[string]any{
		"title":           options.Title,
		"sidebarTitle":    options.SidebarTitle,
		"description":     options.Description,
		"type":            "docs",
		"draft":           false,
		"stdlibGenerated": true,
	}

	if options.Kind != "" {
		frontMatter["stdlibKind"] = options.Kind
	}

	if len(options.Aliases) > 0 {
		frontMatter["aliases"] = options.Aliases
	}

	if options.Weight != 0 {
		frontMatter["weight"] = options.Weight
	}

	if options.SidebarHidden {
		frontMatter["sidebarHidden"] = true
	}

	if options.StdlibLandingHide {
		frontMatter["stdlibLandingHidden"] = true
	}

	body := options.Body
	if len(body) == 0 {
		body = []byte(fmt.Sprintf("# `%s`\n\n%s\n\n{{< stdlib-children >}}\n", options.Title, options.Description))
	}

	return writePage(filepath.Join(root, "_index.md"), frontMatter, body)
}

func writeFunction(filename, namespace string, function api.Function, category *functionCategory) error {
	qualified := qualifiedName(namespace, function.Name)
	signatures := sortedSignatures(function.Signatures)
	page := functionPage{
		QualifiedName: qualified,
		FunctionID:    functionAnchor(namespace, function.Name),
		Signatures:    make([]signatureView, 0, len(signatures)),
	}

	if category != nil {
		page.Breadcrumb = &breadcrumbView{
			CategoryTitle: category.Title,
			CategoryURL:   "/docs/standard-library/" + category.ID + "/",
			FunctionName:  function.Name,
		}
	}

	for _, signature := range signatures {
		kind := "fixed"
		if signature.Variadic {
			kind = "variadic"
		}
		kind += "-" + strconv.Itoa(len(signature.Parameters))

		page.Signatures = append(page.Signatures, signatureView{
			ID:                    page.FunctionID + "-signature-" + kind,
			Label:                 signatureLabel(signature),
			Code:                  signatureCode(qualified, signature),
			Variadic:              signature.Variadic,
			DescriptionParagraphs: paragraphs(signature.Description),
			Parameters:            signature.Parameters,
			Return:                signature.Return,
			Throws:                signature.Throws,
			DeprecatedParagraphs:  paragraphs(signature.Deprecated),
		})
	}

	body := bytes.NewBuffer(nil)
	if err := functionTemplate.Execute(body, page); err != nil {
		return fmt.Errorf("render function %s: %w", qualified, err)
	}

	description := functionSummary(function)
	if description == "" {
		description = "API reference for " + qualified + "."
	}

	frontMatter := map[string]any{
		"title":           qualified,
		"sidebarTitle":    function.Name,
		"description":     description,
		"type":            "docs",
		"draft":           false,
		"sidebarHidden":   true,
		"stdlibGenerated": true,
	}

	return writePage(filename, frontMatter, body.Bytes())
}

func writePage(filename string, frontMatter map[string]any, body []byte) error {
	if err := os.MkdirAll(filepath.Dir(filename), 0o755); err != nil {
		return fmt.Errorf("create generated page directory: %w", err)
	}

	content := bytes.NewBufferString("---\n")
	keys := make([]string, 0, len(frontMatter))
	for key := range frontMatter {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	for _, key := range keys {
		value := frontMatter[key]
		switch typed := value.(type) {
		case string:
			fmt.Fprintf(content, "%s: %s\n", key, strconv.Quote(typed))
		case bool:
			fmt.Fprintf(content, "%s: %t\n", key, typed)
		case int:
			fmt.Fprintf(content, "%s: %d\n", key, typed)
		case []string:
			fmt.Fprintf(content, "%s:\n", key)
			for _, item := range typed {
				fmt.Fprintf(content, "  - %s\n", strconv.Quote(item))
			}
		default:
			return fmt.Errorf("render generated front matter %q: unsupported value %T", key, value)
		}
	}

	content.WriteString("---\n")
	content.Write(body)
	if !bytes.HasSuffix(content.Bytes(), []byte("\n")) {
		content.WriteByte('\n')
	}

	if err := os.WriteFile(filename, content.Bytes(), 0o644); err != nil {
		return fmt.Errorf("write generated page %s: %w", filename, err)
	}

	return nil
}

func registerRoute(paths map[string]string, route, identity string) error {
	key := strings.ToLower(filepath.ToSlash(route))
	if previous, exists := paths[key]; exists {
		return fmt.Errorf("cannot render %s at %q: route collides with %s", identity, filepath.ToSlash(route), previous)
	}

	paths[key] = identity

	return nil
}

func sortedNodeKeys(nodes map[string]*namespaceNode) []string {
	keys := make([]string, 0, len(nodes))
	for key := range nodes {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}

func sortedSignatures(signatures []api.Signature) []api.Signature {
	result := append([]api.Signature(nil), signatures...)
	sort.Slice(result, func(i, j int) bool {
		if result[i].Variadic != result[j].Variadic {
			return !result[i].Variadic
		}

		if len(result[i].Parameters) != len(result[j].Parameters) {
			return len(result[i].Parameters) < len(result[j].Parameters)
		}

		return signatureSortKey(result[i]) < signatureSortKey(result[j])
	})

	return result
}

func signatureSortKey(signature api.Signature) string {
	parts := make([]string, 0, 8+len(signature.Parameters)*3+len(signature.Throws)*2)
	parts = append(parts, signature.Description, signature.Deprecated)
	for _, parameter := range signature.Parameters {
		parts = append(parts, parameter.Name, parameter.Type, parameter.Description)
	}

	if signature.Return == nil {
		parts = append(parts, "0")
	} else {
		parts = append(parts, "1", signature.Return.Type, signature.Return.Description)
	}

	for _, thrown := range signature.Throws {
		parts = append(parts, thrown.Error, thrown.Description)
	}

	var key strings.Builder
	for _, part := range parts {
		key.WriteString(strconv.Itoa(len(part)))
		key.WriteByte(':')
		key.WriteString(part)
	}

	return key.String()
}

func functionSummary(function api.Function) string {
	for _, signature := range sortedSignatures(function.Signatures) {
		description := strings.ReplaceAll(signature.Description, "\r\n", "\n")
		for _, line := range strings.Split(description, "\n") {
			if summary := strings.TrimSpace(line); summary != "" {
				return summary
			}
		}
	}

	return ""
}

func qualifiedName(namespace, function string) string {
	if namespace == "" {
		return function
	}

	return namespace + "::" + function
}

func functionAnchor(namespace, function string) string {
	if namespace == "" {
		return "api-function-global-" + function
	}

	return "api-function-named-" + strings.ReplaceAll(namespace, "::", "-") + "-" + function
}

func signatureLabel(signature api.Signature) string {
	if signature.Variadic {
		return "Variadic signature"
	}

	if len(signature.Parameters) == 1 {
		return "1-parameter signature"
	}

	return strconv.Itoa(len(signature.Parameters)) + "-parameter signature"
}

func signatureCode(qualified string, signature api.Signature) string {
	parameters := make([]string, 0, len(signature.Parameters))
	for index, parameter := range signature.Parameters {
		name := parameter.Name
		if signature.Variadic && index == len(signature.Parameters)-1 {
			name += "..."
		}

		parameters = append(parameters, name)
	}

	return qualified + "(" + strings.Join(parameters, ", ") + ")"
}

func paragraphs(value string) [][]string {
	value = strings.TrimSpace(strings.ReplaceAll(value, "\r\n", "\n"))
	if value == "" {
		return nil
	}

	blocks := strings.Split(value, "\n\n")
	result := make([][]string, 0, len(blocks))
	for _, block := range blocks {
		lines := strings.Split(block, "\n")
		paragraph := make([]string, 0, len(lines))
		for _, line := range lines {
			if trimmed := strings.TrimSpace(line); trimmed != "" {
				paragraph = append(paragraph, trimmed)
			}
		}

		if len(paragraph) > 0 {
			result = append(result, paragraph)
		}
	}

	return result
}

func namespaceAliases(namespace string) []string {
	switch namespace {
	case "io::fs":
		return []string{"/docs/stdlib/io-fs/"}
	case "io::net::http":
		return []string{"/docs/standard-library/io/http/", "/docs/stdlib/io-net-http/"}
	case "t":
		return []string{"/docs/standard-library/testing/", "/docs/stdlib/testing/"}
	default:
		return nil
	}
}
