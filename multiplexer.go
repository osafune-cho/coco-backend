package main

import (
	"context"
	"net/http"
	"regexp"
	"strings"
)

type router struct {
	routes []route
}

type route struct {
	method  string
	path    *regexp.Regexp
	handler http.HandlerFunc
}

func NewRouter() *router {
	return &router{}
}

func (r *router) Add(method string, pattern string, handler http.HandlerFunc) *router {
	newRoute := route{
		method,
		regexp.MustCompile("^" + pattern + "$"),
		handler,
	}
	r.routes = append(r.routes, newRoute)
	return r
}

type pathParamCtxKey struct{}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var allow []string
	for _, route := range r.routes {
		matches := route.path.FindStringSubmatch(req.URL.Path)
		if len(matches) > 0 {
			if req.Method != route.method {
				allow = append(allow, route.method)
				continue
			}
			ctx := context.WithValue(req.Context(), pathParamCtxKey{}, matches[1:])
			route.handler(w, req.WithContext(ctx))
			return
		}
	}
	if len(allow) > 0 {
		w.Header().Set("Allow", strings.Join(allow, ", "))
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func PathParam(r *http.Request, index int) string {
	params := r.Context().Value(pathParamCtxKey{}).([]string)
	return params[index]
}
