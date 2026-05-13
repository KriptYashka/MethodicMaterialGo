package main

import (
	"net/http"
	_ "net/http/pprof"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})

	// pprof доступен на localhost:6060/debug/pprof/
	panic(http.ListenAndServe(":8080", nil))
}
