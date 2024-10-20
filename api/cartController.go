package api

import (
	"net/http"
	"time"

	dbSchema "github.com/Adebusy/cartbackendsvc/dataaccess"
	inputschema "github.com/Adebusy/cartbackendsvc/obj"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// CreateCart godoc
// @Summary		Create new  cart user.
// @Description	Create new cart user.
// @Tags			cart
// @Accept			*/*
// @User			json
// @Param user body inputschema.CartObj true "Create new user"
// @Success		200	{object}	inputschema.ResponseMessage
// @Router			/api/cart/CreateCart [post]
func CreateCart(ctx *gin.Context) {

	carObj := &inputschema.CartObj{}
	if doConvert := ctx.ShouldBindJSON(carObj); doConvert != nil {
		ctx.JSON(http.StatusBadRequest, doConvert)
		return
	}

	//do check
	if docheck := validateMe.Struct(carObj); docheck != nil {
		ctx.JSON(http.StatusBadRequest, docheck)
		return
	}

	//do check userID
	doCheckUser := usww.GetUserByUserId(carObj.UserId)
	if doCheckUser.EmailAddress == "" {
		ctx.JSON(http.StatusBadRequest, "UserId does not exist.")
		return
	}

	//do check userID
	if doCheckCreatedById := usww.GetUserByUserId(carObj.CreatedById); doCheckCreatedById.EmailAddress == "" {
		ctx.JSON(http.StatusBadRequest, "CreatedById does not exist.")
		return
	}

	// check CartTypeId
	if docheckCartId := crt.GetCartTypeByCartId(carObj.CartTypeId); docheckCartId.CartTypeName == "" {
		ctx.JSON(http.StatusBadRequest, "CartTypeId does not exist.")
		return
	}

	crts := dbSchema.TblCart{
		UserId:        carObj.UserId,
		CartTypeId:    carObj.CartTypeId,
		CartName:      carObj.CartName,
		Description:   carObj.Description,
		GroupId:       carObj.GroupId,
		CreatedById:   carObj.CreatedById,
		Status:        "1",
		CreatedAt:     time.Now(),
		LastUpdatedBy: carObj.UserId,
	}

	if doCreate := crt.CreateCart(crts); doCreate != 0 {
		crts := dbSchema.TblCartMember{
			RingMasterEmail: doCheckUser.EmailAddress,
			MemberEmail:     doCheckUser.EmailAddress,
			CartId:          doCreate,
			RingStatus:      1,
			DateAdded:       time.Now(),
		}

		if CreateCartMember := crt.CreateCartMember(crts); CreateCartMember != 0 {
			logrus.Info("Added member to cart")
		}
		ctx.JSON(http.StatusOK, "Cart created successfully!!")
		return
	} else {
		ctx.JSON(http.StatusBadRequest, "Cart cannot be created at the monent!!")
		return
	}
}

// CreateCartMember godoc
// @Summary		Create new  Cart Member.
// @Description	Create new Cart Member.
// @Tags			cart
// @Accept			*/*
// @User			json
// @Param user body inputschema.CartUserObj true "Create new cart member"
// @Success		200	{object}	inputschema.ResponseMessage
// @Router			/api/cart/CreateCartMember [post]
func CreateCartMember(ctx *gin.Context) {

	carObj := &inputschema.CartUserObj{}
	if doConvert := ctx.ShouldBindJSON(carObj); doConvert != nil {
		ctx.JSON(http.StatusBadRequest, doConvert)
		return
	}

	//do check
	if docheck := validateMe.Struct(carObj); docheck != nil {
		ctx.JSON(http.StatusBadRequest, docheck)
		return
	}

	//do check userID
	if doCheckUser := usww.GetUserByEmailAddress(carObj.RingMasterEmail); doCheckUser.EmailAddress == "" {
		ctx.JSON(http.StatusBadRequest, "RingMaster Email does not exist.")
		return
	}

	//do check userID
	if doCheckCreatedById := usww.GetUserByEmailAddress(carObj.MemberEmail); doCheckCreatedById.EmailAddress == "" {
		ctx.JSON(http.StatusBadRequest, "Member Email does not exist.")
		return
	}

	//check if initiator is the cart master
	if GetCartDetailsByCartId := crt.GetCartDetailsByCartIdandMastersId(carObj.CartId, carObj.RingMasterEmail); GetCartDetailsByCartId.RingMasterEmail == "" {
		ctx.JSON(http.StatusBadRequest, "This user does not have the permission required to execute this action.")
		return
	}

	crts := dbSchema.TblCartMember{
		RingMasterEmail: carObj.RingMasterEmail,
		MemberEmail:     carObj.MemberEmail,
		CartId:          carObj.CartId,
		RingStatus:      carObj.RingStatus,
		DateAdded:       time.Now(),
	}

	if doCreate := crt.CreateCartMember(crts); doCreate != 0 {
		ctx.JSON(http.StatusOK, "Cart created successfully!!")
		return
	} else {
		ctx.JSON(http.StatusBadRequest, "Cart cannot be created at the monent!!")
		return
	}
}

// RemoveUserFromCart godoc
// @Summary		Remove user from cart.
// @Description	Remove user from cart.
// @Tags			cart
// @Accept			*/*
// @User			json
// @Param user body inputschema.RemoveUserFromCartObj true "Remove member from cart"
// @Success		200	{object}	inputschema.ResponseMessage
// @Router			/api/cart/RemoveUserFromCart [post]
func RemoveUserFromCart(ctx *gin.Context) {

	resp := inputschema.ResponseMessage{}
	requestObj := &inputschema.RemoveUserFromCartObj{}
	if doConvert := ctx.ShouldBindJSON(requestObj); doConvert != nil {
		resp.ResponseCode = "01"
		resp.ResponseMessage = doConvert.Error()
		ctx.JSON(http.StatusBadRequest, doConvert)
		return
	}

	//do check
	if docheck := validateMe.Struct(requestObj); docheck != nil {
		resp.ResponseCode = "01"
		resp.ResponseMessage = "Ring master email does not exist."
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	//do check userID
	if doCheckUser := usww.GetUserByEmailAddress(requestObj.RingMasterEmail); doCheckUser.EmailAddress == "" {
		resp.ResponseCode = "01"
		resp.ResponseMessage = "Ring master email does not exist."
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	//do check userID
	if doCheckCreatedById := usww.GetUserByEmailAddress(requestObj.MemberEmail); doCheckCreatedById.EmailAddress == "" {
		resp.ResponseCode = "01"
		resp.ResponseMessage = "Member email does not exist."
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	//check if initiator is the cart master
	if GetCartDetailsByCartId := crt.GetCartDetailsByCartIdandMastersId(requestObj.CartId, requestObj.RingMasterEmail); GetCartDetailsByCartId.RingMasterEmail == "" {
		resp.ResponseCode = "01"
		resp.ResponseMessage = "This user does not have the permission required to execute this action."
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	//delete user from cart
	if doCreate := crt.RemoveUserFromCart(requestObj.CartId, requestObj.RingMasterEmail, requestObj.MemberEmail); doCreate == nil {
		resp.ResponseCode = "00"
		resp.ResponseMessage = "User removed from cart successfully!!"
		ctx.JSON(http.StatusOK, resp)
		return
	} else {
		resp.ResponseCode = "01"
		resp.ResponseMessage = "User cannot be removed at the monent!!"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
}

// CloseCart godoc
// @Summary		Close  Cart.
// @Description	Close Cart.
// @Tags			cart
// @Accept			*/*
// @User			json
// @Param user body inputschema.CloseCartObj true "Close cart"
// @Success		200	{object}	inputschema.ResponseMessage
// @Router			/api/cart/CloseCart [put]
func CloseCart(ctx *gin.Context) {

	carObj := &inputschema.CloseCartObj{}
	if doConvert := ctx.ShouldBindJSON(carObj); doConvert != nil {
		ctx.JSON(http.StatusBadRequest, doConvert)
		return
	}

	//do check
	if docheck := validateMe.Struct(carObj); docheck != nil {
		ctx.JSON(http.StatusBadRequest, docheck)
		return
	}

	//do check userID
	if doCheckUser := usww.GetUserByEmailAddress(carObj.RingMasterEmail); doCheckUser.EmailAddress == "" {
		ctx.JSON(http.StatusBadRequest, "RingMaster Email does not exist.")
		return
	}

	//check if initiator is the cart master
	if GetCartDetailsByCartId := crt.GetCartDetailsByCartIdandMastersId(carObj.CartId, carObj.RingMasterEmail); GetCartDetailsByCartId.RingMasterEmail == "" {
		ctx.JSON(http.StatusBadRequest, "This user does not have the permission required to execute this action.")
		return
	}

	// update cart
	if doCreate := crt.CloseCart(carObj.CartId); doCreate != 0 {
		ctx.JSON(http.StatusOK, "Cart closed successfully!!")
		return
	} else {
		ctx.JSON(http.StatusBadRequest, "Cart cannot be close at the monent!!")
		return
	}
}
