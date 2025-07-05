package server

import (
	"github.com/gin-gonic/gin"
)

type Server struct {
	engine  *gin.Engine
	handler *Handler
}

func NewServer(handler *Handler) *Server {
	r := gin.Default()

	r.Static("/static", "./static") 

	r.GET("/order/:order_uid", handler.GetOrder())

	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	return &Server{
		engine:  r,
		handler: handler,
	}
}

func (s *Server) Run(addr string) error {
	return s.engine.Run(addr)
}
