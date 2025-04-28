package api

import (
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"time"

	dbSchema "github.com/Adebusy/cartbackendsvc/dataaccess"
	inpuschema "github.com/Adebusy/cartbackendsvc/obj"
	psg "github.com/Adebusy/cartbackendsvc/postgresql"
	"github.com/Adebusy/cartbackendsvc/utilities"
	"github.com/EDDYCJY/go-gin-example/pkg/upload"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var (
	getdb      = psg.GetDB()
	usww       = dbSchema.ConneectDeal(getdb)
	crtItem    = dbSchema.ConnectCartItem(getdb)
	prd        = dbSchema.ConnectProduct(getdb)
	tit        = dbSchema.ConTitle(getdb)
	client     = dbSchema.ConnectClient(getdb)
	grp        = dbSchema.ConnectGroup(getdb)
	validateMe = validator.New()
)

// SignUp godoc
// @Summary		SignUp new user cart user.
// @Description	SignUp new user cart user.
// @Tags			user
// @Accept			*/*
// @User			json
// @Param user body inpuschema.SignUp true "SignUp new user"
// @Success		200	{object}	dbSchema.User
// @Router			/api/user/SignUp [post]
func SignUp(ctx *gin.Context) {
	currentTime := time.Now()
	usww := dbSchema.ConneectDeal(psg.GetDB())
	reqIn := &inpuschema.SignUp{}
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

	enc := hex.EncodeToString([]byte(reqIn.Password))

	if CheckEmailExist := usww.GetUserByEmailAddress(reqIn.Email); CheckEmailExist.UserName != "" {
		ctx.JSON(http.StatusBadRequest, "User with email address "+reqIn.Email+" already exist!!")
		return
	}

	if CheckMobile := usww.GetUserByMobileNumber(reqIn.MobileNumber); CheckMobile.UserName != "" {
		ctx.JSON(http.StatusBadRequest, "User with mobile number "+reqIn.MobileNumber+" already exist!!")
		return
	}

	doCreate := usww.SignUp(reqIn.Email, reqIn.MobileNumber, enc, currentTime.Format("01-02-2006"))
	logrus.Info(doCreate)
	ctx.JSON(http.StatusOK, doCreate)
}

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
		CreatedAt: time.Now().Format("01-02-2006")}

	doCreate := usww.CreateUser(req)
	logrus.Info(doCreate)
	ctx.JSON(http.StatusOK, doCreate)
}

// CompleteSignUp godoc
// @Summary		CompleteSignUp user signup.
// @Description	CompleteSignUp user signup.
// @Tags			user
// @Accept			*/*
// @User			json
// @Param Authorization header string true "Authorization token"
// @Param clientName header string true "registered client name"
// @Security BearerAuth
// @securityDefinitions.basic BearerAuth
// @Param user body inpuschema.CompleteSignUp true "CompleteSignUp user signup"
// @Success		200	{object}	dbSchema.ResponseMessage
// @Router			/api/user/CompleteSignUp [post]
func CompleteSignUp(ctx *gin.Context) {
	if !ValidateClient(ctx) {
		return
	}
	usww := dbSchema.ConneectDeal(psg.GetDB())
	reqIn := &inpuschema.CompleteSignUp{}
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

	req := dbSchema.CompleteSignUpReq{EmailAddress: reqIn.EmailAddress, TitleId: strconv.Itoa(reqIn.TitleId),
		FirstName: reqIn.FirstName,
		UserName:  reqIn.UserName, NickName: reqIn.NickName,
		LastName:     reqIn.LastName,
		Gender:       reqIn.Gender,
		AgeRange:     reqIn.AgeRange,
		Status:       1,
		MobileNumber: reqIn.MobileNumber,
		CreatedAt:    time.Now().Format("01-02-2006")}

	if CheckEmailExist := usww.GetUserByEmailAddress(req.EmailAddress); CheckEmailExist.EmailAddress == "" {
		ctx.JSON(http.StatusBadRequest, "User with email address "+req.EmailAddress+" does not exist!!")
		return
	}

	if CheckMobile := usww.GetUserByMobileNumber(req.MobileNumber); CheckMobile.MobileNumber == "" {
		ctx.JSON(http.StatusBadRequest, "User with mobile number "+req.MobileNumber+" does not exist!!")
		return
	}

	doCreate := usww.UpdateUserRecord(req)
	logrus.Info(doCreate)
	Response := &inpuschema.ResponseMessage{ResponseCode: "00",
		ResponseMessage: doCreate,
	}

	ctx.JSON(http.StatusOK, Response)
}

// GetUserByEmailAddress create new user
// @Summary		Get user by email address new cart user.
// @Description	Get user by email address new cart user.
// @Tags			user
// @Param EmailAddress path string true "User email address"
// @Produce json
// @Accept			*/*
// @User			json
// @Param Authorization header string true "Authorization token"
// @Param clientName header string true "registered client name"
// @Security BearerAuth
// @securityDefinitions.basic BearerAuth
// @Success		200	{object}	inpuschema.UserResponse
// @Router			/api/user/GetUserByEmailAddress/{EmailAddress} [get]
func GetUserByEmailAddress(ctx *gin.Context) {
	if !ValidateClient(ctx) {
		return
	}
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
// @Param Authorization header string true "Authorization token"
// @Param clientName header string true "registered client name"
// @Security BearerAuth
// @securityDefinitions.basic BearerAuth
// @Success		200	{object}	inpuschema.CartObj
// @Router			/api/user/GetUserByMobile/{MobileNumber} [get]
func GetUserByMobile(ctx *gin.Context) {
	if !ValidateClient(ctx) {
		return
	}
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
// @Produce json
// @Accept			*/*
// @User			json
// @Success		200	{object}	inpuschema.UserResponse
// @Router			/api/user/LogIn/{UserName}/{Password} [get]
func LogIn(ctx *gin.Context) {
	getUSerobj := dbSchema.User{}
	userRespose := &inpuschema.UserResponse{}
	UserName := ctx.Param("UserName")
	Password := ctx.Param("Password")
	enc := hex.EncodeToString([]byte(Password))

	if utilities.IsEmailValid(UserName) {
		getUSerobj = usww.GetUserByEmailAddress(UserName)
	} else if utilities.IsNumberValid(UserName) {
		getUSerobj = usww.GetUserByMobileNumber(UserName)
	} else {
		logAction := fmt.Sprintf("Incorrect username %s", UserName)
		logrus.Info(logAction)
		ctx.JSON(http.StatusBadRequest, logAction)
		return
	}

	if getUSerobj.EmailAddress != "" || getUSerobj.MobileNumber != "" {
		if getUSerobj.Password == enc {
			userRespose.TitleId = getUSerobj.TitleId
			userRespose.UserName = getUSerobj.UserName
			userRespose.NickName = getUSerobj.NickName
			userRespose.FirstName = getUSerobj.FirstName
			userRespose.LastName = getUSerobj.LastName
			userRespose.Email = getUSerobj.EmailAddress
			userRespose.MobileNumber = getUSerobj.MobileNumber
			userRespose.Status = getUSerobj.Status
			userRespose.Gender = getUSerobj.Gender
			userRespose.Location = getUSerobj.Location
			userRespose.CreatedAt = getUSerobj.CreatedAt
			userRespose.Id = uint(getUSerobj.Id)
			if newToken := CreateOrGetToken(userRespose.Email); newToken != "" {
				userRespose.Token = newToken
			}
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

// LogInWithMobileNumber for exiting user
// @Summary		Log user In with mobile number and password.
// @Description	Log user In with mobile number and password.
// @Tags			user
// @Param MobileNumber path string true "MobileNumber"
// @Param Password path string true "Password"
// @Produce json
// @Accept			*/*
// @User			json
// @Success		200	{object}	inpuschema.UserResponse
// @Router			/api/user/LogInWithMobileNumber/{MobileNumber}/{Password} [get]
func LogInWithMobileNumber(ctx *gin.Context) {
	// @Param Authorization header string true "Authorization token"
	// @Param clientName header string true "registered client name"
	// @Security BearerAuth
	// @securityDefinitions.basic BearerAuth
	// if !ValidateClient(ctx) {
	// 	return
	// }
	userRespose := &inpuschema.UserResponse{}
	MobileNumber := ctx.Param("MobileNumber")
	Password := ctx.Param("Password")
	password, _ := utilities.HashPassword(Password)

	if getUSer := usww.GetUserByMobileNumber(MobileNumber); getUSer.EmailAddress != "" {
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
			logrus.Info(fmt.Sprintf("LogIn for user with LogInWithMobileNumber %s", MobileNumber))
			ctx.JSON(http.StatusOK, userRespose)
			return
		} else {
			logAction := fmt.Sprintf("Incorrect password %s", MobileNumber)
			logrus.Info(logAction)
			ctx.JSON(http.StatusBadRequest, logAction)
			return
		}
	} else {
		logAction := fmt.Sprintf("Incorrect username %s", MobileNumber)
		logrus.Info(logAction)
		ctx.JSON(http.StatusBadRequest, logAction)
		return
	}
}

// LogInWithEmailAddress for exiting user
// @Summary		Log user In with email address and password.
// @Description	Log user In with email address and password.
// @Tags			user
// @Param EmailAddress path string true "EmailAddress"
// @Param Password path string true "Password"
// @Produce json
// @Accept			*/*
// @User			json
// @Success		200	{object}	inpuschema.UserResponse
// @Router			/api/user/LogInWithEmailAddress/{EmailAddress}/{Password} [get]
func LogInWithEmailAddress(ctx *gin.Context) {
	userRespose := &inpuschema.UserResponse{}
	EmailAddress := ctx.Param("EmailAddress")
	Password := ctx.Param("Password")
	password, _ := utilities.HashPassword(Password)

	if getUSer := usww.GetUserByEmailAddress(EmailAddress); getUSer.EmailAddress != "" {
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
			logrus.Info(fmt.Sprintf("LogIn for user with LogInWithEmailAddress %s", EmailAddress))
			ctx.JSON(http.StatusOK, userRespose)
			return
		} else {
			logAction := fmt.Sprintf("Incorrect password %s", EmailAddress)
			logrus.Info(logAction)
			ctx.JSON(http.StatusBadRequest, logAction)
			return
		}
	} else {
		logAction := fmt.Sprintf("Incorrect username %s", EmailAddress)
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
// @Param Authorization header string true "Authorization token"
// @Param clientName header string true "registered client name"
// @Security BearerAuth
// @securityDefinitions.basic BearerAuth
// @Success		200	{string}	string "Email sent successfully!!"
// @Failure		400		{string} string	"Unable to send email at the monent!!"
// @Router			/api/user/SendEmail [post]
func SendEmail(ctx *gin.Context) {
	if !ValidateClient(ctx) {
		return
	}
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

// // LogIn exiting user In
// // @Summary		Log user In with username and password.
// // @Description	Log user In with username and password.
// // @Tags			user
// // @Param UserName path string true "Username"
// // @Param Password path string true "Password"
// // @Produce json
// // @Accept			*/*
// // @User			json
// // @Success		200	{object}	inpuschema.UserResponse
// // @Router			/api/user/LogIn/{UserName}/{Password} [get]
// func LogIn(ctx *gin.Context) {
// 	// @Param Authorization header string true "Authorization token"
// 	// @Param clientName header string true "registered client name"
// 	// @Security BearerAuth
// 	// @securityDefinitions.basic BearerAuth
// 	// if !ValidateClient(ctx) {
// 	// 	return
// 	// }

// 	var getUSer dbSchema.User

// 	userRespose := &inpuschema.UserResponse{}
// 	UserName := ctx.Param("UserName")
// 	Password := ctx.Param("Password")
// 	password, _ := utilities.HashPassword(Password)

// 	if utilities.IsEmailValid(UserName) {

// 	}

// 	if getUSer := usww.GetUserByUsername(UserName); getUSer.EmailAddress != "" {
// 		if utilities.CheckPasswordHash(Password, password) {
// 			userRespose.TitleId = getUSer.TitleId
// 			userRespose.UserName = getUSer.UserName
// 			userRespose.NickName = getUSer.NickName
// 			userRespose.FirstName = getUSer.FirstName
// 			userRespose.LastName = getUSer.LastName
// 			userRespose.Email = getUSer.EmailAddress
// 			userRespose.MobileNumber = getUSer.MobileNumber
// 			userRespose.Status = getUSer.Status
// 			userRespose.Gender = getUSer.Gender
// 			userRespose.Location = getUSer.Location
// 			userRespose.CreatedAt = getUSer.CreatedAt
// 			if token, err := utilities.CreateToken(UserName); err.Error() == "" {
// 				userRespose.Token = token
// 			}

// 			logrus.Info(fmt.Sprintf("LogIn for user %s", UserName))
// 			ctx.JSON(http.StatusOK, userRespose)
// 			return
// 		} else {
// 			logAction := fmt.Sprintf("Incorrect password %s", UserName)
// 			logrus.Info(logAction)
// 			ctx.JSON(http.StatusBadRequest, logAction)
// 			return
// 		}
// 	} else {
// 		logAction := fmt.Sprintf("Incorrect username %s", UserName)
// 		logrus.Info(logAction)
// 		ctx.JSON(http.StatusBadRequest, logAction)
// 		return
// 	}
// }

// @Summary Import Image
// @Produce  json
// @Param image formData file true "Image File"
// @Success 200 {object} string
// @Failure 500 {object} string
// @Router /api/user/UploadImage [post]
func UploadImage(ctx *gin.Context) {
	file, image, err := ctx.Request.FormFile("image")
	if err != nil {
		logrus.Warn(err)
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if image == nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	imageName := upload.GetImageName(image.Filename)
	fullPath := upload.GetImageFullPath()
	//savePath := upload.GetImagePath()
	// src := fullPath + imageName

	if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(file) {
		ctx.JSON(http.StatusBadRequest, "ERROR_UPLOAD_CHECK_IMAGE_FORMAT")
		return
	}

	err = upload.CheckImage(fullPath)
	if err != nil {
		logrus.Warn(err)
		ctx.JSON(http.StatusInternalServerError, "ERROR_UPLOAD_CHECK_IMAGE_FAIL")
		return
	}

	// if err := c.SaveUploadedFile(image, src); err != nil {
	// 	logrus.Warn(err)
	// 	ctx.JSON(http.StatusInternalServerError, "ERROR_UPLOAD_SAVE_IMAGE_FAIL")
	// 	return
	// }

	// appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
	// 	"image_url":      upload.GetImageFullUrl(imageName),
	// 	"image_save_url": savePath + imageName,
	// })
}
