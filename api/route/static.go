package route

import (
	"fmt"
	"net/http"
	"os"
)

type StaticRoute struct {
}

func (route *StaticRoute) Register(parent *http.ServeMux) {
	parent.HandleFunc("/", route.root)
	parent.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("/root/public"))))
}

func (route *StaticRoute) root(w http.ResponseWriter, req *http.Request) {
	data, err := os.ReadFile("/root/public/index.html")
	if err != nil {
		fmt.Fprintf(os.Stderr, "[StaticRoute::root] File read error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}
