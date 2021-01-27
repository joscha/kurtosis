/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package main

import (
	"flag"
	"fmt"
	"github.com/kurtosis-tech/kurtosis/commons/logrus_log_levels"
	"github.com/kurtosis-tech/kurtosis/commons/volume_naming_consts"
	"github.com/kurtosis-tech/kurtosis/initializer/initializer_container_constants"
	"github.com/palantir/stacktrace"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"regexp"
	"strings"
	"text/template"
)

const (
	// "Enum" of actions
	StoreTrue Action =  iota
	StoreValue

	flagPrefix = "--"

	failureExitCode = 1

	productionReleaseVersionPatternStr = `^[0-9]+\.[0-9]+$`
)

// Struct containing data specifically about the flag args to the wrapper script, bundled in such a way
//  as to be easy for the template to consume it
type WrapperFlagArgParsingData struct {
	Flag string
	Variable string
	DoStoreTrue bool
	DoStoreValue bool
}

// Package of data that will be used to fill the template
type WrapperTemplateData struct {
	DefaultValues map[string]string

	KurtosisCoreVersion string

	// Certain things in kurtosis.sh should only be enabled when we're generating a version that will be released to the
	//  world (and not a version that's being generated for running a local dev branch of Kurtosis Core)
	IsProductionRelease bool

	FlagArgParsingData []WrapperFlagArgParsingData

	NumPositionalArgs int

	// Mapping of index in Bash's argument array (e.g. "${1}") -> the positional arg that will receive the value
	PositionalArgAssignment map[int]string

	OneLinerHelpText string
	LinewiseHelpText []string

	// Bash date format to use when formatting the timestamp in the suite execution volume name
	VolumeTimestampDateFormat string
}

type Action int

// Definition of a wrapper arg, which will get parsed to generate the template data
type WrapperArg struct {
	// If empty, this is a positional arg
	Flag string

	// The Bash Variable that the value will be stored to
	Variable string

	// The default Variable that the Bash Variable will be assigned to, if not present
	DefaultVal string

	// The Variable HelpText
	HelpText string

	// The Action taken (only relevant for Flag variables; position variables will always store the value)
	Action Action
}

var wrapperArgs = []WrapperArg{
	{
		Flag:       "--custom-params",
		Variable:   "custom_params_json",
		DefaultVal: "{}",
		HelpText:   "JSON string containing arbitrary data that will be passed as-is to your testsuite, so it can modify its behaviour based on input",
		Action:     StoreValue,
	},
	{
		Flag:       "--client-id",
		Variable:   "client_id",
		DefaultVal: "",
		HelpText:   "An OAuth client ID which is needed for running Kurtosis in CI, and should be left empty when running Kurtosis on a local machine",
		Action:     StoreValue,
	},
	{
		Flag:       "--client-secret",
		Variable:   "client_secret",
		DefaultVal: "",
		HelpText:   "An OAuth client secret which is needed for running Kurtosis in CI, and should be left empty when running Kurtosis on a local machine",
		Action:     StoreValue,
	},
	{
		Flag:       "--help",
		Variable:   "show_help",
		DefaultVal: "false",
		HelpText:   "Display this message",
		Action:     StoreTrue,
	},
	{
		Flag:       "--kurtosis-log-level",
		Variable:   "kurtosis_log_level",
		DefaultVal: "info",
		HelpText: fmt.Sprintf(
			"The log level that all output generated by the Kurtosis framework itself should log at (%v)",
			strings.Join(logrus_log_levels.GetAcceptableLogLevelStrs(), "|"),
		),
		Action:     StoreValue,
	},
	{
		Flag:       "--list",
		Variable:   "do_list",
		DefaultVal: "false",
		HelpText:   "Rather than running the tests, lists the tests available to run",
		Action:     StoreTrue,
	},
	{
		Flag:       "--parallelism",
		Variable:   "parallelism",
		DefaultVal: "4",
		HelpText:   "The number of texts to execute in parallel",
		Action:     StoreValue,
	},
	{
		Flag:       "--tests",
		Variable:   "test_names",
		DefaultVal: "",
		HelpText:   "List of test names to run, separated by '" + initializer_container_constants.TestNameArgSeparator + "' (default or empty: run all tests)",
		Action:     StoreValue,
	},
	{
		// TODO Rename this to --suite-log-level to be less of a pain to type
		Flag:       "--test-suite-log-level",
		Variable:   "test_suite_log_level",
		DefaultVal: "info",
		HelpText: fmt.Sprintf("A string that will be passed as-is to the test suite container to indicate what log level the test suite container should output at; this string should be meaningful to the test suite container because Kurtosis won't know what logging framework the testsuite uses"),
		Action:     StoreValue,
	},
	{
		Variable:   "test_suite_image",
		HelpText:   "The Docker image containing the testsuite to execute",
	},
}

// Fills the Bash wrapper script template with the appropriate variables
func main()  {
	productionReleaseVersionPattern := regexp.MustCompile(productionReleaseVersionPatternStr)

	kurtosisCoreVersion, templateFilepath, outputFilepath, err := parseAndValidateFlags()
	if err != nil {
		logrus.Errorf("An error occurred parsing & validating the flags: %v", err)
		os.Exit(failureExitCode)
	}

	if err := validateWrapperArgs(wrapperArgs); err != nil {
		logrus.Errorf("An error occurred validating the wrapper args: %v", err)
		os.Exit(failureExitCode)
	}

	// For some reason, the template name has to match teh basename of the file:
	//  https://stackoverflow.com/questions/49043292/error-template-is-an-incomplete-or-empty-template
	templateFilename := path.Base(templateFilepath)
	tmpl, err := template.New(templateFilename).ParseFiles(templateFilepath)
	if err != nil {
		logrus.Errorf("An error occurred parsing the Bash template: %v", err)
		os.Exit(failureExitCode)
	}

	data, err := generateTemplateData(wrapperArgs, kurtosisCoreVersion, productionReleaseVersionPattern)
	if err != nil {
		logrus.Errorf("An error occurred generating the template data: %v", err)
		os.Exit(failureExitCode)
	}

	fp, err := os.Create(outputFilepath)
	if err != nil {
		logrus.Errorf("An error occurred opening the output file for writing: %v", err)
		os.Exit(failureExitCode)
	}
	defer fp.Close()

	if err := tmpl.Execute(fp, data); err != nil {
		logrus.Errorf("An error occurred filling the template: %v", err)
		os.Exit(failureExitCode)
	}
}

// ===================================================================================================
//                                        Private helper functions
// ===================================================================================================
// TODO Replace this entire thing by:
//  1. Moving the docker_flag_parser out of initializer and into commmons (it can now be generalized)
//  2. Using that now-generalized class here
func parseAndValidateFlags() (kurtosisCoreVersion, templateFilepath, outputFilepath string, err error) {
	kurtosisCoreVersionArg := flag.String(
		"kurtosis-core-version",
		"",
		"Version of Kurtosis core to generate the wrapper script with",
	)
	templateFilepathArg := flag.String(
		"template",
		"",
		"Filepath containing Bash template file that will get rendered into the Kurtosis wrapper script",
	)
	outputFilepathArg := flag.String(
		"output",
		"",
		"Output filepath to write the rendered template to",
	)
	flag.Parse()

	if *kurtosisCoreVersionArg == "" {
		return "", "", "", stacktrace.NewError("Kurtosis Core version arg is required")
	}
	if *templateFilepathArg == "" {
		return "", "", "", stacktrace.NewError("Template filepath arg is required")
	}
	if *outputFilepathArg == "" {
		return "", "", "", stacktrace.NewError("Output filepath arg is required")
	}

	return *kurtosisCoreVersionArg, *templateFilepathArg, *outputFilepathArg, nil
}

/*
Gets the text that an arg should have on the one-liner representation (e.g. "--parallelism parallelism")
 */
func getOneLinerText(arg WrapperArg) (string, error) {
	isFlagArg := arg.Flag != ""

	var onelinerText string
	if isFlagArg {
		if (arg.Action == StoreValue) {
			onelinerText = fmt.Sprintf("%v %v", arg.Flag, arg.Variable)
		} else if (arg.Action == StoreTrue) {
			onelinerText = arg.Flag
		} else {
			return "", stacktrace.NewError("Unrecognized arg Action '%v'; this is a code bug", arg.Action)
		}
	} else {
		onelinerText = arg.Variable
	}
	return onelinerText, nil
}

func padStringToLength(str string, desiredLength int) string {
	numSpacesToAdd := desiredLength - len(str)
	return str + strings.Repeat(" ", numSpacesToAdd)
}

func validateWrapperArgs(args []WrapperArg) error {
	// Validate arguments
	seenFlags := map[string]bool{}
	for _, arg := range args {
		isFlagArg := arg.Flag != ""
		if isFlagArg {
			if _, found := seenFlags[arg.Flag]; found {
				return stacktrace.NewError(
					"Duplicate flag '%v'",
					arg.Flag)
			}
			if !strings.HasPrefix(arg.Flag, flagPrefix) {
				return stacktrace.NewError(
					"Flag '%v' must start with flag prefix '%v'",
					arg.Flag,
					flagPrefix)
			}
		} else {
			if arg.DefaultVal != "" {
				return stacktrace.NewError(
					"Positional argument '%v' cannot have default value",
					arg.Variable)
			}
		}
	}
	return nil
}

// Get the longest one-liner text (which will become the title of the argument in linewise display)
func getLongestOneLinerLength(args []WrapperArg) (int, error) {
	longestOneLinerText := 0
	for _, arg := range args {
		oneLinerText, err := getOneLinerText(arg)
		if err != nil {
			return 0, stacktrace.Propagate(
				err,
				"An error occurred getting oneliner text for arg with variable '%v' while calculating max pad length",
				arg.Variable)
		}
		if len(oneLinerText) > longestOneLinerText {
			longestOneLinerText = len(oneLinerText)
		}
	}
	return longestOneLinerText, nil
}

func generateTemplateData(args []WrapperArg, kurtosisCoreVersion string, productionReleaseVersionPattern *regexp.Regexp) (*WrapperTemplateData, error) {
	longestOneLinerLength, err := getLongestOneLinerLength(args)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred getting the length of the longest one-liner text")
	}

	allFlagArgParsingData := []WrapperFlagArgParsingData{}
	defaultValues := map[string]string{}
	positionalArgAssignment := map[int]string{}
	positionalArgAssignmentIndex := 0

	flagArgsOnelinerFragments := []string{}
	positionalArgsOnelinerFragments := []string{}
	flagArgsLinewiseHelptext := []string{}
	positionalArgsLinewiseHelptext := []string{}
	for _, arg := range args {
		defaultValues[arg.Variable] = arg.DefaultVal

		oneLinerText, err := getOneLinerText(arg)
		if err != nil {
			return nil, stacktrace.Propagate(err, "An error occurred getting oneliner text for arg with variable '%v' while building template data", arg.Variable)
		}
		paddedOneLinerText := padStringToLength(oneLinerText, longestOneLinerLength + 3)

		isFlagArg := arg.Flag != ""
		if isFlagArg {
			flagArgParsingData := WrapperFlagArgParsingData{
				Flag:         arg.Flag,
				Variable:     arg.Variable,
				DoStoreTrue:  arg.Action == StoreTrue,
				DoStoreValue: arg.Action == StoreValue,
			}
			allFlagArgParsingData = append(allFlagArgParsingData, flagArgParsingData)

			flagArgsOnelinerFragments = append(
				flagArgsOnelinerFragments,
				fmt.Sprintf("[%v]", oneLinerText),
			)

			var linewiseText string
			if arg.DefaultVal != "" && arg.Action != StoreTrue {
				linewiseText = fmt.Sprintf(
					"%v%v (default: %v)",
					paddedOneLinerText,
					arg.HelpText,
					arg.DefaultVal,
				)
			} else {
				linewiseText = fmt.Sprintf(
					"%v%v",
					paddedOneLinerText,
					arg.HelpText,
				)
			}
			flagArgsLinewiseHelptext = append(
				flagArgsLinewiseHelptext,
				fmt.Sprintf("   %v", linewiseText),
			)
		} else {
			// Bash's argument list starts at 1, so we add 1
			positionalArgAssignment[positionalArgAssignmentIndex + 1] = arg.Variable
			positionalArgAssignmentIndex++

			positionalArgsOnelinerFragments = append(
				positionalArgsOnelinerFragments,
				oneLinerText,
			)
			linewiseText := fmt.Sprintf(
				"%v%v",
				paddedOneLinerText,
				arg.HelpText,
			)
			positionalArgsLinewiseHelptext = append(
				positionalArgsLinewiseHelptext,
				fmt.Sprintf("   %v", linewiseText),
			)
		}
	}

	flagArgsOneliner := strings.Join(flagArgsOnelinerFragments, " ")
	positionalArgsOneliner := strings.Join(positionalArgsOnelinerFragments, " ")
	combinedOneliner := fmt.Sprintf("%v %v", flagArgsOneliner, positionalArgsOneliner)

	combinedLinewiseHelptext := append(flagArgsLinewiseHelptext, positionalArgsLinewiseHelptext...)

	return &WrapperTemplateData{
		DefaultValues:      defaultValues,
		FlagArgParsingData: allFlagArgParsingData,
		KurtosisCoreVersion: kurtosisCoreVersion,
		IsProductionRelease: productionReleaseVersionPattern.Match([]byte(kurtosisCoreVersion)),
		NumPositionalArgs: len(positionalArgAssignment),
		PositionalArgAssignment: positionalArgAssignment,
		OneLinerHelpText:   combinedOneliner,
		LinewiseHelpText:   combinedLinewiseHelptext,
		VolumeTimestampDateFormat: volume_naming_consts.BashTimestampFormat,
	}, nil
}



