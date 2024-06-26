# bump_version

This is a tool for bumping version numbers in Go files.

## Installation

For the moment, you'll need a working Go installation.

```
go install github.com/kevinburke/bump_version@latest
go install github.com/kevinburke/bump_version/current_version@latest
go install github.com/kevinburke/bump_version/next_version@latest
```

That will install the `bump_version` binary to your `$GOPATH`.

## Usage

```
bump_version <major|minor|patch> <filename>
```

This will:

1. Look for a `const` named `version`, `VERSION`, or `Version` in that file.
   Here's an example:

    ```go
    package main

    const VERSION = "0.2.1"
    ```

2. Apply the version bump - `bump_version major` will increment the major
version number, `bump_version minor` will increment the middle version number,
`bump_version patch` will increment the last version number. If your version is
"0.3" and you ask for `bump_version minor`, the new version will be "0.4".

3. Write the new file to disk, with the bumped version.

4. Add the file with `git add <filename>`.

5. Add a commit with the message "x.y.z" (`git commit -m "<new_version>"`)

6. Tag the new version.

If any of these steps fail, `bump_version` will abort.

#### current_version

This program will retrieve the current version from a Go source file.

```
# current_version main.go
0.6.0
```

#### next_version

This program will retrieve the current version from a Go source file, and then
print out the result of incrementing it (without actually making any changes).

```
# next_version patch main.go
0.6.1
```


## Notes

The VERSION in the Go file should be a string in one of these formats: "3",
"0.3", "0.3.4". Any prefixes like "v" or suffixes like "0.3.3-beta" will be
stripped or generate an error.

- `"v0.1"` - parse error, no prefixes allowed.
- `bump_version("0.1", "minor")` -> "0.2"
- `bump_version("0.1", "patch")` -> "0.1.1"
- `bump_version("0.1", "major")` -> "1.1"
- `bump_version("0.1-beta", "major")` -> "1.1"
- `bump_version("devel", "major")` -> parse error.

To add a prefix to the Git tag, add the `--tag-prefix` flag, for example
`--tag-prefix=v` will generate a Git tag that looks like "v1.2.3".

We use the VERSION in code exclusively - any existing git tags are ignored.

Alan Shreve would like to note that you probably shouldn't store version
numbers in code - instead, check in `const VERSION = "devel"`, then build your
project via:

```
go build -ldflags="-X main.VERSION=0.2"
```

Which you are welcome to do!
