package api

func Path(version string, path string) string {
	return "/api/" + version + path
}
