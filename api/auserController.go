package api

import (
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"

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
// @Param Authorization header string true "Authorization token"
// @Param clientName header string true "registered client name"
// @Success		200	{object}	dbSchema.User
// @Router			/api/user/CreateNewUser [post]
func CreateNewUser(ctx *gin.Context) {
	// @Security BearerAuth
	// @securityDefinitions.basic BearerAuth
	// if !ValidateClient(ctx) {
	// 	return
	// }
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
// @Param Authorization header string true "Authorization token"
// @Param clientName header string true "registered client name"
// @Produce json
// @Accept			*/*
// @User			json
// @Success		200	{object}	inpuschema.UserResponse
// @Router			/api/user/GetUserByEmailAddress/{EmailAddress} [get]
func GetUserByEmailAddress(ctx *gin.Context) {
	// @Security BearerAuth
	// @securityDefinitions.basic BearerAuth
	// if !ValidateClient(ctx) {
	// 	return
	// }
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
// @Param Authorization header string true "Authorization token"
// @Param clientName header string true "registered client name"
// @Produce json
// @Accept			*/*
// @User			json
// @Success		200	{object}	inpuschema.UserResponse
// @Router			/api/user/GetUserByMobile/{MobileNumber} [get]
func GetUserByMobile(ctx *gin.Context) {
	// @Security BearerAuth
	// @securityDefinitions.basic BearerAuth
	// if !ValidateClient(ctx) {
	// 	return
	// }
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
		userRespose.Gender = getUSer.Gender
		userRespose.Location = getUSer.Location
		userRespose.AgeRange = getUSer.AgeRange
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
// @Param Authorization header string true "Authorization token"
// @Param clientName header string true "registered client name"
// @Produce json
// @Accept			*/*
// @User			json
// @Success		200	{object}	inpuschema.UserResponse
// @Router			/api/user/LogIn/{UserName}/{Password} [get]
func LogIn(ctx *gin.Context) {
	// @Security BearerAuth
	// @securityDefinitions.basic BearerAuth
	// if !ValidateClient(ctx) {
	// 	return
	// }
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
			userRespose.Gender = getUSer.Gender
			userRespose.Location = getUSer.Location
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

// SendEmail godoc
// @Summary		Send Email.
// @Description	Send Email.
// @Tags			user
// @Accept			*/*
// @User			json
// @Param user body inpuschema.EmailObj true "Send Email"
// @Success		200	{string}	string "Email sent successfully!!"
// @Failure		400		{string} string	"Unable to send email at the monent!!"
// @Router			/api/user/SendEmail [post]
func SendEmail(ctx *gin.Context) {
	// if !ValidateClient(ctx) {
	// 	return
	// }
	// @Param Authorization header string true "Authorization token"
	// @Param clientName header string true "registered client name"
	// @Security BearerAuth
	// @securityDefinitions.basic BearerAuth
	reqIn := &inpuschema.EmailObj{}
	if err := ctx.ShouldBindJSON(reqIn); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		fmt.Println(err.Error())
		return
	}

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	sender := os.Getenv("SMTP_SENDER")
	recipient := reqIn.ToEmail

	from := "From: " + sender + "\n"
	to := "To: " + recipient + "\n"
	subject := "Subject: Digital cart update\n"
	body := reqIn.MailBody
	message := []byte(from + to + subject + "\n" + body)
	auth := smtp.PlainAuth("", username, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, sender, []string{recipient}, message)
	if err != nil {
		log.Fatalf("Failed to send email: %v", err)
		resp := fmt.Sprintf("Error sending email: %v", err)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	logAction := fmt.Sprintf("SendEmail to %v", recipient)
	logrus.Info(logAction)
	ctx.JSON(http.StatusOK, "Email sent successfully!!!")
}
