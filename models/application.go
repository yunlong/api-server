package models

type App struct {
	Uuid      		string 			`json:"uuid"`
	Commit      	string			`json:"commit,omitempty"`
	ContainerId 	string			`json:"containerId,omitempty"`
	Port 			string 			`json:"port,omitempty"`
	ImageId     	string			`json:"imageId,omitempty"`
	Latest 			string 			`json:"latest,omitempty"`
}

type AppUpdate struct {
	ImageId 	string 			`json:"imageId"`
	Port 		string 			`json:"port"`
	Environment	[]Environment   `json:"environments"`
}