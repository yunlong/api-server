package controllers

import (
	"encoding/json"
	"net/http"
	//"log"

	//"github.com/deviceMP/api-server/models"
	//"time"
)

type Connected struct {
	Connect 	bool 	`json:"connect"`
}

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	connected := Connected{Connect: true}
    json.NewEncoder(w).Encode(&connected)
}

func CheckAllDevice() {
	/*for {
		var devices []models.Device
		db.Find(&devices)

		for _, v := range devices {
			log.Println("--------->", v.IsOnline)
		}

		time.Sleep(10 * time.Second)
	}*/


}