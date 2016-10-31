package models

type App struct {
	Uuid      		string 		`json:"uuid"`
	Commit      	string		`json:"commit,omitempty"`
	ContainerId 	string		`json:"containerId,omitempty"`
	Env         	string		`json:"env,omitempty"`
	ImageId     	string		`json:"imageId,omitempty"`
	Latest 			string 		`json:"latest,omitempty"`
}