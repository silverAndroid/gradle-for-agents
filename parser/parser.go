package parser

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

type Parser struct {
	showWarnings bool
	buffer       []byte

	errors   []string
	warnings []string

	inErrorContext bool
	errorContext   []string
	errorTask      string

	warningRegex  *regexp.Regexp
	taskFailRegex *regexp.Regexp
	errorRegex    *regexp.Regexp
}

func NewParser(showWarnings bool) *Parser {
	return &Parser{
		showWarnings:  showWarnings,
		warningRegex:  regexp.MustCompile(`(?i)^(w:\s|warning:)`),
		taskFailRegex: regexp.MustCompile(`> Task (:[^\s]+) FAILED`),
		errorRegex:    regexp.MustCompile(`(?i)^(e:\s|error:|FAILURE: Build failed)`),
	}
}

func (p *Parser) Write(b []byte) (n int, err error) {
	p.buffer = append(p.buffer, b...)
	for {
		idx := bytes.IndexByte(p.buffer, '\n')
		if idx == -1 {
			break
		}
		line := string(p.buffer[:idx])
		p.processLine(line)
		p.buffer = p.buffer[idx+1:]
	}
	return len(b), nil
}

func (p *Parser) processLine(line string) {
	stripped := strings.TrimRight(line, "\r")

	if p.inErrorContext {
		p.errorContext = append(p.errorContext, stripped)
		if len(p.errorContext) >= 5 {
			p.commitError()
		}
		return
	}

	if match := p.taskFailRegex.FindStringSubmatch(stripped); match != nil {
		p.inErrorContext = true
		p.errorTask = match[1]
		p.errorContext = append(p.errorContext, stripped)
		return
	}

	if p.errorRegex.MatchString(stripped) {
		p.inErrorContext = true
		p.errorTask = "General"
		p.errorContext = append(p.errorContext, stripped)
		return
	}

	if p.warningRegex.MatchString(stripped) {
		p.warnings = append(p.warnings, stripped)
	}
}

func (p *Parser) commitError() {
	if len(p.errorContext) > 0 {
		msg := fmt.Sprintf("ERROR [Task %s]:\n", p.errorTask)
		msg += strings.Join(p.errorContext, "\n")
		p.errors = append(p.errors, msg)
	}
	p.inErrorContext = false
	p.errorContext = nil
	p.errorTask = ""
}

func (p *Parser) PrintSummary(exitCode int, logFile string) {
	// flush remaining
	if len(p.buffer) > 0 {
		p.processLine(string(p.buffer))
	}
	if p.inErrorContext {
		p.commitError()
	}

	if exitCode == 0 {
		if p.showWarnings && len(p.warnings) > 0 {
			for _, w := range p.warnings {
				fmt.Printf("WARNING: %s\n", w)
			}
			fmt.Println()
		}
		fmt.Printf("SUCCESS: Build completed. Full logs at: %s\n", logFile)
	} else {
		fmt.Println("\n--- TOON OUTPUT ---")
		if len(p.errors) > 0 {
			for _, e := range p.errors {
				fmt.Println(e)
				fmt.Println("...")
			}
		} else {
			fmt.Println("ERROR: Build failed, but no explicit error lines were caught.")
		}
		fmt.Printf("\nFAILURE: Build failed with exit code %d. Full logs at: %s\n", exitCode, logFile)
	}
}

func (p *Parser) HasErrors() bool {
	return len(p.errors) > 0 || p.inErrorContext
}
