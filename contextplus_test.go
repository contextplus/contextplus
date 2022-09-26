package contextplus_test

import (
	"context"
	"testing"
	"time"

	"github.com/contextplus/contextplus"
	"gopkg.in/go-playground/assert.v1"
)

type ctxtype string

func TestWithOnlyValue(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		stub func(*testing.T, context.Context)
	}{
		{
			name: "normal",
			args: args{
				ctx: context.WithValue(context.TODO(), ctxtype("k"), "world"),
			},
			stub: func(t *testing.T, ctx context.Context) {
				assert.Equal(t, ctx.Value(ctxtype("k")).(string), "world")
			},
		},
		{
			name: "nil Deadline()",
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
					_ = cancel
					return ctx
				}(),
			},
			stub: func(t *testing.T, ctx context.Context) {
				_, ok := ctx.Deadline()
				assert.Equal(t, ok, false)
			},
		},
		{
			name: "nil Done()",
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithTimeout(context.TODO(), time.Microsecond)
					_ = cancel
					return ctx
				}(),
			},
			stub: func(t *testing.T, ctx context.Context) {
				assert.Equal(t, ctx.Done(), nil)
			},
		},
		{
			name: "nil Err()",
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
					cancel()
					return ctx
				}(),
			},
			stub: func(t *testing.T, ctx context.Context) {
				assert.Equal(t, ctx.Err(), nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := contextplus.WithOnlyValue(tt.args.ctx)
			tt.stub(t, got)
		})
	}
}

func TestWithoutCancel(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		stub func(*testing.T, context.Context)
	}{
		{
			name: "normal",
			args: args{
				ctx: context.WithValue(context.TODO(), ctxtype("k"), "world"),
			},
			stub: func(t *testing.T, ctx context.Context) {
				assert.Equal(t, ctx.Value(ctxtype("k")).(string), "world")
			},
		},
		{
			name: "after withCancel",
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithCancel(context.TODO())
					_ = cancel
					return ctx
				}(),
			},
			stub: func(t *testing.T, ctx context.Context) {
				_, ok := ctx.Deadline()
				assert.Equal(t, ok, false)
			},
		},
		{
			name: "after withDeadline",
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithDeadline(context.TODO(), time.Date(3000, 1, 1, 1, 1, 1, 1, time.Local))
					_ = cancel
					return ctx
				}(),
			},
			stub: func(t *testing.T, ctx context.Context) {
				deadline, ok := ctx.Deadline()
				assert.Equal(t, ok, true)
				assert.Equal(t, deadline, time.Date(3000, 1, 1, 1, 1, 1, 1, time.Local))
			},
		},
		{
			name: "after cancel and withtimeout context",
			args: args{
				ctx: func() context.Context {
					now := time.Now()
					ctx := context.WithValue(context.Background(), ctxtype("t"), now)
					ctx, cancel := context.WithDeadline(ctx, now.Add(1000*time.Millisecond))
					go func() {
						time.Sleep(100 * time.Millisecond)
						cancel()
					}()
					return ctx
				}(),
			},
			stub: func(t *testing.T, ctx context.Context) {
				<-ctx.Done()
				duration := time.Since(ctx.Value(ctxtype("t")).(time.Time))
				assert.Equal(t, ctx.Err(), context.DeadlineExceeded)
				if duration < 200*time.Millisecond {
					t.Error("timeout cancel ")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := contextplus.WithoutCancel(tt.args.ctx)
			tt.stub(t, got)
		})
	}
}
