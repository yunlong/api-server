package controllers

import (
	"encoding/json"
	"net/http"
	"io"
	"io/ioutil"

	"github.com/deviceMP/api-server/models"
	"github.com/gorilla/mux"
	"strconv"
)

func GetApp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	OrgId := vars["orgId"]

	var app models.App

	if OrgIdInt, err := strconv.Atoi(OrgId); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&err)
	} else {
		db.Where(models.App{AppId:OrgIdInt}).First(&app)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&app)
	}
}

// Api to update application when commit
// Case demo: Use this api to update apps and agent get app had updated and do update strategy on device.
// Only change commit value because if agent seen the commit had changed, it will be pull latest image from docker hub and do update.
// TODO: When user commit code, need to trigger and build, push image to docker hub done --> call api UpdateApp.
func UpdateApp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	OrgId := vars["orgId"]

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

		if OrgIdInt, err := strconv.Atoi(OrgId); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&err)
			return
		} else {
			var appUpdate models.App
			rowUpdated := db.Model(&appUpdate).Where(models.App{AppId: OrgIdInt}).UpdateColumn(models.App{Commit: app.Commit}).RowsAffected

			if rowUpdated > 0 {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
	}
}

