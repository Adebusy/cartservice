package dataaccess

import (
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CartConnect struct {
	DbGorm *gorm.DB
}

type TblCart struct {
	Id            int       `json:"Id" gorm:"unique;primaryKey;autoIncrement"`
	UserId        int       `gorm:"column:UserId"`
	CartTypeId    int       `gorm:"column:CartTypeId"`
	CartName      string    `gorm:"column:CartName"`
	Description   string    `gorm:"column:Description"`
	GroupId       int       `gorm:"column:GroupId"`
	CreatedById   int       `gorm:"column:CreatedById"`
	Status        string    `gorm:"column:Status"`
	CreatedAt     time.Time `gorm:"column:CreatedAt"`
	LastUpdatedBy int       `gorm:"column:LastUpdatedBy"`
}

type TblCartMember struct {
	Id              int       `json:"Id" gorm:"unique;primaryKey;autoIncrement"`
	RingMasterEmail string    `json:"RingMasterEmail" validate:"omitempty"`
	MemberEmail     string    `json:"MemberEmail" validate:"omitempty"`
	CartId          int       `json:"CartId" validate:"omitempty"`
	RingStatus      int       `json:"RingStatus" validate:"omitempty"`
	DateAdded       time.Time `json:"DateAdded" validate:"omitempty"`
}

type CartUser struct {
	Id              int       `gorm:"column:Id"`
	RingMasterEmail string    `gorm:"column:RingMasterEmail"`
	MemberEmail     string    `gorm:"column:MemberEmail"`
	CartId          int       `gorm:"column:CartId"`
	RingStatus      int       `gorm:"column:RingStatus"`
	DateAdded       time.Time `gorm:"column:DateAdded"`
}

type CartUserIn struct {
	RingMasterEmail string    `gorm:"column:RingMasterEmail"`
	MemberEmail     string    `gorm:"column:MemberEmail"`
	CartId          int       `gorm:"column:CartId"`
	RingStatus      int       `gorm:"column:RingStatus"`
	DateAdded       time.Time `gorm:"column:DateAdded"`
}

type CartType struct {
	Id           int       `gorm:"column:Id"`
	CartTypeName string    `gorm:"column:CartTypeName"`
	DateAdded    time.Time `gorm:"column:DateAdded"`
}

type TblCartType struct {
	Id           int       `gorm:"column:Id"`
	CartTypeName string    `gorm:"column:CartTypeName"`
	Status       int       `gorm:"column:Status"`
	DateAdded    time.Time `gorm:"column:DateAdded"`
}

func (cn DbConnect) CreateCart(crt TblCart) int {

	dogetcart := cn.DbGorm.Table("TblCart").Create(&crt)
	if dogetcart.Error != nil {
		logrus.Error(dogetcart)
		return 0
	} else {
		return crt.Id
	}
}

func (cn DbConnect) CloseCart(cartId int) int {
	err := cn.DbGorm.Table("TblCart").Where("\"Id\"=?", cartId).Update("Status", "0")
	if err.Error != nil {
		logrus.Error(err)
		return 0
	}
	return cartId
}
func (cn DbConnect) CreateCartMember(crt TblCartMember) int {
	dogetcart := cn.DbGorm.Table("TblCartMember").Create(&crt)
	if dogetcart.Error != nil {
		logrus.Error(dogetcart)
		return 0
	} else {
		return crt.Id
	}
}
func (cn DbConnect) CreateCartMemberIn(crt TblCartMember) int {

	dogetcart := cn.DbGorm.Create(&crt)
	if dogetcart.Error != nil {
		logrus.Error(dogetcart)
		return 0
	} else {
		return crt.Id
	}
}

func (cn DbConnect) GetCartTypeByCartId(CartTypeId int) CartType {
	CartType := CartType{}
	cn.DbGorm.Table("TblCartType").Select("Id", "CartTypeName", "DateAdded").Where("\"Id\"=? and \"Status\"=?", CartTypeId, 1).First(&CartType)
	return CartType
}

func (cn DbConnect) GetCartByCartId(CartId int) TblCart {
	res := TblCart{}
	cn.DbGorm.Table("TblCart").Select("UserId", "CartTypeId", "CartName", "Description", "GroupId", "CreatedById", "Status", "CreatedAt", "LastUpdatedBy").Where("\"Id\"=?", CartId).First(&res)
	return res
}

func (cn DbConnect) GetCartByCartIdAndMemberId(CartId, cartMemberId int) TblCart {
	res := TblCart{}
	cn.DbGorm.Table("TblCart").Select("UserId", "CartTypeId", "CartName", "Description", "GroupId", "CreatedById", "Status", "CreatedAt", "LastUpdatedBy").Where("\"Id\"=? and \"UserId\"=?", CartId, cartMemberId).First(&res)
	return res
}

func (cn DbConnect) GetCartDetailsByCartIdandMastersId(CartId int, masterEmail string) TblCartMember {
	tblCartMember := TblCartMember{}
	cn.DbGorm.Select("Id", "RingMasterEmail", "MemberEmail", "CartId", "RingStatus", "DateAdded").Where("\"CartId\" = ? and \"RingMasterEmail\" =? ", CartId, masterEmail).First(&tblCartMember)
	return tblCartMember
}

func (cn DbConnect) RemoveUserFromCart(CartId int, masterEmail string, UserEmail string) error {
	tblCartMember := TblCartMember{}
	dodelete := cn.DbGorm.Where("\"CartId\"=? and \"RingMasterEmail\"=? and \"MemberEmail\"=? ", CartId, masterEmail, UserEmail).Delete(&tblCartMember).Error
	return dodelete
}

func (cn DbConnect) GetCartByUserId(userId int) TblCart {
	res := TblCart{}
	cn.DbGorm.Debug().Select("Id", "UserId", "CartTypeId", "CartName", "Description", "GroupId", "CreatedById", "Status", "CreatedAt", "LastUpdatedBy").Where("\"UserId\"=?", userId).First(&res)
	return res
}

func (cn DbConnect) GetCartByUserEmail(email string) TblCart {
	res := TblCart{}
	cn.DbGorm.Debug().Select("Id", "UserId", "CartTypeId", "CartName", "Description", "GroupId", "CreatedById", "Status", "CreatedAt", "LastUpdatedBy").Where("\"email\"=?", email).First(&res)
	return res
}
