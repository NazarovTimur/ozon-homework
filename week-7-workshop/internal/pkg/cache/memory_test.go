package cache

import (
	"context"
	"gitlab.ozon.dev/14/week-7-workshop/pkg/logger"
	"go.uber.org/zap"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestLeakyCache_Get(t *testing.T) {
	type fields struct {
		storage sync.Map
		ttl     time.Duration
	}
	type args struct {
		in0 context.Context
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &LeakyCache{
				storage: tt.fields.storage,
				ttl:     tt.fields.ttl,
			}

			config := zap.NewProductionConfig()
			config.OutputPaths = []string{"/dev/null"}
			config.Level.SetLevel(zap.InfoLevel)

			ctx := logger.ToContext(tt.args.in0, logger.NewLogger(config))

			got, err := c.Get(ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}
