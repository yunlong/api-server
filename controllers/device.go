package controllers

import (
    "encoding/json"
    "net/http"
    "time"
    "log"
    "io"
    "io/ioutil"
    "strconv"

    "github.com/deviceMP/api-server/models"
    "github.com/gorilla/mux"
)

func ListDeviceByProject(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    projectId := vars["projectId"]

    var devices []models.Device

    w.Header().Set("Access-Control-Allow-Origin", "*")
    if projectIdInt, err := strconv.Atoi(projectId); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(&err)
    } else {
        w.WriteHeader(http.StatusOK)
        if !db.Where(models.Device{ProjectId: projectIdInt}).Find(&devices).RecordNotFound() {
            json.NewEncoder(w).Encode(&devices)
        }
    }
}

func GetDeviceById(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    projectId := vars["projectId"]
    deviceId := vars["deviceId"]

    var deviceDetail models.Device

    w.Header().Set("Access-Control-Allow-Origin", "*")
    if projectIdInt, err := strconv.Atoi(projectId); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(&err)
    } else {
        if DeviceIdInt, err := strconv.Atoi(deviceId); err != nil {
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(&err)
        } else {
            w.WriteHeader(http.StatusOK)
            if !db.Where(models.Device{ProjectId: projectIdInt, Id: DeviceIdInt}).First(&deviceDetail).RecordNotFound() {
                json.NewEncoder(w).Encode(&deviceDetail)
            }    
        }
    }
}

func RegisterDevice(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    projectId := vars["projectId"]

    if projectIdInt, err := strconv.Atoi(projectId); err != nil {
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

        var project models.Project
        if db.Where(models.Project{Id: projectIdInt}).First(&project).RecordNotFound() {
            w.WriteHeader(http.StatusNotFound)
            return
        } else {
            if clientApiKey[0] != project.ApiKey {
                w.WriteHeader(http.StatusForbidden)
                return
            } else {
                t := time.Now()
                deviceCreate := models.Device {
                    Uuid:       device.Uuid,
                    Name:       device.Name,
                    DeviceType: device.DeviceType,
                    ProjectId:  projectIdInt,
                    IsOnline:   true,
                    LastSeen:   t,
                    PublicIP:   device.PublicIP,
                    IpAddress:  device.IpAddress,
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

                    bodyResponse := models.DeviceReturn{Id: deviceCreate.Id, Name: deviceCreate.Name,ProjectId: deviceCreate.ProjectId,Uuid:deviceCreate.Uuid,DeviceType:deviceCreate.DeviceType}
                    if err := json.NewEncoder(w).Encode(&bodyResponse); err != nil {
                        panic(err)
                    }
                }
            }
        }
    }
}

func UpdateState(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    vars := mux.Vars(r)
    projectId := vars["projectId"]

    if projectIdInt, err := strconv.Atoi(projectId); err != nil {
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
        var project models.Project
        var deviceUpdate models.Device
        if db.Where(models.Project{Id: projectIdInt}).First(&project).RecordNotFound() {
            w.WriteHeader(http.StatusNotFound)
            return
        } else {
            if clientApiKey[0] != project.ApiKey {
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

//Update status online/offline of device
func UpdateStatus(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    vars := mux.Vars(r)
    projectId := vars["projectId"]

    if projectIdInt, err := strconv.Atoi(projectId); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(&err)
    } else {
        clientApiKey := r.URL.Query()["apikey"]
        if len(clientApiKey) == 0 {
            w.WriteHeader(http.StatusNotAcceptable)
            return
        }
        var deviceStatus models.DeviceStatus
        body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
        if err != nil {
            log.Println(err)
            return
        }
        if err := r.Body.Close(); err != nil {
            log.Println(err)
            return
        }
        if err := json.Unmarshal(body, &deviceStatus); err != nil {
            w.Header().Set("Content-Type", "application/json; charset=UTF-8")
            w.WriteHeader(422)
            if err := json.NewEncoder(w).Encode(err); err != nil {
                log.Println(err)
                return
            }
        }

        //Checking app and apikey valid with all request -> common func
        var project models.Project
        var deviceUpdate models.Device
        if db.Where(models.Project{Id: projectIdInt}).First(&project).RecordNotFound() {
            w.WriteHeader(http.StatusNotFound)
            return
        } else {
            if clientApiKey[0] != project.ApiKey {
                w.WriteHeader(http.StatusForbidden)
                return
            } else {
                if db.Where(models.Device{Id: deviceStatus.DeviceId}).First(&deviceUpdate).RecordNotFound() {
                    w.WriteHeader(http.StatusNotFound)
                    return
                } else {
                    deviceUpdate.IsOnline = deviceStatus.Status
                    db.Save(&deviceUpdate)
                    w.WriteHeader(http.StatusOK)
                    return
                }
            }
        }
    }
}

func UpdateProgress(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    projectId := vars["projectId"]

    if projectIdInt, err := strconv.Atoi(projectId); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(&err)
    } else {
        clientApiKey := r.URL.Query()["apikey"]
        if len(clientApiKey) == 0 {
            w.WriteHeader(http.StatusNotAcceptable)
            return
        }
        var deviceState models.DeviceProgress
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
        var project models.Project
        var deviceUpdate models.Device
        if db.Where(models.Project{Id: projectIdInt}).First(&project).RecordNotFound() {
            w.WriteHeader(http.StatusNotFound)
            return
        } else {
            if clientApiKey[0] != project.ApiKey {
                w.WriteHeader(http.StatusForbidden)
                return
            } else {
                if db.Where(models.Device{Id: deviceState.DeviceId}).First(&deviceUpdate).RecordNotFound() {
                    w.WriteHeader(http.StatusNotFound)
                    return
                } else {
                    db.Model(&deviceUpdate).Updates(models.Device{DownloadProgress: deviceState.Progress})
                    w.WriteHeader(http.StatusCreated)
                    return
                }
            }
        }
    }
}

func CheckUpdate(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    projectId := vars["projectId"]
    deviceId := vars["deviceId"]

    projectIdInt, err := strconv.Atoi(projectId); if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    deviceIdInt, err := strconv.Atoi(deviceId); if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    var deviceCheck models.Device
    db.Where(models.Device{ProjectId: projectIdInt, Id: deviceIdInt}).First(&deviceCheck)

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

//Params: OrgId & DeviceId
//QueryString?online=True
func DeviceOnline(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    var deviceOnline models.DeviceOnline
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        log.Println(err)
        return
    }
    if err := r.Body.Close(); err != nil {
        log.Println(err)
        return
    }

    if err := json.Unmarshal(body, &deviceOnline); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(422)
        if err := json.NewEncoder(w).Encode(err); err != nil {
            log.Println(err)
            return
        }
    }

    var deviceCheck models.Device
    rowUpdated := db.Model(&deviceCheck).Where(models.Device{ProjectId: deviceOnline.ProjectId, Id: deviceOnline.Id}).UpdateColumn(models.Device{IsOnline: deviceOnline.IsOnline}).RowsAffected

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    if rowUpdated > 0 {
        w.WriteHeader(http.StatusOK)
    } else {
        w.WriteHeader(http.StatusInternalServerError)
    }
}

func UpdateDeviceName(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    deviceId := vars["deviceId"]

    deviceIdInt, err := strconv.Atoi(deviceId); if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    w.Header().Set("Access-Control-Allow-Origin", "*")
    var deviceEdit models.DeviceEditName
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        log.Println(err)
        return
    }
    if err := r.Body.Close(); err != nil {
        log.Println(err)
        return
    }

    if err := json.Unmarshal(body, &deviceEdit); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(422)
        if err := json.NewEncoder(w).Encode(err); err != nil {
            log.Println(err)
            return
        }
    }

    var deviceUpdate models.Device
    rowUpdated := db.Model(&deviceUpdate).Where(models.Device{Id: deviceIdInt}).UpdateColumn(models.Device{Name: deviceEdit.Name}).RowsAffected

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    if rowUpdated > 0 {
        w.WriteHeader(http.StatusOK)
    } else {
        w.WriteHeader(http.StatusInternalServerError)
    }
}
