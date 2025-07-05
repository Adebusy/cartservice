package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	dbSchema "github.com/Adebusy/cartbackendsvc/dataaccess"
	inputschema "github.com/Adebusy/cartbackendsvc/obj"
	"github.com/Adebusy/cartbackendsvc/utilities"
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
// @Param Authorization header string true "Authorization token"
// @Param clientName header string true "registered client name"
// @Security BearerAuth
// @securityDefinitions.basic BearerAuth
// @Success		200	{object}	inputschema.ResponseMessage
// @Router			/api/cart/CreateCart [post]
func CreateCart(ctx *gin.Context) {
	if !ValidateClient(ctx) {
		return
	}
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
	if docheckCartId := usww.GetCartTypeByCartId(carObj.CartTypeId); docheckCartId.CartTypeName == "" {
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

	if doCreate := usww.CreateCart(crts); doCreate != 0 {
		crts := dbSchema.TblCartMember{
			RingMasterEmail: doCheckUser.EmailAddress,
			MemberEmail:     doCheckUser.EmailAddress,
			CartId:          doCreate,
			RingStatus:      1,
			DateAdded:       time.Now(),
		}

		if CreateCartMember := usww.CreateCartMember(crts); CreateCartMember != 0 {
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
// @Param Authorization header string true "Authorization token"
// @Param clientName header string true "registered client name"
// @Security BearerAuth
// @securityDefinitions.basic BearerAuth
// @Success		200	{object}	inputschema.ResponseMessage
// @Router			/api/cart/CreateCartMember [post]
func CreateCartMember(ctx *gin.Context) {
	if !ValidateClient(ctx) {
		return
	}
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

	//check if initiator is the cart master
	if GetCartDetailsByCartId := usww.GetCartDetailsByCartIdandMastersId(carObj.CartId, carObj.RingMasterEmail); GetCartDetailsByCartId.RingMasterEmail == "" {
		ctx.JSON(http.StatusBadRequest, "This user does not have the permission required to execute this action.")
		return
	}

	//do check userID
	if doCheckCreatedById := usww.GetUserByEmailAddress(carObj.MemberEmail); doCheckCreatedById.EmailAddress == "" {
		utilities.SendEmail(carObj.MemberEmail, "please download the digital cart application from apple store.")
		//ctx.JSON(http.StatusBadRequest, "Member Email does not exist.")
		//return
	}

	crts := dbSchema.TblCartMember{
		RingMasterEmail: carObj.RingMasterEmail,
		MemberEmail:     carObj.MemberEmail,
		CartId:          carObj.CartId,
		RingStatus:      carObj.RingStatus,
		DateAdded:       time.Now(),
	}

	CartByCartId := usww.GetCartByCartId(carObj.CartId)
	if doCreate := usww.CreateCartMember(crts); doCreate != 0 {
		utilities.SendEmail(carObj.MemberEmail, fmt.Sprintf("Hello, You have been invited to join a cart %s. please check", CartByCartId.CartName))
		TAct.CreateAction(dbSchema.TblAction{EmailAddress: carObj.MemberEmail, MobileNumber: "", RequestType: "Notification", Message: fmt.Sprintf("Hello, You have been invited to join a cart %s. please check", CartByCartId.CartName), Status: 1, DateAdded: time.Now()})
		ctx.JSON(http.StatusOK, "A new member added to cart successfully!!")
		return
	} else {
		ctx.JSON(http.StatusBadRequest, "Member cannot not be added to card at the monent!!")
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
// @Param Authorization header string true "Authorization token"
// @Param clientName header string true "registered client name"
// @Security BearerAuth
// @securityDefinitions.basic BearerAuth
// @Success		200	{object}	inputschema.ResponseMessage
// @Router			/api/cart/RemoveUserFromCart [post]
func RemoveUserFromCart(ctx *gin.Context) {
	if !ValidateClient(ctx) {
		return
	}
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
	if GetCartDetailsByCartId := usww.GetCartDetailsByCartIdandMastersId(requestObj.CartId, requestObj.RingMasterEmail); GetCartDetailsByCartId.RingMasterEmail == "" {
		resp.ResponseCode = "01"
		resp.ResponseMessage = "This user does not have the permission required to execute this action."
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	//delete user from cart
	if doCreate := usww.RemoveUserFromCart(requestObj.CartId, requestObj.RingMasterEmail, requestObj.MemberEmail); doCreate == nil {
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
// @Param Authorization header string true "Authorization token"
// @Param clientName header string true "registered client name"
// @Security BearerAuth
// @securityDefinitions.basic BearerAuth
// @Success 200	{string} string	"Cart closed successfully!"
// @Router			/api/cart/CloseCart [put]
func CloseCart(ctx *gin.Context) {
	if !ValidateClient(ctx) {
		return
	}

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
	if GetCartDetailsByCartId := usww.GetCartDetailsByCartIdandMastersId(carObj.CartId, carObj.RingMasterEmail); GetCartDetailsByCartId.RingMasterEmail == "" {
		ctx.JSON(http.StatusBadRequest, "This user does not have the permission required to execute this action.")
		return
	}

	// update cart
	if doCreate := usww.CloseCart(carObj.CartId); doCreate != 0 {
		ctx.JSON(http.StatusOK, "Cart closed successfully!!")
		return
	} else {
		ctx.JSON(http.StatusBadRequest, "Cart cannot be close at the monent!!")
		return
	}
}

// GetCartByUserId Get Cart By User Id
// @Summary		Get Cart By User Id.
// @Description	Get Cart By User Id.
// @Tags			cart
// @Produce json
// @Accept			*/*
// @User			json
// @Param Authorization header string true "Authorization token"
// @Param UserId path string true "User ID"
// @Security BearerAuth
// @securityDefinitions.basic BearerAuth
// @Success		200	{object}	dbSchema.TblCart
// @Router			/api/cart/GetCartByUserId/{UserId} [get]
func GetCartByUserId(ctx *gin.Context) {
	if !ValidateClient(ctx) {
		return
	}
	userId, _ := strconv.Atoi(ctx.Param("UserId"))
	// update cart
	ctx.JSON(http.StatusOK, usww.GetCartByUserId(userId))
}

// GetCartMembersListByCartId Get Cart By Cart Id
// @Summary		Get Cart By Cart Id.
// @Description	Get Cart By Cart Id.
// @Tags			cart
// @Produce json
// @Accept			*/*
// @User			json
// @Param Authorization header string true "Authorization token"
// @Param UserId path string true "Cart ID"
// @Security BearerAuth
// @securityDefinitions.basic BearerAuth
// @Success		200	{object}	dbSchema.TblCart
// @Router			/api/cart/GetCartMembersListByCartId/{CartId} [get]
func GetCartMembersListByCartId(ctx *gin.Context) {
	// if !ValidateClient(ctx) {
	// 	return
	// }
	retlist := []cartMember{}
	cartId, _ := strconv.Atoi(ctx.Param("CartId"))
	// rawget := usww.GetCartMemberByCartId(cartId)
	for _, i := range usww.GetCartMemberByCartId(cartId) {
		if getdet := usww.GetUserByEmailAddress(i.MemberEmail); getdet.EmailAddress != "" {
			retlist = append(retlist, cartMember{Name: getdet.FirstName, Email: getdet.EmailAddress, MobileNumber: getdet.MobileNumber})
		}
	}
	ctx.JSON(http.StatusOK, retlist)
}

type cartMember struct {
	Name         string `json:"Name" validate:"required,email"`
	Email        string `json:"Email" validate:"required,email"`
	MobileNumber string `json:"MobileNumber" validate:"required,min=8"`
}

// GetOpenCartsByUserId Get Carts By User Id
// @Summary		Get Carts By User Id.
// @Description	Get Carts By User Id.
// @Tags			cart
// @Produce json
// @Accept			*/*
// @User			json
// @Param Authorization header string true "Authorization token"
// @Param UserId path string true "User ID"
// @Security BearerAuth
// @securityDefinitions.basic BearerAuth
// @Success		200	{object}	[]dbSchema.TblCart
// @Router			/api/cart/GetOpenCartsByUserId/{UserId} [get]
func GetOpenCartsByUserId(ctx *gin.Context) {
	// if !ValidateClient(ctx) {
	// 	return
	// }
	userId, _ := strconv.Atoi(ctx.Param("UserId"))
	ctx.JSON(http.StatusOK, usww.GetOpenCartsByUserIdandStatus(userId, 1))
}

// GetClosedCartsByUserId Get Closed Carts By User Id
// @Summary		Get Closed Carts By User Id.
// @Description	Get Closed Carts By User Id.
// @Tags			cart
// @Produce json
// @Accept			*/*
// @User			json
// @Param Authorization header string true "Authorization token"
// @Param UserId path string true "User ID"
// @Security BearerAuth
// @securityDefinitions.basic BearerAuth
// @Success		200	{object}	[]dbSchema.TblCart
// @Router			/api/cart/GetClosedCartsByUserId/{UserId} [get]
func GetClosedCartsByUserId(ctx *gin.Context) {
	// if !ValidateClient(ctx) {
	// 	return
	// }

	userId, _ := strconv.Atoi(ctx.Param("UserId"))
	// update cart
	ctx.JSON(http.StatusOK, usww.GetClosedCartsByUserIdandStatus(userId, 0))
}

// GetCartsByUserId Get Carts By User Id
// @Summary		Get Carts By User Id.
// @Description	Get Carts By User Id.
// @Tags			cart
// @Produce json
// @Accept			*/*
// @User			json
// @Param Authorization header string true "Authorization token"
// @Param UserId path string true "User ID"
// @Security BearerAuth
// @securityDefinitions.basic BearerAuth
// @Success		200	{object}	[]dbSchema.TblCart
// @Router			/api/cart/GetCartsByUserId/{UserId} [get]
func GetCartsByUserId(ctx *gin.Context) {
	if !ValidateClient(ctx) {
		return
	}
	userId, _ := strconv.Atoi(ctx.Param("UserId"))
	fmt.Printf("dsdsds %s", strconv.Itoa(userId))
	// update cart
	ctx.JSON(http.StatusOK, usww.GetCartsByUserId(userId))
}

// GetCartByUserEmail Get Cart By User Id
// @Summary		Get Cart By User Email.
// @Description	Get Cart By User Email.
// @Tags			cart
// @Produce json
// @Accept			*/*
// @User			json
// @Param Authorization header string true "Authorization token"
// @Param Email path string true "User Email"
// @Security BearerAuth
// @securityDefinitions.basic BearerAuth
// @Success		200	{object}	dbSchema.TblCart
// @Router			/api/cart/GetCartByUserEmail/{EmailAddress} [get]
func GetCartByUserEmail(ctx *gin.Context) {
	if !ValidateClient(ctx) {
		return
	}
	email := (ctx.Param("EmailAddress"))
	cart := &dbSchema.TblCart{}
	if getUSer := usww.GetUserByEmailAddress(email); getUSer.FirstName != "" {
		if getCart := usww.GetCartByUserId(getUSer.Id); getCart.CartName != "" {
			ctx.JSON(http.StatusOK, getCart)
		} else {
			ctx.JSON(http.StatusBadRequest, cart)
		}
	}
}
