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

	var appGet models.App
	if !db.Where(models.App{Uuid: deviceId}).First(&appGet).RecordNotFound() {
		json.NewEncoder(w).Encode(&appGet)
	}
}

// Api to update application when commit
// Case demo: Use this api to update apps and agent get app had updated and do update strategy on device.
// Only change commit value because if agent seen the commit had changed, it will be pull latest image from docker hub and do update.
// TODO: When user commit code, need to trigger and build, push image to docker hub done --> call api UpdateApp.
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
		var appUpdate models.App
		rowUpdated := db.Model(&appUpdate).Where(models.App{Uuid: app.Uuid}).UpdateColumn(models.App{Commit: app.Commit, ContainerId: app.ContainerId, Env: app.Env, ImageId: app.ImageId}).RowsAffected

		if rowUpdated > 0 {
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(&appUpdate); err != nil {
				panic(err)
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
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
			Uuid:             app.Uuid,
			Commit:      app.Commit,
			ContainerId: app.ContainerId,
			Env:         app.Env,
			ImageId:     app.ImageId,
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