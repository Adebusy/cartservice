package dataaccess

import (
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PConnect struct {
	DbGorm *gorm.DB
}

func ConnectProduct(db *gorm.DB) IProduct {
	return &PConnect{db}
}

type TblProduct struct {
	Id                int       `json:"Id" gorm:"unique;primaryKey;autoIncrement"`
	ProductName       string    `json:"ProductName" validate:"omitempty"`
	ProductDecription string    `json:"ProductDecription" validate:"omitempty"`
	Price             int       `json:"Price" validate:"omitempty"`
	CreatedAt         time.Time `json:"DateAdded" validate:"omitempty"`
}

type IProduct interface {
	CreateProduct(prod TblProduct) int
	GetProductByProductId(productId int) TblProduct
	DeleteProductByProductId(productId int) error
}

func (cn PConnect) CreateProduct(prod TblProduct) int {

	if doInsert := cn.DbGorm.Table("TblProduct").Create(&prod).Error; doInsert != nil {
		return prod.Id
	} else {
		logrus.Error(doInsert)
		return 0
	}
}

func (cn PConnect) GetProductByProductId(productId int) TblProduct {
	prod := TblProduct{}
	cn.DbGorm.Table("TblProduct").Select("Id", "ProductName", "ProductDecription", "Price", "CreatedAt").Where("\"Id\"=?", productId).First(&prod)
	return prod
}

func (cn PConnect) DeleteProductByProductId(productId int) error {
	doDelete := cn.DbGorm.Delete(&TblProduct{}, productId).Error
	return doDelete
}
