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
	if _, err := os.Stat("../studio"); os.IsNotExist(err) {
		t.Skip("Skipping integration test: ../studio directory not found")
	}
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
	if _, err := os.Stat("../totem"); os.IsNotExist(err) {
		t.Skip("Skipping integration test: ../totem directory not found")
	}
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
	if _, err := os.Stat("../totem"); os.IsNotExist(err) {
		t.Skip("Skipping integration test: ../totem directory not found")
	}
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

func TestIntegrationSilentFailureParsing(t *testing.T) {
	bin := buildBinary(t)
	defer os.Remove(bin)

	// Create a temp dir for our mock project
	tempDir, err := os.MkdirTemp("", "gfa-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a mock gradlew script
	mockGradlewPath := filepath.Join(tempDir, "gradlew")
	mockGradlewScript := `#!/bin/sh
echo "> Task :apps:forge:shared:checkKotlinGradlePluginConfigurationErrors"
echo "e: ❌ KMP Dependencies Resolution Failure"
echo "Source set 'commonMain' couldn't resolve dependencies for all target platforms"
echo "BUILD SUCCESSFUL in 2s"
exit 0
`
	if err := os.WriteFile(mockGradlewPath, []byte(mockGradlewScript), 0755); err != nil {
		t.Fatalf("Failed to write mock gradlew: %v", err)
	}

	cmd := exec.Command(bin, "build")
	cmd.Dir = tempDir
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
	if !strings.Contains(outStr, "KMP Dependencies Resolution Failure") {
		t.Errorf("Expected error to mention the resolution failure, got: %s", outStr)
	}
}
