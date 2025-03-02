package dataaccess

import "time"

type TblStatus struct {
	Id         int       `json:"Id" gorm:"unique;primaryKey;autoIncrement"`
	StatusName string    `json:"StatusName" validate:"omitempty"`
	CreatedAt  time.Time `json:"CreatedAt" validate:"omitempty"`
}

type TblRole struct {
	Id        int       `json:"Id" gorm:"unique;primaryKey;autoIncrement"`
	RoleName  string    `json:"RoleName" validate:"omitempty"`
	Status    bool      `json:"Status" validate:"omitempty"`
	CreatedAt time.Time `json:"CreatedAt" validate:"omitempty"`
}

type TblGroupType struct {
	Id            int       `json:"Id" gorm:"unique;primaryKey;autoIncrement"`
	GroupTypeName string    `json:"GroupTypeName"`
	Status        bool      `json:"Status" validate:"omitempty"`
	DateAdded     time.Time `json:"DateAdded"`
}

type TblGroupUser struct {
	Id          int       `json:"Id" gorm:"unique;primaryKey;autoIncrement"`
	GroupName   string    `json:"GroupName"`
	Description string    `json:"Description"`
	UserId      int       `json:"UserId"`
	RoleId      int       `json:"RoleId"`
	Status      bool      `json:"Status" validate:"omitempty"`
	GroupTypeId int       `json:"GroupTypeId"`
	DateAdded   time.Time `json:"DateAdded"`
}

type TblOrderItem struct {
	Id              int       `json:"Id" gorm:"unique;primaryKey;autoIncrement"`
	OrderId         int       `json:"OrderId"`
	ProductId       int       `json:"ProductId"`
	Quantity        int       `json:"Quantity"`
	PriceAtPurchase int       `json:"PriceAtPurchase"`
	DateAdded       time.Time `json:"DateAdded"`
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
