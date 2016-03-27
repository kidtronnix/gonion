# gonion

Golang HTTP handlers abstracted into onions.

## Usage

```go
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

### About

Let me ask the reader a question - what would the ideal http handler look like? 

<img alt="request lifecycle diagram" src="https://docs.google.com/drawings/d/1UslicNjEfqS2rGkINFw3zIK7DPFMSHmA2iszbB0Y6jQ/pub?w=480&amp;h=360" align="right">

An onion right.

Think about it, typically a request lifecycle can be split into distinct stages.
These stages can be tightly coupled, i.e. one stage depends on a value of a previous stage.
Or these stages can be loosely coupled, i.e. one stage has no dependency to the others.

Ok so we need http handlers that can be composed of many individual layers.
These layers should be able to be put together independently of each other as this allows 
for easy testing and flexible composition.
But the layers also need to be able to pass across some concept of request context to allow coupling between stages.

So what about the context, what features should it have?

Well we need a request context to be passed in a way that scales. So no use of global context structs guarded with
mutexes. That approach can introduce lock contention for heavily requested endpoints.

Our context package should be a standard package to allow for creating reusable middleware packages.

On top of this, our context should be safe to pass across service boundaries and processes.
For instance we should be able to safely pass our context to RPC microservices.

The obvious solution to this is google's solution to context, [https://godoc.org/golang.org/x/net/context](https://godoc.org/golang.org/x/net/context).

Put all that together and you have yourself `gonion` - a flexible yet robust way of putting together HTTP handlers.


## Inspiration

The inspirational source of the onion abstraction. A great talk about the benefits of the decorator pattern.

[https://www.youtube.com/watch?v=xyDkyFjzFVc](https://www.youtube.com/watch?v=xyDkyFjzFVc)
