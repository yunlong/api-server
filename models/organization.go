package models

type Org struct {
	Id         int    `sql:"AUTO_INCREMENT",json:"id"`
	Name       string `json:"name"`
	DeviceType string `json:"deviceType"`
	ApiKey     string `json:"apiKey"`
	Commit     string `json:"commit"`
	Repository string `json:"repository"`
}


type OrgConfig struct {
	ApplicationName 	string 	`json:"organization"`
	ApplicationId 		int 	`json:"orgId"`
	ApiKey 				string  `json:"apikey"`
	DeviceType 			string 	`json:"deviceType"`
}