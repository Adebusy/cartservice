package api

import (
	"fmt"
	"net/http"

	"github.com/Adebusy/cartbackendsvc/utilities"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
	logrus.Info("reqBearer")
	logrus.Info(reqBearer)
	logrus.Info("reqBearer")
	if doVerify := utilities.VerifyToken(reqBearer); doVerify != nil {
		ctx.JSON(http.StatusBadRequest, "invalid token")
		return false
	}
	return true
}
