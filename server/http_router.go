package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Router struct {
	routerTree    map[string]*Router
	wildcard      *Router
	uriKey        string
	methodHandler map[string][]http.HandlerFunc
}

func newRouter() *Router {
	return &Router{
		routerTree:    map[string]*Router{},
		methodHandler: map[string][]http.HandlerFunc{},
	}
}

func (s *Server) handle(uri, method string, handlerChain ...http.HandlerFunc) {
	paths := strings.Split(uri, "/")

	r := s.router
	for _, path := range paths {
		if len(path) < 1 {
			continue
		}
		if len(path) > 1 && path[0] == '<' && path[len(path)-1] == '>' {
			if r.wildcard == nil {
				r.wildcard = newRouter()
			}
			r = r.wildcard
			r.uriKey = path[1 : len(path)-1]
		} else {
			if nr, ok := r.routerTree[path]; ok {
				r = nr
			} else {
				nr = newRouter()
				r.routerTree[path] = nr
				r = nr
			}
		}
	}
	r.methodHandler[method] = handlerChain
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	statusCode := 200
	defer func() {
		params := []byte{}
		if ctxP := r.Context().Value("params"); ctxP != nil {
			params, _ = json.Marshal(ctxP.(*RequestParams).m)
		}
		s.log.Trace(
			statusCode,
			r.Method,
			r.RequestURI,
			string(params),
		)
	}()
	defer func() {
		if err := recover(); err != nil && err != ERROR_HANDLER_CHAIN_ABORD {
			s.log.Error(err)
			s.log.Error(r.Method, r.RequestURI)
			s.log.Error(callstack())
			statusCode = 500
			w.WriteHeader(500)
		}
	}()
	uri := strings.Split(r.RequestURI, "?")[0]

	paths := strings.Split(uri, "/")

	uriParams := map[string]interface{}{}
	router := s.router
	for _, path := range paths {
		if len(path) < 1 {
			continue
		}
		if nr, ok := router.routerTree[path]; ok {
			router = nr
		} else {
			if router.wildcard != nil {
				router = router.wildcard
				if v, err := strconv.Atoi(path); err == nil {
					uriParams[router.uriKey] = v
				} else {
					uriParams[router.uriKey] = path
				}
			} else {
				statusCode = 404
				w.WriteHeader(404)
				return
			}
		}
	}
	hc, ok := router.methodHandler[r.Method]
	if !ok {
		statusCode = 405
		w.WriteHeader(405)
		return
	}
	params := &RequestParams{m: map[string]interface{}{}}
	ctx := context.WithValue(r.Context(), "params", params)
	r = r.WithContext(ctx)

	for k, v := range uriParams {
		params.m[k] = v
	}
	for _, h := range hc {
		h(w, r)
	}
}

type RequestParams struct {
	m map[string]interface{}
}
