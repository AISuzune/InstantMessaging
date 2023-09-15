package utils

import (
	"github.com/golang-jwt/jwt/v4"
)

type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
