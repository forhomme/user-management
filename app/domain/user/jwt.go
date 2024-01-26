package user

import (
	"github.com/golang-jwt/jwt"
	"time"
)

type JwtCustomClaims struct {
	UserId   string    `json:"user_id"`
	Email    string    `json:"email"`
	RoleId   int       `json:"role_id"`
	LoginAt  time.Time `json:"login_at"`
	ExpireAt time.Time `json:"expire_at"`
	jwt.StandardClaims
}
