package route

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

type StaticRoute struct {
}

func (route *StaticRoute) Register(parent *http.ServeMux) {
	parent.HandleFunc("/", route.root)
	parent.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("public"))))
}

func (route *StaticRoute) root(w http.ResponseWriter, req *http.Request) {
	path := filepath.Join("public/index.html")

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[StaticRoute::root] File read error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}
