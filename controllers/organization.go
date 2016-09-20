package controllers

import (
    "encoding/json"
    "net/http"
    "log"
    "io"
    "io/ioutil"

    "github.com/deviceMP/api-server/models"
    "github.com/deviceMP/api-server/utils"
)

func ListOrg(w http.ResponseWriter, r *http.Request) {
	var orgs []models.Org
	db.Find(&orgs)

        w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&orgs)
}

func CreateOrg(w http.ResponseWriter, r *http.Request) {
    var org models.Org
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        panic(err)
        return
    }
    if err := r.Body.Close(); err != nil {
        panic(err)
        return
    }
    if err := json.Unmarshal(body, &org); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(422)
        if err := json.NewEncoder(w).Encode(err); err != nil {
            panic(err)
            return
        }
    } else {
        if org.Name == "" || org.DeviceType == "" || org.Repository == "" {
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        orgCreate := models.Org{
            Name:       org.Name,
            DeviceType: org.DeviceType,
            ApiKey:     utils.RandStringRunes(32),
            Commit:     "",
            Repository: org.Repository,
        }

        if dbc := db.Create(&orgCreate); dbc.Error != nil {
            w.WriteHeader(http.StatusBadRequest)
            return
        } else {
            w.Header().Set("Content-Type", "application/json; charset=UTF-8")
            w.WriteHeader(http.StatusCreated)
            if err := json.NewEncoder(w).Encode(&orgCreate); err != nil {
                log.Println(err)
                return
            }
        }
    }
}
