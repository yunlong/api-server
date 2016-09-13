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

func Index(w http.ResponseWriter, r *http.Request) {
	var apps []models.App
	db.Find(&apps)

	json.NewEncoder(w).Encode(apps)
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
        w.WriteHeader(422) // unprocessable entity
        if err := json.NewEncoder(w).Encode(err); err != nil {
            panic(err)
            return
        }
    } else {
        appCreate := models.App{
            Name:       app.Name,
            Devicetype: app.Devicetype,
            Apikey:     utils.RandStringRunes(32),
            Commit:     "",
            Repository: app.Repository,
        }

        if app.Name == "" || app.Devicetype == "" || app.Repository == "" {
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        if dbc := db.Create(&appCreate); dbc.Error != nil {
            w.WriteHeader(http.StatusBadRequest)
            return
        } else {
            w.Header().Set("Content-Type", "application/json; charset=UTF-8")
            w.WriteHeader(http.StatusCreated)
            if err := json.NewEncoder(w).Encode(appCreate); err != nil {
                log.Println(err)
                return
            }
        }
    }



}

/*
func TodoIndex(w http.ResponseWriter, r *http.Request) {
	todos := Todos{
		Todo{Name: "Write presentation"},
		Todo{Name: "Host meetup"},
	}

	if err := json.NewEncoder(w).Encode(todos); err != nil {
		panic(err)
	}
}

func TodoShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoId := vars["todoId"]
	fmt.Fprintf(w, "Todo show: %s\n", todoId)
}
*/
