package rock

import "net/http"

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Params  map[string]string // URL parameters (e.g., :id, *filepath)
}

// Param returns the value of the URL parameter
func (c *Context) Param(key string) string {
	return c.Params[key]
}
