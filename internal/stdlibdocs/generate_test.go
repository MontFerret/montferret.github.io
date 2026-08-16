package stdlibdocs

import (
	"bytes"
	"context"
	"encoding/json"
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
)

const testIndexURL = "https://api.test/index.json"

func TestGenerateRendersPublishedAPI(t *testing.T) {
	client := fixtureClient(t, nil)
	output := filepath.Join(t.TempDir(), "standard-library")

	if err := Generate(context.Background(), client, Options{IndexURL: testIndexURL, Version: "2.0.0-alpha.46", OutputDir: output}); err != nil {
		t.Fatal(err)
	}

	wantRequests := []string{"/index.json", "/artifacts/core-alpha.46.json"}
	if got := client.requests; !reflect.DeepEqual(got, wantRequests) {
		t.Fatalf("requests = %v, want authoritative index href requests %v", got, wantRequests)
	}
	for _, filename := range []string{
		"functions/_index.md",
		"functions/flatten.md",
		"functions/to_number.md",
		"io/_index.md",
		"io/fs/_index.md",
		"io/fs/read.md",
		"io/net/_index.md",
		"io/net/http/_index.md",
		"io/net/http/get.md",
		"t/_index.md",
		"t/eq.md",
		"t/not/_index.md",
		"t/not/eq.md",
	} {
		if _, err := os.Stat(filepath.Join(output, filename)); err != nil {
			t.Errorf("generated page %s: %v", filename, err)
		}
	}

	flatten := readFile(t, filepath.Join(output, "functions", "flatten.md"))
	for _, expected := range []string{
		`id="api-function-global-flatten"`,
		`id="api-function-global-flatten-signature-fixed-1"`,
		`id="api-function-global-flatten-signature-fixed-2"`,
		`flatten(arr)`,
		`flatten(arr, depth)`,
		`<code class="stdlib-api-value-type">Any[]</code>`,
		`<dt>Returns</dt>`,
		`</div><div><dt>Returns</dt>`,
	} {
		if !strings.Contains(flatten, expected) {
			t.Errorf("flatten page does not contain %q", expected)
		}
	}

	toNumber := readFile(t, filepath.Join(output, "functions", "to_number.md"))
	if !strings.Contains(toNumber, "<dt>Throws</dt>") || !strings.Contains(toNumber, "TypeError") {
		t.Error("to_number page does not render thrown errors")
	}

	get := readFile(t, filepath.Join(output, "io", "net", "http", "get.md"))
	for _, expected := range []string{
		`api-function-named-io-net-http-get-signature-variadic`,
		`io::net::http::get(url, options...)`,
		`<span class="stdlib-api-parameter-kind">Variadic</span>`,
		`<strong>Deprecated.</strong>`,
	} {
		if !strings.Contains(get, expected) {
			t.Errorf("namespaced function page does not contain %q", expected)
		}
	}

	functions := readFile(t, filepath.Join(output, "functions", "_index.md"))
	for _, alias := range []string{"/docs/standard-library/arrays/", "/docs/stdlib/types/"} {
		if !strings.Contains(functions, alias) {
			t.Errorf("global functions section does not contain legacy alias %q", alias)
		}
	}
	testing := readFile(t, filepath.Join(output, "t", "_index.md"))
	if !strings.Contains(testing, "/docs/standard-library/testing/") {
		t.Error("t namespace does not contain the Testing compatibility alias")
	}
	filesystem := readFile(t, filepath.Join(output, "io", "fs", "_index.md"))
	if !strings.Contains(filesystem, "/docs/stdlib/io-fs/") {
		t.Error("io::fs namespace does not contain its compatibility alias")
	}
	httpSection := readFile(t, filepath.Join(output, "io", "net", "http", "_index.md"))
	for _, alias := range []string{"/docs/standard-library/io/http/", "/docs/stdlib/io-net-http/"} {
		if !strings.Contains(httpSection, alias) {
			t.Errorf("io::net::http namespace does not contain compatibility alias %q", alias)
		}
	}

	if got, want := flatten, readFile(t, "testdata/golden/flatten.md"); got != want {
		t.Error("generated flatten page does not match its golden file")
	}
}

func TestGenerateIsDeterministicAcrossOrderingAndRepeatedRuns(t *testing.T) {
	firstRoot := filepath.Join(t.TempDir(), "standard-library")
	client := fixtureClient(t, nil)
	options := Options{IndexURL: testIndexURL, Version: "2.0.0-alpha.46", OutputDir: firstRoot}
	if err := Generate(context.Background(), client, options); err != nil {
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
	if err := Generate(context.Background(), fixtureClient(t, map[string]string{"/artifacts/core-alpha.46.json": string(shuffled)}), Options{IndexURL: testIndexURL, Version: "2.0.0-alpha.46", OutputDir: secondRoot}); err != nil {
		t.Fatal(err)
	}
	if got := snapshot(t, secondRoot); !reflect.DeepEqual(got, first) {
		t.Fatal("shuffled source ordering changed generated output")
	}
}

func TestGenerateRejectsInvalidPublication(t *testing.T) {
	validIndex := readFile(t, "testdata/index.json")
	validAPI := readFile(t, "testdata/api.json")
	tests := []struct {
		name      string
		version   string
		overrides map[string]string
		want      string
	}{
		{name: "missing configured version", version: "2.0.0-alpha.99", want: "is not published"},
		{name: "malformed index", version: "2.0.0-alpha.46", overrides: map[string]string{"/index.json": `{`}, want: "parse Ferret API index"},
		{name: "unsupported index schema", version: "2.0.0-alpha.46", overrides: map[string]string{"/index.json": strings.Replace(validIndex, `"schemaVersion": 1`, `"schemaVersion": 2`, 1)}, want: "unsupported Ferret API Reference schema version"},
		{name: "malformed API", version: "2.0.0-alpha.46", overrides: map[string]string{"/artifacts/core-alpha.46.json": `{`}, want: "parse Ferret API version"},
		{name: "unsupported API schema", version: "2.0.0-alpha.46", overrides: map[string]string{"/artifacts/core-alpha.46.json": strings.Replace(validAPI, `"schemaVersion": 1`, `"schemaVersion": 2`, 1)}, want: "unsupported Ferret API Reference schema version"},
		{name: "wrong API id", version: "2.0.0-alpha.46", overrides: map[string]string{"/artifacts/core-alpha.46.json": strings.Replace(validAPI, `"montferret/core"`, `"montferret/other"`, 1)}, want: `want "montferret/core"`},
		{name: "API version mismatch", version: "2.0.0-alpha.46", overrides: map[string]string{"/artifacts/core-alpha.46.json": strings.Replace(validAPI, `"2.0.0-alpha.46"`, `"2.0.0-alpha.45"`, 1)}, want: "want configured version"},
		{name: "cross origin href", version: "2.0.0-alpha.46", overrides: map[string]string{"/index.json": strings.Replace(validIndex, `"./artifacts/core-alpha.46.json"`, `"https://other.test/api.json"`, 1)}, want: "must stay on the Ferret API origin"},
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
	tests := []struct {
		name   string
		client httpClient
		want   string
	}{
		{name: "HTTP status", client: &recordingClient{documents: map[string]string{}}, want: "404 Not Found"},
		{
			name: "cross origin redirect",
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

func TestGeneratePreservesPreviousOutputOnRenderFailure(t *testing.T) {
	output := filepath.Join(t.TempDir(), "standard-library")
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(output, "sentinel"), []byte("previous output"), 0o644); err != nil {
		t.Fatal(err)
	}

	apiData := readFile(t, "testdata/api.json")
	collision := strings.Replace(apiData, `"name": "to_number"`, `"name": "FLATTEN"`, 1)
	err := Generate(context.Background(), fixtureClient(t, map[string]string{"/artifacts/core-alpha.46.json": collision}), Options{IndexURL: testIndexURL, Version: "2.0.0-alpha.46", OutputDir: output})
	if err == nil || !strings.Contains(err.Error(), "route collides") {
		t.Fatalf("Generate error = %v, want route collision", err)
	}
	if got := readFile(t, filepath.Join(output, "sentinel")); got != "previous output" {
		t.Fatalf("previous output = %q", got)
	}
}

type recordingClient struct {
	documents map[string]string
	redirects map[string]string
	requests  []string
}

func fixtureClient(t *testing.T, overrides map[string]string) *recordingClient {
	t.Helper()
	documents := map[string]string{
		"/index.json":                       readFile(t, "testdata/index.json"),
		"/artifacts/core-alpha.46.json":     readFile(t, "testdata/api.json"),
		"/versions/2.0.0-alpha.45/api.json": readFile(t, "testdata/api.json"),
	}
	for path, document := range overrides {
		documents[path] = document
	}

	return &recordingClient{documents: documents}
}

func (client *recordingClient) Do(request *http.Request) (*http.Response, error) {
	client.requests = append(client.requests, request.URL.Path)
	document, ok := client.documents[request.URL.Path]
	statusCode := http.StatusOK
	status := "200 OK"
	if !ok {
		document = "not found"
		statusCode = http.StatusNotFound
		status = "404 Not Found"
	}
	responseRequest := request
	if target := client.redirects[request.URL.Path]; target != "" {
		redirectURL, _ := url.Parse(target)
		responseRequest = request.Clone(request.Context())
		responseRequest.URL = redirectURL
	}

	return &http.Response{
		StatusCode: statusCode,
		Status:     status,
		Body:       io.NopCloser(strings.NewReader(document)),
		Request:    responseRequest,
	}, nil
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
