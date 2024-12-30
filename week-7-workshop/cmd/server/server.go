package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.ozon.dev/14/week-7-workshop/internal/pkg/cache"
	"gitlab.ozon.dev/14/week-7-workshop/internal/pkg/metrics"
	loggerPkg "gitlab.ozon.dev/14/week-7-workshop/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	trace2 "go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"io"
	"math/rand"
	"net/http"
	"time"

	_ "net/http/pprof"
)

func main() {
	ctx := context.Background()

	config := zap.NewProductionConfig()
	config.ErrorOutputPaths = []string{"stdout"}
	config.Level.SetLevel(zap.InfoLevel)

	logger := loggerPkg.NewLogger(config)
	defer logger.Sync()

	const (
		cacheSize = 100
		ttl       = 10 * time.Minute
	)

	exp, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpointURL("http://jaeger:4318"))
	if err != nil {
		panic(err)
	}

	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("workshop"),
			semconv.DeploymentEnvironment("development"),
			semconv.URLFull("jaeger"),
		),
	)
	if err != nil {
		panic(err)
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(r),
	)

	defer func() {
		traceProvider.Shutdown(ctx)
	}()

	otel.SetTracerProvider(traceProvider)

	tracer := otel.GetTracerProvider().Tracer("workshop")

	c := cache.NewCache(cacheSize, ttl)

	http.HandleFunc("POST /set", func(w http.ResponseWriter, r *http.Request) {
		metrics.IncRequestCounter("set")

		defer func(now time.Time) {
			metrics.StoreHandlerDuration("set", time.Since(now))
		}(time.Now())
		ctx := r.Context()

		ctx, span := tracer.Start(
			ctx,
			"POST /set",
			trace2.WithAttributes(attribute.String("handler", "set")),
		)
		defer span.End()

		ms := rand.Intn(1000)
		time.Sleep(time.Duration(ms) * time.Millisecond)

		span.AddEvent("rand time sleep executed")

		var request struct {
			Key   string
			Value string
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))

			return
		}

		err = json.Unmarshal(body, &request)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))

			return
		}

		if len(request.Value) == 0 || len(request.Key) == 0 {
			loggerPkg.Errorw(ctx, "handler: set, request is invalid")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("value or key is empty"))

			return
		}

		loggerPkg.Infow(ctx, "handler: set", "cache-key", request.Key, "cache-value", request.Value)

		err = c.Set(ctx, request.Key, []byte(request.Value))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))

			return
		}
	})

	http.HandleFunc("GET /get", func(w http.ResponseWriter, r *http.Request) {
		metrics.IncRequestCounter("get")

		defer func(now time.Time) {
			metrics.StoreHandlerDuration("get", time.Since(now))
		}(time.Now())

		ms := rand.Intn(1000)
		time.Sleep(time.Duration(ms) * time.Millisecond)

		ctx := r.Context()

		cacheKey := r.URL.Query().Get("key")
		if cacheKey == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("key must be passed"))

			return
		}

		loggerPkg.Infow(ctx, "handler: get", "cache-key", cacheKey)

		value, err := c.Get(r.Context(), cacheKey)
		if err != nil {

			if errors.Is(err, cache.ErrCacheNoValue) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("key not found"))
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))

			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(value))
	})

	http.Handle("GET /metrics", promhttp.Handler())

	loggerPkg.Errorw(ctx, "app bootstrapped")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
