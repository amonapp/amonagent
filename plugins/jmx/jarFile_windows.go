// +build windows

package jmx

import "os"

// TmpJarFile - XXX
var TmpJarFile = os.Getenv("TMP") + "\\amonagent\\mjb.jar"
