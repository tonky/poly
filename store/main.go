package main

import (
    // "context"
    "fmt"
    "log"
    "net/http"
    "github.com/go-chi/chi"
    "github.com/go-chi/chi/middleware"
	"github.com/uber/jaeger-client-go"
    ot "github.com/opentracing/opentracing-go"
    ext "github.com/opentracing/opentracing-go/ext"
    otlog "github.com/opentracing/opentracing-go/log"
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
    closer, err := cfg.InitGlobalTracer("store")

	if err != nil {
		log.Printf("Could not initialize jaeger tracer: %s", err.Error())
		return
	}

	defer closer.Close()

    log.Printf("Global tracer registered: %b?", ot.IsGlobalTracerRegistered())

    r := chi.NewRouter()
    r.Use(middleware.RequestID)
    r.Use(otm.Middleware())

    r.Get("/store/", func(w http.ResponseWriter, r *http.Request) {

        fmt.Printf("Headers: %v\n", r.Header)

        var sp ot.Span

        wireContext, err := ot.GlobalTracer().Extract(
            ot.HTTPHeaders,
            ot.HTTPHeadersCarrier(r.Header))

        if err != nil {
            fmt.Println("Can't extract context: %e", err)
        }

        // Create the span referring to the RPC client if available.
        // If wireContext == nil, a root span will be created.
        fmt.Println("wireContext: ", wireContext)

        sp = ot.StartSpan(
            "store",
            ext.RPCServerOption(wireContext))

        defer sp.Finish()

        // ctx := ot.ContextWithSpan(context.Background(), sp)

        sp.LogFields(
            otlog.String("path", r.URL.Path),
        )

        w.Write([]byte("Hello store!\n"))
    })

    fmt.Println("Starting http server...")

    log.Fatal(http.ListenAndServe(":8080", r))
}
