package models

import "time"

type Device struct {
	Id                   int       `sql:"AUTO_INCREMENT",json:"id"`
	Uuid                 string    `json:"uuid"`
	Name                 string    `json:"name"`
	Devicetype           string    `json:"devicetype"`
	AppId                int    	`json:"appid"`
	Isonline             bool      `json:"isonline"`
	Lastseen             time.Time `json:"last_seen"`
	PublicIP             string    `json:"public_ip"`
	IpAddress            string    `json:"ip_address"`
	Commit               string    `json:"commit"`
	ProvisioningState    string    `json:"provisioning_state"`
	ProvisioningProgress string    `json:"provisioning_progress"`
	DownloadProgress     string    `json:"download_progress"`
}

type DeviceReturn struct {
	Id         int    	`json:"Id,omitempty"`
	Name       string 	`json:"name"`
	Appid      int 	  	`json:"appid"`
	Uuid       string 	`json:"uuid"`
	Devicetype string 	`json:"devicetype"`
}

type DeviceState struct {
	AppId 		int 	`json:"appId"`
	DeviceId 	int 	`json:"deviceId"`
	State 		string  `json:"state"`
}
