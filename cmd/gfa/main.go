package main

import (
	"fmt"
	"os"

	"github.com/silverAndroid/gradle-for-agents/runner"
)

const version = "1.0.0"

func main() {
	args := os.Args[1:]

	showWarnings := false
	passThrough := false
	var gradleArgs []string

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--show-warnings" {
			showWarnings = true
		} else if arg == "--help" || arg == "-h" || arg == "-?" {
			printHelp()
			passThrough = true
			gradleArgs = append(gradleArgs, arg)
		} else if arg == "--version" || arg == "-v" {
			printVersion()
			passThrough = true
			gradleArgs = append(gradleArgs, arg)
		} else {
			gradleArgs = append(gradleArgs, arg)
		}
	}

	exitCode := runner.Run(gradleArgs, showWarnings, passThrough)
	os.Exit(exitCode)
}

func printHelp() {
	fmt.Println("gfa options:")
	fmt.Println("  --show-warnings                    Output warnings in TOON format on successful build.")
	fmt.Println("  --version                          Prints gfa and gradle version information and exits.")
	fmt.Println("  --help                             Shows this help message and gradle's help message.")
	fmt.Println("")
}

func printVersion() {
	fmt.Printf("gfa version %s\n\n", version)
}
