package versions

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const versionsFixture = `# Website release pins.
go: "1.25" # Hugo build toolchain

runtime:
    v2: "2.0.0-alpha.46" # Current runtime
    v1: "1.0.0"

cli:
    go: "1.26.5" # Unrelated CLI build toolchain
    v2: "2.0.0-alpha.40"
    v1: "1.11.1"

lab:
    v2: "2.0.0-alpha.31"

worker:
    v2: "2.0.0-rc.27"
`

func TestUpdateFileResolvesSupportedRepositoriesAndChannels(t *testing.T) {
	tests := []struct {
		name       string
		repository string
		tag        string
		path       string
		previous   string
		version    string
	}{
		{name: "runtime v1", repository: "MontFerret/ferret", tag: "v1.0.1", path: "runtime.v1", previous: "1.0.0", version: "1.0.1"},
		{name: "runtime v2", repository: "MontFerret/ferret", tag: "v2.0.0-alpha.47", path: "runtime.v2", previous: "2.0.0-alpha.46", version: "2.0.0-alpha.47"},
		{name: "CLI v1", repository: "MontFerret/cli", tag: "1.12.0", path: "cli.v1", previous: "1.11.1", version: "1.12.0"},
		{name: "CLI v2", repository: "MontFerret/cli", tag: "2.0.0-alpha.41", path: "cli.v2", previous: "2.0.0-alpha.40", version: "2.0.0-alpha.41"},
		{name: "Lab v2", repository: "MontFerret/lab", tag: "2.0.0-alpha.32", path: "lab.v2", previous: "2.0.0-alpha.31", version: "2.0.0-alpha.32"},
		{name: "Worker v2", repository: "MontFerret/worker", tag: "2.0.0-rc.28", path: "worker.v2", previous: "2.0.0-rc.27", version: "2.0.0-rc.28"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filename := writeVersionFixture(t, versionsFixture)
			result, err := UpdateFile(Request{Repository: test.repository, Tag: test.tag, File: filename})
			if err != nil {
				t.Fatal(err)
			}
			if result.Status != StatusUpdated || result.Path != test.path || result.Previous != test.previous || result.Version != test.version {
				t.Fatalf("UpdateFile result = %+v", result)
			}
			if result.Major != 1 && strings.HasSuffix(test.path, ".v1") {
				t.Fatalf("UpdateFile major = %d, want 1", result.Major)
			}
			if result.Major != 2 && strings.HasSuffix(test.path, ".v2") {
				t.Fatalf("UpdateFile major = %d, want 2", result.Major)
			}
		})
	}
}

func TestUpdateFileAcceptsSemVerReleaseForms(t *testing.T) {
	tests := []struct {
		name     string
		previous string
		tag      string
		want     string
	}{
		{name: "alpha", previous: "2.0.0-alpha.1", tag: "v2.0.0-alpha.2", want: "2.0.0-alpha.2"},
		{name: "beta", previous: "2.0.0-alpha.2", tag: "2.0.0-beta.1", want: "2.0.0-beta.1"},
		{name: "RC", previous: "2.0.0-beta.1", tag: "v2.0.0-rc.1", want: "2.0.0-rc.1"},
		{name: "stable", previous: "2.0.0-rc.1", tag: "2.0.0", want: "2.0.0"},
		{name: "build metadata", previous: "1.9.9", tag: "v2.0.0+linux.amd64", want: "2.0.0+linux.amd64"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filename := writeVersionFixture(t, "cli:\n    v2: \""+test.previous+"\"\n")
			result, err := UpdateFile(Request{Repository: "MontFerret/cli", Tag: test.tag, File: filename})
			if err != nil {
				t.Fatal(err)
			}
			if result.Status != StatusUpdated || result.Version != test.want {
				t.Fatalf("UpdateFile result = %+v", result)
			}
		})
	}
}

func TestUpdateFileRejectsInvalidRequests(t *testing.T) {
	tests := []struct {
		name       string
		repository string
		tag        string
		wantError  string
	}{
		{name: "unknown repository", repository: "MontFerret/unknown", tag: "v2.0.0", wantError: "unknown repository"},
		{name: "unsupported major", repository: "MontFerret/lab", tag: "v1.0.0", wantError: "unsupported repository/major combination"},
		{name: "second leading v", repository: "MontFerret/cli", tag: "vv2.0.0", wantError: "invalid release tag"},
		{name: "uppercase leading V", repository: "MontFerret/cli", tag: "V2.0.0", wantError: "invalid release tag"},
		{name: "leading zero", repository: "MontFerret/cli", tag: "v02.0.0", wantError: "invalid release tag"},
		{name: "partial version", repository: "MontFerret/cli", tag: "v2.0", wantError: "invalid release tag"},
		{name: "empty version", repository: "MontFerret/cli", tag: "", wantError: "invalid release tag"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filename := writeVersionFixture(t, versionsFixture)
			before := readFile(t, filename)
			_, err := UpdateFile(Request{Repository: test.repository, Tag: test.tag, File: filename})
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Fatalf("UpdateFile error = %v, want %q", err, test.wantError)
			}
			if after := readFile(t, filename); after != before {
				t.Fatal("invalid request changed the version file")
			}
		})
	}
}

func TestUpdateFileReturnsUnchangedAndStaleWithoutWriting(t *testing.T) {
	tests := []struct {
		name   string
		tag    string
		status Status
	}{
		{name: "unchanged", tag: "v2.0.0-alpha.40", status: StatusUnchanged},
		{name: "stale prerelease", tag: "v2.0.0-alpha.39", status: StatusStale},
		{name: "stale downgrade", tag: "v1.10.0", status: StatusStale},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filename := writeVersionFixture(t, versionsFixture)
			before := readFile(t, filename)
			result, err := updateFile(Request{Repository: "MontFerret/cli", Tag: test.tag, File: filename}, func(string, []byte, fs.FileMode) error {
				t.Fatal("writer called for no-op update")
				return nil
			})
			if err != nil {
				t.Fatal(err)
			}
			if result.Status != test.status {
				t.Fatalf("UpdateFile status = %q, want %q", result.Status, test.status)
			}
			if after := readFile(t, filename); after != before {
				t.Fatal("no-op update changed the version file")
			}
		})
	}
}

func TestUpdateFileRejectsEqualPrecedenceBuildMetadata(t *testing.T) {
	filename := writeVersionFixture(t, "cli:\n    v2: \"2.0.0+build.1\"\n")
	before := readFile(t, filename)
	_, err := UpdateFile(Request{Repository: "MontFerret/cli", Tag: "v2.0.0+build.2", File: filename})
	if err == nil || !strings.Contains(err.Error(), "equal SemVer precedence") {
		t.Fatalf("UpdateFile error = %v, want build metadata conflict", err)
	}
	if after := readFile(t, filename); after != before {
		t.Fatal("build metadata conflict changed the version file")
	}
}

func TestUpdateFilePreservesYAMLStructureAndUnrelatedValues(t *testing.T) {
	filename := writeVersionFixture(t, versionsFixture)
	if _, err := UpdateFile(Request{Repository: "MontFerret/cli", Tag: "v2.0.0-alpha.41", File: filename}); err != nil {
		t.Fatal(err)
	}
	contents := readFile(t, filename)
	want := strings.Replace(versionsFixture, `    v2: "2.0.0-alpha.40"`, `    v2: "2.0.0-alpha.41"`, 1)
	if contents != want {
		t.Fatalf("updated YAML changed bytes outside the target scalar:\n%s", contents)
	}
}

func TestUpdateFilePreservesCommentOnTargetScalar(t *testing.T) {
	const contents = "cli:\n    v2: \"2.0.0-alpha.9\" # Keep this release note\n"
	filename := writeVersionFixture(t, contents)
	if _, err := UpdateFile(Request{Repository: "MontFerret/cli", Tag: "v2.0.0-alpha.10", File: filename}); err != nil {
		t.Fatal(err)
	}
	want := "cli:\n    v2: \"2.0.0-alpha.10\" # Keep this release note\n"
	if updated := readFile(t, filename); updated != want {
		t.Fatalf("updated YAML = %q, want %q", updated, want)
	}
}

func TestUpdateFileRejectsInvalidVersionFiles(t *testing.T) {
	tests := []struct {
		name      string
		contents  string
		wantError string
	}{
		{name: "malformed YAML", contents: "cli: [\n", wantError: "parse version file"},
		{name: "multiple documents", contents: "cli:\n    v2: \"2.0.0\"\n---\ncli:\n    v2: \"2.0.1\"\n", wantError: "exactly one YAML document"},
		{name: "missing product", contents: "runtime:\n    v2: \"2.0.0\"\n", wantError: "top-level \"cli\" key"},
		{name: "duplicate product", contents: "cli:\n    v2: \"2.0.0\"\ncli:\n    v2: \"2.0.1\"\n", wantError: "top-level \"cli\" key"},
		{name: "missing channel", contents: "cli:\n    v1: \"1.0.0\"\n", wantError: "unsupported repository/major combination"},
		{name: "duplicate channel", contents: "cli:\n    v2: \"2.0.0\"\n    v2: \"2.0.1\"\n", wantError: "exactly one cli.v2 key"},
		{name: "non-string channel", contents: "cli:\n    v2: 2\n", wantError: "double-quoted string scalar"},
		{name: "unquoted string channel", contents: "cli:\n    v2: 2.0.0\n", wantError: "double-quoted string scalar"},
		{name: "invalid stored version", contents: "cli:\n    v2: \"latest\"\n", wantError: "normalized SemVer"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filename := writeVersionFixture(t, test.contents)
			before := readFile(t, filename)
			_, err := UpdateFile(Request{Repository: "MontFerret/cli", Tag: "v2.1.0", File: filename})
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Fatalf("UpdateFile error = %v, want %q", err, test.wantError)
			}
			if after := readFile(t, filename); after != before {
				t.Fatal("invalid version file changed")
			}
		})
	}
}

func TestUpdateFileLeavesOriginalIntactWhenAtomicWriteFails(t *testing.T) {
	filename := writeVersionFixture(t, versionsFixture)
	before := readFile(t, filename)
	wantError := errors.New("injected atomic write failure")
	_, err := updateFile(Request{Repository: "MontFerret/cli", Tag: "v2.0.0-alpha.41", File: filename}, func(string, []byte, fs.FileMode) error {
		return wantError
	})
	if !errors.Is(err, wantError) {
		t.Fatalf("UpdateFile error = %v, want %v", err, wantError)
	}
	if after := readFile(t, filename); after != before {
		t.Fatal("failed atomic write changed the original file")
	}
}

func TestWriteFileAtomicallyPreservesPermissions(t *testing.T) {
	filename := writeVersionFixture(t, "before")
	if err := os.Chmod(filename, 0o640); err != nil {
		t.Fatal(err)
	}
	if err := writeFileAtomically(filename, []byte("after"), 0o640); err != nil {
		t.Fatal(err)
	}
	if contents := readFile(t, filename); contents != "after" {
		t.Fatalf("contents = %q, want after", contents)
	}
	info, err := os.Stat(filename)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o640 {
		t.Fatalf("permissions = %o, want 640", info.Mode().Perm())
	}
}

func writeVersionFixture(t *testing.T, contents string) string {
	t.Helper()
	filename := filepath.Join(t.TempDir(), "versions.yaml")
	if err := os.WriteFile(filename, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
	return filename
}

func readFile(t *testing.T, filename string) string {
	t.Helper()
	contents, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	return string(contents)
}
