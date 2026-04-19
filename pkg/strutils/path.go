package strutils

func APIPath(version string, path string) string {
	return "/api/" + version + path
}
