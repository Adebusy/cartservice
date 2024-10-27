package postgresql

import (
	"fmt"

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
	// env := os.Getenv("ENVIRONMENT")
	// env := "live"
	// SERVER := os.Getenv("DATABASE_SERVER" + "_" + env)
	// USERID := os.Getenv("USERID" + "_" + env)
	// DATABASE := os.Getenv("DATABASE" + "_" + env)
	// PASSWORD := os.Getenv("PASSWORD" + "_" + env)
	// PORT := os.Getenv("PORT" + "_" + env)

	// SERVER := "localhost"
	// PASSWORD := "Password1"
	// DATABASE := "DigitalCartDB"
	// USERID := "postgres"
	// PORT := "5432"

	SERVER := "my-db-postgresql-nyc3-62498-do-user-17863435-0.m.db.ondigitalocean.com"
	PASSWORD := "AVNS_4p8LzBbUn5iE6NeHLQP"
	DATABASE := "cartbackeddb"
	USERID := "cartusr"
	PORT := "25060"

	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=require", USERID, PASSWORD, SERVER, PORT, DATABASE)
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
