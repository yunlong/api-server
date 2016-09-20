package models

type App struct {
	AppId       int 		`json:"appId"`
	Commit      string		`json:"commit"`
	ContainerId string		`json:"containerId"`
	Env         string		`json:"env"`
	ImageId     string		`json:"imageId"`
}