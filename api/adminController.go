package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	dbSchema "github.com/Adebusy/cartbackendsvc/dataaccess"
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
// @Success		200	{object}	string
// @Router			/api/admin/CreateTitle [post]
func CreateTitle(ctx *gin.Context) {
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
// @Success		200	{object}  []dbSchema.TitleResp
// @Router			/api/admin/GetTitles [get]
func GetTitles(ctx *gin.Context) {
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
	ctx.JSON(http.StatusOK, usww.GetAllStatus())
}
