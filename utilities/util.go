package utilities

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func GoDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file" + err.Error())
	}
	return os.Getenv(key)
}

// HashPassword hashes a given password and returns the hashed password or an error
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedBytes), err
}

// CheckPasswordHash verifies the password against the hashed password and returns if it's correct or not
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

var secretKey = []byte("secret-key")

func CreateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

// isEmailValid checks if the email provided is valid by regex.
func IsEmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}

// isEmailValid checks if the email provided is valid by regex.
func IsNumberValid(e string) bool {
	var re = regexp.MustCompile(`^[0-9]+$`)
	if re.MatchString(e) {
		return true
	} else {
		return false
	}
}

// func ValidateClient(ctx *gin.Context) bool {
// 	reqBearer := ctx.GetHeader("Authorization")
// 	if reqBearer == "" {
// 		resp := fmt.Sprintf("Bearer is required!! %s", reqBearer)
// 		ctx.JSON(http.StatusBadRequest, resp)
// 		return false
// 	}

// 	clientName := ctx.GetHeader("clientName")
// 	if clientName == "" {
// 		resp := fmt.Sprintf("Client name is required in the header!! %s", clientName)
// 		ctx.JSON(http.StatusBadRequest, resp)
// 		return false
// 	}
// 	//check client

// 	docheck := dataaccess.GetClientByName(clientName)

// 	reqBearer = reqBearer[len("Bearer "):]
// 	if doVerify := VerifyToken(reqBearer); doVerify != nil {
// 		ctx.JSON(http.StatusBadRequest, "invalid token")
// 		return false
// 	}
// 	return true
// }
