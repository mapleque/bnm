package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (s *Server) SaveParam(r *http.Request, k string, v interface{}) {
	m := r.Context().Value("params").(*RequestParams).m
	if t, ok := m[k]; !ok {
		m[k] = v
	} else {
		s.log.Error("conflict param:", k, v, "use origin value:", t)
	}
}

func (s *Server) GetParamPure(r *http.Request, k string) (interface{}, bool) {
	v, ok := r.Context().Value("params").(*RequestParams).m[k]
	return v, ok
}

func (s *Server) FormParam(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.Response(w, STATUS_INVALID_PARAM, MSG_INVALID_BODY)
		panic(ERROR_HANDLER_CHAIN_ABORD)
	}
	m := map[string]interface{}{}
	if err := json.Unmarshal(body, &m); err != nil {
		s.Response(w, STATUS_INVALID_PARAM, MSG_INVALID_BODY)
		panic(ERROR_HANDLER_CHAIN_ABORD)
	}
	for k, v := range m {
		s.SaveParam(r, k, v)
	}
}

func (s *Server) PagerParam(w http.ResponseWriter, r *http.Request) {
	size, err := strconv.Atoi(r.FormValue("s"))
	if err != nil {
		s.Response(w, STATUS_INVALID_PARAM, MSG_INVALID_PAGER)
		panic(ERROR_HANDLER_CHAIN_ABORD)
	}
	lastId, err := strconv.Atoi(r.FormValue("l"))
	if err != nil {
		lastId = 0
	}
	po := &PagerParam{
		Size:   size,
		LastId: lastId,
	}
	s.SaveParam(r, "pager", po)
}

type HttpParam struct {
	key   string
	value interface{}
	ok    bool
}

func (s *Server) GetParam(r *http.Request, key string) *HttpParam {
	v, ok := r.Context().Value("params").(*RequestParams).m[key]
	return &HttpParam{
		key:   key,
		value: v,
		ok:    ok,
	}
}

func (p *HttpParam) Int() int {
	if p.ok {
		switch v := p.value.(type) {
		case float64:
			return int(v)
		case int:
			return v
		default:
			return 0
		}
	}
	return 0
}

func (p *HttpParam) String() string {
	if p.ok {
		if v, ok := p.value.(string); ok {
			return v
		}
	}
	return ""
}

type PagerParam struct {
	Size   int
	LastId int
}

func (s *Server) GetPager(r *http.Request) *PagerParam {
	v, ok := r.Context().Value("params").(*RequestParams).m["pager"]
	if !ok {
		panic("PagerParam should be use before your GetPager")
	}
	p := v.(*PagerParam)
	return p
}

func (p *PagerParam) Sql() string {
	if p.LastId == 0 {
		return fmt.Sprintf(
			" ORDER BY id DESC LIMIT %d",
			p.Size,
		)
	}
	return fmt.Sprintf(
		" AND id < %d ORDER BY id DESC LIMIT %d",
		p.LastId,
		p.Size,
	)
}
