package api

import (
	"net/http"
	"strconv"
	"time"

	dbSchema "github.com/Adebusy/cartbackendsvc/dataaccess"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// CreateGroup godoc
// @Summary		create group from existing cart.
// @Description	This action can only be performed by the admin.
// @Tags			group
// @Accept			*/*
// @User			json
// @Param user body dbSchema.TblGroupObj true "Create Group"
// @Param Authorization header string true "Authorization token"
// @Param clientName header string true "registered client name"
// @Security BearerAuth
// @securityDefinitions.basic BearerAuth
// @Success		200	{string}	string "Create group was successfully!!"
// @Failure		400		{string} string	"Unable to Create Group at the monent!!"
// @Router			/api/group/CreateGroup [post]
func CreateGroup(ctx *gin.Context) {
	// if !ValidateClient(ctx) {
	// 	return
	// }

	GroupObj := &dbSchema.TblGroupObj{}
	if doConvert := ctx.ShouldBindJSON(GroupObj); doConvert != nil {
		ctx.JSON(http.StatusBadRequest, doConvert)
		return
	}

	//do check
	if docheck := validateMe.Struct(GroupObj); docheck != nil {
		ctx.JSON(http.StatusBadRequest, docheck)
		return
	}

	doCheckCreatedById := usww.GetUserByUserId(GroupObj.UserId)
	if doCheckCreatedById.EmailAddress == "" {
		ctx.JSON(http.StatusBadRequest, "This user has not been created. UserId does not exist")
		return
	}

	if doCheckCartbyCartId := usww.GetCartByCartId(GroupObj.CartId); doCheckCartbyCartId.CartName == "" {
		ctx.JSON(http.StatusBadRequest, "Invalid cart supplied.")
		return
	}

	currentTime := time.Now()
	groupReq := &dbSchema.TblGroupUser{
		GroupName:   GroupObj.GroupName,
		Description: GroupObj.Description,
		UserId:      GroupObj.UserId,
		RoleId:      1,
		Status:      1,
		GroupTypeId: GroupObj.GroupTypeId,
		CartId:      GroupObj.CartId,
		DateAdded:   currentTime,
	}

	//do remove
	if createGoupUser := grp.CreateGroupUser(groupReq); createGoupUser == 0 {
		logrus.Error(createGoupUser)
		ctx.JSON(http.StatusBadRequest, "Service is unable to create or add user to group ATM, Please try again later!!")
		return
	}
	ctx.JSON(http.StatusOK, "Group created and user added successfully !!")
}

// AddUserToCartGroup godoc
// @Summary		Add user group from existing cart.
// @Description	This action can only be performed by the group admin.
// @Tags			group
// @Accept			*/*
// @User			json
// @Param user body dbSchema.TblTeamGroupObj true "Add user to Group"
// @Param Authorization header string true "Authorization token"
// @Param clientName header string true "registered client name"
// @Security BearerAuth
// @securityDefinitions.basic BearerAuth
// @Success		200	{string}	string "Create group was successfully!!"
// @Failure		400		{string} string	"Unable to Create Group at the monent!!"
// @Router			/api/group/AddUserToCartGroup [post]
func AddUserToCartGroup(ctx *gin.Context) {
	// if !ValidateClient(ctx) {
	// 	return
	// }

	GroupObj := &dbSchema.TblTeamGroupObj{}
	if doConvert := ctx.ShouldBindJSON(GroupObj); doConvert != nil {
		ctx.JSON(http.StatusBadRequest, doConvert)
		return
	}

	//do check
	if docheck := validateMe.Struct(GroupObj); docheck != nil {
		ctx.JSON(http.StatusBadRequest, docheck)
		return
	}

	checkIfInitiatorIsAdmin := grp.GetGroupAdminByUserIdAndRoleID(1, GroupObj.AdminId)
	if checkIfInitiatorIsAdmin.GroupName == "" {
		ctx.JSON(http.StatusBadRequest, "Only the group admin can add user to this group.")
		return
	}

	checkIfGroupUserAlreadyAdded := grp.GetGroupAdminByUserIdAndRoleID(2, GroupObj.AdminId)
	if checkIfGroupUserAlreadyAdded.GroupName != "" {
		ctx.JSON(http.StatusBadRequest, "This user had already been added to this group.")
		return
	}

	doCheckCreatedById := usww.GetUserByUserId(GroupObj.UserId)
	if doCheckCreatedById.EmailAddress == "" {
		ctx.JSON(http.StatusBadRequest, "This user has not been created. UserId does not exist")
		return
	}

	if doCheckCartbyCartId := usww.GetCartByCartId(GroupObj.CartId); doCheckCartbyCartId.CartName == "" {
		ctx.JSON(http.StatusBadRequest, "Invalid cart supplied.")
		return
	}

	currentTime := time.Now()
	groupReq := &dbSchema.TblGroupUser{
		GroupName:   GroupObj.GroupName,
		Description: GroupObj.Description,
		UserId:      GroupObj.UserId,
		RoleId:      2,
		Status:      1,
		GroupTypeId: GroupObj.GroupTypeId,
		CartId:      GroupObj.CartId,
		DateAdded:   currentTime,
	}

	//do remove
	if createGoupUser := grp.CreateGroupUser(groupReq); createGoupUser == 0 {
		logrus.Error(createGoupUser)
		ctx.JSON(http.StatusBadRequest, "Service is unable to create or add user to group ATM, Please try again later!!")
		return
	}
	ctx.JSON(http.StatusOK, "User added to the group successfully!!")
}

// GetGroupMemberByCartID Get group member by carID
// @Summary		Get group member by carID.
// @Description	Get group member by carID.
// @Tags			group
// @Produce json
// @Accept			*/*
// @User			json
// @Param Authorization header string true "Authorization token"
// @Param clientName header string true "registered client name"
// @Param CartId path string true "Cart Id"
// @Security BearerAuth
// @securityDefinitions.basic BearerAuth
// @Success		200	{object}	[]dbSchema.TblGroupUser
// @Router			/api/group/GetGroupMemberByCartID/{CartId} [get]
func GetGroupMemberByCartID(ctx *gin.Context) {
	// if !ValidateClient(ctx) {
	// 	return
	// }
	CartId, _ := strconv.Atoi(ctx.Param("CartId"))
	createGoupUser := grp.GetGroupMemberByCartID(CartId)
	if len(createGoupUser) == 0 {
		logrus.Error(createGoupUser)
		ctx.JSON(http.StatusBadRequest, "Service is unable to create or add user to group ATM, Please try again later!!")
		return
	}
	ctx.JSON(http.StatusOK, createGoupUser)
}

// GetGroupByUserID Get group by UserId
// @Summary		Get group by UserId.
// @Description	Get group by UserId.
// @Tags			group
// @Produce json
// @Accept			*/*
// @User			json
// @Param Authorization header string true "Authorization token"
// @Param clientName header string true "registered client name"
// @Param UserId path string true "User Id"
// @Security BearerAuth
// @securityDefinitions.basic BearerAuth
// @Success		200	{object}	[]dbSchema.TblGroupUser
// @Router			/api/group/GetGroupByUserID/{UserId} [get]
func GetGroupByUserID(ctx *gin.Context) {
	// if !ValidateClient(ctx) {
	// 	return
	// }
	userId, _ := strconv.Atoi(ctx.Param("UserId"))
	createGoupUser := grp.GetGroupByUserID(userId)
	if len(createGoupUser) == 0 {
		logrus.Error(createGoupUser)
		ctx.JSON(http.StatusBadRequest, "Service is unable to create or add user to group ATM, Please try again later!!")
		return
	}
	ctx.JSON(http.StatusOK, createGoupUser)
}

// RemoveUserFromCartGroup godoc
// @Summary		Remove user group from existing cart.
// @Description	This action can only be performed by the group admin.
// @Tags			group
// @Accept			*/*
// @User			json
// @Param user body dbSchema.RmoveUserFromGroupObj true "Remove user from Group"
// @Param Authorization header string true "Authorization token"
// @Param clientName header string true "registered client name"
// @Security BearerAuth
// @securityDefinitions.basic BearerAuth
// @Success		200	{string}	string "User remove successfully!!"
// @Failure		400		{string} string	"Unable to remove user from group at the monent!!"
// @Router			/api/group/RemoveUserFromCartGroup [post]
func RemoveUserFromCartGroup(ctx *gin.Context) {
	// if !ValidateClient(ctx) {
	// 	return
	// }

	GroupObj := &dbSchema.RmoveUserFromGroupObj{}
	if doConvert := ctx.ShouldBindJSON(GroupObj); doConvert != nil {
		ctx.JSON(http.StatusBadRequest, doConvert)
		return
	}

	//do check
	if docheck := validateMe.Struct(GroupObj); docheck != nil {
		ctx.JSON(http.StatusBadRequest, docheck)
		return
	}

	checkIfInitiatorIsAdmin := grp.GetGroupAdminByUserIdAndRoleID(1, GroupObj.AdminId)
	if checkIfInitiatorIsAdmin.GroupName == "" {
		ctx.JSON(http.StatusBadRequest, "Only the group admin can remove user from this group.")
		return
	}

	checkIfGroupUserAlreadyAdded := grp.GetGroupAdminByUserIdAndRoleID(2, GroupObj.AdminId)
	if checkIfGroupUserAlreadyAdded.GroupName == "" {
		ctx.JSON(http.StatusBadRequest, "This user has not been added to this group.")
		return
	}

	doCheckCreatedById := usww.GetUserByUserId(GroupObj.UserId)
	if doCheckCreatedById.EmailAddress == "" {
		ctx.JSON(http.StatusBadRequest, "This user has not been created. UserId does not exist")
		return
	}

	if doCheckCartbyCartId := usww.GetCartByCartId(GroupObj.CartId); doCheckCartbyCartId.CartName == "" {
		ctx.JSON(http.StatusBadRequest, "Invalid cart supplied.")
		return
	}
	//do remove
	if removeUserFromGoup := grp.RemoveUserFromGroup(2, GroupObj.CartId, GroupObj.GroupName); removeUserFromGoup == 0 {
		logrus.Error(removeUserFromGoup)
		ctx.JSON(http.StatusBadRequest, "Service is unable to create or add user to group ATM, Please try again later!!")
		return
	}
	ctx.JSON(http.StatusOK, "User added to the group successfully!!")
}
