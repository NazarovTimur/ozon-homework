package cache

import (
	"context"
	"errors"
	"gitlab.ozon.dev/14/week-7-workshop/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"math/rand"
	"sync"
	"time"
)

var (
	ErrCacheNoValue            = errors.New("value not found by key")
	ErrCacheUndefinedValueType = errors.New("value type is not defined")
)

type LeakyCache struct {
	storage sync.Map
	// map[string][]byte
	//mx      sync.Mutex
	ttl time.Duration
}

func NewCache(_ int, ttl time.Duration) *LeakyCache {
	return &LeakyCache{
		storage: sync.Map{},
		ttl:     ttl,
	}
}

func (c *LeakyCache) Set(ctx context.Context, key string, value []byte) error {
	tracer := otel.GetTracerProvider().Tracer("workshop")
	ctx, span := tracer.Start(ctx, "LeakyCache.Set")
	defer span.End()

	ms := rand.Intn(50)
	time.Sleep(time.Duration(ms) * time.Millisecond)

	// Получить span из контекста
	logger.Infow(ctx, "trace_id fetched", "trace_id", trace.SpanFromContext(ctx).SpanContext().TraceID().String())

	c.storage.Store(key, value)

	now := time.Now()
	buf := make([]byte, 0)
	go func() {
		for {
			buf = append(buf, make([]byte, 10<<20)...)
			select {
			default:
				if time.Since(now) > c.ttl {
					c.storage.Delete(key)
					break
				}
			}
			time.Sleep(1 * time.Second)
		}
	}()

	return nil
}

func (c *LeakyCache) Get(_ context.Context, key string) ([]byte, error) {
	value, ok := c.storage.Load(key)

	if !ok {
		return nil, ErrCacheNoValue
	}

	data, ok := value.([]byte)
	if !ok {
		return nil, ErrCacheUndefinedValueType
	}

	return data, nil
}

func (c *LeakyCache) Close() error {
	return nil
}
