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
	functionView struct {
		QualifiedName string
		FunctionID    string
		Signatures    []signatureView
	}

	functionMenuItem struct {
		Label  string
		Anchor string
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
		Parameters            []parameterView
		Return                *returnView
		Throws                []api.Throw
		DeprecatedParagraphs  [][]string
	}

	parameterView struct {
		Name        string
		Type        string
		Description string
	}

	returnView struct {
		Type        string
		Description string
	}

	categoryPage struct {
		Title       string
		Description string
		Functions   []functionView
	}

	sectionOptions struct {
		SidebarTitle string
		Title        string
		Description  string
		Kind         string
		Aliases      []string
		FunctionMenu []functionMenuItem
		Weight       int
		Body         []byte
	}
)

var categoryTemplate = template.Must(template.New("category").Parse(`<h1>{{ .Title }}</h1>
{{ with .Description }}
<p>{{ . }}</p>
{{ end }}

<div class="stdlib-category-reference">
{{ range .Functions }}<section class="stdlib-api-function" aria-labelledby="{{ .FunctionID }}">
  <h2 id="{{ .FunctionID }}">
  <a class="stdlib-api-entity-link" href="#{{ .FunctionID }}" aria-label="{{ .QualifiedName }}" title="Link to {{ .QualifiedName }}">
    <code>{{ .QualifiedName }}</code><span aria-hidden="true">#</span>
  </a>
</h2>
{{ range .Signatures }}  <section class="stdlib-api-signature" aria-labelledby="{{ .ID }}">
  <h3 id="{{ .ID }}">
    <a class="stdlib-api-entity-link" href="#{{ .ID }}">{{ .Label }}<span aria-hidden="true">#</span></a>
  </h3>
  <p class="stdlib-api-signature-code"><code>{{ .Code }}</code></p>
{{ if .DeprecatedParagraphs }}  <div class="stdlib-api-deprecated"><strong>Deprecated.</strong>{{ range .DeprecatedParagraphs }}<p>{{ range $index, $line := . }}{{ if $index }}<br>{{ end }}{{ $line }}{{ end }}</p>{{ end }}</div>
{{ end }}{{ range .DescriptionParagraphs }}  <p>{{ range $index, $line := . }}{{ if $index }}<br>{{ end }}{{ $line }}{{ end }}</p>
{{ end }}  <dl class="stdlib-api-details">
    <div>
      <dt>Parameters</dt>
      <dd>{{ if .Parameters }}<ul>{{ range .Parameters }}<li><span class="stdlib-api-value-heading"><code>{{ .Name }}</code>{{ if .Type }}<code class="stdlib-api-value-type">{{ .Type }}</code>{{ end }}</span>{{ if .Description }}<span>{{ .Description }}</span>{{ end }}</li>{{ end }}</ul>{{ else }}<span>None</span>{{ end }}</dd>
    </div>
    {{- if .Variadic }}<div><dt>Signature</dt><dd><span class="stdlib-api-parameter-kind">Variadic</span></dd></div>{{ end -}}
    {{- with .Return }}<div><dt>Returns</dt><dd class="stdlib-api-value"><span class="stdlib-api-value-heading"><code>{{ .Type }}</code></span><span>{{ .Description }}</span></dd></div>{{ end -}}
    {{- if .Throws }}<div><dt>Throws</dt><dd><ul>{{ range .Throws }}<li class="stdlib-api-value"><span class="stdlib-api-value-heading"><code>{{ .Error }}</code></span><span>{{ .Description }}</span></li>{{ end }}</ul></dd></div>{{ end }}
  </dl>
</section>
{{ end }}</section>
{{ end }}</div>
`))

func renderReference(root string, reference *api.Reference, catalog *apicatalog.Catalog) error {
	if catalog == nil {
		return fmt.Errorf("cannot render Standard Library without an API Catalog")
	}

	paths := make(map[string]string)

	apiFunctions := make(map[functionIdentity]api.Function)
	for _, namespace := range reference.Namespaces {
		for _, function := range namespace.Functions {
			apiFunctions[functionIdentity{Namespace: namespace.Name, Name: function.Name}] = function
		}
	}

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
		}

		if err := writeCategory(filepath.Join(root, category.ID), category, functions, aliases, (index+1)*10); err != nil {
			return err
		}
	}

	return nil
}

func writeCategory(root string, category apicatalog.Category, functions []categorizedFunction, aliases []string, weight int) error {
	page := categoryPage{
		Title:       category.Title,
		Description: category.Description,
		Functions:   make([]functionView, 0, len(functions)),
	}
	functionMenu := make([]functionMenuItem, 0, len(functions))

	anchors := make(map[string]string, len(functions))
	for _, function := range functions {
		view := buildFunctionView(function.Identity, function.Function)
		key := strings.ToLower(view.FunctionID)
		if previous, exists := anchors[key]; exists {
			return fmt.Errorf("cannot render category %q: function anchor %q for %s collides with %s", category.ID, view.FunctionID, function.Identity, previous)
		}

		anchors[key] = function.Identity.String()
		page.Functions = append(page.Functions, view)
		functionMenu = append(functionMenu, functionMenuItem{
			Label:  view.QualifiedName,
			Anchor: view.FunctionID,
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
		FunctionMenu: functionMenu,
		Weight:       weight,
		Body:         body.Bytes(),
	})
}

func writeSection(root string, options sectionOptions) error {
	frontMatter := map[string]any{
		"title":           options.Title,
		"sidebarTitle":    options.SidebarTitle,
		"type":            "docs",
		"draft":           false,
		"stdlibGenerated": true,
	}
	if options.Description != "" {
		frontMatter["description"] = options.Description
	}

	if options.Kind != "" {
		frontMatter["stdlibKind"] = options.Kind
	}

	if len(options.Aliases) > 0 {
		frontMatter["aliases"] = options.Aliases
	}
	if len(options.FunctionMenu) > 0 {
		frontMatter["functionMenu"] = options.FunctionMenu
	}

	if options.Weight != 0 {
		frontMatter["weight"] = options.Weight
	}

	return writePage(filepath.Join(root, "_index.md"), frontMatter, options.Body)
}

func buildFunctionView(identity functionIdentity, function api.Function) functionView {
	qualified := identity.String()
	signatures := sortedSignatures(function.Signatures)
	view := functionView{
		QualifiedName: qualified,
		FunctionID:    functionAnchor(identity.Namespace, identity.Name),
		Signatures:    make([]signatureView, 0, len(signatures)),
	}

	signatureIDs := make([]string, len(signatures))
	signatureIDCounts := make(map[string]int, len(signatures))
	for index, signature := range signatures {
		kind := "fixed"
		if signature.Variadic {
			kind = "variadic"
		}
		kind += "-" + strconv.Itoa(len(signature.Parameters))
		signatureIDs[index] = view.FunctionID + "-signature-" + kind
		signatureIDCounts[signatureIDs[index]]++
	}

	signatureIDIndexes := make(map[string]int, len(signatureIDs))
	for index, signature := range signatures {
		id := signatureIDs[index]
		if signatureIDCounts[id] > 1 {
			signatureIDIndexes[id]++
			id += "-" + strconv.Itoa(signatureIDIndexes[id])
		}

		view.Signatures = append(view.Signatures, signatureView{
			ID:                    id,
			Label:                 signatureLabel(signature),
			Code:                  signatureCode(qualified, signature),
			Variadic:              signature.Variadic,
			DescriptionParagraphs: paragraphs(signature.Description),
			Parameters:            parameterViews(signature.Parameters),
			Return:                returnValueView(signature.Return),
			Throws:                signature.Throws,
			DeprecatedParagraphs:  paragraphs(signature.Deprecated),
		})
	}

	return view
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
		case []functionMenuItem:
			fmt.Fprintf(content, "%s:\n", key)
			for _, item := range typed {
				fmt.Fprintf(content, "  - label: %s\n", strconv.Quote(item.Label))
				fmt.Fprintf(content, "    anchor: %s\n", strconv.Quote(item.Anchor))
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
		parts = append(parts, parameter.Name, typeStructuralKey(parameter.Type), parameter.Description)
	}

	if signature.Return == nil {
		parts = append(parts, "0")
	} else {
		parts = append(parts, "1", typeStructuralKey(signature.Return.Type), signature.Return.Description)
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

func parameterViews(parameters []api.Parameter) []parameterView {
	result := make([]parameterView, len(parameters))
	for index, parameter := range parameters {
		result[index] = parameterView{
			Name:        parameter.Name,
			Type:        formatType(parameter.Type),
			Description: parameter.Description,
		}
	}

	return result
}

func returnValueView(value *api.Return) *returnView {
	if value == nil {
		return nil
	}

	return &returnView{Type: formatType(value.Type), Description: value.Description}
}

func formatType(value *api.Type) string {
	if value == nil {
		return ""
	}

	switch value.Kind {
	case api.TypeKindNamed:
		return value.Name
	case api.TypeKindUnion:
		members := make([]string, len(value.Types))
		for index := range value.Types {
			members[index] = formatType(&value.Types[index])
		}

		return strings.Join(members, " | ")
	case api.TypeKindList:
		return "[" + formatType(value.Element) + "]"
	default:
		return ""
	}
}

func typeStructuralKey(value *api.Type) string {
	if value == nil {
		return "none"
	}

	switch value.Kind {
	case api.TypeKindNamed:
		return "named:" + lengthPrefixed(value.Name)
	case api.TypeKindUnion:
		var key strings.Builder
		key.WriteString("union:")
		key.WriteString(strconv.Itoa(len(value.Types)))
		key.WriteByte(':')
		for index := range value.Types {
			key.WriteString(lengthPrefixed(typeStructuralKey(&value.Types[index])))
		}

		return key.String()
	case api.TypeKindList:
		return "list:" + lengthPrefixed(typeStructuralKey(value.Element))
	default:
		return "invalid:" + lengthPrefixed(string(value.Kind))
	}
}

func lengthPrefixed(value string) string {
	return strconv.Itoa(len(value)) + ":" + value
}

func functionAnchor(namespace, function string) string {
	segments := []string{"global"}
	if namespace == "" {
		return strings.ToLower(strings.Join(append(segments, function), "-"))
	}

	segments = append(strings.Split(namespace, "::"), function)
	return strings.ToLower(strings.Join(segments, "-"))
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
