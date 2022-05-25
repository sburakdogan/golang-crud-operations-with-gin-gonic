package helpers

import (
	"errors"
	"os"
	"time"

	"github.com/burak/product-api/models"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TokenDetails struct {
	Id    primitive.ObjectID
	Email string
	jwt.StandardClaims
}

var secretKey string = os.Getenv("SECRET_KEY")

func GenerateToken(user models.User) (string, error) {
	claims := &TokenDetails{
		Id:    user.Id,
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))
	if err != nil {
		panic(err)
	}

	return token, err
}

func ValidateToken(signedToken string) (err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&TokenDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		},
	)

	if err != nil {
		return err
	}

	claims, ok := token.Claims.(*TokenDetails)
	if !ok {
		return errors.New("Invalid token")
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		return errors.New("Expired token")
	}

	return err
}
