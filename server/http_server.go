package server

import (
	"net/http"
)

type Server struct {
	qiniuAccessKey string
	qiniuSecretKey string
	Wx             *WxSDK
	DB             DBConn

	Host string

	router *Router
	log    Logger
}

func NewServer(db DBConn, wx *WxSDK, host, qiniuAccessKey, qiniuSecretKey string) *Server {
	s := &Server{
		DB:             db,
		Wx:             wx,
		Host:           host,
		qiniuAccessKey: qiniuAccessKey,
		qiniuSecretKey: qiniuSecretKey,
		router:         newRouter(),
		log:            NewStdLogger("[HTTP]"),
	}

	s.initRouter()
	return s
}

func (s *Server) Run() {
	if err := http.ListenAndServe(s.Host, s); err != nil {
		panic(err)
	}
}
