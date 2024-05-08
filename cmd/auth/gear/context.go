package gear

import (
	"dietku-backend/cmd/user/repo"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
	"time"
)

var mySigningKey = []byte("the-secret-of-kalimdor")

type UserClaims struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Email     string             `json:"email" bson:"email"`
	FirstName string             `json:"firstName" bson:"firstName"`
	LastName  string             `json:"lastName" bson:"lastName"`
}

func GenerateToken(user *repo.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID.Hex()
	claims["email"] = user.Email
	claims["firstName"] = user.FirstName
	claims["lastName"] = user.LastName

	expiryDate := time.Now().AddDate(0, 0, 7)
	claims["expiryDate"] = expiryDate
	claims["expiryDateInMillis"] = expiryDate.Unix() * 1000

	generatedToken, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}
	return generatedToken, nil
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return mySigningKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func IsLoggedIn(db *mongo.Database) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header.Get("Authorization")

			loggedIn, err := CheckJWTClaims(db, header)
			if err != nil {
				return echo.ErrUnauthorized
			}

			c.Set("me", loggedIn)
			return next(c)
		}
	}
}

func CheckJWTClaims(db *mongo.Database, header string) (*UserClaims, error) {
	bearer := strings.Split(header, " ")
	if len(bearer) != 2 {
		return nil, errors.New("invalid header")
	}

	if bearer[0] != "Bearer" {
		return nil, errors.New("invalid header")
	}

	token, err := ParseToken(bearer[1])
	if err != nil {
		return nil, errors.New("invalid header")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["id"]

		// check by id
		objectID, err := primitive.ObjectIDFromHex(userID.(string))
		if err != nil {
			return nil, errors.New("invalid header")
		}
		repository := repo.NewUserRepository(db)
		user, err := repository.FindOne(objectID)

		if err != nil {
			return nil, errors.New("invalid header")
		}

		userClaims := &UserClaims{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		}

		return userClaims, nil
	} else {
		return nil, errors.New("invalid header")
	}
}
