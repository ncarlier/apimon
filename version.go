package main

import (
	"fmt"
)

// Version of the app
var Version = "snapshot"

func printVersion() {
	fmt.Printf(`apimon (%s)
Copyright (C) 2018 Nicolas Carlier. All rights reserved.

This work is licensed under the terms of the MIT license.  
For a copy, see <https://opensource.org/licenses/MIT>.
`, Version)
}
