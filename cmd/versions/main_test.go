package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunWritesDeterministicJSON(t *testing.T) {
	filename := writeCommandFixture(t, "cli:\n    v2: \"2.0.0-alpha.40\"\n")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := run([]string{
		"update",
		"--repository", "MontFerret/cli",
		"--tag", "v2.0.0-alpha.41",
		"--file", filename,
		"--format", "json",
	}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("run exit code = %d, stderr = %s", exitCode, stderr.String())
	}
	want := "{\"status\":\"updated\",\"repository\":\"MontFerret/cli\",\"product\":\"cli\",\"major\":2,\"path\":\"cli.v2\",\"previous\":\"2.0.0-alpha.40\",\"version\":\"2.0.0-alpha.41\"}\n"
	if stdout.String() != want {
		t.Fatalf("stdout = %q, want %q", stdout.String(), want)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestRunUsesDefaultVersionFile(t *testing.T) {
	directory := t.TempDir()
	dataDirectory := filepath.Join(directory, "data")
	if err := os.MkdirAll(dataDirectory, 0o755); err != nil {
		t.Fatal(err)
	}
	filename := filepath.Join(dataDirectory, "versions.yaml")
	if err := os.WriteFile(filename, []byte("lab:\n    v2: \"2.0.0-alpha.31\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Chdir(directory)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := run([]string{"update", "--repository", "MontFerret/lab", "--tag", "v2.0.0-alpha.32"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("run exit code = %d, stderr = %s", exitCode, stderr.String())
	}
	if stdout.String() != "updated lab.v2: 2.0.0-alpha.31 -> 2.0.0-alpha.32\n" {
		t.Fatalf("stdout = %q", stdout.String())
	}
	contents, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(contents), `v2: "2.0.0-alpha.32"`) {
		t.Fatalf("default version file was not updated:\n%s", contents)
	}
}

func TestRunKeepsStdoutEmptyOnFailure(t *testing.T) {
	filename := writeCommandFixture(t, "cli:\n    v2: \"2.0.0-alpha.40\"\n")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := run([]string{
		"update",
		"--repository", "MontFerret/unknown",
		"--tag", "v2.0.0",
		"--file", filename,
		"--format", "json",
	}, &stdout, &stderr)
	if exitCode == 0 {
		t.Fatal("run succeeded for an unknown repository")
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
	if !strings.Contains(stderr.String(), "unknown repository") {
		t.Fatalf("stderr = %q, want unknown repository diagnostic", stderr.String())
	}
}

func TestRunReportsStaleReleaseAsSuccess(t *testing.T) {
	filename := writeCommandFixture(t, "worker:\n    v2: \"2.0.0-rc.27\"\n")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := run([]string{
		"update",
		"--repository", "MontFerret/worker",
		"--tag", "v2.0.0-rc.26",
		"--file", filename,
		"--format", "json",
	}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("run exit code = %d, stderr = %s", exitCode, stderr.String())
	}
	if !strings.HasPrefix(stdout.String(), `{"status":"stale"`) {
		t.Fatalf("stdout = %q, want stale JSON", stdout.String())
	}
	if !strings.Contains(stderr.String(), "warning: ignored stale release") {
		t.Fatalf("stderr = %q, want stale warning", stderr.String())
	}
}

func TestRunRejectsMissingArgumentsAndFormats(t *testing.T) {
	tests := [][]string{
		{"update", "--tag", "v2.0.0"},
		{"update", "--repository", "MontFerret/cli"},
		{"update", "--repository", "MontFerret/cli", "--tag", "v2.0.0", "--format", "yaml"},
		{"unknown"},
	}
	for _, arguments := range tests {
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		if exitCode := run(arguments, &stdout, &stderr); exitCode == 0 {
			t.Fatalf("run(%v) succeeded", arguments)
		}
		if stdout.Len() != 0 {
			t.Fatalf("run(%v) stdout = %q, want empty", arguments, stdout.String())
		}
		if stderr.Len() == 0 {
			t.Fatalf("run(%v) did not report an error", arguments)
		}
	}
}

func writeCommandFixture(t *testing.T, contents string) string {
	t.Helper()
	filename := filepath.Join(t.TempDir(), "versions.yaml")
	if err := os.WriteFile(filename, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
	return filename
}
