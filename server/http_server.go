package server

import (
	"fmt"
	"net/http"
)

type Server struct {
	Wx2B *WxSDK
	Wx2C *WxSDK
	DB   DBConn

	Port string

	router *Router
	log    Logger
}

func NewServer(db DBConn, wx2B, wx2C *WxSDK, port int) *Server {
	s := &Server{
		DB:     db,
		Wx2B:   wx2B,
		Wx2C:   wx2C,
		Port:   fmt.Sprintf(":%d", port),
		router: newRouter(),
		log:    NewStdLogger("[HTTP]"),
	}

	s.initRouter()
	return s
}

func (s *Server) Run() {
	http.ListenAndServe(s.Port, s)
}
