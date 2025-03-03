package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	dbSchema "github.com/Adebusy/cartbackendsvc/dataaccess"
	"github.com/Adebusy/cartbackendsvc/obj"
	"github.com/Adebusy/cartbackendsvc/utilities"
)

type TitleObj struct {
	Name   string
	Status int
}

// CreateTitle godoc
// @Summary		Create new Title.
// @Description	Create new Title.
// @Tags			admin
// @Accept			*/*
// @User			json
// @Param user body TitleObj true "Create new title"
// @Param Authorization header string true "Authorization token"
// @Param clientName header string true "registered client name"
// @Security BearerAuth
// @securityDefinitions.basic BearerAuth
// @Success		200	{object}	string
// @Router			/api/admin/CreateTitle [post]
func CreateTitle(ctx *gin.Context) {
	// @Param Authorization header string true "Authorization token"
	// @Param clientName header string true "registered client name"
	// // @Security BearerAuth
	// // @securityDefinitions.basic BearerAuth
	if !ValidateClient(ctx) {
		return
	}
	reqBearer := ctx.GetHeader("Authorization")
	if reqBearer == "" {
		resp := fmt.Sprintf("Bearer is required!! %s", reqBearer)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	reqBearer = reqBearer[len("Bearer "):]

	if doVerify := utilities.VerifyToken(reqBearer); doVerify != nil {
		ctx.JSON(http.StatusBadRequest, "invalid token")
		return
	}

	//validate request body
	title := &TitleObj{}

	if docheck := ctx.ShouldBindJSON(title); docheck != nil {
		ctx.JSON(http.StatusBadRequest, docheck)
	}
	titleobj := dbSchema.TblTitle{
		Name:      title.Name,
		Status:    true,
		CreatedAt: time.Now(),
	}

	if checkTitle := tit.GetTitleByTitleName(titleobj.Name); checkTitle.Name != "" {
		ctx.JSON(http.StatusBadRequest, "Title already exist!!")
		return
	}

	//insert into title name al
	if docreate := tit.CreateTitle(titleobj); docreate == 0 {
		ctx.JSON(http.StatusBadRequest, "unbale to create tiele at the moment.")
		return
	}
	ctx.JSON(http.StatusOK, "Title created successfully!!")
}

// GetTitles godoc
// @Summary		Get all Titles.
// @Description	Get all Title.
// @Tags			admin
// @Accept			*/*
// @User			json
// @Param Authorization header string true "Authorization token"
// @Param clientName header string true "registered client name"
// @Security BearerAuth
// @securityDefinitions.basic BearerAuth
// @Success		200	{object}  []dbSchema.TitleResp
// @Router			/api/admin/GetTitles [get]
func GetTitles(ctx *gin.Context) {
	// @Param Authorization header string true "Authorization token"
	// @Param clientName header string true "registered client name"
	// @Security BearerAuth
	// @securityDefinitions.basic BearerAuth
	if !ValidateClient(ctx) {
		return
	}
	ctx.JSON(http.StatusOK, tit.GetTitles())
}

// GetAllStatus godoc
// @Summary		Get all Status.
// @Description	Get all Status.
// @Tags			admin
// @Accept			*/*
// @User			json
// @Success		200	{object}  []dbSchema.TblStatus
// @Router			/api/admin/GetAllStatus [get]
func GetAllStatus(ctx *gin.Context) {
	// if !ValidateClient(ctx) {
	// 	return
	// }
	ctx.JSON(http.StatusOK, usww.GetAllStatus())
}

// GetToken godoc
// @Summary		Get Token for client.
// @Description	Get Token for client.
// @Param clientname path string true "Registered client name"
// @Tags			admin
// @Accept			*/*
// @User			json
// @Success		200	{object}  obj.TokenResp
// @Router			/api/admin/GetToken/{clientname} [get]
func GetToken(ctx *gin.Context) {
	obj := obj.TokenResp{}
	username := ctx.Param("clientname")
	// if checkClient := client.GetClientByName(username); checkClient.Name == "" {
	// 	ctx.JSON(http.StatusBadRequest, "This client has not been onboarded, please register client and try again.")
	// 	return
	// }

	if respToken, err := utilities.CreateToken(username); err != nil {
		ctx.JSON(http.StatusBadRequest, "unable to get token at the moment.")
		return
	} else {
		obj.Token = respToken
		ctx.JSON(http.StatusOK, obj)
	}
}

// RegisterNewClient godoc
// @Summary		Register New Client.
// @Description	Register New Client.
// @Tags			admin
// @Accept			*/*
// @User			json
// @Param user body dbSchema.ClientRequest true "Create new client"
// @Success		200	{object}	dbSchema.ClientResp
// @Router			/api/admin/RegisterNewClient [post]
func RegisterNewClient(ctx *gin.Context) {
	obj := obj.TokenResp{}
	req := &dbSchema.ClientRequest{}
	if docheck := ctx.ShouldBindJSON(req); docheck != nil {
		ctx.JSON(http.StatusBadRequest, docheck)
	}

	reqClient := dbSchema.TblClient{Name: req.Name,
		Status: 1, Description: "new Client", DateAdded: &time.Time{}}

	if checkClient := client.GetClientByName(req.Name); checkClient.Name != "" {
		if respToken, err := utilities.CreateToken(reqClient.Name); err != nil {
			ctx.JSON(http.StatusBadRequest, "This client already exists but service is unable to generate token at the moment.")
			return
		} else {
			obj.Token = respToken
			ctx.JSON(http.StatusOK, obj)
		}
	} else {
		if doReg := client.RegisterNewClient(reqClient); doReg == "00" {
			if respToken, err := utilities.CreateToken(reqClient.Name); err != nil {
				ctx.JSON(http.StatusBadRequest, "Client is registered but service is unable to generate token at the moment.")
				return
			} else {
				obj.Token = respToken
				ctx.JSON(http.StatusOK, obj)
			}
		}
	}
}
