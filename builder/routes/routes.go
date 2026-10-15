package routes

import "net/http"

func RegisterRoutes() {
	http.HandleFunc("/api/get-version", getVersion)
	http.HandleFunc("/api/get-key", getServerKey)
	http.HandleFunc("/api/gen-build", generateBuild)
}
