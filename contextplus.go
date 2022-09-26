package contextplus

import (
	"context"
	"time"
)

type valueOnlyContext struct{ ctx context.Context }

func (c valueOnlyContext) Value(key interface{}) interface{}     { return c.ctx.Value(key) }
func (valueOnlyContext) Deadline() (deadline time.Time, ok bool) { return }
func (valueOnlyContext) Done() <-chan struct{}                   { return nil }
func (valueOnlyContext) Err() error                              { return nil }

// WithOnlyValue returns a copy of parent with hide `cancel` and `deadline`
func WithOnlyValue(ctx context.Context) context.Context {
	return valueOnlyContext{ctx: ctx}
}

// WithoutCancel returns a copy of parent with hide `cancel` but keep `deadline`
// The returned context's Done channel is closed when the deadline expires
func WithoutCancel(ctx context.Context) context.Context {
	ctx, _ = WithRebirthCancel(ctx)
	return ctx
}

// WithRebirthCancel returns a copy of parent with hide parent `cancel` and keep `deadline`.
// The returned context's Done channel is closed when the deadline expires, when the returned
// cancel function is called
func WithRebirthCancel(ctx context.Context) (context.Context, context.CancelFunc) {
	deadline, ok := ctx.Deadline()
	ctx = WithOnlyValue(ctx)
	if ok {
		return context.WithDeadline(ctx, deadline)
	}
	return ctx, func() {}
}
