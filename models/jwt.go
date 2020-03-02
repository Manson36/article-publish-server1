package models

import (
	"errors"
	"github.com/article-publish-server1/config"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"
)

type AdminUserCustomClaims struct {
	UserID int64 `json:"userID"`
	jwt.StandardClaims
}

func (c *AdminUserCustomClaims) Sign() (string, error) {
	c.StandardClaims = jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24 * time.Duration(config.Web.ExpiresAt)).Unix(),
		IssuedAt:  time.Now().Unix(),
		Id:        strconv.FormatInt(c.UserID, 10),
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, *c).SignedString([]byte(config.Web.JWTSecret))
}

func (c *AdminUserCustomClaims) Parse(tokenStr string) error {
	token, err := jwt.ParseWithClaims(tokenStr, &AdminUserCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Web.JWTSecret), nil
	})

	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(*AdminUserCustomClaims); ok && token.Valid {
		*c = *claims
		return nil
	}

	return errors.New("jwt verify fail")
}
