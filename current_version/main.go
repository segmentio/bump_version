// The current_version binary prints out the value of a version number stored in
// a Go source file.
package main

import (
	"flag"
	"fmt"
	"os"

	bump_version "github.com/segmentio/bump_version/lib"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: current_version <filename>\n")
	flag.PrintDefaults()
}

func _main(flags *flag.FlagSet, cmdArgs []string) int {
	if err := flags.Parse(cmdArgs); err != nil {
		flags.Usage()
		return 2
	}
	args := flags.Args()
	if len(args) != 1 {
		flags.Usage()
		return 2
	}
	filename := args[0]
	version, err := bump_version.GetInFile(filename)
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		return 2
	}
	os.Stdout.WriteString(version.String() + "\n")
	return 0
}

func main() {
	flag.Usage = usage
	os.Exit(_main(flag.CommandLine, os.Args[1:]))
}
