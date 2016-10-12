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
	m.Route{"GetOrg", "GET", "/org/detail/{orgId}", c.GetOrg},
	m.Route{"DConfig", "GET", "/org/configfile/{orgId}", c.DownloadConfig},

	m.Route{"UpdateApp", "PUT", "/app", c.UpdateApp},
	m.Route{"CreateApp", "POST", "/app", c.CreateApp},
	m.Route{"GetApp", "GET", "/app/{deviceId}", c.GetApp},
	m.Route{"DeleteApp", "DELETE", "/app/{deviceId}", c.DeleteApp},

	m.Route{"DeviceOnline", "POST", "/device/online", c.DeviceOnline},
	m.Route{"ListDeviceByOrg", "GET", "/device/{orgId}", c.ListDeviceByOrg},
	m.Route{"RegisterDevice", "POST", "/device/{orgId}", c.RegisterDevice},
	m.Route{"GetDeviceById", "GET", "/device/{orgId}/{deviceId}", c.GetDeviceById},
	m.Route{"UpdateState", "POST", "/device/{orgId}/updatestate", c.UpdateState},
	m.Route{"UpdateProgress", "POST", "/device/{orgId}/updateprogress", c.UpdateProgress},
	m.Route{"CheckUpdate", "POST", "/device/{orgId}/checkupdate/{deviceId}", c.CheckUpdate},
	m.Route{"UpdateDeviceName", "POST", "/device/{orgId}/updatename/{deviceId}", c.UpdateDeviceName},
	//m.Route{"CheckForUpdate", "POST", "/device/{orgId}/checkforupdate/{deviceId}", c.CheckForUpdate},
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

