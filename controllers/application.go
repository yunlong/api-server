package controllers

import (
	"encoding/json"
	"net/http"
	"io"
	"io/ioutil"

	"github.com/deviceMP/api-server/models"
	"github.com/gorilla/mux"
)

func GetApp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	deviceId := vars["deviceId"]

	w.Header().Set("Access-Control-Allow-Origin", "*")

	var appGet []models.App
	if !db.Where(models.App{Uuid: deviceId}).Find(&appGet).RecordNotFound() {
		json.NewEncoder(w).Encode(&appGet)
	}
}

func UpdateApp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var app models.App
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
		return
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
		return
	}
	if err := json.Unmarshal(body, &app); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(&err); err != nil {
			panic(err)
			return
		}
	} else {
		var update models.App
		db.Where(models.App{Uuid: app.Uuid}).First(&update)
		update.Commit = app.Commit
		update.ContainerId = app.ContainerId
		update.Port = app.Port
		update.ImageId = app.ImageId
		db.Save(&update)

		//rowUpdated := db.Model(&appUpdate).UpdateColumn(models.App{Commit: app.Commit, ContainerId: app.ContainerId, Port: app.Port, ImageId: app.ImageId}).RowsAffected
		w.WriteHeader(http.StatusOK)
        return
		
		//return
	}
}

func CreateApp(w http.ResponseWriter, r *http.Request) {
	var app models.App
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
		return
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
		return
	}
	if err := json.Unmarshal(body, &app); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(&err); err != nil {
			panic(err)
			return
		}
	} else {
		appCreate := models.App{
			Uuid:        app.Uuid,
			Commit:      app.Commit,
			ContainerId: app.ContainerId,
			Port:        app.Port,
			ImageId:     app.ImageId,
			Latest: 	 app.Latest,
		}

		if dbc := db.Create(&appCreate); dbc.Error != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(dbc.Error); err != nil {
				panic(err)
			}
			return
		} else {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusCreated)

			if err := json.NewEncoder(w).Encode(&appCreate); err != nil {
				panic(err)
			}
		}
	}
}

func DeleteApp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceId := vars["deviceId"]

	w.Header().Set("Access-Control-Allow-Origin", "*")

	var appDelete models.App
	db.Where(models.App{Uuid: deviceId}).First(&appDelete)
	db.Delete(&appDelete)

	w.WriteHeader(http.StatusOK)
}

func UpdateAppEnv(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)
	deviceUuid := vars["deviceuuid"]

	PushActionAgent(deviceUuid, RestartDeviceApp)
	w.WriteHeader(http.StatusOK)
}

func CheckForUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)
	deviceUuid := vars["deviceuuid"]

	PushActionAgent(deviceUuid, CheckUpdate)
	w.WriteHeader(http.StatusOK)
}

func InstallAppUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)
	deviceUuid := vars["deviceuuid"]

	PushActionAgent(deviceUuid, InstallUpdate)
	w.WriteHeader(http.StatusOK)
}