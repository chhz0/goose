package engines

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type EchoEngineWrapper struct {
	echoEngine *echo.Echo
}

func (e *EchoEngineWrapper) Handler() http.Handler {
	return e.echoEngine
}

func Echo(e *echo.Echo) *EchoEngineWrapper {
	return &EchoEngineWrapper{e}
}
