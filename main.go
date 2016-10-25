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
	m.Route{"UploadOrgImage", "POST", "/orgimage", c.UploadOrgImage},
	m.Route{"GetOrg", "GET", "/org/{orgId}", c.GetOrg},
	m.Route{"DeleteOrg", "POST", "/org/{orgId}", c.DeleteOrg},

	//Project
	m.Route{"ListProjectByOrg", "GET", "/org/{orgId}/project", c.ListProjectByOrg},
	m.Route{"CreateProject", "POST", "/org/{orgId}/project", c.CreateProject},
	m.Route{"UpdateProject", "PUT", "/org/{orgId}/project", c.UpdateProject},
	m.Route{"GetProject", "GET", "/org/{orgId}/project/{projectId}", c.GetProject},
	m.Route{"DConfig", "GET", "/org/{orgId}/download/{projectId}", c.DownloadConfig},
	m.Route{"DeleteProject", "POST", "/org/{orgId}/delete/{projectId}", c.DeleteProject},

	//Device
	m.Route{"DeviceOnline", "POST", "/org/{orgId}/project/{projectId}/device/online", c.DeviceOnline},
	m.Route{"ListDeviceByProject", "GET", "/org/{orgId}/project/{projectId}/device", c.ListDeviceByProject},
	m.Route{"RegisterDevice", "POST", "/org/{orgId}/project/{projectId}/device", c.RegisterDevice},
	m.Route{"GetDeviceById", "GET", "/org/{orgId}/project/{projectId}/device/{deviceId}", c.GetDeviceById},
	m.Route{"UpdateState", "POST", "/org/{orgId}/project/{projectId}/device/{deviceId}/updatestate", c.UpdateState},
	m.Route{"UpdateStatus", "POST", "/org/{orgId}/project/{projectId}/device/{deviceId}/updatestatus", c.UpdateStatus},
	m.Route{"UpdateProgress", "POST", "/org/{orgId}/project/{projectId}/device/{deviceId}/updateprogress", c.UpdateProgress},
	m.Route{"CheckUpdate", "POST", "/org/{orgId}/project/{projectId}/device/{deviceId}/checkupdate", c.CheckUpdate},
	m.Route{"UpdateDeviceName", "POST", "/org/{orgId}/project/{projectId}/device/{deviceId}/updatename", c.UpdateDeviceName},

	//App
	m.Route{"GetApp", "GET", "/device/{deviceId}/app", c.GetApp},
	m.Route{"UpdateApp", "PUT", "/device/{deviceId}/app", c.UpdateApp},
	m.Route{"CreateApp", "POST", "/device/{deviceId}/app", c.CreateApp},
	m.Route{"DeleteApp", "DELETE", "/device/{deviceId}/app", c.DeleteApp},

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

