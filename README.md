# gonion

Golang HTTP handlers abstracted into onions.

## Usage

```golang
func HelloMiddleware(h Handler) Handler {
	return HandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
		ctx = context.WithValue(ctx, "args", []string{"hello"})
		// Pass to next stage of handler
		h.ServeHTTPContext(ctx, rw, req)
	})
}

func WorldMiddleware(h Handler) Handler {
	return HandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
		args := ctx.Value("args").([]string)
		ctx = context.WithValue(ctx, "args", append(args, "world"))
		// Pass to next stage of handler
		h.ServeHTTPContext(ctx, rw, req)
	})
}

func Response(h Handler) Handler {
	return HandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
		args := ctx.Value("args").([]string)
		str := strings.Join(args, " ")
		fmt.Fprint(rw, str)
	})
}

handler := NewHTTPHandler(
	HelloMiddleware,
	WorldMiddleware,
	Response,
)
// Use `handler` in any router that accepts a standard http.Handler

```

### why?

There are a tonne of golang micro http framework packages. So what does this package solve?

Well okay, let me ask the reader a question; what would the ideal http handler look like?

Well an onion right.

Think about it, typically a request lifecycle can be split into a distinct stages.
These stages can be completely coupled, i.e. one stage depends on a value of a previous middleware.
Or these stages can de de-coupled, one stage has no dependency on another stage.

Ok so we need http handlers that can be composed of many individual layers.
These layers should be able to be put together independently of each other as this allows
for decoupling of stages for easy testing and flexible composition.
But the layers also need to be able to pass across some concept of request context to allow coupling between stages.

So what about the context, what features should it have?

Well we need a request context to be passed in a way that scales. So no use of global structs with
mutexes that can introduce lock contention for heavily requested endpoints.

Our context package should be a standard package to allow for creating reusable middleware packages.
On top of this our context should be safe to pass across service boundaries and processes.
So we should be able to safely pass our context to RPC microservices.

Ok so we should just use google's solution to passing context, [https://godoc.org/golang.org/x/net/context](https://godoc.org/golang.org/x/net/context).

Put all that together and you have yourself gonion; a flexible yet robust way of putting together HTTP handlers.


## Initial Inspiration

An excellent talk about the benefits of using the decorator pattern.

[https://www.youtube.com/watch?v=xyDkyFjzFVc](https://www.youtube.com/watch?v=xyDkyFjzFVc)
