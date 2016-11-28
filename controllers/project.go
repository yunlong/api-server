package controllers

import (
	"log"
	"io"
	"os"
	"bufio"
	"strconv"
	"net/http"
	"io/ioutil"
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/deviceMP/api-server/models"
	"github.com/deviceMP/api-server/utils"
)

func ListProjectByOrg(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orgId := vars["orgId"]

	w.Header().Set("Access-Control-Allow-Origin", "*")
	if OrgIdInt, err := strconv.Atoi(orgId); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&err)
	} else {
		var projects []models.Project
		db.Where(models.Project{OrgId:OrgIdInt}).Find(&projects)

		for i, v := range projects {
			var devices []models.Device
			db.Where(models.Device{ProjectId: v.ID}).Find(&devices)
			projects[i].DeviceTotal = len(devices)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&projects)
	}
}

func GetProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orgId := vars["orgId"]
	ProjectId := vars["projectId"]

	w.Header().Set("Access-Control-Allow-Origin", "*")
	if orgIdInt, err := strconv.Atoi(orgId); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&err)
	} else if ProjectIdInt, err := strconv.Atoi(ProjectId); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&err)
	} else {
		var project models.Project
		db.Where(models.Project{OrgId: orgIdInt, ID:ProjectIdInt}).First(&project)
		var devices []models.Device
		db.Where(models.Device{ProjectId: project.ID}).Find(&devices)
		project.DeviceTotal = len(devices)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&project)
	}
}

func CreateProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var project models.Project
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
		return
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
		return
	}
	if err := json.Unmarshal(body, &project); err != nil {
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
			return
		}
	} else {
		if project.Name == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		projectCreate := models.Project{
			OrgId: 			project.OrgId,
			Name:       	project.Name,
			Description: 	project.Description,
			Image:     		project.Image,
			DeviceType: 	project.DeviceType,
			ApiKey:     	utils.RandStringRunes(32),
			Commit:    		project.Commit,
			Port:			project.Port,
			Privileged:		project.Privileged,
			Repository: 	project.Repository,
			Environment: 	project.Environment,
		}

		if dbc := db.Create(&projectCreate); dbc.Error != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusCreated)
			if err := json.NewEncoder(w).Encode(&projectCreate); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func UpdateProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var project models.Project
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
		return
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
		return
	}
	if err := json.Unmarshal(body, &project); err != nil {
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
			return
		}
	} else {
		var projectUpdate models.Project
		db.Model(&projectUpdate).Where(models.Project{ID: project.ID}).UpdateColumn(models.Project{Name: project.Name,
			Description: project.Description,Image:project.Image,Repository:project.Repository, Port:project.Port})
		w.WriteHeader(http.StatusOK)
	}
}

func UpdateProjectApp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var project models.Project
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
		return
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
		return
	}
	if err := json.Unmarshal(body, &project); err != nil {
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
			return
		}
	} else {
		//log.Println(project.Privileged)
		var update models.Project
		db.Where(models.Project{ID: project.ID}).First(&update)
		update.Repository = project.Repository
		update.Port = project.Port
		update.Privileged = project.Privileged
		update.Commit = project.Commit
		db.Save(&update)
		//rowUpdated := db.Where(models.Project{ID: project.ID}).First(&update).UpdateColumns(models.Project{Repository: project.Repository, Port: project.Port, Privileged: project.Privileged}).RowsAffected
		//if rowUpdated > 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		/*} else {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusInternalServerError)
		}*/
	}
}

func DownloadConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)
	ProjectId := vars["projectId"]

	ProjectIdInt, err := strconv.Atoi(ProjectId)
	if err != nil {
		ProjectIdInt = 0
	}

	var project models.Project
	db.Where(models.Project{ID: ProjectIdInt}).First(&project)
	projectConfig := models.ProjectConfig{OrgId: project.OrgId, ProjectId: project.ID, ProjectName: project.Name, ApiKey: project.ApiKey, DeviceType: project.DeviceType}

	var fileName = "config.json"
	var filePath = "/tmp/" + fileName

	f, err := os.Create(filePath)
	handleError(err)

	b, err := json.Marshal(projectConfig)
	handleError(err)

	_, err = f.Write(b)
	handleError(err)
	f.Close()

	file, err := os.Open(filePath)
	handleError(err)

	r4 := bufio.NewReader(file)
	//copy the relevant headers. If you want to preserve the downloaded file name, extract it with go's url parser.
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName+"")
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	w.Header().Set("Content-Length", r.Header.Get("Content-Length"))
	//stream the body to the client without fully loading it into memory
	io.Copy(w, r4)
}

func DeleteProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)
	ProjectId := vars["projectId"]

	if ProjectIdInt, err := strconv.Atoi(ProjectId); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&err)
	} else {
		var project models.Project
		var devices []models.Device
		db.Where(models.Project{ID: ProjectIdInt}).First(&project)
		db.Where(models.Device{ProjectId: project.ID}).Find(&devices)
		for _,v := range devices {
			var apps []models.App
			db.Where(models.App{Uuid: v.Uuid}).Find(&apps)
			for _,v := range apps {
				db.Delete(&v)
			}
			db.Delete(&v)
		}

		db.Delete(&project)
		w.WriteHeader(http.StatusOK)
	}
}

func ListProjectEnv(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectId := vars["projectId"]

	w.Header().Set("Access-Control-Allow-Origin", "*")
	if projectIdInt, err := strconv.Atoi(projectId); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&err)
	} else {
		var listEnv []models.ProjectEnv
		db.Where(models.ProjectEnv{ProjectID:projectIdInt}).Find(&listEnv)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&listEnv)
	}
}

//Add project env variable
//Params ProjectId, key, value
func AddProjectEnv(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	clientApiKey := r.URL.Query()["apikey"]
    if len(clientApiKey) == 0 {
        w.WriteHeader(http.StatusNotAcceptable)
        return
    }

	var newEnv models.ProjectEnv
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        log.Println(err)
        return
    }
    if err := r.Body.Close(); err != nil {
        log.Println(err)
        return
    }
    if err := json.Unmarshal(body, &newEnv); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(422)
        if err := json.NewEncoder(w).Encode(err); err != nil {
            log.Println(err)
            return
        }
    }

    var project models.Project
    if db.Where(models.Project{ID: newEnv.ProjectID}).First(&project).RecordNotFound() {
        w.WriteHeader(http.StatusNotFound)
        return
    } else {
        if clientApiKey[0] != project.ApiKey {
            w.WriteHeader(http.StatusForbidden)
            return
        } else {
        	var oldEnv []models.ProjectEnv
        	db.Where(models.ProjectEnv{ProjectID: project.ID}).Find(&oldEnv)

        	envExist := false
        	for _, v := range oldEnv {
        		if v.Key == newEnv.Key {
        			envExist = true
        		}
        	}

        	if !envExist {
            	if dbc := db.Create(&newEnv); dbc.Error != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				} else {
					w.Header().Set("Content-Type", "application/json; charset=UTF-8")
					w.WriteHeader(http.StatusCreated)
					if err := json.NewEncoder(w).Encode(&newEnv); err != nil {
						log.Println(err)
						return
					}
				}
        	}
        }
    }
}

func UpdateProjectEnv(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	projectId := vars["projectId"]

	clientApiKey := r.URL.Query()["apikey"]
    if len(clientApiKey) == 0 {
        w.WriteHeader(http.StatusNotAcceptable)
        return
    }

	var newEnv models.ProjectEnv
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        log.Println(err)
        return
    }
    if err := r.Body.Close(); err != nil {
        log.Println(err)
        return
    }
    if err := json.Unmarshal(body, &newEnv); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(422)
        if err := json.NewEncoder(w).Encode(err); err != nil {
            log.Println(err)
            return
        }
    }

	if projectIdInt, err := strconv.Atoi(projectId); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&err)
	} else {
		var project models.Project
		if db.Where(models.Project{ID: projectIdInt}).First(&project).RecordNotFound() {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
				if clientApiKey[0] != project.ApiKey {
				w.WriteHeader(http.StatusForbidden)
				return
			} else {
				var oldEnv models.ProjectEnv
				rowUpdated := db.Where(models.ProjectEnv{ID: newEnv.ID}).First(&oldEnv).UpdateColumn(models.ProjectEnv{Value: newEnv.Value}).RowsAffected

				if rowUpdated > 0 {
					w.Header().Set("Content-Type", "application/json; charset=UTF-8")
					w.WriteHeader(http.StatusOK)
				} else {
					w.Header().Set("Content-Type", "application/json; charset=UTF-8")
					w.WriteHeader(http.StatusInternalServerError)
				}
			}
		}
	}
}

func DeleteProjectEnv(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)
	projectId := vars["projectId"]

	clientApiKey := r.URL.Query()["apikey"]
    if len(clientApiKey) == 0 {
        w.WriteHeader(http.StatusNotAcceptable)
        return
    }

	var newEnv models.ProjectEnv
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        log.Println(err)
        return
    }
    if err := r.Body.Close(); err != nil {
        log.Println(err)
        return
    }
    if err := json.Unmarshal(body, &newEnv); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(422)
        if err := json.NewEncoder(w).Encode(err); err != nil {
            log.Println(err)
            return
        }
    }

	if projectIdInt, err := strconv.Atoi(projectId); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&err)
	} else {
		var project models.Project
		if db.Where(models.Project{ID: projectIdInt}).First(&project).RecordNotFound() {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
				if clientApiKey[0] != project.ApiKey {
				w.WriteHeader(http.StatusForbidden)
				return
			} else {
				db.Where(models.ProjectEnv{ID: newEnv.ID}).Delete(&models.ProjectEnv{})
				w.WriteHeader(http.StatusOK)
			}
		}
	}
}

func UpdateProjectAppEnv(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	projectId := vars["projectId"]

	clientApiKey := r.URL.Query()["apikey"]
    if len(clientApiKey) == 0 {
        w.WriteHeader(http.StatusNotAcceptable)
        return
    }

	if projectIdInt, err := strconv.Atoi(projectId); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&err)
	} else {
		var project models.Project
		if db.Where(models.Project{ID: projectIdInt}).First(&project).RecordNotFound() {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
				if clientApiKey[0] != project.ApiKey {
				w.WriteHeader(http.StatusForbidden)
				return
			} else {
				//Update all device running on this project
				var deviceUpdate []models.Device
				db.Where(models.Device{ProjectId: project.ID}).Find(&deviceUpdate)
				for _,v := range deviceUpdate {
					go callAgentUpdateEnv(v.Uuid)
				}

				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusOK)
			}
		}
	}
}

func callAgentUpdateEnv(deviceUuid string){
	PushActionAgent(deviceUuid, RestartDeviceApp)
}