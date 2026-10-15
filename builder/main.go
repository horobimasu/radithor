package main

import (
	"fmt"
	"net/http"

	"radithor-builder/routes"
	"radithor-builder/store"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./builder.html")
	})

	store.InitState()

	routes.RegisterRoutes()

	fmt.Println("listening on http://localhost:" + store.PORT)
	http.ListenAndServe(":"+store.PORT, nil)
}
