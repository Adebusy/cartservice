package dataaccess

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type GConnect struct {
	DbGorm *gorm.DB
}

func ConnectGroup(db *gorm.DB) IGroup {
	return &GConnect{db}
}

type IGroup interface {
	CreateGroupUser(prod *TblGroupUser) int
	GetGroupMemberByCartID(cartId int) []TblGroupUser
	GetGroupByUserID(userId int) []TblGroupUser
	GetGroupAdminByUserIdAndRoleID(roleId, userId int) TblGroupUser
	RemoveUserFromGroup(userId, cartId int, GroupName string) int
}

func (cn GConnect) CreateGroupUser(prod *TblGroupUser) int {

	if doInsert := cn.DbGorm.Debug().Table("TblGroupUser").Create(&prod).Debug().Error; doInsert == nil {
		return prod.Id
	} else {
		logrus.Error(doInsert)
		return 0
	}
}

func (cn GConnect) GetGroupMemberByCartID(cartId int) []TblGroupUser {
	prod := []TblGroupUser{}
	cn.DbGorm.Table("TblGroupUser").Select("Id", "GroupName", "Status", "Description", "UserId", "RoleId", "GroupTypeId", "CartId", "DateAdded").Where("\"CartId\"=?", cartId).Find(&prod)
	return prod
}

func (cn GConnect) GetGroupByUserID(userId int) []TblGroupUser {
	prod := []TblGroupUser{}
	cn.DbGorm.Table("TblGroupUser").Debug().Select("Id", "GroupName", "Status", "Description", "UserId", "RoleId", "GroupTypeId", "CartId", "DateAdded").Where("\"UserId\"=?", userId).Find(&prod)
	return prod
}

func (cn GConnect) GetGroupAdminByUserIdAndRoleID(roleId, userId int) TblGroupUser {
	prod := TblGroupUser{}
	cn.DbGorm.Table("TblGroupUser").Select("Id", "GroupName", "Status", "Description", "UserId", "RoleId", "GroupTypeId", "CartId", "DateAdded").Where("\"UserId\"=? and \"RoleId\"=?", userId, roleId).Debug().Find(&prod)
	return prod
}

func (cn GConnect) RemoveUserFromGroup(roleId, userId int, groupName string) int {
	prod := TblGroupUser{}
	if err := cn.DbGorm.Table("TblGroupUser").Where("\"UserId\"=? and \"RoleId\"=? and \"GroupName\"=?", userId, roleId, groupName).Delete(prod).Error; err != nil {
		fmt.Print(err)
		return 0
	}

	return 1
}
