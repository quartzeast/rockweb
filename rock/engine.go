package rock

import (
	"net/http"
	"slices"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

type routerGroup struct {
	prefix           string
	handlerFuncMap   map[string]HandlerFunc // path -> handler
	handlerMethodMap map[string][]string
}

func (rg *routerGroup) AddRoute(pattern string, handler HandlerFunc) {
	fullPattern := rg.prefix + pattern
	rg.handlerFuncMap[fullPattern] = handler
}

func (rg *routerGroup) ANY(pattern string, handler HandlerFunc) {
	rg.AddRoute(pattern, handler)
	rg.handlerMethodMap["ANY"] = append(rg.handlerMethodMap["ANY"], rg.prefix+pattern)
}

func (rg *routerGroup) GET(pattern string, handler HandlerFunc) {
	rg.AddRoute(pattern, handler)
	rg.handlerMethodMap[http.MethodGet] = append(rg.handlerMethodMap[http.MethodGet], rg.prefix+pattern)
}

func (rg *routerGroup) POST(pattern string, handler HandlerFunc) {
	rg.AddRoute(pattern, handler)
	rg.handlerMethodMap[http.MethodPost] = append(rg.handlerMethodMap[http.MethodPost], rg.prefix+pattern)
}

type router struct {
	routerGroups []*routerGroup
}

func (r *router) Group(prefix string) *routerGroup {
	group := &routerGroup{
		prefix:           prefix,
		handlerFuncMap:   make(map[string]HandlerFunc),
		handlerMethodMap: make(map[string][]string),
	}
	r.routerGroups = append(r.routerGroups, group)
	return group
}

type Engine struct {
	router
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	for _, group := range e.routerGroups {
		for pattern, handler := range group.handlerFuncMap {
			if r.RequestURI == pattern {
				allowedPatterns, ok := group.handlerMethodMap["ANY"]
				if ok {
					if slices.Contains(allowedPatterns, pattern) {
						handler(w, r)
						return
					}
				}

				// 根据 method 进行匹配
				allowedPatterns, ok = group.handlerMethodMap[method]
				if ok {
					if slices.Contains(allowedPatterns, pattern) {
						handler(w, r)
						return
					}
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
		router: router{},
	}
}

func (e *Engine) Run(addr string) error {
	http.Handle("/", e)
	return http.ListenAndServe(addr, nil)
}
