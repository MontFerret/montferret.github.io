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

	for _, filename := range []string{
		"arrays/_index.md",
		"io/_index.md",
		"math/_index.md",
		"strings/_index.md",
		"testing/_index.md",
		"types/_index.md",
		"functions/_index.md",
		"functions/abs.md",
		"functions/append.md",
		"functions/concat.md",
		"functions/flatten.md",
		"functions/to_number.md",
		"io/fs/read.md",
		"io/net/http/get.md",
		"t/eq.md",
		"t/not/eq.md",
	} {
		if _, err := os.Stat(filepath.Join(output, filename)); err != nil {
			t.Errorf("generated page %s: %v", filename, err)
		}
	}

	mathPage := readFile(t, filepath.Join(output, "math", "_index.md"))
	for _, expected := range []string{
		`aliases:`,
		`"/docs/stdlib/math/"`,
		`stdlibKind: "Category"`,
		`weight: 30`,
		`href="/docs/standard-library/functions/abs/"`,
		`abs returns the absolute value of a number.`,
	} {
		if !strings.Contains(mathPage, expected) {
			t.Errorf("Math category does not contain %q", expected)
		}
	}

	flatten := readFile(t, filepath.Join(output, "functions", "flatten.md"))
	for _, expected := range []string{
		`description: "flatten turns an array of arrays into a flat array."`,
		`class="stdlib-api-breadcrumbs"`,
		`href="/docs/standard-library/arrays/"`,
		`id="api-function-global-flatten"`,
		`id="api-function-global-flatten-signature-fixed-1"`,
		`id="api-function-global-flatten-signature-fixed-2"`,
		`flatten(arr)`,
		`flatten(arr, depth)`,
		`<code class="stdlib-api-value-type">Any[]</code>`,
		`<dt>Returns</dt>`,
	} {
		if !strings.Contains(flatten, expected) {
			t.Errorf("flatten page does not contain %q", expected)
		}
	}

	functions := readFile(t, filepath.Join(output, "functions", "_index.md"))
	for _, expected := range []string{
		`sidebarHidden: true`,
		`stdlibLandingHidden: true`,
		`Global functions are organized by category`,
	} {
		if !strings.Contains(functions, expected) {
			t.Errorf("Functions compatibility page does not contain %q", expected)
		}
	}
	if strings.Contains(functions, "stdlib-children") || strings.Contains(functions, "/docs/stdlib/math/") {
		t.Error("Functions compatibility page retained the flat listing or category aliases")
	}

	toNumber := readFile(t, filepath.Join(output, "functions", "to_number.md"))
	if !strings.Contains(toNumber, "<dt>Throws</dt>") || !strings.Contains(toNumber, "TypeError") {
		t.Error("to_number page does not render thrown errors")
	}

	get := readFile(t, filepath.Join(output, "io", "net", "http", "get.md"))
	for _, expected := range []string{
		`class="stdlib-api-breadcrumbs"`,
		`href="/docs/standard-library/io/"`,
		`aria-current="page">io::net::http::get`,
		`api-function-named-io-net-http-get-signature-variadic-2`,
		`io::net::http::get(url, options...)`,
		`<span class="stdlib-api-parameter-kind">Variadic</span>`,
		`<strong>Deprecated.</strong>`,
	} {
		if !strings.Contains(get, expected) {
			t.Errorf("namespaced function page does not contain %q", expected)
		}
	}

	ioCategory := readFile(t, filepath.Join(output, "io", "_index.md"))
	for _, expected := range []string{
		`stdlibKind: "Category"`,
		`weight: 20`,
		`href="/docs/standard-library/io/fs/read/"`,
		`<code>io::fs::read</code>`,
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
		`href="/docs/standard-library/t/eq/"`,
		`<code>t::eq</code>`,
		`href="/docs/standard-library/t/not/eq/"`,
		`<code>t::not::eq</code>`,
		`"/docs/standard-library/t/"`,
		`"/docs/standard-library/t/not/"`,
	} {
		if !strings.Contains(testingCategory, expected) {
			t.Errorf("Testing category does not contain %q", expected)
		}
	}

	for _, filename := range []string{"io/fs/_index.md", "io/net/_index.md", "io/net/http/_index.md", "t/_index.md", "t/not/_index.md"} {
		if _, err := os.Stat(filepath.Join(output, filename)); !os.IsNotExist(err) {
			t.Errorf("generated obsolete namespace section %s: %v", filename, err)
		}
	}

	if got, want := flatten, readFile(t, "testdata/golden/flatten.md"); got != want {
		t.Error("generated flatten page does not match its golden file")
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

func TestGenerateKeepsCanonicalFunctionRoutesAfterRecategorization(t *testing.T) {
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

	abs := readFile(t, filepath.Join(output, "functions", "abs.md"))
	if !strings.Contains(abs, `href="/docs/standard-library/arrays/"`) {
		t.Fatal("recategorized abs page does not use its new category breadcrumb")
	}
	if _, err := os.Stat(filepath.Join(output, "arrays", "abs.md")); !os.IsNotExist(err) {
		t.Fatalf("recategorization changed the canonical abs route: %v", err)
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

func TestFunctionSummaryUsesDeterministicSignatureOrdering(t *testing.T) {
	function := api.Function{Signatures: []api.Signature{
		{Parameters: []api.Parameter{{Name: "value", Type: "String"}}, Description: "String overload."},
		{Parameters: []api.Parameter{{Name: "value", Type: "Number"}}, Description: "Number overload."},
	}}

	if got, want := functionSummary(function), "Number overload."; got != want {
		t.Fatalf("function summary = %q, want %q", got, want)
	}

	reverse(function.Signatures)
	if got, want := functionSummary(function), "Number overload."; got != want {
		t.Fatalf("reversed function summary = %q, want %q", got, want)
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
		{name: "route collision", overrides: map[string]string{testCatalogPath: mutateCatalog(t, func(value *apicatalog.Catalog) { value.Categories[0].ID = "functions" })}, want: "route collides"},
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
