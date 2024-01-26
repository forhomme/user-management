package user

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"user-management/config"
)

type RefreshToken struct {
	Email        string `json:"email" validate:"required,email"`
	RefreshToken string `json:"refresh_token" validate:"required"`
}

var ErrorSigningMethod = errors.New("unexpected signing method")
var ErrorTokenNotValid = errors.New("token not valid")

func (r *RefreshToken) ParsingToken(user *User, cfg *config.Config) error {
	tokenRt, err := jwt.Parse(r.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrorSigningMethod
		}
		return []byte(cfg.SecretKey), nil
	})
	if err != nil {
		return ErrorTokenNotValid
	}

	if rtClaims, ok := tokenRt.Claims.(jwt.MapClaims); ok && tokenRt.Valid {
		if rtClaims["sub"].(string) != user.UserId {
			return ErrorTokenNotValid
		}
		return ErrorTokenNotValid
	}
	return nil
}
