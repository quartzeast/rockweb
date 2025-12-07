package rock

import (
	"net/http"
)

const ANY = "ANY"

type HandlerFunc func(ctx *Context)

type routerGroup struct {
	prefix           string
	handlerFuncMap   map[string]map[string]HandlerFunc // {pattern: {method: handler} }
	handlerMethodMap map[string][]string
}

func (rg *routerGroup) AddRoute(pattern string, method string, handler HandlerFunc) {
	fullPattern := rg.prefix + pattern
	_, ok := rg.handlerFuncMap[fullPattern]
	if !ok {
		rg.handlerFuncMap[fullPattern] = make(map[string]HandlerFunc)
	}

	_, ok = rg.handlerFuncMap[fullPattern][method]
	if ok {
		panic("route already exists: " + method + " " + fullPattern)
	}

	rg.handlerFuncMap[fullPattern][method] = handler
	rg.handlerMethodMap[method] = append(rg.handlerMethodMap[method], fullPattern)
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
		prefix:           prefix,
		handlerFuncMap:   make(map[string]map[string]HandlerFunc),
		handlerMethodMap: make(map[string][]string),
	}
	r.routerGroups = append(r.routerGroups, group)
	return group
}

type Engine struct {
	*router
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	for _, group := range e.routerGroups {
		for pattern, methodMap := range group.handlerFuncMap {
			if r.RequestURI == pattern {
				ctx := &Context{
					Writer:  w,
					Request: r,
				}

				handler, ok := methodMap[ANY]
				if ok {
					handler(ctx)
					return
				}

				handler, ok = methodMap[method]
				if ok {
					handler(ctx)
					return
				}

				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
				return
			}
		}
	}

	http.NotFound(w, r)
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
