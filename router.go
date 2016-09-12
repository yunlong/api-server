package main

import (
	"net/http"
	//"exposeapi/util"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true).PathPrefix("/v1").Subrouter()
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		//handler = util.RequireTokenAuthentication

		router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(handler)

	}

	return router
}
