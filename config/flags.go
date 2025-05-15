package config

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
)

type stringArray []string

func (s *stringArray) String() string {
	return fmt.Sprintf("%v", *s)
}

func (s *stringArray) Set(value string) error {
	*s = append(*s, value)
	return nil
}

var m = map[string]func(){
	"header": handleHeader,
	"apply":  handleApply,
}

var (
	HeaderCmd *flag.FlagSet
	ApplyCmd  *flag.FlagSet
)

var (
	MirrorFlag stringArray
	FilterFlag stringArray
	RotateFlag stringArray
	CropFlag   stringArray
)

var (
	SourceFileName string
	OutputFileName string
)

var OrderedFlags []string

func InitFlags() {
	flag.Usage = func() {
		_, _ = fmt.Fprint(os.Stderr, helpText)
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	cmd, ok := m[os.Args[1]]
	if !ok {
		flag.Usage()
		os.Exit(1)
	}

	cmd()
}

func handleHeader() {
	HeaderCmd = flag.NewFlagSet("header", flag.ExitOnError)
	HeaderCmd.Usage = func() {
		fmt.Print(headerHelpText)
	}
	parseFlags(HeaderCmd)
	if !validateHeader() {
		HeaderCmd.Usage()
		os.Exit(1)
	}
	SourceFileName = HeaderCmd.Args()[0]
}

func handleApply() {
	ApplyCmd = flag.NewFlagSet("apply", flag.ExitOnError)
	ApplyCmd.Var(&MirrorFlag, "mirror", "mirrors the image")
	ApplyCmd.Var(&FilterFlag, "filter", "applies a filter to the image")
	ApplyCmd.Var(&RotateFlag, "rotate", "rotates the image")
	ApplyCmd.Var(&CropFlag, "crop", "crops the image")
	ApplyCmd.Usage = func() {
		fmt.Print(applyHelpText)
	}
	parseFlags(ApplyCmd)
	if !validateApply() {
		ApplyCmd.Usage()
		os.Exit(1)
	}
	SourceFileName = ApplyCmd.Args()[0]
	OutputFileName = ApplyCmd.Args()[1]
	parseOrderedFlags()
}

func parseFlags(cmd *flag.FlagSet) {
	if len(os.Args) < 3 {
		cmd.Usage()
		os.Exit(1)
	}

	err := cmd.Parse(os.Args[2:])
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func parseOrderedFlags() {
	for i := 1; i < len(os.Args[:len(os.Args)-2]); i++ {
		arg := os.Args[i]
		if strings.HasPrefix(arg, "--") || strings.HasPrefix(arg, "-") {
			name := strings.SplitN(arg, "=", 2)[0]
			name = strings.TrimPrefix(name, "--")
			name = strings.TrimPrefix(name, "-")
			OrderedFlags = append(OrderedFlags, name)
		}
	}
}

func validateHeader() bool {
	args := HeaderCmd.Args()

	if len(args) > 1 {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR: too many arguments")
		return false
	}

	if hasFlags(args) {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR: invalid flag")
		return false
	}

	if !strings.HasSuffix(args[0], ".bmp") {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR: invalid file format")
		return false
	}

	return true
}

func validateApply() bool {
	args := ApplyCmd.Args()

	if len(args) < 2 {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR: not enough arguments")
		return false
	}

	if len(args) > 2 {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR: too many arguments")
		return false
	}

	if hasFlags(args) {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR: invalid flags")
		return false
	}

	if !strings.HasSuffix(args[0], ".bmp") || !strings.HasSuffix(args[1], ".bmp") {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR: invalid file format")
		return false
	}

	return true
}

func hasFlags(arr []string) bool {
	return slices.ContainsFunc(arr, func(s string) bool {
		return strings.HasPrefix(s, "-")
	})
}
