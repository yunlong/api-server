package models

type App struct {
	Id 					int     	`sql:"AUTO_INCREMENT",json:"id"`
	Name 				string 		`json:"name"`
	Devicetype 			string 		`json:"devicetype"`
	Apikey 				string 		`json:"apikey"`
	Commit 				string 		`json:"commit"`	
	Repository 			string 		`json:"repository"`
}
