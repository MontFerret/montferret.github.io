package stdlibdocs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/MontFerret/specs/pkg/api"
	apicatalog "github.com/MontFerret/specs/pkg/api/catalog"
)

const (
	testIndexURL    = "https://api.test/index.json"
	testAPIPath     = "/versions/2.0.0-alpha.46/api.json"
	testCatalogPath = "/versions/2.0.0-alpha.46/catalog.json"
)

func TestGenerateRendersPublishedAPICatalog(t *testing.T) {
	client := fixtureClient(t, nil)
	output := filepath.Join(t.TempDir(), "standard-library")

	if err := Generate(context.Background(), client, Options{IndexURL: testIndexURL, Version: "2.0.0-alpha.46", OutputDir: output}); err != nil {
		t.Fatal(err)
	}

	wantRequests := []string{"/index.json", testAPIPath, testCatalogPath}
	if got := client.requests; !reflect.DeepEqual(got, wantRequests) {
		t.Fatalf("requests = %v, want authoritative sibling artifact requests %v", got, wantRequests)
	}

	wantFiles := []string{
		"arrays/_index.md",
		"io/_index.md",
		"math/_index.md",
		"strings/_index.md",
		"testing/_index.md",
		"types/_index.md",
	}
	generated := snapshot(t, output)
	gotFiles := make([]string, 0, len(generated))
	for filename := range generated {
		gotFiles = append(gotFiles, filename)
	}
	sort.Strings(gotFiles)
	if !reflect.DeepEqual(gotFiles, wantFiles) {
		t.Fatalf("generated files = %v, want one page per catalog category: %v", gotFiles, wantFiles)
	}

	catalog, err := apicatalog.Parse([]byte(readFile(t, "testdata/catalog.json")))
	if err != nil {
		t.Fatal(err)
	}
	for _, category := range catalog.Categories {
		page := readFile(t, filepath.Join(output, category.ID, "_index.md"))
		menuStart := strings.Index(page, "functionMenu:\n")
		if menuStart == -1 {
			t.Fatalf("category %q does not contain its right-side function menu", category.ID)
		}
		if strings.Contains(page, "stdlib-category-index") {
			t.Errorf("category %q contains the obsolete compact function index", category.ID)
		}

		previousFunctionSection := -1
		previousMenuItem := -1
		for _, function := range category.Functions {
			identity := functionIdentity{Namespace: function.Namespace, Name: function.Name}
			anchor := functionAnchor(identity.Namespace, identity.Name)
			functionSection := fmt.Sprintf(`<section class="stdlib-api-function" aria-labelledby="%s">`, anchor)
			sectionPosition := strings.Index(page, functionSection)
		
			if sectionPosition == -1 {
				t.Errorf("category %q does not contain an inline section for %s", category.ID, identity)
			} else if sectionPosition <= previousFunctionSection {
				t.Errorf("category %q inline sections do not preserve catalog order at %s", category.ID, identity)
			}
		
			previousFunctionSection = sectionPosition

			menuItem := fmt.Sprintf("  - label: %q\n    anchor: %q", identity.String(), anchor)
			menuPosition := strings.Index(page[menuStart:], menuItem)
		
			if menuPosition == -1 {
				t.Errorf("category %q right-side menu does not link %s to #%s", category.ID, identity, anchor)
			} else if menuPosition <= previousMenuItem {
				t.Errorf("category %q right-side menu does not preserve catalog order at %s", category.ID, identity)
			}
			
			previousMenuItem = menuPosition

			if count := strings.Count(page, fmt.Sprintf(`id="%s"`, anchor)); count != 1 {
				t.Errorf("category %q contains %d headings for anchor %q, want 1", category.ID, count, anchor)
			}
		}
	}

	mathPage := readFile(t, filepath.Join(output, "math", "_index.md"))
	for _, expected := range []string{
		`aliases:`,
		`"/docs/stdlib/math/"`,
		`description: "Mathematical and numeric global functions."`,
		`stdlibKind: "Category"`,
		`weight: 30`,
		`<h1>Math</h1>`,
		`href="#global-abs"`,
		`id="global-abs"`,
		`abs(value)`,
		`abs returns the absolute value of a number.`,
		`<code>value</code><code class="stdlib-api-value-type">Number</code>`,
		`<dt>Returns</dt>`,
	} {
		if !strings.Contains(mathPage, expected) {
			t.Errorf("Math category does not contain %q", expected)
		}
	}
	arrays := readFile(t, filepath.Join(output, "arrays", "_index.md"))
	for _, expected := range []string{
		`href="#global-append"`,
		`href="#global-flatten"`,
		`id="global-flatten"`,
		`id="global-flatten-signature-fixed-1"`,
		`id="global-flatten-signature-fixed-2"`,
		`flatten(arr)`,
		`flatten(arr, depth)`,
		`<code class="stdlib-api-value-type">Any[]</code>`,
		`<dt>Returns</dt>`,
	} {
		if !strings.Contains(arrays, expected) {
			t.Errorf("Arrays category does not contain %q", expected)
		}
	}
	typesPage := readFile(t, filepath.Join(output, "types", "_index.md"))
	if !strings.Contains(typesPage, `id="global-to_number"`) || !strings.Contains(typesPage, "<dt>Throws</dt>") || !strings.Contains(typesPage, "TypeError") {
		t.Error("Types category does not render to_number and its thrown errors")
	}

	ioCategory := readFile(t, filepath.Join(output, "io", "_index.md"))
	for _, expected := range []string{
		`description: "Functions for working with files, networks, and other input/output operations."`,
		`href="#io-fs-read"`,
		`href="#io-net-http-get"`,
		`id="io-fs-read"`,
		`<code>io::fs::read</code>`,
		`id="io-net-http-get"`,
		`<code>io::net::http::get</code>`,
		`io-net-http-get-signature-variadic-2`,
		`io::net::http::get(url, options...)`,
		`<span class="stdlib-api-parameter-kind">Variadic</span>`,
		`<strong>Deprecated.</strong>`,
		`stdlibKind: "Category"`,
		`weight: 20`,
		`read reads from a file.`,
		`"/docs/standard-library/io/fs/"`,
		`"/docs/standard-library/io/net/http/"`,
		`"/docs/stdlib/io-net-http/"`,
	} {
		if !strings.Contains(ioCategory, expected) {
			t.Errorf("I/O category does not contain %q", expected)
		}
	}

	testingCategory := readFile(t, filepath.Join(output, "testing", "_index.md"))
	for _, expected := range []string{
		`href="#t-eq"`,
		`id="t-eq"`,
		`<code>t::eq</code>`,
		`href="#t-not-eq"`,
		`id="t-not-eq"`,
		`<code>t::not::eq</code>`,
		`"/docs/standard-library/t/"`,
		`"/docs/standard-library/t/not/"`,
	} {
		if !strings.Contains(testingCategory, expected) {
			t.Errorf("Testing category does not contain %q", expected)
		}
	}

	if got, want := arrays, readFile(t, "testdata/golden/arrays.md"); got != want {
		t.Errorf("generated Arrays category does not match its golden file:\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestGenerateKeepsCatalogAndAPIMetadataInTheirOwnSections(t *testing.T) {
	catalog := mutateCatalog(t, func(value *apicatalog.Catalog) {
		value.Categories[2].Description = "Catalog-owned Math description."
	})
	reference := mutateReference(t, func(value *api.Reference) {
		value.Namespaces[0].Functions[0].Signatures[0].Description = "API-owned abs description."
	})
	output := filepath.Join(t.TempDir(), "standard-library")
	if err := Generate(context.Background(), fixtureClient(t, map[string]string{
		testAPIPath:     reference,
		testCatalogPath: catalog,
	}), Options{IndexURL: testIndexURL, Version: "2.0.0-alpha.46", OutputDir: output}); err != nil {
		t.Fatal(err)
	}

	mathPage := readFile(t, filepath.Join(output, "math", "_index.md"))
	for _, expected := range []string{
		`description: "Catalog-owned Math description."`,
		`<p>Catalog-owned Math description.</p>`,
		`<p>API-owned abs description.</p>`,
	} {
		if !strings.Contains(mathPage, expected) {
			t.Errorf("Math category does not contain %q", expected)
		}
	}
}

func TestGenerateOmitsEmptyCategoryDescription(t *testing.T) {
	output := filepath.Join(t.TempDir(), "math")
	category := apicatalog.Category{ID: "math", Title: "Math"}
	functions := []categorizedFunction{{
		Identity: functionIdentity{Name: "abs"},
		Function: api.Function{Name: "abs", Signatures: []api.Signature{{Description: "Absolute value."}}},
	}}
	if err := writeCategory(output, category, functions, nil, 10); err != nil {
		t.Fatal(err)
	}

	mathPage := readFile(t, filepath.Join(output, "_index.md"))
	if strings.Contains(mathPage, "description:") || strings.Contains(mathPage, "<p></p>") {
		t.Fatal("Math category rendered an empty catalog description")
	}
	if !strings.Contains(mathPage, "functionMenu:\n  - label: \"abs\"\n    anchor: \"global-abs\"") {
		t.Fatal("single-function Math category does not contain its right-side function menu")
	}
}

func TestGenerateIsDeterministicAcrossAPIOrderingAndRepeatedRuns(t *testing.T) {
	firstRoot := filepath.Join(t.TempDir(), "standard-library")
	options := Options{IndexURL: testIndexURL, Version: "2.0.0-alpha.46", OutputDir: firstRoot}
	if err := Generate(context.Background(), fixtureClient(t, nil), options); err != nil {
		t.Fatal(err)
	}
	first := snapshot(t, firstRoot)

	if err := Generate(context.Background(), fixtureClient(t, nil), options); err != nil {
		t.Fatal(err)
	}
	if got := snapshot(t, firstRoot); !reflect.DeepEqual(got, first) {
		t.Fatal("repeated generation changed output")
	}

	reference, err := api.Parse([]byte(readFile(t, "testdata/api.json")))
	if err != nil {
		t.Fatal(err)
	}
	reverse(reference.Namespaces)
	for namespaceIndex := range reference.Namespaces {
		reverse(reference.Namespaces[namespaceIndex].Functions)
		for functionIndex := range reference.Namespaces[namespaceIndex].Functions {
			reverse(reference.Namespaces[namespaceIndex].Functions[functionIndex].Signatures)
		}
	}
	shuffled, err := json.Marshal(reference)
	if err != nil {
		t.Fatal(err)
	}

	secondRoot := filepath.Join(t.TempDir(), "standard-library")
	client := fixtureClient(t, map[string]string{testAPIPath: string(shuffled)})
	if err := Generate(context.Background(), client, Options{IndexURL: testIndexURL, Version: "2.0.0-alpha.46", OutputDir: secondRoot}); err != nil {
		t.Fatal(err)
	}
	if got := snapshot(t, secondRoot); !reflect.DeepEqual(got, first) {
		t.Fatal("shuffled API ordering changed generated output")
	}
}

func TestGenerateMovesFunctionsWhenCatalogRecategorizesThem(t *testing.T) {
	output := filepath.Join(t.TempDir(), "standard-library")
	catalog := mutateCatalog(t, func(value *apicatalog.Catalog) {
		value.Categories[0].Functions[0], value.Categories[2].Functions[0] = value.Categories[2].Functions[0], value.Categories[0].Functions[0]
	})

	if err := Generate(context.Background(), fixtureClient(t, map[string]string{testCatalogPath: catalog}), Options{
		IndexURL:  testIndexURL,
		Version:   "2.0.0-alpha.46",
		OutputDir: output,
	}); err != nil {
		t.Fatal(err)
	}

	arrays := readFile(t, filepath.Join(output, "arrays", "_index.md"))
	if !strings.Contains(arrays, `href="#global-abs"`) || !strings.Contains(arrays, `id="global-abs"`) {
		t.Fatal("recategorized abs function is not rendered on its new Arrays category page")
	}
	if strings.Contains(arrays, `global-append`) {
		t.Fatal("Arrays category retained append after recategorization")
	}

	mathPage := readFile(t, filepath.Join(output, "math", "_index.md"))
	if !strings.Contains(mathPage, `href="#global-append"`) || strings.Contains(mathPage, `global-abs`) {
		t.Fatal("Math category does not reflect recategorized function membership")
	}
}

func TestGenerateRequiresCatalogAndPreservesPreviousOutput(t *testing.T) {
	client := fixtureClient(t, nil)
	delete(client.documents, testCatalogPath)
	output := filepath.Join(t.TempDir(), "standard-library")
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(output, "sentinel"), []byte("previous output"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := Generate(context.Background(), client, Options{IndexURL: testIndexURL, Version: "2.0.0-alpha.46", OutputDir: output})
	if err == nil || !strings.Contains(err.Error(), "404 Not Found") {
		t.Fatalf("Generate error = %v, want required catalog 404", err)
	}
	if got := readFile(t, filepath.Join(output, "sentinel")); got != "previous output" {
		t.Fatalf("previous output = %q", got)
	}
}

func TestBuildFunctionViewUsesUniqueDeterministicSignatureAnchors(t *testing.T) {
	function := api.Function{Signatures: []api.Signature{
		{Parameters: []api.Parameter{{Name: "value", Type: "String"}}, Description: "String overload."},
		{Parameters: []api.Parameter{{Name: "value", Type: "Number"}}, Description: "Number overload."},
	}}

	want := []string{
		"global-convert-signature-fixed-1-1",
		"global-convert-signature-fixed-1-2",
	}
	view := buildFunctionView(functionIdentity{Name: "convert"}, function)
	got := []string{view.Signatures[0].ID, view.Signatures[1].ID}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("signature anchors = %v, want %v", got, want)
	}

	reverse(function.Signatures)
	view = buildFunctionView(functionIdentity{Name: "convert"}, function)
	got = []string{view.Signatures[0].ID, view.Signatures[1].ID}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("reversed signature anchors = %v, want %v", got, want)
	}
}

func TestFunctionAnchorUsesQualifiedIdentity(t *testing.T) {
	tests := []struct {
		namespace string
		name      string
		want      string
	}{
		{name: "ABS", want: "global-abs"},
		{namespace: "IO::FS", name: "READ", want: "io-fs-read"},
		{namespace: "T::NOT", name: "EQ", want: "t-not-eq"},
	}

	for _, test := range tests {
		if got := functionAnchor(test.namespace, test.name); got != test.want {
			t.Errorf("functionAnchor(%q, %q) = %q, want %q", test.namespace, test.name, got, test.want)
		}
	}
}

func TestGenerateRejectsInvalidPublication(t *testing.T) {
	validIndex := readFile(t, "testdata/index.json")
	validAPI := readFile(t, "testdata/api.json")
	validCatalog := readFile(t, "testdata/catalog.json")
	tests := []struct {
		name      string
		version   string
		overrides map[string]string
		want      string
	}{
		{name: "missing configured version", version: "2.0.0-alpha.99", want: "is not published"},
		{name: "malformed index", version: "2.0.0-alpha.46", overrides: map[string]string{"/index.json": `{`}, want: "parse Ferret API index"},
		{name: "unsupported index schema", version: "2.0.0-alpha.46", overrides: map[string]string{"/index.json": strings.Replace(validIndex, `"schemaVersion": 1`, `"schemaVersion": 2`, 1)}, want: "unsupported Ferret API Reference schema version"},
		{name: "malformed API", version: "2.0.0-alpha.46", overrides: map[string]string{testAPIPath: `{`}, want: "parse Ferret API version"},
		{name: "unsupported API schema", version: "2.0.0-alpha.46", overrides: map[string]string{testAPIPath: strings.Replace(validAPI, `"schemaVersion": 1`, `"schemaVersion": 2`, 1)}, want: "unsupported Ferret API Reference schema version"},
		{name: "wrong API id", version: "2.0.0-alpha.46", overrides: map[string]string{testAPIPath: strings.Replace(validAPI, `"montferret/core"`, `"montferret/other"`, 1)}, want: `want "montferret/core"`},
		{name: "API version mismatch", version: "2.0.0-alpha.46", overrides: map[string]string{testAPIPath: strings.Replace(validAPI, `"2.0.0-alpha.46"`, `"2.0.0-alpha.45"`, 1)}, want: "want configured version"},
		{name: "cross origin href", version: "2.0.0-alpha.46", overrides: map[string]string{"/index.json": strings.Replace(validIndex, `"./versions/2.0.0-alpha.46/api.json"`, `"https://other.test/api.json"`, 1)}, want: "must stay on the Ferret API origin"},
		{name: "malformed catalog", version: "2.0.0-alpha.46", overrides: map[string]string{testCatalogPath: `{`}, want: "parse Ferret API Catalog"},
		{name: "unsupported catalog schema", version: "2.0.0-alpha.46", overrides: map[string]string{testCatalogPath: strings.Replace(validCatalog, `"schemaVersion": 1`, `"schemaVersion": 2`, 1)}, want: "unsupported API Catalog schema version"},
		{name: "wrong catalog id", version: "2.0.0-alpha.46", overrides: map[string]string{testCatalogPath: strings.Replace(validCatalog, `"montferret/core"`, `"montferret/other"`, 1)}, want: "does not match API id"},
		{name: "catalog version mismatch", version: "2.0.0-alpha.46", overrides: map[string]string{testCatalogPath: strings.Replace(validCatalog, `"2.0.0-alpha.46"`, `"2.0.0-alpha.45"`, 1)}, want: "does not match API version"},
		{name: "unknown catalog function", version: "2.0.0-alpha.46", overrides: map[string]string{testCatalogPath: strings.Replace(validCatalog, `"abs"`, `"missing"`, 1)}, want: "unknown function"},
		{name: "unknown catalog namespace", version: "2.0.0-alpha.46", overrides: map[string]string{testCatalogPath: mutateCatalog(t, func(value *apicatalog.Catalog) { value.Categories[1].Functions[0].Namespace = "io::missing" })}, want: "unknown API namespace"},
		{name: "uncategorized function", version: "2.0.0-alpha.46", overrides: map[string]string{testCatalogPath: mutateCatalog(t, func(value *apicatalog.Catalog) { value.Categories = value.Categories[1:] })}, want: "is not assigned"},
		{name: "uncategorized namespaced function", version: "2.0.0-alpha.46", overrides: map[string]string{testCatalogPath: mutateCatalog(t, func(value *apicatalog.Catalog) { value.Categories[1].Functions = value.Categories[1].Functions[1:] })}, want: "io::fs::read"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := filepath.Join(t.TempDir(), "standard-library")
			err := Generate(context.Background(), fixtureClient(t, test.overrides), Options{IndexURL: testIndexURL, Version: test.version, OutputDir: output})
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("Generate error = %v, want error containing %q", err, test.want)
			}
		})
	}
}

func TestGenerateRejectsHTTPAndRedirectFailures(t *testing.T) {
	catalogStatus := fixtureClient(t, nil)
	catalogStatus.statusCodes = map[string]int{testCatalogPath: http.StatusServiceUnavailable}
	catalogRedirect := fixtureClient(t, nil)
	catalogRedirect.redirects = map[string]string{testCatalogPath: "https://other.test/catalog.json"}

	tests := []struct {
		name   string
		client HTTPClient
		want   string
	}{
		{name: "index HTTP status", client: &recordingClient{documents: map[string]string{}}, want: "404 Not Found"},
		{name: "catalog non-404 status", client: catalogStatus, want: "503 Service Unavailable"},
		{name: "catalog cross-origin redirect", client: catalogRedirect, want: "redirected outside"},
		{
			name: "index cross-origin redirect",
			client: &recordingClient{
				documents: map[string]string{"/index.json": readFile(t, "testdata/index.json")},
				redirects: map[string]string{"/index.json": "https://other.test/index.json"},
			},
			want: "redirected outside",
		},
		{
			name:   "oversized response",
			client: &recordingClient{documents: map[string]string{"/index.json": strings.Repeat("x", maxDocumentSize+1)}},
			want:   "exceeded",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := Generate(context.Background(), test.client, Options{IndexURL: testIndexURL, Version: "2.0.0-alpha.46", OutputDir: filepath.Join(t.TempDir(), "standard-library")})
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("Generate error = %v, want error containing %q", err, test.want)
			}
		})
	}
}

func TestGeneratePreservesPreviousOutputOnCatalogAndRenderFailures(t *testing.T) {
	tests := []struct {
		name      string
		overrides map[string]string
		want      string
	}{
		{name: "invalid catalog", overrides: map[string]string{testCatalogPath: `{`}, want: "parse Ferret API Catalog"},
		{name: "route collision", overrides: map[string]string{testCatalogPath: mutateCatalog(t, func(value *apicatalog.Catalog) { value.Categories[0].ID = "t" })}, want: "route collides"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := filepath.Join(t.TempDir(), "standard-library")
			if err := os.MkdirAll(output, 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(output, "sentinel"), []byte("previous output"), 0o644); err != nil {
				t.Fatal(err)
			}

			err := Generate(context.Background(), fixtureClient(t, test.overrides), Options{IndexURL: testIndexURL, Version: "2.0.0-alpha.46", OutputDir: output})
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("Generate error = %v, want %q", err, test.want)
			}
			if got := readFile(t, filepath.Join(output, "sentinel")); got != "previous output" {
				t.Fatalf("previous output = %q", got)
			}
		})
	}
}

type recordingClient struct {
	documents   map[string]string
	redirects   map[string]string
	statusCodes map[string]int
	requests    []string
}

func fixtureClient(t *testing.T, overrides map[string]string) *recordingClient {
	t.Helper()
	documents := map[string]string{
		"/index.json":   readFile(t, "testdata/index.json"),
		testAPIPath:     readFile(t, "testdata/api.json"),
		testCatalogPath: readFile(t, "testdata/catalog.json"),
	}
	for path, document := range overrides {
		documents[path] = document
	}

	return &recordingClient{documents: documents}
}

func (client *recordingClient) Do(request *http.Request) (*http.Response, error) {
	client.requests = append(client.requests, request.URL.Path)
	document, ok := client.documents[request.URL.Path]
	statusCode := client.statusCodes[request.URL.Path]
	if statusCode == 0 {
		if ok {
			statusCode = http.StatusOK
		} else {
			statusCode = http.StatusNotFound
		}
	}
	if !ok {
		document = strings.ToLower(http.StatusText(statusCode))
	}

	responseRequest := request
	if target := client.redirects[request.URL.Path]; target != "" {
		redirectURL, _ := url.Parse(target)
		responseRequest = request.Clone(request.Context())
		responseRequest.URL = redirectURL
	}

	return &http.Response{
		StatusCode: statusCode,
		Status:     fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)),
		Body:       io.NopCloser(strings.NewReader(document)),
		Request:    responseRequest,
	}, nil
}

func mutateCatalog(t *testing.T, mutate func(*apicatalog.Catalog)) string {
	t.Helper()
	catalog, err := apicatalog.Parse([]byte(readFile(t, "testdata/catalog.json")))
	if err != nil {
		t.Fatal(err)
	}
	mutate(catalog)

	data, err := json.Marshal(catalog)
	if err != nil {
		t.Fatal(err)
	}

	return string(data)
}

func mutateReference(t *testing.T, mutate func(*api.Reference)) string {
	t.Helper()
	reference, err := api.Parse([]byte(readFile(t, "testdata/api.json")))
	if err != nil {
		t.Fatal(err)
	}
	mutate(reference)

	data, err := json.Marshal(reference)
	if err != nil {
		t.Fatal(err)
	}

	return string(data)
}

func readFile(t *testing.T, filename string) string {
	t.Helper()
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}

	return string(data)
}

func snapshot(t *testing.T, root string) map[string][]byte {
	t.Helper()
	result := make(map[string][]byte)
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			return nil
		}

		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		result[filepath.ToSlash(relative)] = bytes.Clone(data)

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	return result
}

func reverse[T any](values []T) {
	for left, right := 0, len(values)-1; left < right; left, right = left+1, right-1 {
		values[left], values[right] = values[right], values[left]
	}
}
