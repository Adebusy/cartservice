package dataaccess

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type CartItemObj struct {
	CartId      int    `json:"CartId" validate:"omitempty"`
	ProductId   int    `json:"ProductId" validate:"omitempty"`
	Quantity    int    `json:"Quantity" validate:"omitempty"`
	Description string `json:"Description" validate:"omitempty"`
	UserId      int    `json:"UserId" validate:"omitempty"`
}

type RemoveCartItemObj struct {
	CartId    int `json:"CartId" validate:"omitempty"`
	ProductId int `json:"ProductId" validate:"omitempty"`
	UserId    int `json:"UserId" validate:"omitempty"`
}

type TblCartItem struct {
	Id          int       `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	CartId      int       `json:"CartId" validate:"omitempty"`
	ProductId   int       `json:"ProductId" validate:"omitempty"`
	Quantity    int       `json:"Quantity" validate:"omitempty"`
	Description string    `json:"Description" validate:"omitempty"`
	DateAdded   time.Time `json:"DateAdded" validate:"omitempty"`
	UserId      int       `json:"UserId" validate:"omitempty"`
}

func ConnectCartItem(db *gorm.DB) ICartItem {
	return &DbConnect{db}
}

type ICartItem interface {
	AddItemToCart(CartItem TblCartItem) int
	RemoveItemFromCart(ProductId, CartId, UserId int) error
}

func (cn DbConnect) AddItemToCart(CartItem TblCartItem) int {
	if createCartItem := cn.DbGorm.Table("TblCartItem").Create(&CartItem).Error; createCartItem != nil {
		fmt.Printf("Error %s", createCartItem)
		return 0
	}
	return CartItem.Id
}

func (cn DbConnect) RemoveItemFromCart(ProductId, CartId, UserId int) error {
	doDeleteItem := cn.DbGorm.Where("\"CartId\"=? and \"ProductId\"=? and \"UserId\"=?", CartId, ProductId, UserId).Delete(&TblCartItem{}).Error
	return doDeleteItem
}
