package user

import (
	"github.com/golang-jwt/jwt"
	"time"
	"user-management/app/common/utils"
	"user-management/config"
)

type User struct {
	RoleId   int    `json:"role_id"`
	RoleName string `json:"role_name"`
	UserId   string `json:"user_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *User) IsExist() bool {
	return u.UserId != ""
}

func (u *User) ChangePassword(password string) error {
	hashPassword, err := utils.HashPassword(password)
	if err != nil {
		return errorGeneratePassword
	}
	u.Password = hashPassword
	return nil
}

func (u *User) GenerateToken(cfg *config.Config) (*Token, error) {
	expire := time.Now().Add(cfg.AuthExpire)
	claims := &JwtCustomClaims{
		UserId:   u.UserId,
		Email:    u.Email,
		RoleId:   u.RoleId,
		LoginAt:  time.Now(),
		ExpireAt: expire,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expire.Unix(),
		},
	}

	// access token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(cfg.SecretKey))
	if err != nil {
		return nil, err
	}

	// refresh token
	rtExpire := time.Now().Add(cfg.RefreshTokenExpire)
	rtClaims := jwt.StandardClaims{
		ExpiresAt: rtExpire.Unix(),
		Subject:   u.UserId,
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	rt, err := refreshToken.SignedString([]byte(cfg.SecretKey))
	if err != nil {
		return nil, err
	}

	return &Token{
		AccessToken:  t,
		RefreshToken: rt,
	}, nil
}
