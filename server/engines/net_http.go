package engines

import "net/http"

type GoHttpEngine struct {
	handler http.Handler
}

func (e *GoHttpEngine) Handler() http.Handler {
	return e.handler
}

func NetHttp() *GoHttpEngine {
	return &GoHttpEngine{
		handler: http.NewServeMux(),
	}
}
