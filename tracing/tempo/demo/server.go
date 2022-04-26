package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func initProvider() func() {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithFromEnv(),
	)
	handleErr(err, "failed to create resource")

	traceExporter, err := otlptracehttp.New(ctx)
	handleErr(err, "failed to create trace exporter")

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	otel.SetTextMapPropagator(propagation.TraceContext{})

	return func() {
		handleErr(tracerProvider.Shutdown(ctx), "failed to shutdown TracerProvider")
	}
}

func main() {

	shutdown := initProvider()
	tracer := otel.Tracer("demo")
	defer shutdown()
	mux := mux.NewRouter()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(r.Header))
		ctx, span := tracer.Start(
			ctx,
			"serve hello")
		defer span.End()
		fmt.Println(span.SpanContext().TraceID())

		time.Sleep(time.Second)
	})

	http.ListenAndServe(":8080", mux)

	log.Printf("Done!")
}

func handleErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}
