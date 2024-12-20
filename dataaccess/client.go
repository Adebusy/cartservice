package dataaccess

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

func ConnectClient(db *gorm.DB) IClient {
	return &DbConnect{db}
}

type IClient interface {
	GetClientByName(clientName string) TblClient
	RegisterNewClient(req TblClient) string
}

type TblClient struct {
	Id          int        `json:"Id" validate:"omitempty"`
	Name        string     `json:"Name" validate:"omitempty"`
	Status      int        `json:"Status" validate:"omitempty"`
	Description string     `json:"Description" validate:"omitempty"`
	DateAdded   *time.Time `json:"DateAdded"`
}

type ClientRequest struct {
	Name        string `json:"Name" validate:"omitempty"`
	Description string `json:"Description" validate:"omitempty"`
}

type ClientResp struct {
	Id        int    `json:"Id" validate:"omitempty"`
	Name      string `json:"Name" validate:"omitempty"`
	RespToken string `json:"respToken" validate:"omitempty"`
}

func (cn DbConnect) GetClientByName(clientName string) TblClient {
	res := TblClient{}
	fmt.Println("dasdsas")
	fmt.Printf("name is clientName %s", clientName)
	cn.DbGorm.Table("TblClient").Select("Id", "Name", "Status", "Description", "DateAdded").Where("\"Name\"=? and \"Status\"=1", clientName).First(&res)
	fmt.Printf("name is %s", res.Name)
	return res
}

func (cn DbConnect) RegisterNewClient(req TblClient) string {
	if doinssert := cn.DbGorm.Table("TblClient").Create(&req).Error; doinssert != nil {
		return "01"
	} else {
		return "00"
	}

}
