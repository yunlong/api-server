package models

import "time"

type Device struct {
	Id                   int       `sql:"AUTO_INCREMENT",json:"id"`
	Uuid                 string    `json:"uuid"`
	Name                 string    `json:"name"`
	DeviceType           string    `json:"deviceType"`
	AppId                int       `json:"appId"`
	IsOnline             bool      `json:"isOnline"`
	LastSeen             time.Time `json:"lastSeen"`
	PublicIP             string    `json:"publicIp"`
	IpAddress            string    `json:"ipAddress"`
	Commit               string    `json:"commit"`
	Status 		     	 string    `json:"status"`
	ProvisioningState    string    `json:"provisioningState,omitempty"`
	ProvisioningProgress int       `json:"provisioningProgress,omitempty"`
	DownloadProgress     int       `json:"downloadProgress,omitempty"`
	UpdatePending 	     bool      `json:"updatePending,omitempty"`
	UpdateDownloaded     bool      `json:"updateDownloaded,omitempty"`
	UpdateFailed 	     bool      `json:"updateFailed,omitempty"`
}

type DeviceReturn struct {
	Id         int    	`json:"id,omitempty"`
	Name       string 	`json:"name"`
	AppId      int 	  	`json:"appId"`
	Uuid       string 	`json:"uuid"`
	DeviceType string 	`json:"deviceType"`
}

type DeviceState struct {
	AppId 		int 	`json:"appId"`
	DeviceId 	int 	`json:"deviceId"`
	State 		string  `json:"state"`
}

type DeviceStatus struct {
	AppId 		int 	`json:"appId"`
	DeviceId 	int 	`json:"deviceId"`
	Status 		bool  	`json:"status"`
}

type DeviceProgress struct {
	AppId 		int 	`json:"appId"`
	DeviceId 	int 	`json:"deviceId"`
	Progress 	int 	`json:"progress"`
}

type DeviceOnline struct {
	Id 			int 	`json:"id"`
	AppId 		int 	`json:"appId"`
	IsOnline 	bool  	`json:"isOnline"`
}

type DeviceEditName struct {
	Name 	string 		`json:"name"`
}
