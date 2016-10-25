package models

type Project struct {
	Id                   int     	 `sql:"AUTO_INCREMENT",json:"id"`
	OrgId                int      	 `json:"orgId"`
	Name                 string   	 `json:"name"`
	Description          string   	 `json:"description"`
	Image                string    	 `json:"image"`
	DeviceType 	     string 	 `json:"deviceType"`
	ApiKey     	     string 	 `json:"apiKey"`
	Commit     	     string 	 `json:"commit"`
	Repository 	     string 	 `json:"repository"`
	DeviceTotal 	     int 	 `json:"deviceTotal"`
}

type ProjectConfig struct {
	ProjectName 		string 	`json:"projectName"`
	ProjectId 		int 	`json:"projectId"`
	ApiKey 			string  `json:"apikey"`
	DeviceType 		string 	`json:"deviceType"`
}