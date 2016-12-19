package models

type App struct {	
	ID 				int				`sql:"AUTO_INCREMENT" json:"id"`
	Uuid      		string 			`json:"uuid"`
	Commit      	string			`json:"commit"`
	ContainerId 	string			`json:"containerId"`
	ImageId     	string			`json:"imageId"`
	Port 			string 			`json:"port"`
	Latest 			string 			`json:"latest"`
}

type AppUpdate struct {
	ImageId 	string 			`json:"imageId"`
	Port 		string 			`json:"port"`
	Environment	[]Environment   `json:"environments"`
}

type DataTransfer struct {
	DeviceUuid 	string 		`json:"deviceUuid"`
	EnableLog	bool	 	`json:"enableLog"`
	Logs 		string		`json:"logs"`
	Action 		string 		`json:"action"`
}