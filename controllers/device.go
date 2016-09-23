package controllers

import (
    "encoding/json"
    "net/http"
    "time"
    "log"
    "io"
    "io/ioutil"

    "github.com/deviceMP/api-server/models"
    "github.com/gorilla/mux"
    "strconv"
)

func ListDeviceByApp(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    orgId := vars["orgId"]

    var devices []models.Device

    w.Header().Set("Access-Control-Allow-Origin", "*")
    if OrgIdInt, err := strconv.Atoi(orgId); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(&err)
    } else {
        db.Where(models.Device{AppId: OrgIdInt}).Find(&devices)
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(&devices)
    }
}

func RegisterDevice(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    orgId := vars["orgId"]

    if OrgIdInt, err := strconv.Atoi(orgId); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(&err)
    } else {
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
            w.WriteHeader(422)
            if err := json.NewEncoder(w).Encode(err); err != nil {
                log.Println(err)
                return
            }
        }

        var org models.Org
        if db.Where(models.Org{Id: OrgIdInt}).First(&org).RecordNotFound() {
            w.WriteHeader(http.StatusNotFound)
            return
        } else {
            if clientApiKey[0] != org.ApiKey {
                w.WriteHeader(http.StatusForbidden)
                return
            } else {
                t := time.Now()
                deviceCreate := models.Device {
                    Uuid:       device.Uuid,
                    Name:       device.Name,
                    DeviceType: device.DeviceType,
                    AppId:      OrgIdInt,
                    IsOnline:   true,
                    LastSeen:   t,
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

                    bodyResponse := models.DeviceReturn{Id: deviceCreate.Id, Name: deviceCreate.Name,AppId: deviceCreate.AppId,Uuid:deviceCreate.Uuid,DeviceType:deviceCreate.DeviceType}
                    if err := json.NewEncoder(w).Encode(&bodyResponse); err != nil {
                        panic(err)
                    }
                }
            }
        }
    }
}

func UpdateState(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    orgId := vars["orgId"]

    if OrgIdInt, err := strconv.Atoi(orgId); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(&err)
    } else {
        clientApiKey := r.URL.Query()["apikey"]
        if len(clientApiKey) == 0 {
            w.WriteHeader(http.StatusNotAcceptable)
            return
        }
        var deviceState models.DeviceState
        body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
        if err != nil {
            log.Println(err)
            return
        }
        if err := r.Body.Close(); err != nil {
            log.Println(err)
            return
        }
        if err := json.Unmarshal(body, &deviceState); err != nil {
            w.Header().Set("Content-Type", "application/json; charset=UTF-8")
            w.WriteHeader(422)
            if err := json.NewEncoder(w).Encode(err); err != nil {
                log.Println(err)
                return
            }
        }

        //Checking app and apikey valid with all request -> common func
        var org models.Org
        var deviceUpdate models.Device
        if db.Where(models.Org{Id: OrgIdInt}).First(&org).RecordNotFound() {
            w.WriteHeader(http.StatusNotFound)
            return
        } else {
            if clientApiKey[0] != org.ApiKey {
                w.WriteHeader(http.StatusForbidden)
                return
            } else {
                if db.Where(models.Device{Id: deviceState.DeviceId}).First(&deviceUpdate).RecordNotFound() {
                    w.WriteHeader(http.StatusNotFound)
                    return
                } else {
                    db.Model(&deviceUpdate).Updates(models.Device{Status: deviceState.State})
                    w.WriteHeader(http.StatusCreated)
                    return
                }
            }
        }
    }
}

func CheckUpdate(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    orgId := vars["orgId"]
    deviceId := vars["deviceId"]

    orgIdInt, err := strconv.Atoi(orgId); if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    deviceIdInt, err := strconv.Atoi(deviceId); if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    var deviceCheck models.Device
    db.Where(models.Device{AppId: orgIdInt, Id: deviceIdInt}).First(&deviceCheck)

    result := map[string]interface{}{
        "deviceId":  deviceCheck.Id,
        "status": deviceCheck.Status,
    }

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusCreated)

    if err := json.NewEncoder(w).Encode(&result); err != nil {
        panic(err)
    }

}