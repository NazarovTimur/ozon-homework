# Week 7 Workshop - Observability

## Adding Observability to Cache App 

## Pprof
1. Снять профиль памяти приложения в течение 5 секунд и посмотреть его через http на 8000 порту
```shell
go tool pprof -http=":8000" http://localhost:8080/debug/pprof/heap?seconds=5
```

2. Снять профиль CPU приложения в течение 5 секунд и посмотреть его через http на 8000 порту
```shell
go tool pprof -http=":8000" http://localhost:8080/debug/pprof/profile?seconds=10
```

## Runtime tracing
TRACE:
```shell
# снять трейс
curl -o trace.out http://localhost:8080/debug/pprof/trace?seconds=5

# отобразить трейс
go tool trace trace.out
```

## Полезные ссылки

- OpenTelemetry Instrumentation https://opentelemetry.io/docs/languages/go/instrumentation/
- Prometheus Instrumentation https://prometheus.io/docs/guides/go-application/
- Zap Docs https://pkg.go.dev/go.uber.org/zap
- pprof Docs https://pkg.go.dev/net/http/pprof
- Runtime Tracing Docs https://pkg.go.dev/cmd/trace 
- Diagnostics Go programs https://go.dev/doc/diagnostics 
- Profiling Go programs https://go.dev/blog/pprof
- How to read pprof https://github.com/google/pprof/blob/main/doc/README.md
