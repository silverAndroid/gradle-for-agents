package main_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func buildBinary(t *testing.T) string {
	t.Helper()
	cmd := exec.Command("go", "build", "-o", "gfa-test-bin", "./cmd/gfa")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build gfa binary: %v", err)
	}
	absPath, err := filepath.Abs("gfa-test-bin")
	if err != nil {
		t.Fatalf("Failed to get abs path: %v", err)
	}
	return absPath
}

func TestIntegrationStudio(t *testing.T) {
	bin := buildBinary(t)
	defer os.Remove(bin)

	cmd := exec.Command(bin, "tasks")
	cmd.Dir = "../studio"
	out, err := cmd.CombinedOutput()
	
	if err != nil {
		t.Fatalf("gfa failed on studio: %v\nOutput: %s", err, out)
	}
	if !strings.Contains(string(out), "SUCCESS: Build completed.") {
		t.Errorf("Expected success output from gfa, got: %s", out)
	}
}

func TestIntegrationTotem(t *testing.T) {
	bin := buildBinary(t)
	defer os.Remove(bin)

	cmd := exec.Command(bin, "tasks")
	cmd.Dir = "../totem"
	out, err := cmd.CombinedOutput()
	
	if err != nil {
		t.Fatalf("gfa failed on totem: %v\nOutput: %s", err, out)
	}
	if !strings.Contains(string(out), "SUCCESS: Build completed.") {
		t.Errorf("Expected success output from gfa, got: %s", out)
	}
}

func TestIntegrationErrorParsing(t *testing.T) {
	bin := buildBinary(t)
	defer os.Remove(bin)

	// We run a command that we know will fail in totem, e.g. a non-existent task
	cmd := exec.Command(bin, "thisTaskDoesNotExist")
	cmd.Dir = "../totem"
	out, err := cmd.CombinedOutput()
	
	if err == nil {
		t.Fatalf("Expected gfa to fail, but it succeeded.\nOutput: %s", out)
	}
	outStr := string(out)
	if !strings.Contains(outStr, "FAILURE: Build failed with exit code 1.") {
		t.Errorf("Expected FAILURE output from gfa, got: %s", outStr)
	}
	if !strings.Contains(outStr, "--- TOON OUTPUT ---") {
		t.Errorf("Expected TOON OUTPUT section, got: %s", outStr)
	}
	if !strings.Contains(outStr, "Task 'thisTaskDoesNotExist' not found") {
		t.Errorf("Expected error to mention the missing task, got: %s", outStr)
	}
}
