package _interface

import (
	"encoding/json"
	"fmt"
	"github.com/forhomme/app-base/errs"
	"github.com/forhomme/app-base/usecase/controller"
	"github.com/forhomme/app-base/usecase/logger"
	"net/http"
	domain2 "user-management/app/user/domain"
	"user-management/app/user/usecase"
	"user-management/config"
	"user-management/utils"
)

type httpTransport struct {
	cfg    *config.Config
	logger logger.Logger
	iPort  usecase.InputPorts
}

func NewHttpTransport(cfg *config.Config, logger logger.Logger, iPort usecase.InputPorts) ControllerPorts {
	return &httpTransport{
		cfg:    cfg,
		logger: logger,
		iPort:  iPort,
	}
}

func (h *httpTransport) InsertAuditTrail(c controller.Context, request, response, userId interface{}, err error) {
	requestJson, _ := json.Marshal(request)
	responseJson, _ := json.Marshal(response)
	if err != nil {
		responseJson = []byte(err.Error())
	}
	err = h.iPort.InsertAuditTrail(&domain2.AuditTrail{
		UserId:   userId.(string),
		Menu:     c.Path(),
		Method:   http.MethodPost,
		Request:  string(requestJson),
		Response: string(responseJson),
	})
	if err != nil {
		err = fmt.Errorf("error insert audit trail: %w", err)
		h.logger.Error(err)
	}
	return
}

func (h *httpTransport) SignUp(c controller.Context) (err error) {
	request := new(domain2.SignUpRequest)
	response := new(domain2.SignUpResponse)
	defer func() {
		h.InsertAuditTrail(c, request.Email, response, request.Email, err)
	}()

	if err = c.ReadRequest(request); err != nil {
		h.logger.Error(err)
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}
	response, err = h.iPort.SignUp(request)
	if err != nil {
		h.logger.Error(err)
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *httpTransport) Login(c controller.Context) (err error) {
	request := new(domain2.LoginRequest)
	response := new(domain2.LoginResponse)
	defer func() {
		h.InsertAuditTrail(c, request.Email, response, request.Email, err)
	}()

	if err = c.ReadRequest(request); err != nil {
		h.logger.Error(err)
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}

	response, err = h.iPort.Login(request)
	if err != nil {
		h.logger.Error(err)
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}

	return c.JSON(http.StatusOK, response)
}

func (h *httpTransport) RefreshToken(c controller.Context) (err error) {
	request := new(domain2.RefreshTokenRequest)
	response := new(domain2.RefreshTokenResponse)
	defer func() {
		h.InsertAuditTrail(c, request, response, request.Email, err)
	}()

	if err = c.ReadRequest(request); err != nil {
		h.logger.Error(err)
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}

	response, err = h.iPort.RefreshToken(request)
	if err != nil {
		h.logger.Error(err)
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}

	return c.JSON(http.StatusOK, response)
}

func (h *httpTransport) ChangePassword(c controller.Context) (err error) {
	request := new(domain2.ChangePasswordRequest)
	defer func() {
		h.InsertAuditTrail(c, request, nil, c.GetAuthUser().UserId, err)
	}()

	if err = c.ReadRequest(request); err != nil {
		h.logger.Error(err)
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}

	err = h.iPort.ChangePassword(request)
	if err != nil {
		h.logger.Error(err)
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}

	return c.NoContent(http.StatusOK)
}

func (h *httpTransport) GetMenu(c controller.Context) (err error) {
	request := new(domain2.MenuRequest)
	response := new(domain2.MenuResponse)
	defer func() {
		h.InsertAuditTrail(c, request, response, c.GetAuthUser().UserId, err)
	}()

	if err = c.ReadRequest(request); err != nil {
		h.logger.Error(err)
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}

	response, err = h.iPort.GetMenuByUserId(request)
	if err != nil {
		h.logger.Error(err)
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}

	return c.JSON(http.StatusOK, response)
}
