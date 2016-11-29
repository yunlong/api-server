package controllers

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/deviceMP/api-server/models"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "127.0.0.1:8080", "http service address")
var upgrader = websocket.Upgrader{} // use default options

var deviceMap map[string]string
var deviceMapOffline map[string]bool
var deviceCurrentStatus map[string]bool

const (
	CheckUpdate 		= "CheckUpdate"
	RestartDeviceApp	= "RestartDeviceApp"
	InstallUpdate 		= "InstallUpdate"
	UpdateEnvironment 	= "UpdateEnvironment"
	RestartDevice 		= "RestartDevice"
	ShutdownDevice 		= "ShutdownDevice"
	StopApplication 	= "StopApplication"
	StartApplication 	= "StartApplication"
	RestartApplication 	= "RestartApplication"

	Message       = "Pong"
	StopCharacter = "\r\n\r\n"
)

func DeviceDBSyncUp() {
	var devices []models.Device
	db.Find(&devices)
	deviceMap = make(map[string]string)
	deviceMapOffline = make(map[string]bool)
	deviceCurrentStatus = make(map[string]bool)

	for _,v := range devices {
		if _, ok := deviceMap[v.Uuid]; !ok {
			var newDevice string = "None"
			
			deviceMap[v.Uuid] = newDevice
			deviceMapOffline[v.Uuid] = v.IsOnline
			deviceCurrentStatus[v.Uuid] = v.IsOnline
		}
	}
}

func ConnectivityListen(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		//deviceAction := deviceMap[string(message)]
		CheckDeviceOnline(string(message))
		err = c.WriteMessage(mt, []byte(deviceMap[string(message)]))
		if err != nil {
			log.Println("write:", err)
			break
		} else {
			deviceMap[string(message)] = "None"
		}
	}
}

func SyncDB() {
	//Synchronize with database
	go func() {
		for {
			DeviceDBSyncUp()
			time.Sleep(time.Second * 10)
		}
	}()

	go func() {
		for {
			SetDeviceOffline()
			time.Sleep(time.Second * 2)
		}
	}()
}

//Check package from client...
func CheckDeviceOnline(deviceUuid string) {
	//Check with deviceMapOffline -> 1.Get current state, 2. Compare update
	deviceMapOffline[deviceUuid] = true

	//Check neu device exist
	if dev, ok := deviceMapOffline[deviceUuid]; ok {
		if deviceCurrentStatus[deviceUuid] != dev {
			var device models.Device
			if !db.Where(models.Device{Uuid: deviceUuid}).First(&device).RecordNotFound() {
				log.Println("call to update on db")
				device.IsOnline = deviceMapOffline[deviceUuid]
				db.Save(&device)
				deviceCurrentStatus[deviceUuid] = deviceMapOffline[deviceUuid]
			}
		}	
	} else {
		deviceMapOffline[deviceUuid] = true
	}
}

//Interval check in 3s, If device IsOnline = false, update in db
func SetDeviceOffline() {
	for uuid, status := range deviceMapOffline {
		if deviceCurrentStatus[uuid] != status {
			var device models.Device
			if !db.Where(models.Device{Uuid: uuid}).First(&device).RecordNotFound() {
				log.Println("call to set offline")
				deviceCurrentStatus[uuid] = status
				deviceMapOffline[uuid] = status
    			device.IsOnline = status
    			db.Save(&device)
			}
		}
	}

	for k, _ := range deviceMapOffline {
		deviceMapOffline[k] = false
	}
}

func PushActionAgent(deviceUuid, action string) {
	deviceMap[deviceUuid] = action
}