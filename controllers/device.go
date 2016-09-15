package controllers

import (
    "encoding/json"
    "net/http"
    "time"
    "log"
    "io"
    "io/ioutil"

    "github.com/deviceMP/api-server/models"
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

    var application models.App
    if db.Where(models.App{Id: device.AppId}).First(&application).RecordNotFound() {
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

func UpdateState(w http.ResponseWriter, r *http.Request) {
    clientApiKey := r.URL.Query()["apikey"]
    if len(clientApiKey) == 0 {
        w.WriteHeader(http.StatusNotAcceptable)
        return
    }
    var devicestate models.DeviceState
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        log.Println(err)
        return
    }
    if err := r.Body.Close(); err != nil {
        log.Println(err)
        return
    }
    if err := json.Unmarshal(body, &devicestate); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(422) // unprocessable entity
        if err := json.NewEncoder(w).Encode(err); err != nil {
            log.Println(err)
            return
        }
    }
    
    //Checking app and apikey valid with all request -> common func
    var application models.App
    var deviceUpdate models.Device

    if db.Where(models.App{Id: devicestate.AppId}).First(&application).RecordNotFound() {
        w.WriteHeader(http.StatusNotFound)
        return
    } else {
        if clientApiKey[0] != application.Apikey {
            w.WriteHeader(http.StatusForbidden)
            return
        } else {
            if db.Where(models.Device{Id: devicestate.DeviceId}).First(&deviceUpdate).RecordNotFound() {
                w.WriteHeader(http.StatusNotFound)
                return
            } else {
                db.Model(&deviceUpdate).Updates(models.Device{ProvisioningState: devicestate.State})
                w.WriteHeader(http.StatusCreated)
                return
            }
        }
    }
}