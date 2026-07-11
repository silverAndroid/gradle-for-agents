package parser

import (
	"strings"
	"testing"
)

func TestParserSuccessNoWarnings(t *testing.T) {
	parser := NewParser(false)
	_, _ = parser.Write([]byte("BUILD SUCCESSFUL in 2s\n"))
	
	if len(parser.errors) > 0 || len(parser.warnings) > 0 {
		t.Errorf("Expected 0 errors and warnings, got %d errors, %d warnings", len(parser.errors), len(parser.warnings))
	}
}

func TestParserWarnings(t *testing.T) {
	parser := NewParser(true)
	_, _ = parser.Write([]byte("w: some warning\nwarning: another one\nBUILD SUCCESSFUL\n"))
	
	if len(parser.warnings) != 2 {
		t.Errorf("Expected 2 warnings, got %d", len(parser.warnings))
	}
}

func TestParserTaskFailure(t *testing.T) {
	parser := NewParser(true)
	input := `> Task :app:compileDebugJavaWithJavac FAILED
/path/to/file.java:12: error: cannot find symbol
        Foo foo = new Foo();
        ^
  symbol:   class Foo
  location: class MainActivity
5 errors
`
	_, _ = parser.Write([]byte(input))
	// force commit
	if parser.inErrorContext {
		parser.commitError()
	}
	
	if len(parser.errors) != 1 {
		t.Fatalf("Expected 1 error, got %d", len(parser.errors))
	}
	
	if !strings.Contains(parser.errors[0], ":app:compileDebugJavaWithJavac") {
		t.Errorf("Expected error to contain task name, got %s", parser.errors[0])
	}
	if !strings.Contains(parser.errors[0], "cannot find symbol") {
		t.Errorf("Expected error to contain error text, got %s", parser.errors[0])
	}
}

func TestParserGeneralFailure(t *testing.T) {
	parser := NewParser(true)
	input := `FAILURE: Build failed with an exception.

* What went wrong:
A problem occurred evaluating root project 'dummy'.
> Could not find method compile() for arguments...
`
	_, _ = parser.Write([]byte(input))
	// force commit
	if parser.inErrorContext {
		parser.commitError()
	}
	
	if len(parser.errors) != 1 {
		t.Fatalf("Expected 1 error, got %d", len(parser.errors))
	}
	if !strings.Contains(parser.errors[0], "General") {
		t.Errorf("Expected error to contain General task, got %s", parser.errors[0])
	}
	if !strings.Contains(parser.errors[0], "What went wrong:") {
		t.Errorf("Expected error to contain error text, got %s", parser.errors[0])
	}
}
