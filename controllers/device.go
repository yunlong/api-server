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
    ProjectId := vars["projectId"]

    var devices []models.Device

    w.Header().Set("Access-Control-Allow-Origin", "*")
    if ProjectIdInt, err := strconv.Atoi(ProjectId); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(&err)
    } else {
        w.WriteHeader(http.StatusOK)
        if !db.Where(models.Device{ProjectId: ProjectIdInt}).Find(&devices).RecordNotFound() {
            json.NewEncoder(w).Encode(&devices)
        }
    }
}

func GetDeviceById(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    projectId := vars["projectId"]
    deviceID := vars["deviceId"]

    var deviceDetail models.Device

    w.Header().Set("Access-Control-Allow-Origin", "*")
    if projectIdInt, err := strconv.Atoi(projectId); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(&err)
    } else {
        if deviceIDInt, err := strconv.Atoi(deviceID); err != nil {
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(&err)
        } else {
            w.WriteHeader(http.StatusOK)
            if !db.Where(models.Device{ProjectId: projectIdInt, ID: deviceIDInt}).First(&deviceDetail).RecordNotFound() {
                json.NewEncoder(w).Encode(&deviceDetail)
            }    
        }
    }
}

func RegisterDevice(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")

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
    if db.Where(models.Project{ID: device.ProjectId}).First(&project).RecordNotFound() {
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
                ProjectId:  device.ProjectId,
                IsOnline:   true,
                LastSeen:   t,
                Commit:     project.Commit,
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

                //bodyResponse := models.DeviceReturn{Id: deviceCreate.ID, Name: deviceCreate.Name,ProjectId: deviceCreate.ProjectId,Uuid:deviceCreate.Uuid,DeviceType:deviceCreate.DeviceType}
                var projectEnv []models.ProjectEnv
                db.Where(models.ProjectEnv{ProjectID: project.ID}).Find(&projectEnv)
                var penvs []models.Environment
                for _, v := range projectEnv {
                    var env models.Environment
                    env.Name = v.Key
                    env.Value = v.Value
                    penvs = append(penvs, env)
                }
                bodyResponse := models.RegisterSuccess{DeviceId: deviceCreate.ID, Image: project.Repository, Port: project.Port, Privileged: project.Privileged, Environments: penvs, Commit: project.Commit}
                if err := json.NewEncoder(w).Encode(&bodyResponse); err != nil {
                    panic(err)
                }
            }
        }
    }
}

func UpdateState(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")

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

    //Checking app and apikey valID with all request -> common func
    var project models.Project
    var deviceUpdate models.Device
    if db.Where(models.Project{ID: deviceState.ProjectId}).First(&project).RecordNotFound() {
        w.WriteHeader(http.StatusNotFound)
        return
    } else {
        if clientApiKey[0] != project.ApiKey {
            w.WriteHeader(http.StatusForbidden)
            return
        } else {
            if db.Where(models.Device{ID: deviceState.DeviceId}).First(&deviceUpdate).RecordNotFound() {
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

//Update status online/offline of device
func UpdateStatus(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")

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

    //Checking app and apikey valID with all request -> common func
    var project models.Project
    var deviceUpdate models.Device
    if db.Where(models.Project{ID: deviceStatus.ProjectId}).First(&project).RecordNotFound() {
        w.WriteHeader(http.StatusNotFound)
        return
    } else {
        if clientApiKey[0] != project.ApiKey {
            w.WriteHeader(http.StatusForbidden)
            return
        } else {
            if db.Where(models.Device{ID: deviceStatus.DeviceId}).First(&deviceUpdate).RecordNotFound() {
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

func UpdateProgress(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")

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
    //log.Println(deviceState)
    //Checking app and apikey valID with all request -> common func
    var project models.Project
    var deviceUpdate models.Device
    if db.Where(models.Project{ID: deviceState.ProjectId}).First(&project).RecordNotFound() {
        w.WriteHeader(http.StatusNotFound)
        return
    } else {
        if clientApiKey[0] != project.ApiKey {
            w.WriteHeader(http.StatusForbidden)
            return
        } else {
            if db.Where(models.Device{ID: deviceState.DeviceId}).First(&deviceUpdate).RecordNotFound() {
                w.WriteHeader(http.StatusNotFound)
                return
            } else {
                deviceUpdate.DownloadProgress = deviceState.Progress
                db.Save(&deviceUpdate)
                w.WriteHeader(http.StatusOK)
                return
            }
        }
    }
}

func CheckAppUpdate(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    vars := mux.Vars(r)
    ProjectId := vars["projectId"]
    deviceID := vars["deviceId"]

    ProjectIdInt, err := strconv.Atoi(ProjectId); if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    deviceIDInt, err := strconv.Atoi(deviceID); if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    var deviceCheck models.Device
    db.Where(models.Device{ProjectId: ProjectIdInt, ID: deviceIDInt}).First(&deviceCheck)

    log.Println(deviceMap[deviceCheck.Uuid])
    /*result := map[string]interface{}{
        "deviceID":  deviceCheck.ID,
        "status": deviceCheck.Status,
    }

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusCreated)

    if err := json.NewEncoder(w).Encode(&result); err != nil {
        panic(err)
    }*/
}

//Params: OrgID & DeviceID
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
    rowUpdated := db.Model(&deviceCheck).Where(models.Device{ProjectId: deviceOnline.ProjectId, ID: deviceOnline.Id}).UpdateColumn(models.Device{IsOnline: deviceOnline.IsOnline}).RowsAffected

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    if rowUpdated > 0 {
        w.WriteHeader(http.StatusOK)
    } else {
        w.WriteHeader(http.StatusInternalServerError)
    }
}

func UpdateDeviceName(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    vars := mux.Vars(r)
    deviceID := vars["deviceId"]

    deviceIDInt, err := strconv.Atoi(deviceID); if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

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
    rowUpdated := db.Model(&deviceUpdate).Where(models.Device{ID: deviceIDInt}).UpdateColumn(models.Device{Name: deviceEdit.Name}).RowsAffected

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    if rowUpdated > 0 {
        w.WriteHeader(http.StatusOK)
    } else {
        w.WriteHeader(http.StatusInternalServerError)
    }
}

func UpdateDeviceVersion(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")

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

    var project models.Project
    var deviceUpdate models.Device
    if db.Where(models.Project{ID: deviceState.ProjectId}).First(&project).RecordNotFound() {
        w.WriteHeader(http.StatusNotFound)
        return
    } else {
        if clientApiKey[0] != project.ApiKey {
            w.WriteHeader(http.StatusForbidden)
            return
        } else {
            if db.Where(models.Device{ID: deviceState.DeviceId}).First(&deviceUpdate).RecordNotFound() {
                w.WriteHeader(http.StatusNotFound)
                return
            } else {
                deviceUpdate.Commit = deviceState.State
                db.Save(&deviceUpdate)
                w.WriteHeader(http.StatusOK)
                if err := json.NewEncoder(w).Encode(http.StatusOK); err != nil {
                    log.Println(err)
                }
                return
            }
        }
    }
}

func UpdateLatestVersion(w http.ResponseWriter, r *http.Request) {
    log.Println("FUCKKKKK")
    w.Header().Set("Access-Control-Allow-Origin", "*")

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

    var project models.Project
    var deviceUpdate models.Device
    if db.Where(models.Project{ID: deviceState.ProjectId}).First(&project).RecordNotFound() {
        w.WriteHeader(http.StatusNotFound)
        return
    } else {
        if clientApiKey[0] != project.ApiKey {
            w.WriteHeader(http.StatusForbidden)
            return
        } else {
            if db.Where(models.Device{ID: deviceState.DeviceId}).First(&deviceUpdate).RecordNotFound() {
                w.WriteHeader(http.StatusNotFound)
                return
            } else {
                var appUpdate models.App
                rowUpdated := db.Model(&appUpdate).Where(models.App{Uuid: deviceUpdate.Uuid}).UpdateColumn(models.App{Latest: deviceState.State}).RowsAffected

                w.Header().Set("Content-Type", "application/json; charset=UTF-8")
                if rowUpdated > 0 {
                    w.WriteHeader(http.StatusOK)
                } else {
                    w.WriteHeader(http.StatusInternalServerError)
                }
                if err := json.NewEncoder(w).Encode(http.StatusOK); err != nil {
                    log.Println(err)
                }
                return
            }
        }
    }
}