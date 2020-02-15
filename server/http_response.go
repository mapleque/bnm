package server

import (
	"encoding/json"
	"net/http"
)

type HttpResponse struct {
	Status  int         `json:"status"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
}

func (s *Server) Response(w http.ResponseWriter, status int, err interface{}) {
	if rerr, ok := err.(error); ok {
		err = rerr.Error()
	}
	ret, _ := json.Marshal(HttpResponse{
		Status:  status,
		Message: err,
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write(ret)
}

func (s *Server) Success(w http.ResponseWriter, data interface{}) {
	ret, _ := json.Marshal(HttpResponse{
		Status:  STATUS_SUCCESS,
		Message: nil,
		Data:    data,
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write(ret)
}
