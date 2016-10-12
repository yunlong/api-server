package controllers

import (
    "encoding/json"
    "net/http"
    "log"
    "io"
    "io/ioutil"
    "bufio"
    "os"
    "strconv"
    "fmt"

    "github.com/deviceMP/api-server/models"
    "github.com/deviceMP/api-server/utils"
    "github.com/gorilla/mux"
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
        db.Where(models.Org{Id: OrgIdInt}).First(&org)

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

func DownloadConfig(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    vars := mux.Vars(r)
    orgId := vars["orgId"]

    OrgIdInt, err := strconv.Atoi(orgId)
    if err != nil {
        OrgIdInt = 0
    }

    var org models.Org
    db.Where(models.Org{Id: OrgIdInt}).First(&org)
    orgConfig := models.OrgConfig{ApplicationId: org.Id, ApplicationName: org.Name, ApiKey: org.ApiKey, DeviceType: org.DeviceType}


    var fileName = "config.json"
    var filePath = "/tmp/" + fileName

    f, err := os.Create(filePath)
    handleError(err)

    b, err := json.Marshal(orgConfig)
    handleError(err)

    _, err = f.Write(b)
    handleError(err)
    f.Close()

    file, err := os.Open(filePath)
    handleError(err)

    r4 := bufio.NewReader(file)
    //copy the relevant headers. If you want to preserve the downloaded file name, extract it with go's url parser.
    w.Header().Set("Content-Disposition", "attachment; filename="+fileName+"")
    w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
    w.Header().Set("Content-Length", r.Header.Get("Content-Length"))
    //stream the body to the client without fully loading it into memory
    io.Copy(w, r4)
}

func handleError(e error) {
    if e != nil {
        fmt.Println(e)
    }
}