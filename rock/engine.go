package rock

import (
	"net/http"
)

const ANY = "ANY"

type HandlerFunc func(ctx *Context)

type routerGroup struct {
	prefix   string
	roots    map[string]*node       // {method: root node}
	handlers map[string]HandlerFunc // {pattern: handler}
}

func (rg *routerGroup) AddRoute(pattern string, method string, handler HandlerFunc) {
	fullPattern := rg.prefix + pattern
	parts := parsePath(fullPattern)

	key := method + "-" + fullPattern
	_, ok := rg.handlers[key]
	if ok {
		panic("route already exists: " + method + " " + fullPattern)
	}

	// Initialize root node for this method if not exists
	if _, ok := rg.roots[method]; !ok {
		rg.roots[method] = &node{}
	}

	// Insert pattern into the tree
	rg.roots[method].insert(fullPattern, parts, 0)
	rg.handlers[key] = handler
}

func (rg *routerGroup) ANY(pattern string, handler HandlerFunc) {
	rg.AddRoute(pattern, ANY, handler)
}

func (rg *routerGroup) GET(pattern string, handler HandlerFunc) {
	rg.AddRoute(pattern, http.MethodGet, handler)
}

func (rg *routerGroup) POST(pattern string, handler HandlerFunc) {
	rg.AddRoute(pattern, http.MethodPost, handler)
}

func (rg *routerGroup) DELETE(pattern string, handler HandlerFunc) {
	rg.AddRoute(pattern, http.MethodDelete, handler)
}

func (rg *routerGroup) PUT(pattern string, handler HandlerFunc) {
	rg.AddRoute(pattern, http.MethodPut, handler)
}

func (rg *routerGroup) PATCH(pattern string, handler HandlerFunc) {
	rg.AddRoute(pattern, http.MethodPatch, handler)
}

func (rg *routerGroup) OPTIONS(pattern string, handler HandlerFunc) {
	rg.AddRoute(pattern, http.MethodOptions, handler)
}

func (rg *routerGroup) HEAD(pattern string, handler HandlerFunc) {
	rg.AddRoute(pattern, http.MethodHead, handler)
}

type router struct {
	routerGroups []*routerGroup
}

func (r *router) Group(prefix string) *routerGroup {
	group := &routerGroup{
		prefix:   prefix,
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
	r.routerGroups = append(r.routerGroups, group)
	return group
}

type Engine struct {
	*router
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	path := r.URL.Path
	parts := parsePath(path)

	ctx := &Context{
		Writer:  w,
		Request: r,
		Params:  make(map[string]string),
	}

	for _, group := range e.routerGroups {
		// Try to match with ANY method first
		if root, ok := group.roots[ANY]; ok {
			if n := root.search(parts, 0); n != nil {
				// Extract params
				e.extractParams(ctx, n.pattern, parts)
				key := ANY + "-" + n.pattern
				if handler, ok := group.handlers[key]; ok {
					handler(ctx)
					return
				}
			}
		}

		// Try to match with specific method
		if root, ok := group.roots[method]; ok {
			if n := root.search(parts, 0); n != nil {
				// Extract params
				e.extractParams(ctx, n.pattern, parts)
				key := method + "-" + n.pattern
				if handler, ok := group.handlers[key]; ok {
					handler(ctx)
					return
				}
			}
		}
	}

	http.NotFound(w, r)
}

// extractParams extracts path parameters from the matched route
func (e *Engine) extractParams(ctx *Context, pattern string, parts []string) {
	patternParts := parsePath(pattern)

	for i, part := range patternParts {
		if len(part) > 0 && part[0] == ':' {
			// Named parameter
			if i < len(parts) {
				paramName := part[1:]
				ctx.Params[paramName] = parts[i]
			}
		} else if len(part) > 0 && part[0] == '*' {
			// Catch-all parameter
			paramName := part[1:]
			if i < len(parts) {
				// Join remaining parts
				remaining := parts[i:]
				value := ""
				for j, p := range remaining {
					if j > 0 {
						value += "/"
					}
					value += p
				}
				ctx.Params[paramName] = value
			}
			break
		}
	}
}

func New() *Engine {
	return &Engine{
		router: &router{},
	}
}

func (e *Engine) Run(addr string) error {
	http.Handle("/", e)
	return http.ListenAndServe(addr, nil)
}
