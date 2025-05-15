package api

import (
	"net/http"
	"strconv"
	"time"

	dbSchema "github.com/Adebusy/cartbackendsvc/dataaccess"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// AddItemToCart godoc
// @Summary		Add item to an existing cart.
// @Description	Add item to an existing cart.
// @Tags			cart
// @Accept			*/*
// @User			json
// @Param user body dbSchema.CartItemObj true "Add item to cart"
// @Param Authorization header string true "Authorization token"
// @Param clientName header string true "registered client name"
// @Security BearerAuth
// @securityDefinitions.basic BearerAuth
// @Success		200	{string}	string "Item added to cart successfully!!"
// @Failure		400		{string} string	"Unable to add item to cart at the monent!!"
// @Router			/api/cart/AddItemToCart [post]
func AddItemToCart(ctx *gin.Context) {
	if !ValidateClient(ctx) {
		return
	}
	cartItemObj := &dbSchema.CartItemObj{}
	if doConvert := ctx.ShouldBindJSON(cartItemObj); doConvert != nil {
		ctx.JSON(http.StatusBadRequest, doConvert)
		return
	}

	//do check
	if docheck := validateMe.Struct(cartItemObj); docheck != nil {
		ctx.JSON(http.StatusBadRequest, docheck)
		return
	}
	//do check userID
	if doCheckCreatedById := usww.GetUserByUserId(cartItemObj.UserId); doCheckCreatedById.EmailAddress == "" {
		ctx.JSON(http.StatusBadRequest, "This user has not been created. UserId does not exist")
		return
	}

	// check CartTypeId ( crt.GetCartByCartId)
	if docheckCartId := usww.GetCartByCartId(cartItemObj.CartId); docheckCartId.CartName == "" {
		ctx.JSON(http.StatusBadRequest, "CartId does not exist.")
		return
	}

	//check user has been added to cart team cartUser (crt.GetCartByCartIdAndMemberId)
	if docheckcartdetail := usww.GetCartByCartIdAndMemberId(cartItemObj.CartId, cartItemObj.UserId); docheckcartdetail.CartName == "" {
		ctx.JSON(http.StatusBadRequest, "this user has not been added to this cart group.")
		return
	}

	crts := dbSchema.TblCartItem{
		UserId:      cartItemObj.UserId,
		CartId:      cartItemObj.CartId,
		Name:        cartItemObj.Name,
		Description: cartItemObj.Description,
		Quantity:    cartItemObj.Quantity,
		DateAdded:   time.Now(),
	}

	if doCreate := crtItem.AddItemToCart(crts); doCreate > 0 {
		ctx.JSON(http.StatusOK, "Item added to cart successfully!!")
		return
	} else {
		ctx.JSON(http.StatusBadRequest, "Unable to add item to cart at the monent!!")
		return
	}
}

// RemoveItemFromCart godoc
// @Summary		Remove item  from existing cart.
// @Description	This action can only be performed by the cart master.
// @Tags			cart
// @Accept			*/*
// @User			json
// @Param user body dbSchema.RemoveCartItemObj true "Remove item from cart"
// @Param Authorization header string true "Authorization token"
// @Param clientName header string true "registered client name"
// @Security BearerAuth
// @securityDefinitions.basic BearerAuth
// @Success		200	{string}	string "Item removed from cart successfully!!"
// @Failure		400		{string} string	"Unable to remove item to cart at the monent!!"
// @Router			/api/cart/RemoveItemFromCart [post]
func RemoveItemFromCart(ctx *gin.Context) {
	// if !ValidateClient(ctx) {
	// 	return
	// }
	RemoveCartItemObj := &dbSchema.RemoveCartItemObj{}
	if doConvert := ctx.ShouldBindJSON(RemoveCartItemObj); doConvert != nil {
		ctx.JSON(http.StatusBadRequest, doConvert)
		return
	}

	//do check
	if docheck := validateMe.Struct(RemoveCartItemObj); docheck != nil {
		ctx.JSON(http.StatusBadRequest, docheck)
		return
	}

	// check Check product to be removed
	if GetProductByProductId := prd.GetProductByProductId(RemoveCartItemObj.ProductId); GetProductByProductId.ProductName == "" {
		ctx.JSON(http.StatusBadRequest, "Product does not exist.")
		return
	}

	doCheckCreatedById := usww.GetUserByUserId(RemoveCartItemObj.UserId)
	if doCheckCreatedById.EmailAddress == "" {
		ctx.JSON(http.StatusBadRequest, "This user has not been created. UserId does not exist")
		return
	}

	//do check if user can delete Item from cart (crt.GetCartDetailsByCartIdandMastersId)
	if docheckRemovalPower := usww.GetCartDetailsByCartIdandMastersId(RemoveCartItemObj.CartId, doCheckCreatedById.EmailAddress); docheckRemovalPower.RingMasterEmail == "" {
		ctx.JSON(http.StatusBadRequest, "This user does not have the permission required to execute this action.")
		return
	}

	//do remove
	if doDeleteProduct := crtItem.RemoveItemFromCart(RemoveCartItemObj.ProductId, RemoveCartItemObj.CartId, RemoveCartItemObj.UserId); doDeleteProduct != nil {
		logrus.Error(doDeleteProduct)
		ctx.JSON(http.StatusBadRequest, "Service is unable to remove this product at the moment, Please try again later!!")
		return
	}
	ctx.JSON(http.StatusOK, "Product delete successfully !!")
}

// GetCartItemsByUserId Get Cart Item By User Id
// @Summary		Get Cart By User Id.
// @Description	Get Cart By User Id.
// @Tags			cart
// @Produce json
// @Accept			*/*
// @User			json
// @Param Authorization header string true "Authorization token"
// @Param clientName header string true "registered client name"
// @Param UserId path string true "User Id"
// @Security BearerAuth
// @securityDefinitions.basic BearerAuth
// @Success		200	{object}	[]dbSchema.TblCartItem
// @Router			/api/cart/GetCartItemsByUserId/{UserId} [get]
func GetCartItemsByUserId(ctx *gin.Context) {
	if !ValidateClient(ctx) {
		return
	}
	userId, _ := strconv.Atoi(ctx.Param("UserId"))
	ctx.JSON(http.StatusOK, usww.GetCartItemsByUserId(userId))
}

// GetCartItemsByCartId Get Cart Item By Cart Id
// @Summary		Get Cart By Cart Id.
// @Description	Get Cart By Cart Id.
// @Tags			cart
// @Produce json
// @Accept			*/*
// @User			json
// @Param Authorization header string true "Authorization token"
// @Param clientName header string true "registered client name"
// @Param CartId path string true "Cart ID"
// @Security BearerAuth
// @securityDefinitions.basic BearerAuth
// @Success		200	{object}	[]dbSchema.TblCartItem
// @Router			/api/cart/GetCartItemsByCartId/{CartId} [get]
func GetCartItemsByCartId(ctx *gin.Context) {
	if !ValidateClient(ctx) {
		return
	}
	cartId, _ := strconv.Atoi(ctx.Param("CartId"))
	ctx.JSON(http.StatusOK, usww.GetCartItemsByCartId(cartId))
}
