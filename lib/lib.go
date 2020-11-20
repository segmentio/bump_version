// Package bump_version makes it easy to parse and update version numbers in Go
// source files.
package bump_version

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"strconv"
	"strings"
)

const VERSION = "2.1"

type VersionType string

const Major = VersionType("major")
const Minor = VersionType("minor")
const Patch = VersionType("patch")

func ValidVersionType(vtype VersionType) bool {
	switch vtype {
	case Major, Minor, Patch:
		return true
	default:
		return false
	}
}

type Version struct {
	Major int64
	// May be "-1" to signify that this version field is unused.
	Minor int64
	Patch int64
}

func (v Version) String() string {
	if v.Patch >= 0 {
		return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	}
	if v.Minor >= 0 {
		return fmt.Sprintf("%d.%d", v.Major, v.Minor)
	}
	if v.Major >= 0 {
		return fmt.Sprintf("%d", v.Major)
	}
	return "%!s(INVALID_VERSION)"
}

// ParseVersion parses a version string of the forms "2", "2.3", or "0.10.11".
// Any information after the third number ("2.0.0-beta") is discarded. Very
// little effort is taken to validate the input.
//
// If a field is omitted from the string version (e.g. "0.2"), it's stored in
// the Version string as the integer -1.
func Parse(version string) (Version, error) {
	if len(version) == 0 {
		return Version{}, errors.New("bump_version: empty version string")
	}

	parts := strings.SplitN(version, ".", 3)
	if len(parts) == 1 {
		major, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return Version{}, err
		}
		return Version{
			Major: major,
			Minor: -1,
			Patch: -1,
		}, nil
	}
	if len(parts) == 2 {
		major, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return Version{}, err
		}
		minor, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return Version{}, err
		}
		return Version{
			Major: major,
			Minor: minor,
			Patch: -1,
		}, nil
	}
	if len(parts) == 3 {
		major, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return Version{}, err
		}
		minor, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return Version{}, err
		}
		patchParts := strings.SplitN(parts[2], "-", 2)
		patch, err := strconv.ParseInt(patchParts[0], 10, 64)
		if err != nil {
			return Version{}, err
		}
		return Version{
			Major: major,
			Minor: minor,
			Patch: patch,
		}, nil
	}
	return Version{}, fmt.Errorf("bump_version: invalid version string: %q", version)
}

// Bump increments the version number by the given vtype (major/minor/patch).
// Bump panics if vtype is not a known VersionType.
func Bump(version Version, vtype VersionType) Version {
	switch vtype {
	case Major:
		version.Major++
		if version.Minor != -1 {
			version.Minor = 0
		}
		if version.Patch != -1 {
			version.Patch = 0
		}
		return version
	case Minor:
		if version.Minor == -1 {
			version.Minor = 0
		}
		if version.Patch != -1 {
			version.Patch = 0
		}
		version.Minor++
		return version
	case Patch:
		if version.Patch == -1 {
			version.Patch = 0
		}
		version.Patch++
		return version
	default:
		panic(fmt.Sprintf("bump_version: invalid version type: %s", vtype))
	}
}

// changeVersion takes a basic literal representing a string version, and
// increments the version number per the given VersionType.
func changeVersion(vtype VersionType, value string) (Version, error) {
	versionNoQuotes := strings.Replace(value, "\"", "", -1)
	version, err := Parse(versionNoQuotes)
	if err != nil {
		return Version{}, err
	}
	return Bump(version, vtype), nil
}

func findBasicLit(file *ast.File) (*ast.BasicLit, error) {
	for _, decl := range file.Decls {
		switch gd := decl.(type) {
		case *ast.GenDecl:
			if gd.Tok != token.CONST {
				continue
			}
			spec, _ := gd.Specs[0].(*ast.ValueSpec)
			if strings.ToUpper(spec.Names[0].Name) == "VERSION" {
				value, ok := spec.Values[0].(*ast.BasicLit)
				if !ok || value.Kind != token.STRING {
					return nil, fmt.Errorf("bump_version: VERSION constant is not a string, was %#v", value.Value)
				}
				return value, nil
			}
		default:
			continue
		}
	}
	return nil, errors.New("bump_version: No version const found")
}

func writeFile(filename string, fset *token.FileSet, file *ast.File) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	cfg := printer.Config{Mode: printer.UseSpaces | printer.TabIndent, Tabwidth: 8}
	return cfg.Fprint(f, fset, file)
}

var errNoChanges = errors.New("bump_version: no changes made")

func changeInFile(filename string, f func(*ast.BasicLit) error) error {
	fset := token.NewFileSet()
	parsedFile, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	lit, err := findBasicLit(parsedFile)
	if err != nil {
		return fmt.Errorf("bump_version: no Version const found in %s", filename)
	}
	funcErr := f(lit)
	if funcErr != nil && funcErr != errNoChanges {
		return err
	}
	if funcErr == errNoChanges {
		return nil
	}
	writeErr := writeFile(filename, fset, parsedFile)
	return writeErr
}

// GetInFile returns the version found in the file.
func GetInFile(filename string) (Version, error) {
	var version Version
	err := changeInFile(filename, func(lit *ast.BasicLit) error {
		versionNoQuotes := strings.Replace(lit.Value, "\"", "", -1)
		var err error
		version, err = Parse(versionNoQuotes)
		if err != nil {
			return err
		}
		return errNoChanges
	})
	return version, err
}

// SetInFile sets the version in filename to newVersion.
func SetInFile(newVersion Version, filename string) error {
	return changeInFile(filename, func(lit *ast.BasicLit) error {
		lit.Value = fmt.Sprintf("%q", newVersion.String())
		return nil
	})
}

// BumpInFile finds a constant named VERSION, version, or Version in the file
// with the given filename, increments the version per the given VersionType,
// and writes the file back to disk. Returns the incremented Version object.
func BumpInFile(vtype VersionType, filename string) (Version, error) {
	var version Version
	err := changeInFile(filename, func(lit *ast.BasicLit) error {
		var err error
		version, err = changeVersion(vtype, lit.Value)
		if err != nil {
			return err
		}
		lit.Value = fmt.Sprintf("%q", version.String())
		return nil
	})
	return version, err
}

// Less reports whether i is a lower version number than j.
func Less(i, j Version) bool {
	if i.Major < j.Major {
		return true
	}
	if i.Major > j.Major {
		return false
	}
	if i.Minor < j.Minor {
		return true
	}
	if i.Minor > j.Minor {
		return false
	}
	return i.Patch < j.Patch
}
