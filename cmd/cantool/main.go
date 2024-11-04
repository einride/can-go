package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/scanner"

	"github.com/alecthomas/kingpin/v2"
	"github.com/fatih/color"
	"go.einride.tech/can/internal/generate"
	"go.einride.tech/can/pkg/dbc"
	"go.einride.tech/can/pkg/dbc/analysis"
	"go.einride.tech/can/pkg/dbc/analysis/passes/definitiontypeorder"
	"go.einride.tech/can/pkg/dbc/analysis/passes/intervals"
	"go.einride.tech/can/pkg/dbc/analysis/passes/lineendings"
	"go.einride.tech/can/pkg/dbc/analysis/passes/messagenames"
	"go.einride.tech/can/pkg/dbc/analysis/passes/multiplexedsignals"
	"go.einride.tech/can/pkg/dbc/analysis/passes/newsymbols"
	"go.einride.tech/can/pkg/dbc/analysis/passes/nodereferences"
	"go.einride.tech/can/pkg/dbc/analysis/passes/noreservedsignals"
	"go.einride.tech/can/pkg/dbc/analysis/passes/requireddefinitions"
	"go.einride.tech/can/pkg/dbc/analysis/passes/signalbounds"
	"go.einride.tech/can/pkg/dbc/analysis/passes/signalnames"
	"go.einride.tech/can/pkg/dbc/analysis/passes/singletondefinitions"
	"go.einride.tech/can/pkg/dbc/analysis/passes/siunits"
	"go.einride.tech/can/pkg/dbc/analysis/passes/uniquemessageids"
	"go.einride.tech/can/pkg/dbc/analysis/passes/uniquenodenames"
	"go.einride.tech/can/pkg/dbc/analysis/passes/uniquesignalnames"
	"go.einride.tech/can/pkg/dbc/analysis/passes/unitsuffixes"
	"go.einride.tech/can/pkg/dbc/analysis/passes/valuedescriptions"
	"go.einride.tech/can/pkg/dbc/analysis/passes/version"
)

func main() {
	app := kingpin.New("cantool", "CAN tool for Go programmers")
	generateCommand(app)
	lintCommand(app)
	kingpin.MustParse(app.Parse(os.Args[1:]))
}

func generateCommand(app *kingpin.Application) {
	command := app.Command("generate", "generate CAN messages")
	inputDir := command.
		Arg("input-dir", "input directory").
		Required().
		ExistingDir()
	outputDir := command.
		Arg("output-dir", "output directory").
		Required().
		String()
	allowedMessageIds := command.
		Arg("allowed-message-ids", "optional filter of message-ids to compile").
		Uint32List()
	command.Action(func(_ *kingpin.ParseContext) error {
		return filepath.Walk(*inputDir, func(p string, i os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if i.IsDir() || filepath.Ext(p) != ".dbc" {
				return nil
			}
			relPath, err := filepath.Rel(*inputDir, p)
			if err != nil {
				return err
			}
			outputFile := relPath + ".go"
			outputPath := filepath.Join(*outputDir, outputFile)
			return genGo(p, outputPath, *allowedMessageIds)
		})
	})
}

func lintCommand(app *kingpin.Application) {
	command := app.Command("lint", "lint DBC files")
	fileOrDir := command.
		Arg("file-or-dir", "DBC file or directory").
		Required().
		ExistingFileOrDir()
	command.Action(func(_ *kingpin.ParseContext) error {
		filesToLint, err := resolveFileOrDirectory(*fileOrDir)
		if err != nil {
			return err
		}
		var hasFailed bool
		for _, lintFile := range filesToLint {
			f, err := os.Open(lintFile)
			if err != nil {
				return err
			}
			source, err := io.ReadAll(f)
			if err != nil {
				return err
			}
			p := dbc.NewParser(f.Name(), source)
			if err := p.Parse(); err != nil {
				printError(source, err.Position(), err.Reason(), "parse")
				continue
			}
			for _, a := range analyzers() {
				pass := &analysis.Pass{
					Analyzer: a,
					File:     p.File(),
				}
				if err := a.Run(pass); err != nil {
					return err
				}
				hasFailed = hasFailed || len(pass.Diagnostics) > 0
				for _, d := range pass.Diagnostics {
					printError(source, d.Pos, d.Message, a.Name)
				}
			}
		}
		if hasFailed {
			return errors.New("one or more lint errors")
		}
		return nil
	})
}

func analyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		// TODO: Re-evaluate if we want boolprefix.Analyzer(), since it creates a lot of churn in vendor schemas
		definitiontypeorder.Analyzer(),
		intervals.Analyzer(),
		lineendings.Analyzer(),
		messagenames.Analyzer(),
		multiplexedsignals.Analyzer(),
		newsymbols.Analyzer(),
		nodereferences.Analyzer(),
		noreservedsignals.Analyzer(),
		requireddefinitions.Analyzer(),
		signalbounds.Analyzer(),
		signalnames.Analyzer(),
		singletondefinitions.Analyzer(),
		siunits.Analyzer(),
		uniquemessageids.Analyzer(),
		uniquenodenames.Analyzer(),
		uniquesignalnames.Analyzer(),
		unitsuffixes.Analyzer(),
		valuedescriptions.Analyzer(),
		version.Analyzer(),
	}
}

func genGo(inputFile, outputFile string, allowedMessageIds []uint32) error {
	if err := os.MkdirAll(filepath.Dir(outputFile), 0o755); err != nil {
		return err
	}
	input, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}
	result, err := generate.Compile(inputFile, input, generate.WithAllowedMessageIds(allowedMessageIds))
	if err != nil {
		return err
	}
	for _, warning := range result.Warnings {
		return warning
	}
	output, err := generate.Database(result.Database)
	if err != nil {
		return err
	}
	if err := os.WriteFile(outputFile, output, 0o600); err != nil {
		return err
	}
	fmt.Println("wrote:", outputFile)
	return nil
}

func resolveFileOrDirectory(fileOrDirectory string) ([]string, error) {
	fileInfo, err := os.Stat(fileOrDirectory)
	if err != nil {
		return nil, err
	}
	if !fileInfo.IsDir() {
		return []string{fileOrDirectory}, nil
	}
	var files []string
	if err := filepath.Walk(fileOrDirectory, func(path string, info os.FileInfo, _ error) error {
		if !info.IsDir() && filepath.Ext(path) == ".dbc" {
			files = append(files, path)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return files, nil
}

func printError(source []byte, pos scanner.Position, msg, name string) {
	fmt.Printf("\n%s: %s (%s)\n", pos, color.RedString("%s", msg), name)
	fmt.Printf("%s\n", getSourceLine(source, pos))
	fmt.Printf("%s\n", caretAtPosition(pos))
}

func getSourceLine(source []byte, pos scanner.Position) []byte {
	lineStart := pos.Offset
	for lineStart > 0 && source[lineStart-1] != '\n' {
		lineStart--
	}
	lineEnd := pos.Offset
	for lineEnd < len(source) && source[lineEnd] != '\n' {
		lineEnd++
	}
	return source[lineStart:lineEnd]
}

func caretAtPosition(pos scanner.Position) string {
	return strings.Repeat(" ", pos.Column-1) + color.YellowString("^")
}
