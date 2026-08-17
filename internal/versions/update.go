// Package versions updates the website's pinned Ferret tool versions.
package versions

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver/v3"
	"gopkg.in/yaml.v3"
)

// Status describes the result of an update request.
type Status string

const (
	StatusUpdated   Status = "updated"
	StatusUnchanged Status = "unchanged"
	StatusStale     Status = "stale"
)

// Request identifies a release and the website version file to update.
type Request struct {
	Repository string
	Tag        string
	File       string
}

// Result is the stable machine-readable result of an update request.
type Result struct {
	Status     Status `json:"status"`
	Repository string `json:"repository"`
	Product    string `json:"product"`
	Major      uint64 `json:"major"`
	Path       string `json:"path"`
	Previous   string `json:"previous"`
	Version    string `json:"version"`
}

var repositoryProducts = map[string]string{
	"MontFerret/ferret": "runtime",
	"MontFerret/cli":    "cli",
	"MontFerret/lab":    "lab",
	"MontFerret/worker": "worker",
}

// UpdateFile validates request and atomically advances its version channel.
func UpdateFile(request Request) (Result, error) {
	return updateFile(request, writeFileAtomically)
}

type fileWriter func(string, []byte, fs.FileMode) error

func updateFile(request Request, write fileWriter) (Result, error) {
	product, ok := repositoryProducts[request.Repository]
	if !ok {
		return Result{}, fmt.Errorf("unknown repository %q", request.Repository)
	}

	normalized, releaseVersion, err := parseReleaseTag(request.Tag)
	if err != nil {
		return Result{}, err
	}

	channel := fmt.Sprintf("v%d", releaseVersion.Major())
	path := product + "." + channel
	contents, err := os.ReadFile(request.File)
	if err != nil {
		return Result{}, fmt.Errorf("read version file %q: %w", request.File, err)
	}

	document, err := parseDocument(contents)
	if err != nil {
		return Result{}, fmt.Errorf("parse version file %q: %w", request.File, err)
	}

	target, err := findVersionNode(document, product, channel)
	if err != nil {
		return Result{}, fmt.Errorf("validate version file %q: %w", request.File, err)
	}

	previous := target.Value
	currentVersion, err := semver.StrictNewVersion(previous)
	if err != nil {
		return Result{}, fmt.Errorf("validate version file %q: %s must contain a normalized SemVer version: %w", request.File, path, err)
	}

	result := Result{
		Repository: request.Repository,
		Product:    product,
		Major:      releaseVersion.Major(),
		Path:       path,
		Previous:   previous,
		Version:    normalized,
	}

	if normalized == previous {
		result.Status = StatusUnchanged

		return result, nil
	}

	comparison := releaseVersion.Compare(currentVersion)
	if comparison < 0 {
		result.Status = StatusStale

		return result, nil
	}

	if comparison == 0 {
		return Result{}, fmt.Errorf("release %q and stored version %q have equal SemVer precedence but different build metadata", normalized, previous)
	}

	target.Value = normalized
	replacement, err := serializeScalar(target)
	if err != nil {
		return Result{}, fmt.Errorf("serialize version file %q: %w", request.File, err)
	}
	updatedContents, err := replaceDoubleQuotedScalar(contents, target, replacement)
	if err != nil {
		return Result{}, fmt.Errorf("update version file %q: %w", request.File, err)
	}

	info, err := os.Stat(request.File)
	if err != nil {
		return Result{}, fmt.Errorf("stat version file %q: %w", request.File, err)
	}

	if err := write(request.File, updatedContents, info.Mode()); err != nil {
		return Result{}, fmt.Errorf("write version file %q: %w", request.File, err)
	}

	result.Status = StatusUpdated

	return result, nil
}

func parseReleaseTag(tag string) (string, *semver.Version, error) {
	normalized := strings.TrimPrefix(tag, "v")

	version, err := semver.StrictNewVersion(normalized)
	if err != nil {
		return "", nil, fmt.Errorf("invalid release tag %q: %w", tag, err)
	}

	return normalized, version, nil
}

func parseDocument(contents []byte) (*yaml.Node, error) {
	decoder := yaml.NewDecoder(bytes.NewReader(contents))
	var document yaml.Node
	if err := decoder.Decode(&document); err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("document is empty")
		}

		return nil, err
	}

	var trailing yaml.Node
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			return nil, fmt.Errorf("expected exactly one YAML document")
		}

		return nil, err
	}

	if document.Kind != yaml.DocumentNode || len(document.Content) != 1 || document.Content[0].Kind != yaml.MappingNode {
		return nil, fmt.Errorf("root must be a mapping")
	}

	return &document, nil
}

func findVersionNode(document *yaml.Node, product, channel string) (*yaml.Node, error) {
	root := document.Content[0]
	products := mappingValues(root, product)

	if len(products) != 1 {
		return nil, fmt.Errorf("expected exactly one top-level %q key, found %d", product, len(products))
	}

	productNode := products[0]
	if productNode.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("%s must be a mapping", product)
	}

	versions := mappingValues(productNode, channel)
	if len(versions) == 0 {
		return nil, fmt.Errorf("unsupported repository/major combination: %s.%s does not exist", product, channel)
	}

	if len(versions) != 1 {
		return nil, fmt.Errorf("expected exactly one %s.%s key, found %d", product, channel, len(versions))
	}

	versionNode := versions[0]
	if versionNode.Kind != yaml.ScalarNode || versionNode.Tag != "!!str" || versionNode.Style != yaml.DoubleQuotedStyle {
		return nil, fmt.Errorf("%s.%s must be a double-quoted string scalar", product, channel)
	}

	return versionNode, nil
}

func serializeScalar(node *yaml.Node) ([]byte, error) {
	scalar := yaml.Node{
		Kind:  yaml.ScalarNode,
		Style: node.Style,
		Tag:   node.Tag,
		Value: node.Value,
	}
	var output bytes.Buffer
	encoder := yaml.NewEncoder(&output)
	encoder.SetIndent(4)
	if err := encoder.Encode(&scalar); err != nil {
		return nil, err
	}
	if err := encoder.Close(); err != nil {
		return nil, err
	}

	serialized := bytes.TrimSuffix(output.Bytes(), []byte("\n"))
	if len(serialized) < 2 || serialized[0] != '"' || serialized[len(serialized)-1] != '"' || bytes.ContainsAny(serialized, "\r\n") {
		return nil, fmt.Errorf("expected one double-quoted scalar")
	}

	return append([]byte(nil), serialized...), nil
}

func replaceDoubleQuotedScalar(contents []byte, node *yaml.Node, replacement []byte) ([]byte, error) {
	start, err := nodeOffset(contents, node.Line, node.Column)
	if err != nil {
		return nil, err
	}
	if start >= len(contents) || contents[start] != '"' {
		return nil, fmt.Errorf("expected double-quoted scalar at line %d, column %d", node.Line, node.Column)
	}

	escaped := false
	for end := start + 1; end < len(contents); end++ {
		switch contents[end] {
		case '\n', '\r':
			return nil, fmt.Errorf("double-quoted version scalar must not span lines")
		case '\\':
			escaped = !escaped
		case '"':
			if escaped {
				escaped = false
				continue
			}

			updated := make([]byte, 0, len(contents)-(end+1-start)+len(replacement))
			updated = append(updated, contents[:start]...)
			updated = append(updated, replacement...)
			updated = append(updated, contents[end+1:]...)
			return updated, nil
		default:
			escaped = false
		}
	}

	return nil, fmt.Errorf("unterminated double-quoted version scalar")
}

func nodeOffset(contents []byte, line, column int) (int, error) {
	if line < 1 || column < 1 {
		return 0, fmt.Errorf("invalid scalar position")
	}

	offset := 0
	for currentLine := 1; currentLine < line; currentLine++ {
		newline := bytes.IndexByte(contents[offset:], '\n')
		if newline < 0 {
			return 0, fmt.Errorf("scalar line %d is outside the document", line)
		}
		offset += newline + 1
	}
	offset += column - 1
	if offset > len(contents) {
		return 0, fmt.Errorf("scalar column %d is outside line %d", column, line)
	}

	return offset, nil
}

func mappingValues(mapping *yaml.Node, key string) []*yaml.Node {
	var values []*yaml.Node
	for index := 0; index+1 < len(mapping.Content); index += 2 {
		candidate := mapping.Content[index]

		if candidate.Kind == yaml.ScalarNode && candidate.Value == key {
			values = append(values, mapping.Content[index+1])
		}
	}

	return values
}

func writeFileAtomically(filename string, contents []byte, mode fs.FileMode) (err error) {
	temporary, err := os.CreateTemp(filepath.Dir(filename), "."+filepath.Base(filename)+".tmp-*")
	if err != nil {
		return err
	}

	temporaryName := temporary.Name()
	defer func() {
		if temporary != nil {
			_ = temporary.Close()
		}
		_ = os.Remove(temporaryName)
	}()

	if err := temporary.Chmod(mode.Perm()); err != nil {
		return err
	}

	if _, err := temporary.Write(contents); err != nil {
		return err
	}

	if err := temporary.Sync(); err != nil {
		return err
	}

	if err := temporary.Close(); err != nil {
		return err
	}

	temporary = nil

	return os.Rename(temporaryName, filename)
}
