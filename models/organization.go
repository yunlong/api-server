package models

type Org struct {
	Id         	int    `sql:"AUTO_INCREMENT",json:"id"`
	Name       	string `json:"name"`
	Description 	string `json:"description"`
	Image 		string `json:"image"`
}