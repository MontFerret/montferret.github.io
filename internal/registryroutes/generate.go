package registryroutes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const schemaVersion = 1

var routeSegmentPattern = regexp.MustCompile(`^[A-Za-z0-9](?:[A-Za-z0-9._-]*[A-Za-z0-9])?$`)

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type rootDocument struct {
	SchemaVersion int               `json:"schemaVersion"`
	Artifacts     map[string]string `json:"artifacts"`
}

type moduleCatalog struct {
	SchemaVersion int                 `json:"schemaVersion"`
	Modules       []moduleCatalogItem `json:"modules"`
}

type moduleCatalogItem struct {
	ID   string `json:"id"`
	Href string `json:"href"`
}

type moduleDocument struct {
	SchemaVersion int                     `json:"schemaVersion"`
	ID            string                  `json:"id"`
	Versions      []moduleDocumentVersion `json:"versions"`
}

type moduleDocumentVersion struct {
	Version string `json:"version"`
	Href    string `json:"href"`
}

// Generate enumerates the live Registry and copies the built Registry shell to
// every published module and version route. Call it after Pagefind so duplicate
// shells are not added to global search.
func Generate(ctx context.Context, client httpClient, baseURL, shellPath, outputRoot string) error {
	base, err := url.Parse(baseURL)
	if err != nil || base.Scheme != "https" && base.Scheme != "http" || base.Host == "" {
		return fmt.Errorf("invalid Registry base URL %q", baseURL)
	}

	shell, err := os.ReadFile(shellPath)
	if err != nil {
		return fmt.Errorf("read Registry route shell: %w", err)
	}
	deepRouteShell := bytes.Replace(shell, []byte(" data-pagefind-body"), []byte(` data-pagefind-ignore="all"`), 1)
	if bytes.Equal(deepRouteShell, shell) {
		return fmt.Errorf("read Registry route shell: missing Pagefind body marker")
	}

	var root rootDocument
	if err := getJSON(ctx, client, base, base, &root); err != nil {
		return fmt.Errorf("enumerate Registry: %w", err)
	}
	if root.SchemaVersion != schemaVersion || root.Artifacts["modules"] == "" {
		return fmt.Errorf("enumerate Registry: root document is missing the v1 modules artifact")
	}

	modulesURL, err := resolveArtifact(base, root.Artifacts["modules"])
	if err != nil {
		return fmt.Errorf("enumerate Registry modules: %w", err)
	}

	var catalog moduleCatalog
	if err := getJSON(ctx, client, base, modulesURL, &catalog); err != nil {
		return fmt.Errorf("enumerate Registry modules: %w", err)
	}
	if catalog.SchemaVersion != schemaVersion {
		return fmt.Errorf("enumerate Registry modules: unsupported schema version %d", catalog.SchemaVersion)
	}

	for _, entry := range catalog.Modules {
		owner, name, err := parseModuleID(entry.ID)
		if err != nil {
			return fmt.Errorf("enumerate Registry module %q: %w", entry.ID, err)
		}

		moduleURL, err := resolveArtifact(modulesURL, entry.Href)
		if err != nil {
			return fmt.Errorf("enumerate Registry module %q: %w", entry.ID, err)
		}

		var module moduleDocument
		if err := getJSON(ctx, client, base, moduleURL, &module); err != nil {
			return fmt.Errorf("enumerate Registry module %q: %w", entry.ID, err)
		}
		if module.SchemaVersion != schemaVersion || module.ID != entry.ID {
			return fmt.Errorf("enumerate Registry module %q: document identity or schema version does not match", entry.ID)
		}

		if err := writeShell(outputRoot, deepRouteShell, owner, name); err != nil {
			return err
		}
		for _, version := range module.Versions {
			if !validRouteSegment(version.Version) || version.Href == "" {
				return fmt.Errorf("enumerate Registry module %q: invalid version %q", entry.ID, version.Version)
			}
			if _, err := resolveArtifact(moduleURL, version.Href); err != nil {
				return fmt.Errorf("enumerate Registry module %q version %q: %w", entry.ID, version.Version, err)
			}
			if err := writeShell(outputRoot, deepRouteShell, owner, name, version.Version); err != nil {
				return err
			}
		}
	}

	return nil
}

func getJSON(ctx context.Context, client httpClient, registryBase, endpoint *url.URL, destination any) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return err
	}
	request.Header.Set("Accept", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("GET %s returned %s", endpoint, response.Status)
	}
	if response.Request != nil && response.Request.URL != nil && !sameOrigin(registryBase, response.Request.URL) {
		return fmt.Errorf("GET %s redirected outside the Registry origin", endpoint)
	}

	decoder := json.NewDecoder(io.LimitReader(response.Body, 4<<20))
	if err := decoder.Decode(destination); err != nil {
		return fmt.Errorf("decode %s: %w", endpoint, err)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return fmt.Errorf("decode %s: trailing JSON data", endpoint)
	}
	return nil
}

func resolveArtifact(base *url.URL, href string) (*url.URL, error) {
	reference, err := url.Parse(href)
	if err != nil {
		return nil, fmt.Errorf("invalid artifact URL %q", href)
	}
	resolved := base.ResolveReference(reference)
	if !sameOrigin(base, resolved) || resolved.User != nil || resolved.RawQuery != "" || resolved.Fragment != "" {
		return nil, fmt.Errorf("artifact URL %q must stay on the Registry origin", href)
	}
	return resolved, nil
}

func sameOrigin(left, right *url.URL) bool {
	return strings.EqualFold(left.Scheme, right.Scheme) && strings.EqualFold(left.Host, right.Host)
}

func parseModuleID(id string) (string, string, error) {
	parts := strings.Split(id, "/")
	if len(parts) != 2 || !validRouteSegment(parts[0]) || !validRouteSegment(parts[1]) {
		return "", "", fmt.Errorf("module ID must contain two safe route segments")
	}
	return parts[0], parts[1], nil
}

func validRouteSegment(value string) bool {
	return value != "." && value != ".." && routeSegmentPattern.MatchString(value)
}

func writeShell(outputRoot string, shell []byte, segments ...string) error {
	targetDirectory := filepath.Join(append([]string{outputRoot, "registry"}, segments...)...)
	root, err := filepath.Abs(outputRoot)
	if err != nil {
		return err
	}
	target, err := filepath.Abs(targetDirectory)
	if err != nil {
		return err
	}
	if target != filepath.Join(root, "registry") && !strings.HasPrefix(target, filepath.Join(root, "registry")+string(filepath.Separator)) {
		return fmt.Errorf("refusing to write Registry shell outside output root")
	}
	if err := os.MkdirAll(target, 0o755); err != nil {
		return fmt.Errorf("create Registry route %q: %w", strings.Join(segments, "/"), err)
	}
	if err := os.WriteFile(filepath.Join(target, "index.html"), shell, 0o644); err != nil {
		return fmt.Errorf("write Registry route %q: %w", strings.Join(segments, "/"), err)
	}
	return nil
}
