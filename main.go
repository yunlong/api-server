package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	m "github.com/deviceMP/api-server/models"
	c "github.com/deviceMP/api-server/controllers"
)

var Routes = m.Routes{
	m.Route{"ListOrg", "GET", "/org", c.ListOrg},
	m.Route{"CreateOrg", "POST", "/org", c.CreateOrg},
	m.Route{"GetOrg", "GET", "/org/{orgId}", c.GetOrg},
	m.Route{"DownloadConfig", "GET", "/org/{orgId}/configfile", c.DownloadConfig},

	m.Route{"GetApp", "GET", "/org/app/{orgId}", c.GetApp},
	m.Route{"DeleteApp", "DELETE", "/org/app/{orgId}", c.DeleteApp},
	m.Route{"UpdateApp", "PUT", "/org/app/{orgId}", c.UpdateApp},
	m.Route{"CreateApp", "POST", "/org/app/{orgId}", c.CreateApp},

	m.Route{"DeviceOnline", "POST", "/device/online", c.DeviceOnline},
	m.Route{"ListDeviceByApp", "GET", "/device/{orgId}", c.ListDeviceByApp},
	m.Route{"RegisterDevice", "POST", "/device/{orgId}", c.RegisterDevice},
	m.Route{"GetDeviceById", "GET", "/device/{orgId}/{deviceId}", c.GetDeviceById},
	m.Route{"UpdateState", "POST", "/device/{orgId}/updatestate", c.UpdateState},
	m.Route{"CheckUpdate", "POST", "/device/{orgId}/checkupdate/{deviceId}", c.CheckUpdate},
}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

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

