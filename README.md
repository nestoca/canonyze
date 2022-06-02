# yamlcanonizer

CLI for standardizing a set of yaml documents into a canonical form suitable for comparison against another set

# Installation

## Installing via brew on MacOS (recommended)

```
$ brew tap nestoca/yamlcanonizer
$ brew install yamlcanonizer
```

## Downloading binary

Download and install latest release for your platform from the GitHub releases page.

Make sure the binary is accessible via your `$PATH`.

## Building and installing from source

This approach requires that you replace both mentions of version number with your desired version in the following command:

$ go install -ldflags "-s -w" github.com/nestoca/yamlcanonizer/cmd/yamlcanonizer@v0.1.2
