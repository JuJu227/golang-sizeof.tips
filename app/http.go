package app

import (
	"net/http"
	"strings"

	"github.com/gophergala/golang-sizeof.tips/internal/bindata"
)

func bindHttpHandlers() {
	fileServer := http.NewServeMux()
	fileServer.Handle("/", useCustom404(http.FileServer(http.Dir("pub/"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "."):
			fileServer.ServeHTTP(w, r)
			return
		case r.URL.Path != "/":
			write404(w)
			return
		}
		pageHandler(w, r)
	})
}

func write404(w http.ResponseWriter) {
	w.Write([]byte("gala not found!"))
}

type hijack404 struct {
	http.ResponseWriter
}

func (h *hijack404) WriteHeader(code int) {
	if code == 404 {
		write404(h.ResponseWriter)
		panic(h)
	}
	h.ResponseWriter.WriteHeader(code)
}

func useCustom404(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hijack := &hijack404{w}
		defer func() {
			if p := recover(); p != nil {
				if p == hijack {
					return
				}
				panic(p)
			}
		}()
		handler.ServeHTTP(hijack, r)
	})
}

func pageHandler(w http.ResponseWriter, _ *http.Request) {
	data, _ := bindata.Asset("templs/index.tmpl")
	w.Write(data)
}
