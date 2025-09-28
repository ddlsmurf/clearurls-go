// Command line tool to run ClearURLs rules or generate go code from them
package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/ddlsmurf/clearurls-go/clearurls"
)

func commandGenerate(source, destrinationFile string) error {
	providers, err := clearurls.GetProvidersFromSourceArgument(source)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Got %d providers\n", len(providers))
	if _, err := clearurls.Compile(providers); err != nil {
		return fmt.Errorf("Couldn't compile providers. %w", err)
	}
	goSource := clearurls.GenerateGoSourceCodeForProviders(providers)
	fmt.Fprintf(os.Stderr, "Writing %d bytes to %q\n", len(goSource), destrinationFile)
	if destrinationFile == "-" {
		fmt.Println(goSource)
		return nil
	}
	return os.WriteFile(destrinationFile, []byte(goSource), 0644)
}

func runCleanURLTests(providers []clearurls.RunnableProvider) int {
	failed := 0
	runCleanURLTest := func(url, expected string) {
		fmt.Fprintf(os.Stderr, "\n- Test        %q\n", url)
		cleaned, err := clearurls.ClearURL(providers, url, false)
		if err != nil {
			panic(err)
		}
		if expected != cleaned {
			fmt.Fprintf(os.Stderr, "   FAIL  Got: %q\n    Expected: %q\n", cleaned, expected)
			failed++
		} else {
			fmt.Fprintf(os.Stderr, "   pass  Got: %q\n", cleaned)
		}
	}
	runCleanURLTest("https://amazon.com?zoup=com&keywords=truc", "https://amazon.com?zoup=com")
	runCleanURLTest("https://amazon.com?zoup=com&keywords=truc#bidule=truc&keywords=ohno", "https://amazon.com?zoup=com#bidule=truc")
	runCleanURLTest("https://indeed.com?zoup=com&yclid=truc", "https://indeed.com?zoup=com")
	runCleanURLTest("https://indeed.com?yclid=truc", "https://indeed.com")
	runCleanURLTest("https://indeed.com/rc/clk?from=com&keywords=truc", "https://indeed.com/rc/clk?from=com&keywords=truc") // exception
	runCleanURLTest("https://google.com/plop?adurl=https%3A%2F%2Famazon.com%3Fzoup%3Dcom", "https://amazon.com?zoup=com")
	runCleanURLTest("https://google.com/plop?adurl=https%3A%2F%2Famazon.com%3Fzoup%3Dcom%26keywords%3Dtruc", "https://amazon.com?zoup=com")
	if failed > 0 {
		fmt.Fprintf(os.Stderr, "\n  => %d failed\n", failed)
	}
	return failed
}

func commandMiniTests(source string) error {
	providers, err := clearurls.GetProvidersFromSourceArgument(source)
	if err != nil {
		return err
	}
	wasCompiled := providers[0].IsCompiled()
	fmt.Fprintf(os.Stderr, "Got %d providers (compiled: %v)\n", len(providers), wasCompiled)
	fails := runCleanURLTests(providers)
	if !wasCompiled {
		compiledProviders, err := clearurls.Compile(providers)
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "\nRe-running tests on compiled providers (%v)\n\n", compiledProviders[0].IsCompiled())
		fails += runCleanURLTests(compiledProviders)
		if clearurls.GenerateGoSourceCodeForProviders(providers) != clearurls.GenerateGoSourceCodeForProviders(compiledProviders) {
			fmt.Fprintf(os.Stderr, "\nFail: Source compiled from json and from compiled providers didn't match")
			fails++
		} else {
			fmt.Fprintf(os.Stderr, "\nPass: Checked go source generation from json and compiled\n")
		}
	}
	providers, err = clearurls.HardcodedProviders()
	if err != nil {
		return fmt.Errorf("Error retreiving hardcoded providers: %w", err)
	}
	if providers == nil {
		fmt.Fprintf(os.Stderr, "\nSkip: Hardcoded providers not included in this version\n")
	} else {
		fmt.Fprintf(os.Stderr, "\nGot %d hardcoded providers (compiled: %v)\n", len(providers), wasCompiled)
		fails += runCleanURLTests(providers)
	}
	if fails > 0 {
		return fmt.Errorf("Got %d fail(s)\n", fails)
	}
	return nil
}

func readStdinByLine(gotLine func(line string) error) error {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			if err := gotLine(line); err != nil {
				return err
			}
		}
	}
	return scanner.Err()
}

func commandClean(source, urlToClean string, includeReferralMarketingParams bool) error {
	providers, err := clearurls.GetProvidersFromSourceArgument(source)
	if err != nil {
		return err
	}
	processLine := func(line string) error {
		cleaned, err := clearurls.ClearURL(providers, line, includeReferralMarketingParams)
		if err != nil {
			return err
		}
		fmt.Println(cleaned)
		return nil
	}
	if urlToClean == "-" {
		readStdinByLine(processLine)
	} else {
		processLine(urlToClean)
	}
	return nil
}

type commandType struct {
	name           string
	argsHelp, help string
	minArgs        int
	maxArgs        int
	run            func(args []string) error
}

type InvalidArgumentsError struct{ msg string }

func (e *InvalidArgumentsError) Error() string { return e.msg }
func NewInvalidArgumentsError(format string, a ...any) error {
	return &InvalidArgumentsError{msg: fmt.Sprintf(format, a...)}
}

var ArgsErrorJustPrintHelp = NewInvalidArgumentsError("Help command called")
var commands = []*commandType{
	{
		name:     "clean",
		argsHelp: "<source> <url or '-'>",
		help: "" +
			"Apply ClearURL process to an URL.\n" +
			"  - `source` can be the same as `generate`, or `hardcoded` if available\n" +
			"  - `url` can be the url to clean, or `-` to process each line on stdin\n",
		minArgs: 2,
		maxArgs: 2,
		run:     func(args []string) error { return commandClean(args[0], args[1], false) },
	},
	{
		name:     "cleanKeepingReferrals",
		argsHelp: "<source> <url or '-'>",
		help: "" +
			"Same as `clean` but does not remove parameters matching a `referralMarketing`\n",
		minArgs: 2,
		maxArgs: 2,
		run:     func(args []string) error { return commandClean(args[0], args[1], true) },
	},
	{
		name:     "generate",
		argsHelp: "<source> <destination_file>",
		help: "" +
			"Download CleanURL's JSON and generate hardoded data in GO source.\n" +
			"  - `source` can be '{github,gitlab}[:path_to_cache_file[:max_age_in_minutes]]'",
		minArgs: 2,
		maxArgs: 2,
		run:     func(args []string) error { return commandGenerate(args[0], args[1]) },
	},
	{
		name:     "mini_tests", // Some mini unit-ish tests to run on real data
		argsHelp: "[source]",
		help:     "hidden",
		minArgs:  1,
		maxArgs:  1,
		run:      func(args []string) error { return commandMiniTests(args[0]) },
	},
	{
		name:    "help",
		help:    "hidden",
		maxArgs: -1,
		run:     func(args []string) error { return ArgsErrorJustPrintHelp },
	},
}

func printUsage() {
	slices.SortFunc(commands[:], func(a, b *commandType) int {
		return strings.Compare(a.name, b.name)
	})
	fmt.Fprintf(os.Stderr, `
  Go package to work with https://docs.clearurls.xyz/1.27.3/ data for removing
  tracking and similar URL parameters

  Commands:

`)
	for _, command := range commands {
		if command.help != "hidden" {
			fmt.Fprintf(os.Stderr, "\t%s %s\n\t\t%s\n\n", command.name, command.argsHelp, strings.ReplaceAll(command.help, "\n", "\n\t\t"))
		}
	}
}

func runCommandFromArguments(args []string) error {
	if len(args) < 2 {
		return NewInvalidArgumentsError("Missing command name")
	}
	if slices.Contains(args, "--help") {
		return ArgsErrorJustPrintHelp
	}
	commandName := args[1]
	args = args[2:]
	for _, command := range commands {
		if commandName != command.name {
			continue
		}
		if command.minArgs >= 0 && len(args) < command.minArgs {
			return NewInvalidArgumentsError("Not enough arguments (got %d, need %d)", len(args), command.minArgs)
		}
		if command.maxArgs >= 0 && len(args) > command.maxArgs {
			return NewInvalidArgumentsError("Too many arguments (got %d, max %d)", len(args), command.maxArgs)
		}
		// fmt.Fprintf(os.Stderr, "Running command %q with arguments %q\n", commandName, args)
		return command.run(args)
	}
	return NewInvalidArgumentsError("Invalid command name %q", commandName)
}

func processFinalResult(err error) {
	if err == nil {
		return
	}
	if _, ok := err.(*InvalidArgumentsError); ok {
		if err != ArgsErrorJustPrintHelp {
			fmt.Fprintf(os.Stderr, "\nError: %s\n\nUsage:\n", err)
		}
		printUsage()
		if err == ArgsErrorJustPrintHelp {
			return
		}
	} else {
		fmt.Fprintf(os.Stderr, "\nError: %s\n", err)
	}
	os.Exit(1)
}

func main() {
	clearurls.Verbose = true
	processFinalResult(runCommandFromArguments(os.Args))
}
