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
)

type (
	namespaceNode struct {
		Name      string
		Segment   string
		Functions []api.Function
		Children  map[string]*namespaceNode
	}

	functionPage struct {
		QualifiedName string
		FunctionID    string
		Signatures    []signatureView
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
<h1 id="{{ .FunctionID }}">
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
)

func renderReference(root string, reference *api.Reference) error {
	namespaces := append([]api.Namespace(nil), reference.Namespaces...)
	sort.Slice(namespaces, func(i, j int) bool { return namespaces[i].Name < namespaces[j].Name })

	paths := make(map[string]string)
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
    
		node.Functions = functions
	}

	if len(rootNode.Functions) > 0 {
		functionsRoot := filepath.Join(root, "functions")
		if err := registerRoute(paths, "functions", "global functions"); err != nil {
			return err
		}
	
		if err := writeSection(functionsRoot, "Functions", "Global functions", "Functions in the global Ferret namespace.", legacyFunctionAliases); err != nil {
			return err
		}
	
		if err := writeFunctions(functionsRoot, "", rootNode.Functions, paths, "functions"); err != nil {
			return err
		}
	}

	segments := sortedNodeKeys(rootNode.Children)
	for _, segment := range segments {
		if err := writeNamespace(root, rootNode.Children[segment], paths, nil); err != nil {
			return err
		}
	}

	return nil
}

func writeNamespace(root string, node *namespaceNode, paths map[string]string, parentSegments []string) error {
	segments := append(append([]string(nil), parentSegments...), node.Segment)
	relative := filepath.Join(segments...)
	if err := registerRoute(paths, relative, "namespace "+node.Name); err != nil {
		return err
	}

	aliases := namespaceAliases(node.Name)
	if err := writeSection(filepath.Join(root, relative), node.Segment, node.Name, "Functions in the "+node.Name+" Ferret namespace.", aliases); err != nil {
		return err
	}
	
	if err := writeFunctions(filepath.Join(root, relative), node.Name, node.Functions, paths, relative); err != nil {
		return err
	}

	for _, segment := range sortedNodeKeys(node.Children) {
		if err := writeNamespace(root, node.Children[segment], paths, segments); err != nil {
			return err
		}
	}

	return nil
}

func writeFunctions(root, namespace string, functions []api.Function, paths map[string]string, parentRoute string) error {
	for _, function := range functions {
		route := filepath.Join(parentRoute, function.Name)
		qualified := qualifiedName(namespace, function.Name)
	
		if err := registerRoute(paths, route, "function "+qualified); err != nil {
			return err
		}
		
		if err := writeFunction(filepath.Join(root, function.Name+".md"), namespace, function); err != nil {
			return err
		}
	}

	return nil
}

func writeSection(root, sidebarTitle, title, description string, aliases []string) error {
	frontMatter := map[string]any{
		"title":           title,
		"sidebarTitle":    sidebarTitle,
		"description":     description,
		"type":            "docs",
		"draft":           false,
		"stdlibGenerated": true,
	}
	
	if len(aliases) > 0 {
		frontMatter["aliases"] = aliases
	}
	
	body := fmt.Sprintf("# `%s`\n\n%s\n\n{{< stdlib-children >}}\n", title, description)

	return writePage(filepath.Join(root, "_index.md"), frontMatter, []byte(body))
}

func writeFunction(filename, namespace string, function api.Function) error {
	qualified := qualifiedName(namespace, function.Name)
	signatures := append([]api.Signature(nil), function.Signatures...)
	
	sort.Slice(signatures, func(i, j int) bool {
		if signatures[i].Variadic != signatures[j].Variadic {
			return !signatures[i].Variadic
		}
    
		return len(signatures[i].Parameters) < len(signatures[j].Parameters)
	})

	page := functionPage{
		QualifiedName: qualified,
		FunctionID:    functionAnchor(namespace, function.Name),
		Signatures:    make([]signatureView, 0, len(signatures)),
	}
	
	for _, signature := range signatures {
		kind := "fixed-" + strconv.Itoa(len(signature.Parameters))
		
		if signature.Variadic {
			kind = "variadic"
		}
		
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
	
	frontMatter := map[string]any{
		"title":           qualified,
		"sidebarTitle":    function.Name,
		"description":     "API reference for " + qualified + ".",
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
