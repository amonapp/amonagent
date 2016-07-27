package pluginhelper

import (
	"io"
	"log"
	"os"
	"path"
	"strings"
)

// WritePluginConfig - XXX
func WritePluginConfig(name string, value string) error {

	var PluginGlobalConfigPath = path.Join("/tmp", "plugins-enabled")

	PluginDirCleanupErr := os.RemoveAll(PluginGlobalConfigPath)

	if PluginDirCleanupErr != nil {
		log.Fatalf("testing/pluginhelper.go line: 19 - Can't cleanup plugin dir: %v", PluginDirCleanupErr)
	}

	PluginDirErr := os.MkdirAll(PluginGlobalConfigPath, os.ModePerm)

	if PluginDirErr != nil {
		log.Fatalf("testing/pluginhelper.go line: 25 -  Can't create config directory: %v", PluginDirErr)
	}

	var pluginConfig = path.Join(PluginGlobalConfigPath, strings.Join([]string{name, "conf"}, "."))

	configFile, err := os.OpenFile(pluginConfig, os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		log.Fatalf("testing/pluginhelper.go line: 33 -  Can't create config file: %v", err)
	}

	_, writeErr := io.WriteString(configFile, value)

	if writeErr != nil {
		log.Fatalf("testing/pluginhelper.go line: 39 -  Can't write to file: %v", writeErr)
	}

	configFile.Close()

	// Use only to check the contents of the file

	// f, err := os.Open(pluginConfig)

	// if err != nil {
	// 	fmt.Println(err)
	// }

	// scanner := bufio.NewScanner(f)
	// for scanner.Scan() {
	// 	fmt.Println(scanner.Text())
	// }
	// f.Close()

	return nil
}
