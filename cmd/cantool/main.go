package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/scanner"

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
	"go.einride.tech/can/pkg/dbc/analysis/passes/uniquenodenames"
	"go.einride.tech/can/pkg/dbc/analysis/passes/uniquesignalnames"
	"go.einride.tech/can/pkg/dbc/analysis/passes/unitsuffixes"
	"go.einride.tech/can/pkg/dbc/analysis/passes/valuedescriptions"
	"go.einride.tech/can/pkg/dbc/analysis/passes/version"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	app := kingpin.New("cantool", "CAN tool for Go programmers")
	generateCommand(app)
	lintCommand(app)
	kingpin.MustParse(app.Parse(os.Args[1:]))
}

type SignalFilters map[string]map[string]*bool

func generateCommand(app *kingpin.Application) {
	command := app.Command("generate", "generate CAN messages")
	inputFileOrDir := command.
		Arg("input-file-or-dir", "input directory").
		Required().
		ExistingFileOrDir()
	outputDir := command.
		Arg("output-dir", "output directory").
		Required().
		String()
	filter := command.
		Flag("filter", "comma-separated list of filters (<message>:[<signal>])").
		String()
	filterFile := command.
		Flag("filter-file", "path to file containing messages to include, one per line").
		File()

	command.Action(func(c *kingpin.ParseContext) error {
		var signalFilters SignalFilters
		if *filter != "" {
			signalFilters = make(map[string]map[string]*bool)
			for _, e := range strings.Split(*filter, ",") {
				err := parseFilter(e, signalFilters)
				if err != nil {
					return err
				}
			}
		} else if *filterFile != nil {
			signalFilters = make(map[string]map[string]*bool)
			scan := bufio.NewScanner(*filterFile)
			for scan.Scan() {
				err := parseFilter(scan.Text(), signalFilters)
				if err != nil {
					return err
				}
			}
			if scan.Err() != nil {
				fmt.Println("Failed to parse message file", scan.Err())
			}
		}
		if (signalFilters != nil) {
			fmt.Println("Using filters (case-insensitive):")
			for msg, sigs := range signalFilters {
				fmt.Printf("\t%s:\n", msg)
				for sig := range sigs {
					fmt.Printf("\t\t%s\n", sig)
				}
			}
			defer func() {
				for msg, sigs := range signalFilters {
					for sig, found := range sigs {
						if !*found {
							fmt.Printf("Warning: no signal found matching filter '%s:%s'\n", msg, sig)
						}
					}
				}
			}()
		}

		return filepath.Walk(*inputFileOrDir, func(p string, i os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if i.IsDir() || filepath.Ext(p) != ".dbc" {
				return nil
			}
			var relPath string
			if *inputFileOrDir == p {
				relPath = i.Name()
			} else {
				relPath, err = filepath.Rel(*inputFileOrDir, p)
				if err != nil {
					return err
				}
			}
			outputFile := relPath + ".go"
			outputPath := filepath.Join(*outputDir, outputFile)
			return genGo(p, outputPath, signalFilters)
		})
	})
}
func boolPtr(b bool) *bool {
	return &b
}
func parseFilter(entry string, filters SignalFilters) error {
	pieces := strings.Split(entry, ":")
	if len(pieces) > 2 {
		return errors.New(fmt.Sprintf("Invalid filter entry: '%s', format is <message>[:<signal>]", entry))
	}
	message := strings.ToLower(pieces[0])
	signalSet, ok := filters[message]
	if !ok {
		signalSet = make(map[string]*bool)
	}
	if len(pieces) == 2 {
		signal := strings.ToLower(pieces[1])
		signalSet[signal] = boolPtr(false)
	}
	filters[message] = signalSet
	return nil
}

func lintCommand(app *kingpin.Application) {
	command := app.Command("lint", "lint DBC files")
	fileOrDir := command.
		Arg("file-or-dir", "DBC file or directory").
		Required().
		ExistingFileOrDir()
	command.Action(func(context *kingpin.ParseContext) error {
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
		uniquenodenames.Analyzer(),
		uniquesignalnames.Analyzer(),
		unitsuffixes.Analyzer(),
		valuedescriptions.Analyzer(),
		version.Analyzer(),
	}
}

func genGo(inputFile, outputFile string, filters SignalFilters) error {
	if err := os.MkdirAll(filepath.Dir(outputFile), 0o755); err != nil {
		return err
	}
	input, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}
	result, err := generate.Compile(inputFile, input)
	if err != nil {
		return err
	}
	for _, warning := range result.Warnings {
		return warning
	}

	if (filters != nil) {
		skips := make(map[string][]string)
		// Filter in-place for only the messages and signals matching the filter
		allMessages := result.Database.Messages
		result.Database.Messages = result.Database.Messages[:0]
		for _, msg := range allMessages {
			if signalSet, msgMatch := filters[strings.ToLower(msg.Name)]; msgMatch {
				allSignals := msg.Signals
				msg.Signals = msg.Signals[:0]
				for _, sig := range allSignals {
					if sigFound, sigMatch := signalSet[strings.ToLower(sig.Name)]; sigMatch {
						*sigFound = true
						msg.Signals = append(msg.Signals, sig)
					} else {
						skips[msg.Name] = append(skips[msg.Name], sig.Name)
					}
				}
				result.Database.Messages = append(result.Database.Messages, msg)
			} else {
				skips[msg.Name] = make([]string, 0)
			}
		}
		if len(skips) > 0 {
			fmt.Printf("The following messages/signals in %s were ignored due to filtering:\n", inputFile)
			sortedMsgs := make([]string, 0, len(skips))
			for msg := range skips {
				sortedMsgs = append(sortedMsgs, msg)
			}
			sort.Strings(sortedMsgs)
			for _, msg := range sortedMsgs {
				sigs := skips[msg]
				if len(sigs) > 0 {
					sort.Strings(sigs)
					for _, sig := range sigs {
						fmt.Printf("\t%s:%s\n", msg, sig)
					}
				} else {
					fmt.Printf("\t%s:*\n", msg)
				}
			}
		}
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
	if err := filepath.Walk(fileOrDirectory, func(path string, info os.FileInfo, err error) error {
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
