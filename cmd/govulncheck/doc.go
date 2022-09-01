// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Govulncheck reports known vulnerabilities that affect Go code. It uses static
analysis of source code or a binary's symbol table to narrow down reports to
only those that potentially affect the application. For more information about
the API behind govulncheck, see https://go.dev/security/vuln/vulncheck.

By default, govulncheck uses the Go vulnerability database at
https://vuln.go.dev. Set the GOVULNDB environment variable to specify a
different database.  The database must follow the specification at
https://go.dev/security/vuln/database.

Govulncheck requires Go version 1.18 or higher to run.

# Usage

To analyze source code, run govulncheck from the module directory, using the
same package path syntax that the go command uses:

	$ cd my-module
	$ govulncheck ./...

If no vulnerabilities are found, govulncheck will display a short message. If
there are vulnerabilities, each is displayed briefly, with a summary of a call
stack.

The call stack summary shows in brief how the package calls a vulnerable
function. For example, it might say

	main.go:[line]:[column]: mypackage.main calls golang.org/x/text/language.Parse

For a more detailed call path that resembles Go panic stack traces, use the -v
flag.

To control which files are processed, use the -tags flag to provide a
comma-separated list of build tags, and the -tests flag to indicate that test
files should be included.

To run govulncheck on a compiled binary, pass it the path to the binary file:

	$ govulncheck $HOME/go/bin/my-go-program

Govulncheck uses the binary's symbol information to find mentions of vulnerable
functions. Its output and exit codes are as described above, except that
without source it cannot produce call stacks.

# Other Modes

A few flags control govulncheck's output. Regardless of output, govulncheck
exits successfully if there are no vulnerabilities, and exits unsuccessfully if
there are.

The -v flag outputs more information about call stacks when run on source. It
has no effect when run on a binary.

The -json flag outputs a JSON object with vulnerability information. The output
corresponds to the type golang.org/x/vuln/vulncheck.Result.

# Weaknesses

Govulncheck is built on top of golang.org/x/vuln/vulncheck library and thus
shares its limitations described at
https://go.dev/security/vulncheck#limitations.
*/
package main
