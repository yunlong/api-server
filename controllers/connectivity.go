package controllers

import (
	"encoding/json"
	"net/http"
	"log"

	"github.com/deviceMP/api-server/models"
	"time"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(true); err != nil {
		log.Println(err)
		return
	}
}

func CheckAllDevice() {
	for {
		var devices []models.Device
		db.Find(&devices)

		for _, v := range devices {
			log.Println("--------->", v.Isonline)
		}

		time.Sleep(10 * time.Second)
	}


}