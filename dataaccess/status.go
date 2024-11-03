package dataaccess

import "time"

type TblStatus struct {
	Id         int       `json:"Id" gorm:"unique;primaryKey;autoIncrement"`
	StatusName string    `json:"StatusName" validate:"omitempty"`
	CreatedAt  time.Time `json:"CreatedAt" validate:"omitempty"`
}

type StatusResp struct {
	Id         int       `json:"Id" gorm:"unique;primaryKey;autoIncrement"`
	StatusName string    `json:"StatusName" validate:"omitempty"`
	CreatedAt  time.Time `json:"CreatedAt" validate:"omitempty"`
}

func (cn DbConnect) GetAllStatus() []TblStatus {
	stat := []TblStatus{}
	cn.DbGorm.Select("Id", "StatusName", "CreatedAt").Find(&stat)
	return stat
}
