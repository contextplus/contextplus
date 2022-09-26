[![Go Report Card](https://goreportcard.com/badge/github.com/contextplus/contextplus)](https://goreportcard.com/report/github.com/contextplus/contextplus)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/contextplus/contextplus)](https://pkg.go.dev/github.com/contextplus/contextplus)
[![Go package](https://github.com/contextplus/contextplus/actions/workflows/test.yaml/badge.svg?branch=main)](https://github.com/contextplus/contextplus/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/contextplus/contextplus/branch/main/graph/badge.svg?token=L8E15TDQFP)](https://codecov.io/gh/contextplus/contextplus)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)

# contextplus

Package contextplus provide more easy to use functions for contexts.

# Use

```
go get -u github.com/contextplus/contextplus
```
# Example
This is how an http.Handler should run a goroutine that need values from the context. Pretend to use the middleware that timeout after one second.

```golang
func myMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx, cancel := context.WithTimeout(time.Second)
		defer cancel()
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
```

```golang
func handle(w http.ResponseWriter, r *http.Request) {
	asyncCtx := contextplus.WithOnlyValue(r.Context())
	go func() {
		// will not cancel if timeout
		// will not cancel if call 'cancel'
		asyncTask(asyncCtx)
	}()
}
```

```golang
func handle(w http.ResponseWriter, r *http.Request) {
	asyncCtx := contextplus.WithoutCancel(r.Context())
	go func() {
		// will cancel if timeout
		// will not cancel if call 'cancel'
		asyncTask(asyncCtx)
	}()
}
```
