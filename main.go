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
	m.Route{"UpdateProject", "POST", "/org/{orgId}/updateproject", c.UpdateProject},
	m.Route{"GetProject", "GET", "/org/{orgId}/project/{projectId}", c.GetProject},
	m.Route{"DConfig", "GET", "/org/{orgId}/download/{projectId}", c.DownloadConfig},
	m.Route{"DeleteProject", "POST", "/org/{orgId}/delete/{projectId}", c.DeleteProject},
	m.Route{"ListProjectEnv", "GET", "/org/{orgId}/project/{projectId}/env", c.ListProjectEnv},
	m.Route{"AddProjectEnv", "POST", "/org/{orgId}/project/{projectId}/env", c.AddProjectEnv},
	m.Route{"DeleteProjectEnv", "POST", "/deleteenv/{envId}", c.DeleteProjectEnv},

	//Device
	m.Route{"DeviceOnline", "POST", "/org/{orgId}/project/{projectId}/device/online", c.DeviceOnline},
	m.Route{"ListDeviceByProject", "GET", "/org/{orgId}/project/{projectId}/device", c.ListDeviceByProject},
	m.Route{"RegisterDevice", "POST", "/registerdevice", c.RegisterDevice},
	m.Route{"GetDeviceById", "GET", "/org/{orgId}/project/{projectId}/device/{deviceId}", c.GetDeviceById},
	m.Route{"UpdateState", "POST", "/updatestate", c.UpdateState},
	m.Route{"UpdateStatus", "POST", "/updatestatus", c.UpdateStatus},
	m.Route{"UpdateProgress", "POST", "/updateprogress", c.UpdateProgress},
	m.Route{"CheckUpdate", "POST", "/checkupdate/{projectId}/{deviceId}", c.CheckUpdate},
	m.Route{"UpdateDeviceName", "POST", "/updatename/{deviceId}", c.UpdateDeviceName},

	//App
	m.Route{"GetApp", "GET", "/device/{deviceId}/app", c.GetApp},
	m.Route{"UpdateApp", "POST", "/updateapp", c.UpdateApp},
	m.Route{"CreateApp", "POST", "/createapp", c.CreateApp},
	m.Route{"DeleteApp", "DELETE", "/deleteapp", c.DeleteApp},

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

