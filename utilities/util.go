package utilities

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/rand"
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

func DeactivateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Minute * 0).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func Logout(token, username string) {
	// ttl := time.Now().Add(time.Hour * 0).Unix()

	// if err != nil {
	// 	c.JSON(http.StatusUnauthorized, "unauthorized")
	// 	return
	// }

	// deleted, delErr := DeleteAuth(au.AccessUuid)

	// if delErr != nil || deleted == 0 {
	// 	c.JSON(http.StatusUnauthorized, "unauthorized")
	// 	return
	// }

	//c.JSON(http.StatusOK, "Successfully logged out")
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

func SendEmail(toEmail, mailBody string) string {

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	sender := os.Getenv("SMTP_SENDER")
	recipient := toEmail
	from := "From: " + sender + "\n"
	to := "To: " + recipient + "\n"
	subject := "Subject: Digital cart update\n"
	body := mailBody
	message := []byte(from + to + subject + "\n" + body)
	auth := smtp.PlainAuth("", username, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, sender, []string{recipient}, message)
	if err != nil {
		log.Fatalf("Failed to send email: %v", err.Error())
		return "01"
	} else {
		return "00"
	}
}

const (
	letterBytes  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialBytes = "!@#$%^&*()_+-=[]{}\\|;':\",.<>/?`~"
	numBytes     = "0123456789"
)

func TempPassword(length int, useLetters bool, useSpecial bool, useNum bool) string {
	b := make([]byte, length)
	for i := range b {
		if useLetters {
			b[i] = letterBytes[rand.Intn(len(letterBytes))]
		} else if useSpecial {
			b[i] = specialBytes[rand.Intn(len(specialBytes))]
		} else if useNum {
			b[i] = numBytes[rand.Intn(len(numBytes))]
		}
	}
	return string(b)
}
