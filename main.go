package main

import (
	"log"
	"net/http"

	"github.com/deviceMP/api-server/controllers"
	"github.com/gorilla/mux"
	m "github.com/deviceMP/api-server/models"
	c "github.com/deviceMP/api-server/controllers"
)

var Routes = m.Routes{
	m.Route{"Ping","GET","/ping",c.Ping,},

	m.Route{"Index","GET","/app",c.Index,},
	m.Route{"CreateApp","POST","/app",c.CreateApp,},

	m.Route{"RegisterDevice","POST","/device",c.RegisterDevice,},
}

func main() {
	go controllers.CheckAllDevice()

	router := NewRouter()
    	log.Fatal(http.ListenAndServe(":8080", router))
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true).PathPrefix("/v1").Subrouter()
	for _, route := range Routes {
		var handler http.Handler

		handler = route.HandlerFunc
		//handler = util.RequireTokenAuthentication

		router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(handler)

	}

	return router
}

