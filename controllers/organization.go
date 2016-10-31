package controllers

import (
    "encoding/json"
    "net/http"
    "log"
    "io/ioutil"
    "strconv"
    "fmt"

    "github.com/deviceMP/api-server/models"
    //"github.com/deviceMP/api-server/utils"
    "github.com/gorilla/mux"
    "io"
    "github.com/deviceMP/api-server/utils"
    "os"
)

func ListOrg(w http.ResponseWriter, r *http.Request) {
    var orgs []models.Org
    db.Find(&orgs)

    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(&orgs)
}

func GetOrg(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    orgId := vars["orgId"]

    w.Header().Set("Access-Control-Allow-Origin", "*")
    if OrgIdInt, err := strconv.Atoi(orgId); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(&err)
    } else {
        var org models.Org
        db.Where(models.Org{ID: OrgIdInt}).First(&org)

        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(&org)
    }
}

func CreateOrg(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")

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
        w.WriteHeader(422)
        if err := json.NewEncoder(w).Encode(err); err != nil {
            panic(err)
            return
        }
    } else {
        if org.Name == "" {
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        orgCreate := models.Org{
            Name:           org.Name,
            Description:    org.Description,
            Image:          org.Image,
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

func UploadOrgImage(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")

    //log.Println(r.Body)
    file, handler, err := r.FormFile("file")
    if err != nil {
        fmt.Println(err)
    }
    if file != nil {
        data, err := ioutil.ReadAll(file)
        if err != nil {
            fmt.Println(err)
        }

        reName := utils.RenameImage(handler.Filename, 6)
        imageFolderPath := os.Getenv("HOME") + "/" + "api-images"
        if _, err := os.Stat(imageFolderPath); os.IsNotExist(err) {
            os.Mkdir(imageFolderPath, 0777)
        }

        err = ioutil.WriteFile(imageFolderPath + "/" + reName, data, 0777)
        if err != nil {
            fmt.Println(err)
        }
        reData := map[string]interface{}{
            "filename": reName,
        }

        w.WriteHeader(http.StatusOK)
        if err := json.NewEncoder(w).Encode(&reData); err != nil {
            log.Println(err)
            return
        }
    }
}

func DeleteOrg(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    vars := mux.Vars(r)
    orgId := vars["orgId"]

    if OrgIdInt, err := strconv.Atoi(orgId); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(&err)
    } else {
        var org models.Org
        var devices []models.Device
        db.Where(models.Org{ID: OrgIdInt}).First(&org)
        db.Where(models.Device{ProjectId: org.ID}).Find(&devices)
        for _,v := range devices {
            var apps []models.App
            db.Where(models.App{Uuid: v.Uuid}).Find(&apps)
            for _,v := range apps {
                db.Delete(&v)
            }
            db.Delete(&v)
        }

        db.Delete(&org)
        w.WriteHeader(http.StatusOK)
    }
}

func handleError(e error) {
    if e != nil {
        fmt.Println(e)
    }
}