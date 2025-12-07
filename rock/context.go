package rock

import "net/http"

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
}
