package models

type Org struct {
	ID         	int    `sql:"AUTO_INCREMENT" json:"id"`
	Name       	string `json:"name"`
	Description string `json:"description"`
	Image 		string `json:"image"`
}