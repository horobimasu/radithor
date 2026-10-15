package routes

import (
	"fmt"
	"net/http"

	"radithor-builder/store"
)

func getVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s %s", store.SERVER_VERSION, store.SITE_VERSION)
}
