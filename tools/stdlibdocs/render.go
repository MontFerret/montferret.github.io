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

	categorizedFunction struct {
		Identity functionIdentity
		Function api.Function
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
	if catalog == nil {
		return fmt.Errorf("cannot render Standard Library without an API Catalog")
	}

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

	apiFunctions := make(map[functionIdentity]api.Function)
	for _, namespace := range reference.Namespaces {
		for _, function := range namespace.Functions {
			apiFunctions[functionIdentity{Namespace: namespace.Name, Name: function.Name}] = function
		}
	}

	functionCategories := make(map[functionIdentity]functionCategory, len(apiFunctions))
	for index, category := range catalog.Categories {
		if err := registerRoute(paths, category.ID, "category "+category.ID); err != nil {
			return err
		}

		aliases := categoryAliases(category.ID)
		for _, alias := range aliases {
			if err := registerRoute(paths, alias, "alias for category "+category.ID); err != nil {
				return err
			}
		}

		functions := make([]categorizedFunction, 0, len(category.Functions))
		for _, reference := range category.Functions {
			identity := functionIdentity{Namespace: reference.Namespace, Name: reference.Name}
			function, exists := apiFunctions[identity]
			if !exists {
				return fmt.Errorf("cannot render category %q: function %q is missing", category.ID, identity)
			}

			functions = append(functions, categorizedFunction{Identity: identity, Function: function})
			functionCategories[identity] = functionCategory{ID: category.ID, Title: category.Title}
		}

		if err := writeCategory(filepath.Join(root, category.ID), category, functions, aliases, (index+1)*10); err != nil {
			return err
		}
	}

	identities := make([]functionIdentity, 0, len(apiFunctions))
	for identity := range apiFunctions {
		identities = append(identities, identity)
	}
	sort.Slice(identities, func(i, j int) bool {
		if identities[i].Namespace != identities[j].Namespace {
			return identities[i].Namespace < identities[j].Namespace
		}

		return identities[i].Name < identities[j].Name
	})

	for _, identity := range identities {
		category, exists := functionCategories[identity]
		if !exists {
			return fmt.Errorf("cannot render function %q without a category", identity)
		}

		route := functionRoute(identity)
		if err := registerRoute(paths, route, "function "+identity.String()); err != nil {
			return err
		}

		filename := filepath.Join(root, filepath.FromSlash(route)) + ".md"
		if err := writeFunction(filename, identity, apiFunctions[identity], &category); err != nil {
			return err
		}
	}

	return nil
}

func writeCategory(root string, category apicatalog.Category, functions []categorizedFunction, aliases []string, weight int) error {
	page := categoryPage{
		Title:       category.Title,
		Description: category.Description,
		Functions:   make([]categoryFunctionView, 0, len(functions)),
	}

	for _, function := range functions {
		page.Functions = append(page.Functions, categoryFunctionView{
			Name:    function.Identity.String(),
			URL:     "/docs/standard-library/" + functionRoute(function.Identity) + "/",
			Summary: functionSummary(function.Function),
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
		Kind:         "Category",
		Aliases:      aliases,
		Weight:       weight,
		Body:         body.Bytes(),
	})
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

func writeFunction(filename string, identity functionIdentity, function api.Function, category *functionCategory) error {
	qualified := identity.String()
	signatures := sortedSignatures(function.Signatures)
	page := functionPage{
		QualifiedName: qualified,
		FunctionID:    functionAnchor(identity.Namespace, identity.Name),
		Signatures:    make([]signatureView, 0, len(signatures)),
	}

	if category != nil {
		page.Breadcrumb = &breadcrumbView{
			CategoryTitle: category.Title,
			CategoryURL:   "/docs/standard-library/" + category.ID + "/",
			FunctionName:  identity.String(),
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
		"title":               qualified,
		"sidebarTitle":        function.Name,
		"description":         description,
		"type":                "docs",
		"draft":               false,
		"sidebarHidden":       true,
		"stdlibLandingHidden": true,
		"stdlibGenerated":     true,
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
	normalized := filepath.ToSlash(route)
	if !strings.HasPrefix(normalized, "/") {
		normalized = "/docs/standard-library/" + strings.Trim(normalized, "/") + "/"
	}
	key := strings.ToLower(normalized)
	if previous, exists := paths[key]; exists {
		return fmt.Errorf("cannot render %s at %q: route collides with %s", identity, normalized, previous)
	}

	paths[key] = identity

	return nil
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

func functionRoute(identity functionIdentity) string {
	if identity.Namespace == "" {
		return filepath.ToSlash(filepath.Join("functions", identity.Name))
	}

	segments := strings.Split(identity.Namespace, "::")
	segments = append(segments, identity.Name)

	return filepath.ToSlash(filepath.Join(segments...))
}

func categoryAliases(categoryID string) []string {
	aliases := []string{"/docs/stdlib/" + categoryID + "/"}
	switch categoryID {
	case "io":
		aliases = append(aliases,
			"/docs/standard-library/io/fs/",
			"/docs/stdlib/io-fs/",
			"/docs/standard-library/io/net/",
			"/docs/standard-library/io/net/http/",
			"/docs/standard-library/io/http/",
			"/docs/stdlib/io-net-http/",
		)
	case "testing":
		aliases = append(aliases,
			"/docs/standard-library/t/",
			"/docs/standard-library/t/not/",
		)
	}

	return aliases
}
