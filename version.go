package main

import (
	"fmt"
)

// Version of the app
var Version = ""

// GitCommit hash
var GitCommit = "HEAD"

func printVersion() {
	version := Version
	if version == "" {
		version = GitCommit
	}
	fmt.Printf(`apimon (%s)
Copyright (C) 2018 Nicolas Carlier. All rights reserved.

This work is licensed under the terms of the MIT license.  
For a copy, see <https://opensource.org/licenses/MIT>.
`, Version)
}
