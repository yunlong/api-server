package models

type App struct {
	AppId       int 		`json:"appId"`
	Commit      string		`json:"commit,omitempty"`
	ContainerId string		`json:"containerId,omitempty"`
	Env         string		`json:"env,omitempty"`
	ImageId     string		`json:"imageId,omitempty"`
}