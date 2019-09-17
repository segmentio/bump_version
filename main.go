// The bump_version binary makes it easy to increment version constants in a Go
// source file.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	bump_version "github.com/kevinburke/bump_version/lib"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: bump_version [--version=<version>] [<major|minor|patch>] <filename>\n")
	flag.PrintDefaults()
}

// runCommand execs the given command and exits if it fails.
func runCommand(binary string, args ...string) {
	out, err := exec.Command(binary, args...).CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when running command: %s.\nOutput was:\n%s", err.Error(), string(out))
		os.Exit(2)
	}
}

func _main(flags *flag.FlagSet, cmdArgs []string) int {
	var vsn = flags.String("version", "", "Set this version in the file (don't increment whatever version is present)")
	if err := flags.Parse(cmdArgs); err != nil {
		flags.Usage()
		return 2
	}
	args := flags.Args()
	var filename string
	var version bump_version.Version
	if *vsn != "" {
		// no "minor"
		if len(args) != 1 {
			flags.Usage()
			return 2
		}
		var err error
		version, err = bump_version.Parse(*vsn)
		if err != nil {
			os.Stderr.WriteString(err.Error())
			return 2
		}
		filename = args[0]
		setErr := bump_version.SetInFile(version, filename)
		if setErr != nil {
			os.Stderr.WriteString(setErr.Error() + "\n")
			return 2
		}
	} else {
		if len(args) != 2 {
			flags.Usage()
			return 2
		}
		versionTypeStr := args[0]
		filename = args[1]

		var err error
		version, err = bump_version.BumpInFile(bump_version.VersionType(versionTypeStr), filename)
		if err != nil {
			os.Stderr.WriteString(err.Error() + "\n")
			return 2
		}
	}
	runCommand("git", "add", filename)
	runCommand("git", "commit", "-m", version.String())
	runCommand("git", "tag", version.String(), "--annotate", "--message", version.String())
	os.Stdout.WriteString(version.String() + "\n")
	return 0
}

func main() {
	flag.Usage = usage
	os.Exit(_main(flag.CommandLine, os.Args[1:]))
}
