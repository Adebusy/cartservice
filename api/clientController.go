package api

import (
	"fmt"
	"net/http"
	"time"

	dbSchema "github.com/Adebusy/cartbackendsvc/dataaccess"
	"github.com/Adebusy/cartbackendsvc/utilities"
	"github.com/gin-gonic/gin"
)

func ValidateClient(ctx *gin.Context) bool {
	reqBearer := ctx.GetHeader("Authorization")
	if reqBearer == "" {
		resp := fmt.Sprintf("Bearer is required!! %s", reqBearer)
		ctx.JSON(http.StatusBadRequest, resp)
		return false
	}

	clientName := ctx.GetHeader("clientName")
	if clientName == "" {
		resp := fmt.Sprintf("Client name is required in the header!! %s", clientName)
		ctx.JSON(http.StatusBadRequest, resp)
		return false
	}

	//check client
	if docheck := client.GetClientByName(clientName); docheck.Name == "" {
		resp := fmt.Sprintf("Client %s is not registered!!", clientName)
		ctx.JSON(http.StatusBadRequest, resp)
		return false
	}

	reqBearer = reqBearer[len("Bearer "):]
	if doVerify := utilities.VerifyToken(reqBearer); doVerify != nil {
		ctx.JSON(http.StatusBadRequest, "invalid token")
		return false
	}
	return true
}

func CreateOrGetToken(username string) string {
	var respToken string
	res := dbSchema.TblClient{
		Name:        username,
		Status:      1,
		Description: username,
		DateAdded:   &time.Time{},
	}

	return CreateOrUpdateToken(username, res, respToken)
}

func CreateOrUpdateToken(username string, res dbSchema.TblClient, respToken string) string {
	fmt.Printf("checkClient.Name is")
	if checkClient := client.GetClientByName(username); checkClient.Name == "" {
		fmt.Printf("checkClient.Name is %s", checkClient.Name)
		if doReg := client.RegisterNewClient(res); doReg == "00" {
			if respToken, err := utilities.CreateToken(username); err != nil {
				return respToken
			} else {
				return respToken
			}
		}
	} else {
		if respToken, err := utilities.CreateToken(username); err != nil {
			return respToken
		} else {
			return respToken
		}
	}
	return respToken
}
