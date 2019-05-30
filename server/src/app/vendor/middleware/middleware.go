package middleware

import "net/http"

// Handler - Type define for method chaining
type Handler func(http.Handler) http.Handler

// Container - for middleware chaining
type Container struct {
	Handlers []Handler
}

// Chain - chain factory
func Chain(handlers ...Handler) Container {
	return Container{append(([]Handler)(nil), handlers...)}
}

// Then - Then dunction
func (c Container) Then(h http.Handler) http.Handler {
	if h == nil {
		h = http.DefaultServeMux
	}

	for i := range c.Handlers {
		h = c.Handlers[len(c.Handlers)-1-i](h)
	}

	return h
}

// ThenFunc provides all the guarantees of Then.
func (c Container) ThenFunc(fn http.Handler) http.Handler {
	if fn == nil {
		return c.Then(nil)
	}
	return c.Then(fn)
}

// Append extends a chain, adding the specified constructors
func (c Container) Append(handlers ...Handler) Container {
	newCons := make([]Handler, 0, len(c.Handlers)+len(handlers))
	newCons = append(newCons, c.Handlers...)
	newCons = append(newCons, handlers...)

	return Container{newCons}
}

// Extend extends a chain by adding the specified chain
func (c Container) Extend(middlewareChain Container) Container {
	return c.Append(middlewareChain.Handlers...)
}
