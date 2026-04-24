package system_setting

import (
	"os"
	"strings"
)

func GetFrontendBaseURL() string {
	return strings.TrimRight(strings.TrimSpace(os.Getenv("FRONTEND_BASE_URL")), "/")
}

func GetPublicWebBaseURL() string {
	if frontendBaseURL := GetFrontendBaseURL(); frontendBaseURL != "" {
		return frontendBaseURL
	}
	return strings.TrimRight(strings.TrimSpace(ServerAddress), "/")
}

func BuildPublicWebURL(path string) string {
	baseURL := GetPublicWebBaseURL()
	if baseURL == "" {
		return path
	}
	if path == "" {
		return baseURL
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return baseURL + path
}
