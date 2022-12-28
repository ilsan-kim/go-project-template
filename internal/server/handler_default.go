package server

import (
	"github.com/labstack/echo/v4"
)

type hello struct {
	Hello string `json:"hello"`
}

func (s *Server) helloHandler(c echo.Context) (err error) {
	resp := &hello{"world"}
	return c.JSON(200, resp)
}
