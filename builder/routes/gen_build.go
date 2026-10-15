package routes

import (
	"net/http"

	"radithor-builder/build"
	"radithor-builder/utils"
)

func generateBuild(w http.ResponseWriter, r *http.Request) {
	utils.NewLogChunk()

	build.GenerateRadithorBuild(w, r)

	w.WriteHeader(http.StatusOK)
}
