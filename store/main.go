package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/go-chi/chi"
    "github.com/go-chi/chi/middleware"
)

func main() {
    r := chi.NewRouter()
    r.Use(middleware.RequestID)

    r.Get("/store/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello store!\n"))
    })

    fmt.Println("Starting http server...")

    log.Fatal(http.ListenAndServe(":8080", r))
}
