package registryroutes

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateCreatesModuleAndVersionShells(t *testing.T) {
	client := registryFixtureClient{
		"/":                                   `{"schemaVersion":1,"artifacts":{"modules":"./modules/index.json"}}`,
		"/modules/index.json":                 `{"schemaVersion":1,"modules":[{"id":"montferret/html","href":"./montferret/html/index.json"}]}`,
		"/modules/montferret/html/index.json": `{"schemaVersion":1,"id":"montferret/html","description":"kept additive","versions":[{"version":"1.0.0-rc.21","href":"./versions/1.0.0-rc.21/index.json"}]}`,
	}

	output := t.TempDir()
	shellPath := filepath.Join(output, "registry", "index.html")
	writeTestFile(t, shellPath, []byte(`<main data-pagefind-body>registry shell</main>`))

	if err := Generate(context.Background(), client, "https://registry.test/", shellPath, output); err != nil {
		t.Fatal(err)
	}

	for _, relative := range []string{
		"registry/montferret/html/index.html",
		"registry/montferret/html/1.0.0-rc.21/index.html",
	} {
		contents, err := os.ReadFile(filepath.Join(output, relative))
		if err != nil {
			t.Errorf("read %s: %v", relative, err)
			continue
		}
		if string(contents) != `<main data-pagefind-ignore="all">registry shell</main>` {
			t.Errorf("%s was not marked as excluded from Pagefind", relative)
		}
	}
}

func TestGenerateRejectsUnsafeCatalogData(t *testing.T) {
	tests := map[string]map[string]string{
		"cross-origin artifact": {
			"/": `{"schemaVersion":1,"artifacts":{"modules":"https://example.org/modules.json"}}`,
		},
		"artifact query": {
			"/": `{"schemaVersion":1,"artifacts":{"modules":"/modules/index.json?download=1"}}`,
		},
		"trailing JSON data": {
			"/": `{"schemaVersion":1,"artifacts":{"modules":"/modules/index.json"}} {}`,
		},
		"unsafe module ID": {
			"/":                   `{"schemaVersion":1,"artifacts":{"modules":"/modules/index.json"}}`,
			"/modules/index.json": `{"schemaVersion":1,"modules":[{"id":"../escape","href":"/module.json"}]}`,
		},
		"unsafe version": {
			"/":                   `{"schemaVersion":1,"artifacts":{"modules":"/modules/index.json"}}`,
			"/modules/index.json": `{"schemaVersion":1,"modules":[{"id":"montferret/html","href":"/module.json"}]}`,
			"/module.json":        `{"schemaVersion":1,"id":"montferret/html","versions":[{"version":"../escape","href":"/version.json"}]}`,
		},
	}

	for name, documents := range tests {
		t.Run(name, func(t *testing.T) {
			output := t.TempDir()
			shellPath := filepath.Join(output, "registry", "index.html")
			writeTestFile(t, shellPath, []byte(`<main data-pagefind-body>registry shell</main>`))
			err := Generate(context.Background(), registryFixtureClient(documents), "https://registry.test/", shellPath, output)
			if err == nil {
				t.Fatal("Generate succeeded for unsafe Registry data")
			}
			if _, statErr := os.Stat(filepath.Join(output, "escape")); !os.IsNotExist(statErr) {
				t.Fatalf("unsafe route escaped output: %v", statErr)
			}
		})
	}
}

func TestGenerateFailsWhenRegistryCannotBeEnumerated(t *testing.T) {
	output := t.TempDir()
	shellPath := filepath.Join(output, "registry", "index.html")
	writeTestFile(t, shellPath, []byte(`<main data-pagefind-body>registry shell</main>`))
	err := Generate(context.Background(), registryFixtureClient{}, "https://registry.test/", shellPath, output)
	if err == nil || !strings.Contains(err.Error(), "404") {
		t.Fatalf("Generate error = %v, want Registry enumeration failure", err)
	}
}

type registryFixtureClient map[string]string

func (client registryFixtureClient) Do(request *http.Request) (*http.Response, error) {
	document, ok := client[request.URL.Path]
	statusCode := http.StatusOK
	status := "200 OK"
	if !ok {
		document = "not found"
		statusCode = http.StatusNotFound
		status = "404 Not Found"
	}
	return &http.Response{
		StatusCode: statusCode,
		Status:     status,
		Body:       io.NopCloser(strings.NewReader(document)),
		Request:    request,
	}, nil
}

func writeTestFile(t *testing.T, filename string, contents []byte) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(filename), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filename, contents, 0o644); err != nil {
		t.Fatal(err)
	}
}
