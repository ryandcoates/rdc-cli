package sys

import (
	"os"
	"runtime"
	"strings"
)

type OSType string

const (
	Windows OSType = "windows"
	MacOS   OSType = "macos"
	Linux   OSType = "linux"
	WSL     OSType = "wsl"
	Unknown OSType = "unknown"
)

func GetOSType() OSType {
	switch runtime.GOOS {
	case "windows":
		return Windows
	case "darwin":
		return MacOS
	case "linux":
		if isWSL() {
			return WSL
		}
		return Linux
	default:
		return Unknown
	}
}

func isWSL() bool {
	data, err := os.ReadFile("/proc/version")
	if err != nil {
		return false
	}
	content := strings.ToLower(string(data))
	return strings.Contains(content, "microsoft") || strings.Contains(content, "wsl")
}

func GetArch() string {
	return runtime.GOARCH
}
