package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/taylormonacelli/goldbug"
	"github.com/taylormonacelli/ourbmw"
)

var (
	verbose   bool
	logFormat string
	rootDir   string
	oldStr    string
	newStr    string
)

func main() {
	flag.StringVar(&rootDir, "rootDir", "/path/to/your/directory", "Root directory to process")
	flag.StringVar(&oldStr, "oldStr", "bluefashion", "String to replace")
	flag.StringVar(&newStr, "newStr", "{{ cookiecutter.project_slug }}", "Replacement string")

	flag.BoolVar(&verbose, "verbose", false, "Enable verbose output")
	flag.BoolVar(&verbose, "v", false, "Enable verbose output (shorthand)")
	flag.StringVar(&logFormat, "log-format", "", "Log format (text or json)")

	flag.Parse()

	if verbose || logFormat != "" {
		if logFormat == "json" {
			goldbug.SetDefaultLoggerJson(slog.LevelDebug)
		} else {
			goldbug.SetDefaultLoggerText(slog.LevelDebug)
		}
	}

	code := ourbmw.Main(rootDir, oldStr, newStr)
	os.Exit(code)
}
