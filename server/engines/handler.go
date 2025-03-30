package engines

import "net/http"

type Handler interface {
	Handler() http.Handler
}
