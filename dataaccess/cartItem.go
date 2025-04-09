package dataaccess

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type CartItemObj struct {
	CartId      int    `json:"CartId" validate:"omitempty"`
	Name        string `json:"Name" validate:"omitempty"`
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
	Name        string    `json:"Name" validate:"omitempty"`
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
	GetCartItemsByUserId(UserId int) []TblCartItem
	GetCartItemsByCartId(CartId int) []TblCartItem
}

func (cn DbConnect) AddItemToCart(CartItem TblCartItem) int {
	if createCartItem := cn.DbGorm.Debug().Table("TblCartItem").Create(&CartItem).Error; createCartItem != nil {
		fmt.Printf("Error %s", createCartItem)
		return 0
	}
	return CartItem.Id
}

func (cn DbConnect) RemoveItemFromCart(ProductId, CartId, UserId int) error {
	doDeleteItem := cn.DbGorm.Debug().Where("\"CartId\"=? and \"ProductId\"=? and \"UserId\"=?", CartId, ProductId, UserId).Delete(&TblCartItem{}).Error
	return doDeleteItem
}

func (cn DbConnect) GetCartItemsByUserId(UserId int) []TblCartItem {
	resp := []TblCartItem{}
	cn.DbGorm.Select([]string{"Id", "Name", "CartId", "ProductId", "Quantity", "Description", "DateAdded", "UserId"}).Where("\"UserId\"=?", UserId).Find(&resp)
	return resp
}

func (cn DbConnect) GetCartItemsByCartId(CartId int) []TblCartItem {
	resp := []TblCartItem{}
	cn.DbGorm.Select([]string{"Id", "Name", "CartId", "ProductId", "Quantity", "Description", "DateAdded", "UserId"}).Where("\"CartId\"=?", CartId).Find(&resp)
	return resp
}
