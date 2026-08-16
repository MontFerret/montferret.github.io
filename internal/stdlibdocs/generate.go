package stdlibdocs

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/MontFerret/specs/pkg/api"
	apicatalog "github.com/MontFerret/specs/pkg/api/catalog"
)

const (
	// IndexURL is the published Ferret Core API Reference discovery document.
	IndexURL        = "https://ferretlang.org/ferret/index.json"
	moduleID        = "montferret/core"
	maxDocumentSize = 4 << 20
)

type (
	httpClient interface {
		Do(*http.Request) (*http.Response, error)
	}

	// Options selects one published Ferret API and its generated Hugo content root.
	Options struct {
		IndexURL  string
		Version   string
		OutputDir string
	}
)

// Generate loads one exact published Ferret Core API and atomically renders its
// unversioned Hugo content tree.
func Generate(ctx context.Context, client httpClient, options Options) error {
	if client == nil {
		return fmt.Errorf("generate Standard Library: HTTP client is required")
	}

	if strings.TrimSpace(options.Version) == "" {
		return fmt.Errorf("generate Standard Library: configured Ferret version is required")
	}

	if strings.TrimSpace(options.OutputDir) == "" {
		return fmt.Errorf("generate Standard Library: output directory is required")
	}

	indexURL := options.IndexURL
	if indexURL == "" {
		indexURL = IndexURL
	}

	indexEndpoint, err := parseIndexURL(indexURL)
	if err != nil {
		return fmt.Errorf("generate Standard Library: %w", err)
	}

	indexData, err := getDocument(ctx, client, indexEndpoint, indexEndpoint)
	if err != nil {
		return fmt.Errorf("load Ferret API index: %w", err)
	}

	index, err := api.ParseIndex(indexData)
	if err != nil {
		return fmt.Errorf("parse Ferret API index: %w", err)
	}

	entry, ok := exactVersion(index, options.Version)
	if !ok {
		return fmt.Errorf("configured Ferret API version %q is not published", options.Version)
	}

	apiEndpoint, err := resolveArtifact(indexEndpoint, entry.Href)
	if err != nil {
		return fmt.Errorf("resolve Ferret API version %q: %w", options.Version, err)
	}

	referenceData, err := getDocument(ctx, client, indexEndpoint, apiEndpoint)
	if err != nil {
		return fmt.Errorf("load Ferret API version %q: %w", options.Version, err)
	}

	reference, err := api.Parse(referenceData)
	if err != nil {
		return fmt.Errorf("parse Ferret API version %q: %w", options.Version, err)
	}

	if reference.ID != moduleID {
		return fmt.Errorf("Ferret API id is %q, want %q", reference.ID, moduleID)
	}

	if reference.Version != options.Version {
		return fmt.Errorf("Ferret API version is %q, want configured version %q", reference.Version, options.Version)
	}

	catalogEndpoint, err := resolveCatalogArtifact(apiEndpoint)
	if err != nil {
		return fmt.Errorf("resolve Ferret API Catalog version %q: %w", options.Version, err)
	}

	catalogData, err := getDocument(ctx, client, indexEndpoint, catalogEndpoint)
	if err != nil {
		var statusErr *httpStatusError
		if errors.As(err, &statusErr) && statusErr.StatusCode == http.StatusNotFound {
			return renderAtomic(options.OutputDir, reference, nil)
		}

		return fmt.Errorf("load Ferret API Catalog version %q: %w", options.Version, err)
	}

	catalog, err := apicatalog.Parse(catalogData)
	if err != nil {
		return fmt.Errorf("parse Ferret API Catalog version %q: %w", options.Version, err)
	}

	if err := validateCatalogAgainstReference(reference, catalog); err != nil {
		return fmt.Errorf("validate Ferret API Reference and Catalog version %q: %w", options.Version, err)
	}

	return renderAtomic(options.OutputDir, reference, catalog)
}

func parseIndexURL(value string) (*url.URL, error) {
	endpoint, err := url.Parse(value)
	if err != nil {
		return nil, fmt.Errorf("invalid Ferret API index URL %q: %w", value, err)
	}

	if endpoint.Scheme != "https" && endpoint.Scheme != "http" || endpoint.Host == "" || endpoint.User != nil || endpoint.RawQuery != "" || endpoint.Fragment != "" {
		return nil, fmt.Errorf("invalid Ferret API index URL %q", value)
	}

	return endpoint, nil
}

func exactVersion(index *api.Index, version string) (api.IndexVersion, bool) {
	for _, entry := range index.Versions {
		if entry.Version == version {
			return entry, true
		}
	}

	return api.IndexVersion{}, false
}

func getDocument(ctx context.Context, client httpClient, origin, endpoint *url.URL) ([]byte, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Accept", "application/json")
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.Request != nil && response.Request.URL != nil && !sameOrigin(origin, response.Request.URL) {
		return nil, fmt.Errorf("GET %s redirected outside the Ferret API origin", endpoint)
	}

	if response.StatusCode != http.StatusOK {
		return nil, &httpStatusError{Endpoint: endpoint.String(), StatusCode: response.StatusCode, Status: response.Status}
	}

	data, err := io.ReadAll(io.LimitReader(response.Body, maxDocumentSize+1))
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", endpoint, err)
	}

	if len(data) > maxDocumentSize {
		return nil, fmt.Errorf("GET %s exceeded the %d byte limit", endpoint, maxDocumentSize)
	}

	return data, nil
}

func resolveCatalogArtifact(apiEndpoint *url.URL) (*url.URL, error) {
	resolved := apiEndpoint.ResolveReference(&url.URL{Path: "catalog.json"})
	if !sameOrigin(apiEndpoint, resolved) || resolved.User != nil || resolved.RawQuery != "" || resolved.Fragment != "" {
		return nil, fmt.Errorf("catalog URL derived from %q must stay on the Ferret API origin", apiEndpoint)
	}

	return resolved, nil
}

func resolveArtifact(indexURL *url.URL, href string) (*url.URL, error) {
	reference, err := url.Parse(href)
	if err != nil {
		return nil, fmt.Errorf("invalid artifact URL %q: %w", href, err)
	}

	resolved := indexURL.ResolveReference(reference)
	if !sameOrigin(indexURL, resolved) || resolved.User != nil || resolved.RawQuery != "" || resolved.Fragment != "" {
		return nil, fmt.Errorf("artifact URL %q must stay on the Ferret API origin", href)
	}

	return resolved, nil
}

func sameOrigin(left, right *url.URL) bool {
	return strings.EqualFold(left.Scheme, right.Scheme) && strings.EqualFold(left.Host, right.Host)
}

func renderAtomic(outputDir string, reference *api.Reference, catalog *apicatalog.Catalog) error {
	parent := filepath.Dir(outputDir)
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return fmt.Errorf("create generated content parent: %w", err)
	}

	staging, err := os.MkdirTemp(parent, ".stdlib-staging-")
	if err != nil {
		return fmt.Errorf("create generated content staging directory: %w", err)
	}

	defer os.RemoveAll(staging)

	if err := renderReference(staging, reference, catalog); err != nil {
		return err
	}

	backup := staging + "-previous"
	hadPrevious := false
	if _, err := os.Lstat(outputDir); err == nil {
		if err := os.Rename(outputDir, backup); err != nil {
			return fmt.Errorf("stage previous Standard Library content: %w", err)
		}

		hadPrevious = true
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("inspect previous Standard Library content: %w", err)
	}

	if err := os.Rename(staging, outputDir); err != nil {
		if hadPrevious {
			_ = os.Rename(backup, outputDir)
		}

		return fmt.Errorf("publish generated Standard Library content: %w", err)
	}

	if hadPrevious {
		if err := os.RemoveAll(backup); err != nil {
			return fmt.Errorf("remove previous Standard Library content: %w", err)
		}
	}

	return nil
}
