package main

import (
    "fmt"
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
    "github.com/go-chi/chi"
    "github.com/go-chi/chi/middleware"
	"github.com/uber/jaeger-client-go"
    ot "github.com/opentracing/opentracing-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
    otm "github.com/improbable-eng/go-httpwares/tracing/opentracing" 
)

func main() {
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
            LocalAgentHostPort: "jaeger-agent:6831",
		},
	}

	// Initialize tracer with a logger and a metrics factory
    closer, err := cfg.InitGlobalTracer("api-gateway")

	if err != nil {
		log.Printf("Could not initialize jaeger tracer: %s", err.Error())
		return
	}

	defer closer.Close()

    log.Printf("Global tracer registered: %b?", ot.IsGlobalTracerRegistered())

    r := chi.NewRouter()
    r.Use(middleware.RequestID)
    r.Use(otm.Middleware())

    r.Get("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello gw!\n"))
    })

    r.HandleFunc("/store*", StoreProxyHandler)

    fmt.Println("Starting http server...")

    log.Fatal(http.ListenAndServe(":8080", r))
}

func StoreProxyHandler(w http.ResponseWriter, r *http.Request) {
    sp := ot.StartSpan("main")
    defer sp.Finish()

    ot.GlobalTracer().Inject(
        sp.Context(),
        ot.HTTPHeaders,
        ot.HTTPHeadersCarrier(r.Header))

	url, _ := url.Parse("http://store-service")

	proxy := httputil.ReverseProxy{Director: func(r *http.Request) {
		r.URL.Scheme = url.Scheme
		r.URL.Host = url.Host
		r.URL.Path = url.Path + r.URL.Path
		r.Host = url.Host
	}}

    fmt.Printf("Proxying to store: %s\n", r.URL.Path)

    fmt.Printf("Headers: %v\n", r.Header)

	proxy.ServeHTTP(w, r)
}
