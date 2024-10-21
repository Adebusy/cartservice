package postgresql

import (
	"fmt"
	"os"

	dbSchema "github.com/Adebusy/cartbackendsvc/dataaccess"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DbGorm *gorm.DB
var err error

func GetDB() *gorm.DB {
	if loadEnv := godotenv.Load(); loadEnv != nil {
		ret := fmt.Sprintf("Unable to load environment variable. %s", loadEnv.Error())
		fmt.Println(ret)
	}
	env := os.Getenv("ENVIRONMENT")
	SERVER := os.Getenv("DATABASE_SERVER" + "_" + env)
	USERID := os.Getenv("USERID" + "_" + env)
	DATABASE := os.Getenv("DATABASE" + "_" + env)
	PASSWORD := os.Getenv("PASSWORD" + "_" + env)
	PORT := os.Getenv("PORT" + "_" + env)
	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", SERVER, USERID, PASSWORD, DATABASE, PORT)
	DbGorm, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{NamingStrategy: schema.NamingStrategy{
		SingularTable: true, NoLowerCase: true,
	}})
	if err != nil {
		panic("failed to connect database")
	}
	DbGorm.AutoMigrate(&dbSchema.TblCart{})
	DbGorm.AutoMigrate(&dbSchema.TblTitle{})
	DbGorm.AutoMigrate(&dbSchema.TblCartItem{})
	DbGorm.AutoMigrate(&dbSchema.TblCartMember{})
	return DbGorm
}
