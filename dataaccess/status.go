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
	Status      int       `json:"Status" validate:"omitempty"`
	Description string    `json:"Description"`
	UserId      int       `json:"UserId"`
	RoleId      int       `json:"RoleId"`
	GroupTypeId int       `json:"GroupTypeId"`
	CartId      int       `json:"CartId"`
	DateAdded   time.Time `json:"DateAdded"`
}

type TblGroupObj struct {
	GroupName   string `json:"GroupName"`
	Description string `json:"Description"`
	UserId      int    `json:"UserId"`
	// RoleId      int    `json:"RoleId"`
	GroupTypeId int `json:"GroupTypeId"`
	CartId      int `json:"CartId"`
}

type TblTeamGroupObj struct {
	GroupName   string `json:"GroupName"`
	Description string `json:"Description"`
	UserId      int    `json:"UserId"`
	GroupTypeId int    `json:"GroupTypeId"`
	CartId      int    `json:"CartId"`
	AdminId     int    `json:"AdminId"`
}

type RmoveUserFromGroupObj struct {
	GroupName   string `json:"GroupName"`
	AdminId     int    `json:"AdminId"`
	GroupTypeId int    `json:"GroupTypeId"`
	CartId      int    `json:"CartId"`
	UserId      int    `json:"UserId"`
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
