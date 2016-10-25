package controllers

import (
	"encoding/json"
	"net/http"
	"log"
	"io"
	"io/ioutil"
	"bufio"
	"os"
	"strconv"

	"github.com/deviceMP/api-server/models"
	"github.com/gorilla/mux"
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
			db.Where(models.Device{ProjectId: v.Id}).Find(&devices)
			projects[i].DeviceTotal = len(devices)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&projects)
	}
}

func GetProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orgId := vars["orgId"]
	projectId := vars["projectId"]

	w.Header().Set("Access-Control-Allow-Origin", "*")
	if orgIdInt, err := strconv.Atoi(orgId); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&err)
	} else if projectIdInt, err := strconv.Atoi(projectId); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&err)
	} else {
		var project models.Project
		db.Where(models.Project{OrgId: orgIdInt, Id:projectIdInt}).First(&project)
		var devices []models.Device
		db.Where(models.Device{ProjectId: project.Id}).Find(&devices)
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
			OrgId: 		project.OrgId,
			Name:       	project.Name,
			Description: 	project.Description,
			Image:     	project.Image,
			DeviceType: 	project.DeviceType,
			ApiKey:     	utils.RandStringRunes(32),
			Commit:    	"",
			Repository: 	project.Repository,
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
		if project.Name == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var projectUpdate models.Project
		rowUpdated := db.Model(&projectUpdate).Where(models.Project{Id: project.Id}).UpdateColumn(models.Project{Name: project.Name,
			Description: project.Description,Image:project.Image,Repository:project.Repository}).RowsAffected

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if rowUpdated > 0 {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func DownloadConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)
	projectId := vars["projectId"]

	projectIdInt, err := strconv.Atoi(projectId)
	if err != nil {
		projectIdInt = 0
	}

	var project models.Project
	db.Where(models.Project{Id: projectIdInt}).First(&project)
	projectConfig := models.ProjectConfig{ProjectId: project.Id, ProjectName: project.Name, ApiKey: project.ApiKey, DeviceType: project.DeviceType}

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
	orgId := vars["orgId"]

	if OrgIdInt, err := strconv.Atoi(orgId); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&err)
	} else {
		var org models.Org
		var devices []models.Device
		db.Where(models.Org{Id: OrgIdInt}).First(&org)
		db.Where(models.Device{ProjectId: org.Id}).Find(&devices)
		for _,v := range devices {
			var apps []models.App
			db.Where(models.App{Uuid: v.Uuid}).Find(&apps)
			for _,v := range apps {
				db.Delete(&v)
			}
			db.Delete(&v)
		}

		db.Delete(&org)
		w.WriteHeader(http.StatusOK)
	}
}