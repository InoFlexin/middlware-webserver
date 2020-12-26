package thttp

import (
	"net/http"
)

type HttpHandler struct{}

func ListenAndServe(port string, handler HttpHandler) {
	
}

func(h HttpHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// data := []byte("hello world")
	// res.Write(data)
}