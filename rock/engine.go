package rock

import "net/http"

type HandlerFunc func(http.ResponseWriter, *http.Request)

type routerGroup struct {
	prefix         string
	handlerFuncMap map[string]HandlerFunc // path -> handler
}

func (rg *routerGroup) AddRoute(pattern string, handler HandlerFunc) {
	fullPattern := rg.prefix + pattern
	rg.handlerFuncMap[fullPattern] = handler
}

type router struct {
	routerGroups []*routerGroup
}

func (r *router) Group(prefix string) *routerGroup {
	group := &routerGroup{
		prefix:         prefix,
		handlerFuncMap: make(map[string]HandlerFunc),
	}
	r.routerGroups = append(r.routerGroups, group)
	return group
}

type Engine struct {
	router
}

func New() *Engine {
	return &Engine{
		router: router{},
	}
}

func (e *Engine) Run(addr string) error {
	for _, group := range e.routerGroups {
		for pattern, handler := range group.handlerFuncMap {
			http.HandleFunc(pattern, handler)
		}
	}

	return http.ListenAndServe(addr, nil)
}
