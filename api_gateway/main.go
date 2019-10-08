package main

import (
    "fmt"
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
    "github.com/go-chi/chi"
    "github.com/go-chi/chi/middleware"
)

func main() {
    r := chi.NewRouter()
    r.Use(middleware.RequestID)

    r.Get("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello world!\n"))
    })

    r.HandleFunc("/store*", StoreProxyHandler)

    fmt.Println("Starting http server...")

    log.Fatal(http.ListenAndServe(":8080", r))
}

func StoreProxyHandler(w http.ResponseWriter, r *http.Request) {
	url, _ := url.Parse("http://store/")

	proxy := httputil.ReverseProxy{Director: func(r *http.Request) {
		r.URL.Scheme = url.Scheme
		r.URL.Host = url.Host
		r.URL.Path = url.Path + r.URL.Path
		r.Host = url.Host
	}}

    fmt.Printf("Proxying to store: %s\n", r.URL.Path)

	// req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))

	proxy.ServeHTTP(w, r)
}
