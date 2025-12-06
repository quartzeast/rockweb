package rock

import "net/http"

type HandlerFunc func(http.ResponseWriter, *http.Request)

type router struct {
	handlerFuncMap map[string]HandlerFunc // path -> handler
}

func (r *router) AddRoute(name string, handler HandlerFunc) {
	if r.handlerFuncMap == nil {
		r.handlerFuncMap = make(map[string]HandlerFunc)
	}
	r.handlerFuncMap[name] = handler
}

type Engine struct {
	router
}

func New() *Engine {
	return &Engine{
		router: router{
			handlerFuncMap: make(map[string]HandlerFunc),
		},
	}
}

func (e *Engine) Run(addr string) error {
	for pattern, handler := range e.handlerFuncMap {
		http.HandleFunc(pattern, handler)
	}

	return http.ListenAndServe(addr, nil)
}
