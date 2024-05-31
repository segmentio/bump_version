// The next_version binary prints out the value of the next version, using the
// current version stored in a Go source file as a base.
package main

import (
	"flag"
	"fmt"
	"os"

	bump_version "github.com/segmentio/bump_version/lib"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: next_version [<major|minor|patch>] <filename>\n")
	flag.PrintDefaults()
}

func _main(flags *flag.FlagSet, cmdArgs []string) int {
	if err := flags.Parse(cmdArgs); err != nil {
		flags.Usage()
		return 2
	}
	args := flags.Args()
	if len(args) != 2 {
		flags.Usage()
		return 2
	}
	filename := args[1]
	version, err := bump_version.GetInFile(filename)
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		return 2
	}
	versionType := bump_version.VersionType(args[0])
	if !bump_version.ValidVersionType(versionType) {
		fmt.Fprintf(os.Stderr, "invalid version type (want major/minor/patch): %q\n", versionType)
		return 2
	}
	newVersion := bump_version.Bump(version, versionType)
	os.Stdout.WriteString(newVersion.String() + "\n")
	return 0
}

func main() {
	flag.Usage = usage
	os.Exit(_main(flag.CommandLine, os.Args[1:]))
}
