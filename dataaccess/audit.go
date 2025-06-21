package dataaccess

import "time"

type logAction struct {
	Id          int       `json:"Id" gorm:"unique;primaryKey;autoIncrement"`
	Function    string    `gorm:"column:Function"`
	Requestbody string    `gorm:"column:Requestbody"`
	CreatedAt   time.Time `gorm:"column:CreatedAt"`
	RequestedBy string    `gorm:"column:RequestedBy"`
}
