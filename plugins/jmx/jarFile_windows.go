// +build windows

package jmx

import "os"

// JarFile - XXX
var JarFile = os.Getenv("TMP") + "\\amonagent\\mjb.jar"
