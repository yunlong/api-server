package controllers

import (
    "encoding/json"
    //"fmt"
    "net/http"
    "time"
    "log"
    "strconv"
    "io"
    "io/ioutil"

    "cli-client/models"
)

func RegisterDevice(w http.ResponseWriter, r *http.Request) {
    clientApiKey := r.URL.Query()["apikey"]
    if len(clientApiKey) == 0 {
        w.WriteHeader(http.StatusNotAcceptable)
        return
    }
    var device models.Device
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        log.Println(err)
        return
    }
    if err := r.Body.Close(); err != nil {
        log.Println(err)
        return
    }
    if err := json.Unmarshal(body, &device); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(422) // unprocessable entity
        if err := json.NewEncoder(w).Encode(err); err != nil {
            log.Println(err)
            return
        }
    }

    //Get apikey in application to compare with client request
    appIdInt, errc := strconv.Atoi(device.AppId)
    if errc != nil {
        appIdInt = 0
        log.Println(errc)
    }

    var application models.App
    if db.Where(models.App{Id: appIdInt}).First(&application).RecordNotFound() {
        w.WriteHeader(http.StatusNotFound)
        return
    } else {
        if clientApiKey[0] != application.Apikey {
            w.WriteHeader(http.StatusForbidden)
            return
        } else {
            t := time.Now()
            deviceCreate := models.Device {
                Uuid:       device.Uuid,
                Name:       device.Name,
                Devicetype: device.Devicetype,
                AppId:      device.AppId,
                Isonline:   true,
                Lastseen:   t,
                PublicIP:   device.PublicIP,
            }

            if dbc := db.Create(&deviceCreate); dbc.Error != nil {
                w.Header().Set("Content-Type", "application/json; charset=UTF-8")
                w.WriteHeader(http.StatusBadRequest)
                if err := json.NewEncoder(w).Encode(dbc.Error); err != nil {
                    panic(err)
                }
                return
            } else {
                w.Header().Set("Content-Type", "application/json; charset=UTF-8")
                w.WriteHeader(http.StatusCreated)

                bodyResponse := models.DeviceReturn{Id: deviceCreate.Id, Name: deviceCreate.Name,Appid: deviceCreate.AppId,Uuid:deviceCreate.Uuid,Devicetype:deviceCreate.Devicetype}
                if err := json.NewEncoder(w).Encode(&bodyResponse); err != nil {
                    panic(err)
                }
            }
        }
    }
}