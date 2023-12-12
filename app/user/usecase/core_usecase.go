package usecase

import (
	"fmt"
	"github.com/forhomme/app-base/usecase/logger"
	"github.com/forhomme/app-base/usecase/storage"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"time"
	domain2 "user-management/app/user/domain"
	"user-management/app/user/infrastructure"
	"user-management/config"
	"user-management/utils"
)

type coreUsecase struct {
	cfg     *config.Config
	logger  logger.Logger
	dPorts  infrastructure.DatabasePorts
	storage storage.Storage
}

func NewCoreUsecase(cfg *config.Config, logger logger.Logger, dPorts infrastructure.DatabasePorts, storage storage.Storage) InputPorts {
	return &coreUsecase{
		cfg:     cfg,
		logger:  logger,
		dPorts:  dPorts,
		storage: storage,
	}
}

func (c *coreUsecase) InsertAuditTrail(request *domain2.AuditTrail) error {
	return c.dPorts.InsertAuditTrail(request)
}

func (c *coreUsecase) SignUp(request *domain2.SignUpRequest) (*domain2.SignUpResponse, error) {
	hashPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		c.logger.Error(fmt.Errorf("error hash password: %w", err))
		return nil, err
	}

	dataUser, err := c.dPorts.InsertUser(&domain2.UsersDatabase{
		RoleId:   request.RoleId,
		UserId:   uuid.New().String(),
		Email:    request.Email,
		Password: hashPassword,
	})
	if err != nil {
		c.logger.Error(fmt.Errorf("error insert users : %w", err))
		return nil, err
	}

	out, err := generateToken(dataUser, c.cfg)
	if err != nil {
		c.logger.Error(fmt.Errorf("error create token : %w", err))
		return nil, err
	}
	return out, nil
}

func (c *coreUsecase) Login(request *domain2.LoginRequest) (*domain2.LoginResponse, error) {
	dataUser, err := c.dPorts.GetUserByEmail(request.Email)
	if err != nil {
		c.logger.Error(fmt.Errorf("error get data users: %w", err))
		return nil, err
	}
	if dataUser == nil || dataUser == (&domain2.UsersDatabase{}) {
		err = fmt.Errorf("users not found")
		return nil, err
	}

	if !utils.CheckPasswordHash(request.Password, dataUser.Password) {
		err = fmt.Errorf("password invalid")
		return nil, err
	}

	out, err := generateToken(dataUser, c.cfg)
	if err != nil {
		c.logger.Error(fmt.Errorf("error create token : %w", err))
		return nil, err
	}
	return out, nil
}

func (c *coreUsecase) RefreshToken(request *domain2.RefreshTokenRequest) (*domain2.RefreshTokenResponse, error) {
	tokenRt, err := jwt.Parse(request.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(c.cfg.SecretKey), nil
	})

	if rtClaims, ok := tokenRt.Claims.(jwt.MapClaims); ok && tokenRt.Valid {
		if rtClaims["sub"].(string) == request.Email {
			dataUser, errRt := c.dPorts.GetUserByEmail(request.Email)
			if errRt != nil {
				c.logger.Error(errRt)
				return nil, errRt
			}
			newTokenPair, errRt := generateToken(dataUser, c.cfg)
			if errRt != nil {
				c.logger.Error(fmt.Errorf("error create token : %w", errRt))
				return nil, errRt
			}

			return newTokenPair, nil
		}
		return nil, fmt.Errorf("token not valid")
	}

	return nil, err
}

func (c *coreUsecase) ChangePassword(request *domain2.ChangePasswordRequest) error {
	dataUser, err := c.dPorts.GetUserById(request.UserId)
	if err != nil {
		c.logger.Error(fmt.Errorf("error get data users: %w", err))
		return err
	}
	if !utils.CheckPasswordHash(request.OldPassword, dataUser.Password) {
		err = fmt.Errorf("old password invalid")
		return err
	}

	hashPassword, err := utils.HashPassword(request.NewPassword)
	if err != nil {
		c.logger.Error(fmt.Errorf("error hash password: %w", err))
		return err
	}
	dataUser.Password = hashPassword
	err = c.dPorts.UpdateUserById(dataUser.UserId, dataUser)
	if err != nil {
		c.logger.Error(fmt.Errorf("error update data users: %w", err))
		return err
	}
	return nil
}

func (c *coreUsecase) GetMenuByUserId(request *domain2.MenuRequest) (*domain2.MenuResponse, error) {
	dataUser, err := c.dPorts.GetUserById(request.UserId)
	if err != nil {
		c.logger.Error(fmt.Errorf("error get data users: %w", err))
		return nil, err
	}

	menus, err := c.dPorts.GetUserMenu(dataUser.RoleId)
	if err != nil {
		c.logger.Error(fmt.Errorf("error get users menu: %w", err))
		return nil, err
	}

	return &domain2.MenuResponse{
		RoleId:   dataUser.RoleId,
		UserId:   dataUser.UserId,
		RoleName: dataUser.RoleName,
		Menus:    menus,
	}, nil
}

func (c *coreUsecase) GetAllCourseCategory() ([]*domain2.CourseCategoryDatabase, error) {
	return c.dPorts.GetAllCourseCategory()
}

func generateToken(dataUser *domain2.UsersDatabase, cfg *config.Config) (*domain2.LoginResponse, error) {
	expire := time.Now().Add(cfg.AuthExpire)
	claims := &domain2.JwtCustomClaims{
		UserId:   dataUser.UserId,
		Email:    dataUser.Email,
		RoleId:   dataUser.RoleId,
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
		Subject:   dataUser.UserId,
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	rt, err := refreshToken.SignedString([]byte(cfg.SecretKey))
	if err != nil {
		return nil, err
	}

	return &domain2.LoginResponse{
		AccessToken:  t,
		RefreshToken: rt,
	}, nil
}
