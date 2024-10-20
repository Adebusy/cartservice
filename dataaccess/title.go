package dataaccess

import (
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ConnectTitle struct {
	DbGorm *gorm.DB
}

func ConTitle(db *gorm.DB) ITitle {
	return &ConnectTitle{db}
}

type TblTitle struct {
	Id        int       `json:"Id" gorm:"unique;primaryKey;autoIncrement"`
	Name      string    `json:"Name" validate:"omitempty"`
	Status    bool      `json:"Status" validate:"omitempty"`
	CreatedAt time.Time `json:"CreatedAt" validate:"omitempty"`
}

type TitleResp struct {
	Id        int       `json:"Id" gorm:"unique;primaryKey;autoIncrement"`
	Name      string    `json:"Name" validate:"omitempty"`
	Status    bool      `json:"Status" validate:"omitempty"`
	CreatedAt time.Time `json:"CreatedAt" validate:"omitempty"`
}

// type TblCart struct {
// 	Id        int       `json:"Id" gorm:"unique;primaryKey;autoIncrement"`
// 	Name      string    `json:"Name" validate:"omitempty"`
// 	Status    bool      `json:"Status" validate:"omitempty"`
// 	CreatedAt time.Time `json:"CreatedAt" validate:"omitempty"`
// }

type ITitle interface {
	CreateTitle(prod TblTitle) int
	GetTitleByTitleId(TitleId int) TblTitle
	GetTitleByTitleName(TitleId string) TblTitle
	DeleteTitleByTitleId(titleId int) error
	GetTitles() []TblTitle
}

func (cn ConnectTitle) CreateTitle(prod TblTitle) int {

	if doInsert := cn.DbGorm.Table("TblTitle").Create(&prod).Error; doInsert == nil {
		return prod.Id
	} else {
		logrus.Error(doInsert)
		return 0
	}
}

func (cn ConnectTitle) GetTitleByTitleId(TitleId int) TblTitle {
	title := TblTitle{}
	cn.DbGorm.Table("TblTitle").Select("Id", "Name", "Status", "CreatedAt").Where("\"Id\"=?", TitleId).First(&title)
	return title
}

func (cn ConnectTitle) GetTitleByTitleName(TitleId string) TblTitle {
	title := TblTitle{}
	cn.DbGorm.Table("TblTitle").Select("Id", "Name", "Status", "CreatedAt").Where("\"Name\"=?", TitleId).First(&title)
	return title
}

func (cn ConnectTitle) DeleteTitleByTitleId(titleId int) error {
	doDelete := cn.DbGorm.Delete(&TblTitle{}, titleId).Error
	return doDelete
}
func (cn ConnectTitle) GetTitles() []TblTitle {
	req := []TblTitle{}
	cn.DbGorm.Select("Id", "Name", "Status").Find(&req)
	return req
}
