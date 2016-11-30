package models

import "time"

type Log struct {
	ID                   int	`sql:"AUTO_INCREMENT" json:"id"`
	DeviceId 	     	 int	`json:"deviceId"`
	Status               string	`json:"status"`
	Name                 string	`json:"name"`
	Version              string	`json:"version"`
	CreateAt             time.Time	`json:"createAt"`
	CompleteAt           time.Time	`json:"completeAt"`
	Link                 string	`json:"link"`
}