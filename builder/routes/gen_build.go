package routes

import (
	"net/http"
	"time"

	"radithor-builder/build"
	"radithor-builder/store"
	"radithor-builder/utils"
)

func generateBuild(w http.ResponseWriter, r *http.Request) {
	utils.NewLogChunk()

	store.BuildId = time.Now().UnixNano()

	build.GenerateRadithorBuild(w, r)

	w.WriteHeader(http.StatusOK)
}
