package dataaccess

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	TitleId      string `gorm:"column:TitleId"`
	UserName     string `gorm:"column:UserName"`
	NickName     string `gorm:"column:NickName"`
	FirstName    string `gorm:"column:FirstName"`
	LastName     string `gorm:"column:LastName"`
	EmailAddress string `gorm:"column:EmailAddress"`
	MobileNumber string `gorm:"column:MobileNumber"`
	Password     string `gorm:"column:Password"`
	Status       string `gorm:"column:Status"`
	CreatedAt    string `gorm:"column:CreatedAt"`
}

type TblUser struct {
	Id           int       `json:"Id" gorm:"unique;primaryKey;autoIncrement"`
	TitleId      int       `gorm:"column:TitleId"`
	UserName     string    `gorm:"column:UserName"`
	NickName     string    `gorm:"column:NickName"`
	FirstName    string    `gorm:"column:FirstName"`
	LastName     string    `gorm:"column:LastName"`
	EmailAddress string    `gorm:"column:EmailAddress"`
	MobileNumber string    `gorm:"column:MobileNumber"`
	Password     string    `gorm:"column:Password"`
	Status       int       `gorm:"column:Status"`
	CreatedAt    time.Time `gorm:"column:CreatedAt"`
}

type DbConnect struct {
	DbGorm *gorm.DB
}

func ConneectDeal(db *gorm.DB) Iuser {
	return &DbConnect{db}
}

type Iuser interface {
	CreateUser(usr *User) string
	GetUserByEmailAddress(EmailAddress string) User
	GetUserByEmailUsername(EmailAddress string) User
	GetUserByMobileNumber(MobileNumber string) User
	LoginUser(UserName, Password string) User
	GetUserByUserId(UserId int) User

	CreateCart(crt TblCart) int
	GetCartByCartId(CartId int) TblCart
	GetCartTypeByCartId(CartTypeId int) CartType
	CreateCartMember(cusr TblCartMember) int
	GetCartByCartIdAndMemberId(CartId, cartMemberId int) TblCart
	GetCartDetailsByCartIdandMastersId(CartId int, masterEmail string) TblCartMember
	CreateCartMemberIn(crt TblCartMember) int
	RemoveUserFromCart(CartId int, masterEmail string, UserEmail string) error
	CloseCart(cartId int) int
}

func (cn DbConnect) CreateUser(usr *User) string {
	if doinssert := cn.DbGorm.Table("TblUser").Create(&usr).Error; doinssert != nil {
		return "Unable to create user at the moment!!"
	} else {
		return "User created successfully!!"
	}
}

func (cn DbConnect) GetUserByUserId(UserId int) User {
	res := User{}
	cn.DbGorm.Table("TblUser").Select("TitleId", "UserName", "NickName", "FirstName", "LastName", "EmailAddress", "MobileNumber", "Status", "CreatedAt").Where("\"Id\"=?", UserId).First(&res)
	return res
}

func (cn DbConnect) GetUserByEmailAddress(EmailAddress string) User {
	res := User{}
	cn.DbGorm.Table("TblUser").Select("TitleId", "UserName", "NickName", "FirstName", "LastName", "EmailAddress", "MobileNumber", "Status", "CreatedAt").Where("\"EmailAddress\"=?", EmailAddress).First(&res)
	return res
}

func (cn DbConnect) GetUserByMobileNumber(MobileNumber string) User {
	res := User{}
	cn.DbGorm.Table("TblUser").Select("TitleId", "UserName", "NickName", "FirstName", "LastName", "EmailAddress", "MobileNumber", "Status", "CreatedAt").Where("\"MobileNumber\"=?", MobileNumber).First(&res)
	return res
}

func (cn DbConnect) GetUserByEmailUsername(username string) User {
	res := User{}
	cn.DbGorm.Table("TblUser").Select("TitleId", "UserName", "NickName", "FirstName", "LastName", "EmailAddress", "MobileNumber", "Status", "CreatedAt", "Password").Where("\"UserName\"=?", username).First(&res)
	return res
}

func (cn DbConnect) LoginUser(UserName, Password string) User {
	res := User{}
	cn.DbGorm.Table("TblUser").Select("Id", "TitleId", "UserName", "NickName", "FirstName", "LastName", "EmailAddress", "MobileNumber", "Status", "CreatedAt", "Password").Where("\"UserName\"=? and \"Password\"=?", UserName, Password).First(&res)
	return res
}
