package src

import (
	"fmt"
	"net/http"
	"strings"
)

type Router struct {
	rules map[string]map[string]http.HandlerFunc
}

// Router instance constructor
func NewRouter() *Router {
	return &Router{
		rules: make(map[string]map[string]http.HandlerFunc),
	}
}

//Find if a Handler is registered
func (r *Router) FindHandler(path string, method string) (http.HandlerFunc, bool, bool) {
	var rootPath string
	pathSplit := strings.Split(path, "/")
	if len(pathSplit) > 2 {
		rootPath = fmt.Sprintf("/%s/", pathSplit[1])
	} else {
		rootPath = fmt.Sprintf("/%s", pathSplit[1])
	}
	_, exist := r.rules[rootPath]
	handler, methodExist := r.rules[rootPath][method]
	return handler, methodExist, exist
}

//Find if a Path is registered
func (r *Router) FindPath(path string) bool {
	_, exist := r.rules[path]
	return exist
}

func (r *Router) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	handler, methodExist, exist := r.FindHandler(request.URL.Path, request.Method)
	if !exist {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if !methodExist {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	handler(w, request)
}
