package models

type Project struct {
	ID              int     	`sql:"AUTO_INCREMENT" json:"id"`
	OrgId           int      	`json:"orgId"`
	Name            string   	`json:"name"`
	Description     string   	`json:"description"`
	Image           string    	`json:"image"`
	DeviceType 	    string 	 	`json:"deviceType"`
	ApiKey     	    string 	 	`json:"apiKey"`
	Commit     	    string 	 	`json:"commit"`
	Repository 	    string 	 	`json:"repository"`
	Port 			string 		`json:"port"`
	DeviceTotal 	int 	 	`json:"deviceTotal"`
	Environment		[]ProjectEnv `json:"environment"`
}

type ProjectConfig struct {
	ProjectName 	string 	`json:"projectName"`
	ProjectId 		int 	`json:"projectId"`
	ApiKey 			string  `json:"apikey"`
	DeviceType 		string 	`json:"deviceType"`
}

type ProjectEnv struct {
	ID 			int 	`sql:"AUTO_INCREMENT" json:"id"`
    ProjectID  	int     `gorm:"index" json:"projectId"`
    Key 		string 	`json:"key"`
	Value 		string 	`json:"value"`
}