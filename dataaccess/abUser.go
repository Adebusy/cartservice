package dataaccess

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type User struct {
	Id           int    `gorm:"column:Id"`
	TitleId      string `gorm:"column:TitleId"`
	UserName     string `gorm:"column:UserName"`
	NickName     string `gorm:"column:NickName"`
	FirstName    string `gorm:"column:FirstName"`
	LastName     string `gorm:"column:LastName"`
	EmailAddress string `gorm:"column:EmailAddress"`
	MobileNumber string `gorm:"column:MobileNumber"`
	Gender       string `gorm:"column:Gender"`
	Location     string `gorm:"column:Location"`
	AgeRange     string `gorm:"column:AgeRange"`
	Password     string `gorm:"column:Password"`
	Status       string `gorm:"column:Status"`
	CreatedAt    string `gorm:"column:CreatedAt"`
}

type CompleteSignUpReq struct {
	EmailAddress string `gorm:"column:EmailAddress"`
	TitleId      string `gorm:"column:TitleId"`
	UserName     string `gorm:"column:UserName"`
	NickName     string `gorm:"column:NickName"`
	FirstName    string `gorm:"column:FirstName"`
	LastName     string `gorm:"column:LastName"`
	MobileNumber string `gorm:"column:MobileNumber"`
	Gender       string `gorm:"column:Gender"`
	AgeRange     string `gorm:"column:AgeRange"`
	Status       int    `gorm:"column:Status"`
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
	Gender       string    `gorm:"column:Gender"`
	Location     string    `gorm:"column:Location"`
	AgeRange     string    `gorm:"column:AgeRange"`
	Password     string    `gorm:"column:Password"`
	Status       int       `gorm:"column:Status"`
	CreatedAt    time.Time `gorm:"column:CreatedAt"`
}

type ResponseMessage struct {
	ResponseCode    string
	ResponseMessage string
}

type DbConnect struct {
	DbGorm *gorm.DB
}

func ConneectDeal(db *gorm.DB) Iuser {
	return &DbConnect{db}
}

type Iuser interface {
	CreateUser(usr *User) string
	UpdateUserRecord(usr CompleteSignUpReq) string
	SignUp(emailAddress, mobileNumber, password, createdAt string) string
	GetUserByEmailAddress(EmailAddress string) User
	GetUserByUsername(EmailAddress string) User
	GetUserByMobileNumber(MobileNumber string) User
	LoginUser(UserName, Password string) User
	GetUserByUserId(UserId int) User

	CreateCart(crt TblCart) int
	GetCartByCartId(CartId int) TblCart
	// GetCartsByCartId(CartId int) []TblCart
	GetCartTypeByCartId(CartTypeId int) CartType
	CreateCartMember(cusr TblCartMember) int
	GetCartByCartIdAndMemberId(CartId, cartMemberId int) TblCart
	GetCartDetailsByCartIdandMastersId(CartId int, masterEmail string) TblCartMember
	CreateCartMemberIn(crt TblCartMember) int
	RemoveUserFromCart(CartId int, masterEmail, UserEmail string) error
	CloseCart(cartId int) int
	GetCartByUserId(cartId int) TblCart

	GetCartByUserEmail(email string) TblCart

	GetCartsByUserId(cartId int) []TblCart

	GetAllStatus() []TblStatus

	GetCartItemsByUserId(userId int) []TblCartItem

	GetCartItemsByCartId(cartId int) []TblCartItem
}

func (cn DbConnect) CreateUser(usr *User) string {
	if doinssert := cn.DbGorm.Table("TblUser").Create(&usr).Error; doinssert != nil {
		logrus.Error(doinssert)
		return "Unable to create user at the moment!!"
	} else {
		return "User created successfully!!"
	}
}

func (cn DbConnect) UpdateUserRecord(usr CompleteSignUpReq) string {

	if doinssertupdate := cn.DbGorm.Table("TblUser").Debug().Where("\"EmailAddress\"=? or \"MobileNumber\"=?", usr.EmailAddress, usr.MobileNumber).Updates(&usr).Error; doinssertupdate != nil {
		logrus.Error(doinssertupdate)
		return "Unable to create user at the moment!!"
	} else {
		logrus.Error(fmt.Sprintf("UpdateUserRecord for %s", usr.EmailAddress))
		return "User created successfully!!"
	}
}

func (cn DbConnect) SignUp(emailAddress, mobileNumber, password, createdAt string) string {

	if doinssert := cn.DbGorm.Table("TblUser").Select("FirstName", "LastName", "EmailAddress", "MobileNumber", "Password", "Status", "CreatedAt").Create(map[string]interface{}{"FirstName": "", "LastName": "", "EmailAddress": emailAddress, "MobileNumber": mobileNumber, "Password": password, "Status": "4", "CreatedAt": createdAt}).Error; doinssert != nil {
		logrus.Error(doinssert)
		return "Unable to create create sign up at the moment!!"
	} else {
		return "User signed up successfully!!"
	}
}

func (cn DbConnect) GetUserByUserId(UserId int) User {
	res := User{}
	cn.DbGorm.Table("TblUser").Select("TitleId", "UserName", "NickName", "FirstName", "LastName", "EmailAddress", "MobileNumber", "Gender", "Location", "AgeRange", "Password", "Status", "CreatedAt").Where("\"Id\"=?", UserId).First(&res)
	return res
}

func (cn DbConnect) GetUserByEmailAddress(EmailAddress string) User {
	res := User{}
	cn.DbGorm.Table("TblUser").Select("Id", "TitleId", "UserName", "NickName", "FirstName", "LastName", "EmailAddress", "MobileNumber", "Gender", "Location", "AgeRange", "Password", "Status", "CreatedAt").Where("\"EmailAddress\"=?", EmailAddress).First(&res)
	return res
}

func (cn DbConnect) GetUserByMobileNumber(MobileNumber string) User {
	res := User{}
	cn.DbGorm.Table("TblUser").Select("Id", "TitleId", "UserName", "NickName", "FirstName", "LastName", "EmailAddress", "MobileNumber", "Gender", "Location", "AgeRange", "Password", "Status", "CreatedAt").Where("\"MobileNumber\"=?", MobileNumber).First(&res)
	return res
}

func (cn DbConnect) GetUserByUsername(username string) User {
	res := User{}
	cn.DbGorm.Table("TblUser").Select("Id", "TitleId", "UserName", "NickName", "FirstName", "LastName", "EmailAddress", "MobileNumber", "Gender", "Location", "AgeRange", "Password", "Status", "CreatedAt", "Password").Where("\"UserName\"=?", username).First(&res)
	return res
}

func (cn DbConnect) LoginUser(UserName, Password string) User {
	res := User{}
	cn.DbGorm.Table("TblUser").Select("Id", "TitleId", "UserName", "NickName", "FirstName", "LastName", "EmailAddress", "MobileNumber", "Gender", "Location", "AgeRange", "Status", "CreatedAt", "Password").Where("\"UserName\"=? and \"Password\"=?", UserName, Password).First(&res)
	return res
}
