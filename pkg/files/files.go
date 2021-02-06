package files

import (
	"fmt"
	"os"

	"github.com/hiromaily/go-google-home/pkg/config"
)

const installPath = "/usr/local/etc/google-home/"

// GetConfigPath returns toml file path
func GetConfigPath(tomlPath string) string {
	if tomlPath != "" && isExist(tomlPath) {
		return tomlPath
	}
	// gh.toml
	if installed := getInstallPath("toml"); isExist(installed) {
		return installed
	}

	envFile := config.GetEnvConfPath()
	if envFile != "" && isExist(envFile) {
		return envFile
	}
	return ""
}

func isExist(file string) bool {
	_, err := os.OpenFile(file, os.O_RDONLY, 0)
	if err != nil {
		if os.IsNotExist(err) {
			return false // file is not existing
		}
		return false // error occurred somehow, e.g. permission error
	}
	return true
}

func getInstallPath(ext string) string {
	return fmt.Sprintf("%s%s.%s", installPath, os.Args[0], ext)
}
