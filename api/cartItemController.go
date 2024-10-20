package api

import (
	"net/http"
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
// @Success		200	{string}	string "Item added to cart successfully!!"
// @Failure		400		{string} string	"Unable to add item to cart at the monent!!"
// @Router			/api/cart/AddItemToCart [post]
func AddItemToCart(ctx *gin.Context) {

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

	// check CartTypeId
	if docheckCartId := crt.GetCartByCartId(cartItemObj.CartId); docheckCartId.CartName == "" {
		ctx.JSON(http.StatusBadRequest, "CartId does not exist.")
		return
	}

	//check user has been added to cart team cartUser
	if docheckcartdetail := crt.GetCartByCartIdAndMemberId(cartItemObj.CartId, cartItemObj.UserId); docheckcartdetail.CartName == "" {
		ctx.JSON(http.StatusBadRequest, "this user has not been added to this cart group.")
		return
	}

	crts := dbSchema.TblCartItem{
		UserId:      cartItemObj.UserId,
		CartId:      cartItemObj.CartId,
		ProductId:   cartItemObj.ProductId,
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
// @Success		200	{string}	string "Item removed from cart successfully!!"
// @Failure		400		{string} string	"Unable to remove item to cart at the monent!!"
// @Router			/api/cart/RemoveItemFromCart [post]
func RemoveItemFromCart(ctx *gin.Context) {

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

	//do check if user can delete Item from cart
	if docheckRemovalPower := crt.GetCartDetailsByCartIdandMastersId(RemoveCartItemObj.CartId, doCheckCreatedById.EmailAddress); docheckRemovalPower.RingMasterEmail == "" {
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
