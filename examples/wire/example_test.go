package wire_test

import (
	"bufio"
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestDocumentedPrograms(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	binaries := t.TempDir()
	for _, program := range []string{"host", "client"} {
		cmd := exec.CommandContext(ctx, "go", "build", "-o", filepath.Join(binaries, program), "./"+program)
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("build %s: %v\n%s", program, err, output)
		}
	}

	host := exec.CommandContext(ctx, filepath.Join(binaries, "host"), "-listen", "127.0.0.1:0")
	var hostErrors bytes.Buffer
	host.Stderr = &hostErrors
	stdout, err := host.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := host.Start(); err != nil {
		t.Fatal(err)
	}
	waited := false
	t.Cleanup(func() {
		if !waited {
			_ = host.Process.Kill()
			_ = host.Wait()
		}
	})

	ready := make(chan string, 1)
	go func() {
		scanner := bufio.NewScanner(stdout)
		if scanner.Scan() {
			ready <- scanner.Text()
		} else {
			ready <- ""
		}
	}()
	var address string
	select {
	case line := <-ready:
		if !strings.HasPrefix(line, "Listening on 127.0.0.1:") {
			t.Fatalf("unexpected host readiness: %q", line)
		}
		address = strings.TrimPrefix(line, "Listening on ")
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}

	for _, name := range []string{"Ada", "Grace"} {
		t.Run(name, func(t *testing.T) {
			cmd := exec.CommandContext(ctx, filepath.Join(binaries, "client"), "-address", address, "-name", name)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("client: %v\n%s", err, output)
			}
			want := "Content-Type: application/json\n\"Hello, " + name + "!\"\n"
			if string(output) != want {
				t.Fatalf("output = %q, want %q", output, want)
			}
		})
	}

	t.Run("expired context", func(t *testing.T) {
		cmd := exec.CommandContext(ctx, filepath.Join(binaries, "client"), "-address", address, "-timeout", "0s")
		output, err := cmd.CombinedOutput()
		if err == nil || !bytes.Contains(output, []byte("context deadline exceeded")) {
			t.Fatalf("expired client = %v, %s", err, output)
		}
	})

	if err := host.Process.Signal(os.Interrupt); err != nil {
		t.Fatal(err)
	}
	err = host.Wait()
	waited = true
	if err != nil || hostErrors.Len() != 0 {
		t.Fatalf("host shutdown = %v, %s", err, hostErrors.String())
	}

	t.Run("unavailable endpoint", func(t *testing.T) {
		cmd := exec.CommandContext(ctx, filepath.Join(binaries, "client"), "-address", address, "-timeout", "200ms")
		output, err := cmd.CombinedOutput()
		if err == nil || len(output) == 0 || bytes.Contains(output, []byte("Content-Type:")) {
			t.Fatalf("unavailable client = %v, %s", err, output)
		}
	})
}
