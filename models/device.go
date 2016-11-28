package models

import "time"

type Device struct {
	ID                   int       `sql:"AUTO_INCREMENT" json:"id"`
	Uuid                 string    `json:"uuid"`
	Name                 string    `json:"name"`
	DeviceType           string    `json:"deviceType"`
	ProjectId            int       `json:"projectId"`
	IsOnline             bool      `json:"isOnline"`
	LastSeen             time.Time `json:"lastSeen"`
	PublicIP             string    `json:"publicIp"`
	IpAddress            string    `json:"ipAddress"`
	Commit               string    `json:"commit"`
	Status 		     	 string    `json:"status"`
	ProvisioningState    string    `json:"provisioningState,omitempty"`
	ProvisioningProgress int       `json:"provisioningProgress"`
	DownloadProgress     int       `json:"downloadProgress"`
	UpdatePending 	     bool      `json:"updatePending,omitempty"`
	UpdateDownloaded     bool      `json:"updateDownloaded,omitempty"`
	UpdateFailed 	     bool      `json:"updateFailed,omitempty"`
	Environment			 []DeviceEnv `json:"environment"`
}

type DeviceReturn struct {
	Id         int    	`json:"id,omitempty"`
	Name       string 	`json:"name"`
	ProjectId  int 	  	`json:"projectId"`
	Uuid       string 	`json:"uuid"`
	DeviceType string 	`json:"deviceType"`
}

type DeviceState struct {
	ProjectId 	int 	`json:"projectId"`
	DeviceId 	int 	`json:"deviceId"`
	State 		string  `json:"state"`
}

type DeviceStatus struct {
	ProjectId 	int 	`json:"projectId"`
	DeviceId 	int 	`json:"deviceId"`
	Status 		bool  	`json:"status"`
}

type DeviceProgress struct {
	ProjectId 	int 	`json:"projectId"`
	DeviceId 	int 	`json:"deviceId"`
	Progress 	int 	`json:"progress"`
}

type DeviceOnline struct {
	Id 			int 	`json:"id"`
	ProjectId 	int 	`json:"projectId"`
	IsOnline 	bool  	`json:"isOnline"`
}

type DeviceEditName struct {
	Name 	string 		`json:"name"`
}

type DeviceEnv struct {
	ID 			int 	`sql:"AUTO_INCREMENT" json:"id"`
    DeviceID  	int     `gorm:"index" json:"deviceId"`
    Key 		string 	`json:"key"`
	Value 		string 	`json:"value"`
}

type RegisterSuccess struct {
	DeviceId 	int		`json:"deviceId"`
	Image 		string	`json:"image"`
	Port 		string 	`json:"port"`
	Privileged  bool 	`json:"privileged"`
	Environments	[]Environment `json:"environments"`
	RegisterAt 	int		`json:"registerAt"`
	Commit 		string 	`json:"commit"`
}

type Environment struct {
	Name   	string `json:"name"`
	Value   string `json:"value"`
}