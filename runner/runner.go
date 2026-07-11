package runner

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"runtime"

	"github.com/silverAndroid/gradle-for-agents/parser"
)

func getGradleExecutable() (string, error) {
	gradlewName := "gradlew"
	if runtime.GOOS == "windows" {
		gradlewName = "gradlew.bat"
	}
	
	if _, err := os.Stat(gradlewName); err == nil {
		if absPath, err := filepath.Abs(gradlewName); err == nil {
			return absPath, nil
		}
	}
	
	if path, err := exec.LookPath("gradle"); err == nil {
		return path, nil
	}
	
	return "", fmt.Errorf("could not find local '%s' or global 'gradle' command", gradlewName)
}

func Run(gradleArgs []string, showWarnings bool, passThrough bool) int {
	gradleExec, err := getGradleExecutable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		return 1
	}

	timestamp := time.Now().Format("20060102-150405")
	logDir := filepath.Join(os.TempDir(), fmt.Sprintf("gfa-%s", timestamp))
	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to create log directory: %v\n", err)
		return 1
	}

	logFile := filepath.Join(logDir, "full_output.log")
	f, err := os.Create(logFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to create log file: %v\n", err)
		return 1
	}
	defer f.Close()

	cmd := exec.Command(gradleExec, gradleArgs...)

	if passThrough {
		multiWriter := io.MultiWriter(os.Stdout, f)
		cmd.Stdout = multiWriter
		cmd.Stderr = multiWriter

		if err := cmd.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Failed to start Gradle: %v\n", err)
			return 1
		}

		err = cmd.Wait()
		exitCode := 0
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else if err != nil {
			exitCode = 1
		}
		return exitCode
	}

	toonParser := parser.NewParser(showWarnings)

	multiWriter := io.MultiWriter(f, toonParser)
	cmd.Stdout = multiWriter
	cmd.Stderr = multiWriter

	fmt.Printf("INFO: Logging full output to: %s\n", logFile)

	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				fmt.Print(".")
			case <-done:
				return
			}
		}
	}()

	if err := cmd.Start(); err != nil {
		close(done)
		fmt.Fprintf(os.Stderr, "ERROR: Failed to start Gradle: %v\n", err)
		return 1
	}

	err = cmd.Wait()
	close(done)
	fmt.Println("\nINFO: Build finished, generating summary...")

	exitCode := 0
	if exitError, ok := err.(*exec.ExitError); ok {
		exitCode = exitError.ExitCode()
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		exitCode = 1
	}

	// Print TOON output
	toonParser.PrintSummary(exitCode, logFile)

	return exitCode
}
