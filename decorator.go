package gonion

type Decorator func(Handler) Handler

// Decorate will decorate a context handler with a slice of passed decorators
func Decorate(c Handler, ds ...Decorator) Handler {
	decorated := c
	for i := 1; i <= len(ds); i++ {
		j := len(ds) - i
		decorate := ds[j]
		decorated = decorate(decorated)
	}

	return decorated
}
