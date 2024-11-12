package api

import (
	"fmt"
	"net/http"

	dbSchema "github.com/Adebusy/cartbackendsvc/dataaccess"
	inpuschema "github.com/Adebusy/cartbackendsvc/obj"
	psg "github.com/Adebusy/cartbackendsvc/postgresql"
	"github.com/Adebusy/cartbackendsvc/utilities"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var (
	usww = dbSchema.ConneectDeal(psg.GetDB())
	//crt        = dbSchema.ConnectCart(psg.GetDB())
	crtItem    = dbSchema.ConnectCartItem(psg.GetDB())
	prd        = dbSchema.ConnectProduct(psg.GetDB())
	tit        = dbSchema.ConTitle(psg.GetDB())
	client     = dbSchema.ConnectClient(psg.GetDB())
	validateMe = validator.New()
)

// CreateNewUser godoc
// @Summary		Create new user cart user.
// @Description	Create new user cart user.
// @Tags			user
// @Accept			*/*
// @User			json
// @Param user body inpuschema.UserObj true "Create new user"
// @Success		200	{object}	dbSchema.User
// @Router			/api/user/CreateNewUser [post]
func CreateNewUser(ctx *gin.Context) {
	usww := dbSchema.ConneectDeal(psg.GetDB())
	reqIn := &inpuschema.UserObj{}
	if err := ctx.ShouldBindJSON(reqIn); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		fmt.Println(err.Error())
		return
	}

	//validate request
	if validateObj := validateMe.Struct(reqIn); validateObj != nil {
		ctx.JSON(http.StatusBadRequest, validateObj.Error())
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(reqIn.Password), 8)
	req := &dbSchema.User{TitleId: reqIn.TitleId, FirstName: reqIn.FirstName, UserName: reqIn.UserName, NickName: reqIn.NickName,
		LastName: reqIn.LastName, EmailAddress: reqIn.Email,
		MobileNumber: reqIn.MobileNumber, Status: reqIn.Status, Password: string(hashedPassword),
		CreatedAt: "2024-09-15"}

	if CheckEmailExist := usww.GetUserByEmailAddress(req.EmailAddress); CheckEmailExist.UserName != "" {
		ctx.JSON(http.StatusBadRequest, "User with email address "+req.EmailAddress+" already exist!!")
		return
	}

	if CheckMobile := usww.GetUserByMobileNumber(req.MobileNumber); CheckMobile.UserName != "" {
		ctx.JSON(http.StatusBadRequest, "User with mobile number "+req.MobileNumber+" already exist!!")
		return
	}

	doCreate := usww.CreateUser(req)
	logrus.Info(doCreate)
	ctx.JSON(http.StatusOK, doCreate)
}

// GetUserByEmailAddress create new user
// @Summary		Get user by email address new cart user.
// @Description	Get user by email address new cart user.
// @Tags			user
// @Param EmailAddress path string true "User email address"
// @Produce json
// @Accept			*/*
// @User			json
// @Success		200	{object}	inpuschema.UserResponse
// @Router			/api/user/GetUserByEmailAddress/{EmailAddress} [get]
func GetUserByEmailAddress(ctx *gin.Context) {
	requestEmail := ctx.Param("EmailAddress")
	getUSer := usww.GetUserByEmailAddress(requestEmail)
	logAction := fmt.Sprintf("GetUserByEmailAddress %v", requestEmail)
	logrus.Info(logAction)
	ctx.JSON(http.StatusOK, getUSer)
}

// GetUserByMobile existing user destails by mobile number
// @Summary		existing user destails by mobile number.
// @Description	existing user destails by mobile number.
// @Tags			user
// @Param MobileNumber path string true "User mobile number"
// @Produce json
// @Accept			*/*
// @User			json
// @Success		200	{object}	inpuschema.UserResponse
// @Router			/api/user/GetUserByMobile/{MobileNumber} [get]
func GetUserByMobile(ctx *gin.Context) {
	userRespose := &inpuschema.UserResponse{}
	requestMobile := ctx.Param("MobileNumber")
	if getUSer := usww.GetUserByMobileNumber(requestMobile); getUSer.FirstName != "" {
		userRespose.TitleId = getUSer.TitleId
		userRespose.UserName = getUSer.UserName
		userRespose.NickName = getUSer.NickName
		userRespose.FirstName = getUSer.FirstName
		userRespose.LastName = getUSer.LastName
		userRespose.Email = getUSer.EmailAddress
		userRespose.MobileNumber = getUSer.MobileNumber
		userRespose.Status = getUSer.Status
		userRespose.CreatedAt = getUSer.CreatedAt
		ctx.JSON(http.StatusOK, userRespose)
		return
	}

	logAction := fmt.Sprintf("GetUserByMobile %v", requestMobile)
	logrus.Info(logAction)
	ctx.JSON(http.StatusBadRequest, userRespose)
}

// LogIn exiting user In
// @Summary		Log user In with username and password.
// @Description	Log user In with username and password.
// @Tags			user
// @Param UserName path string true "Username"
// @Param Password path string true "Password"
// @Produce json
// @Accept			*/*
// @User			json
// @Success		200	{object}	inpuschema.UserResponse
// @Router			/api/user/LogIn/{UserName}/{Password} [get]
func LogIn(ctx *gin.Context) {
	userRespose := &inpuschema.UserResponse{}
	UserName := ctx.Param("UserName")
	Password := ctx.Param("Password")
	password, _ := utilities.HashPassword(Password)

	if getUSer := usww.GetUserByEmailUsername(UserName); getUSer.FirstName != "" {
		if utilities.CheckPasswordHash(Password, password) {
			userRespose.TitleId = getUSer.TitleId
			userRespose.UserName = getUSer.UserName
			userRespose.NickName = getUSer.NickName
			userRespose.FirstName = getUSer.FirstName
			userRespose.LastName = getUSer.LastName
			userRespose.Email = getUSer.EmailAddress
			userRespose.MobileNumber = getUSer.MobileNumber
			userRespose.Status = getUSer.Status
			userRespose.CreatedAt = getUSer.CreatedAt
			logrus.Info(fmt.Sprintf("LogIn for user %s", UserName))
			ctx.JSON(http.StatusOK, userRespose)
			return
		} else {
			logAction := fmt.Sprintf("Incorrect password %s", UserName)
			logrus.Info(logAction)
			ctx.JSON(http.StatusBadRequest, logAction)
			return
		}
	} else {
		logAction := fmt.Sprintf("Incorrect username %s", UserName)
		logrus.Info(logAction)
		ctx.JSON(http.StatusBadRequest, logAction)
		return
	}
}
