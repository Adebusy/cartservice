package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mssql"

	api "github.com/Adebusy/cartbackendsvc/api"

	"github.com/Adebusy/cartbackendsvc/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/sirupsen/logrus"

	_ "github.com/gofiber/swagger"

	"github.com/gin-contrib/cors"
)

func OptionMessage(c *gin.Context) {
	//c.Header("Access-Control-Allow-Origin", "http://localhost:8080")
	c.Header("Access-Control-Allow-Origin", "https://jellyfish-app-gz2qc.ondigitalocean.app")
	c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT")
	// @host			https://jellyfish-app-gz2qc.ondigitalocean.app
	// @host			localhost:8080
}

// @title			Cart Backend service
// @version		1.0
// @description	This service is meant to manage Cart request.
// @termsOfService	http://swagger.io/terms/
// @contact.name	Alao Adebisi
// @contact.email	alao.adebusy@gmail.com
// @license.name	Cart Manager Concept
// @license.url	https://github.com/MartinHeinz/go-project-blueprint/blob/master/LICENSE
// @host			https://jellyfish-app-gz2qc.ondigitalocean.app
// @BasePath		/
// @schemes		http
func main() {

	docs.SwaggerInfo.Title = "Cart Backend service API"
	docs.SwaggerInfo.Description = "This service is meant to manage Cart request"
	docs.SwaggerInfo.Version = "1.0"
	//docs.SwaggerInfo.Host = "localhost" + ":8080"
	docs.SwaggerInfo.Host = "jellyfish-app-gz2qc.ondigitalocean.app"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"https"}

	svc := gin.Default()
	svc.Use(cors.Default())
	svc.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"https://jellyfish-app-gz2qc.ondigitalocean.app"}, // List of allowed origins
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},
	}))
	//url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	url := ginSwagger.URL("https://jellyfish-app-gz2qc.ondigitalocean.app/swagger/doc.json")

	svc.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	svc.GET("/health", testSvc)
	svc.POST("/api/user/CreateNewUser", api.CreateNewUser)                             //done
	svc.GET("api/user/GetUserByEmailAddress/:EmailAddress", api.GetUserByEmailAddress) //done
	svc.GET("api/user/LogIn/:UserName/:Password", api.LogIn)                           //done
	svc.GET("api/user/GetUserByMobile/:MobileNumber", api.GetUserByMobile)             //done
	//svc.POST("api/user/UploadUserPicture", api.TestSevc)                               // to be done later
	svc.POST("api/user/SendEmail", api.SendEmail) //done
	svc.POST("api/user/LogOutUser", testSvc)      //to be done later

	svc.POST("api/cart/CreateCart", api.CreateCart)                 //done
	svc.POST("api/cart/AddUserToCart", api.CreateCartMember)        //done
	svc.POST("api/cart/AddItemToCart", api.AddItemToCart)           //done
	svc.POST("api/cart/CreateCartType", testSvc)                    //to be done later
	svc.POST("api/cart/RemoveItemFromCart", api.RemoveItemFromCart) //done
	svc.PUT("api/cart/CloseCart", api.CloseCart)                    //done
	svc.POST("api/cart/RemoveUserFromCart", api.RemoveUserFromCart) //done

	svc.POST("api/admin/CreateTitle", api.CreateTitle)       //done
	svc.GET("api/admin/GetTitles", api.GetTitles)            //done
	svc.GET("api/admin/GetAllStatus", api.GetAllStatus)      //done
	svc.GET("/api/admin/GetToken/:clientname", api.GetToken) //done
	svc.POST("api/admin/RegisterNewClient", api.RegisterNewClient)

	svc.POST("api/group/CreateGroupType", testSvc)
	svc.POST("api/group/DeleteGroupType", testSvc)

	svc.POST("api/admin/CreateRole", testSvc) //test service
	svc.POST("api/admin/DeleteRole", testSvc)
	svc.GET("api/admin/GetAllRoles", testSvc)
	svc.GET("api/admin/GetRoleByRoleId", testSvc)
	svc.POST("api/admin/CreateProduct", testSvc)
	svc.GET("api/admin/GetProductById", testSvc)
	svc.GET("api/admin/GetProducts", testSvc)
	svc.POST("api/admin/DeleteProduct", testSvc)
	svc.POST("api/admin/CreateStatus", testSvc)
	svc.DELETE("api/admin/DeleteStatus", testSvc)
	svc.Run(":8080")
}

// testSvc test service
// @Summary		Show the status of server.
// @Description	get the status of server.
// @Tags			root
// @Accept			*/*
// @User			json
// @Success		200	{object}	map[string]interface{}
// @Router			/testSvc [get]
func testSvc(ctx *gin.Context) {
	logrus.Info("testing service")
	fmt.Println("i am running")
	ctx.JSON(http.StatusOK, "good to go")
}
