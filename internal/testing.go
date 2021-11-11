package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func GetTestConfigFile() string {
	_, path, _, _ := runtime.Caller(0)
	dir := filepath.Dir(path)
	dirs := strings.Split(dir, string(os.PathSeparator))
	dir = strings.Join(dirs[0: len(dirs)-1], string(os.PathSeparator))
	configFilePath := dir + "/example.toml"
	fmt.Printf("config file path: %s\n", configFilePath)
	return configFilePath
}
