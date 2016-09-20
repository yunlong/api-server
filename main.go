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
	m.Route{"Ping", "GET", "/ping", c.Ping},

	m.Route{"ListOrg", "GET", "/org", c.ListOrg},
	m.Route{"CreateOrg", "POST", "/org", c.CreateOrg},

	m.Route{"GetApp", "GET", "/org/app/{orgId}", c.GetApp},
	m.Route{"UpdateApp", "POST", "/org/app/{orgId}", c.UpdateApp},

	m.Route{"ListDeviceByApp", "GET", "/device/{orgId}", c.ListDeviceByApp},
	m.Route{"RegisterDevice", "POST", "/device/{orgId}", c.RegisterDevice},
	m.Route{"UpdateState", "POST", "/device/{orgId}/updatestate", c.UpdateState},
}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

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

