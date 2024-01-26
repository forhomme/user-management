package _interface

import (
	"context"
	"github.com/forhomme/app-base/errs"
	"github.com/forhomme/app-base/usecase/controller"
	"github.com/forhomme/app-base/usecase/logger"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"user-management/app/common/utils"
	"user-management/app/domain/user"
	"user-management/app/usecase"
	"user-management/app/usecase/command"
	"user-management/app/usecase/query"
	"user-management/config"
)

type HttpServer struct {
	cfg    *config.Config
	logger logger.Logger
	app    usecase.Application
}

func NewHttpServer(cfg *config.Config, log logger.Logger, app usecase.Application) HttpServer {
	return HttpServer{cfg: cfg, logger: log, app: app}
}

func (h HttpServer) SignUp(c controller.Context) (err error) {
	request := new(user.SignUp)
	response := new(user.Token)
	if err = c.ReadRequest(request); err != nil {
		h.logger.Error(err)
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}

	response, err = h.app.Queries.SignUp.Handle(context.TODO(), request)
	if err != nil {
		h.logger.Error(err)
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}
	return c.JSON(http.StatusCreated, utils.ParseResponse(http.StatusOK, "success", response))
}

func (h HttpServer) Login(c controller.Context) (err error) {
	request := new(user.Login)
	response := new(user.Token)
	if err = c.ReadRequest(request); err != nil {
		h.logger.Error(err)
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}

	response, err = h.app.Queries.Login.Handle(context.TODO(), request)
	if err != nil {
		h.logger.Error(err)
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}
	return c.JSON(http.StatusOK, utils.ParseResponse(http.StatusOK, "success", response))
}

func (h HttpServer) RefreshToken(c controller.Context) (err error) {
	request := new(user.RefreshToken)
	response := new(user.Token)
	if err = c.ReadRequest(request); err != nil {
		h.logger.Error(err)
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}

	response, err = h.app.Queries.RefreshToken.Handle(context.TODO(), request)
	if err != nil {
		h.logger.Error(err)
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}
	return c.JSON(http.StatusOK, utils.ParseResponse(http.StatusOK, "success", response))
}

func (h HttpServer) ChangePassword(c controller.Context) (err error) {
	request := new(user.ChangePassword)
	if err = c.ReadRequest(request); err != nil {
		h.logger.Error(err)
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}

	err = h.app.Commands.ChangePassword.Handle(context.TODO(), request)
	if err != nil {
		h.logger.Error(err)
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}
	return c.JSON(http.StatusOK, utils.ParseResponse(http.StatusOK, "success", nil))
}

func (h HttpServer) InsertCategory(c controller.Context) (err error) {
	request := new(Category)
	if err = c.ReadRequest(request); err != nil {
		h.logger.Error(err)
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}

	var cmd command.AddCategory
	err = mapstructure.Decode(request, &cmd)
	if err != nil {
		h.logger.Error(err)
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}

	err = h.app.Commands.AddCategory.Handle(context.TODO(), &cmd)
	if err != nil {
		h.logger.Error(err)
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}

	return c.JSON(http.StatusCreated, utils.ParseResponse(http.StatusCreated, "success", nil))
}

func (h HttpServer) GetAllCategory(c controller.Context) (err error) {
	response := new(AllCategory)
	categories, err := h.app.Queries.GetAllCategories.Handle(context.TODO(), &query.Category{})
	if err != nil {
		h.logger.Error(err)
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}

	err = mapstructure.Decode(categories, &response)
	if err != nil {
		h.logger.Error(err)
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}

	return c.JSON(http.StatusCreated, utils.ParseResponse(http.StatusOK, "success", response))
}

func (h HttpServer) GetCourses(c controller.Context) (err error) {
	request := new(GetCourses)
	response := new(AllCourse)
	if err = c.ReadRequest(request); err != nil {
		h.logger.Error(err)
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}

	user, err := newUserFromContext(c)
	if err != nil {
		h.logger.Error(err)
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}

	q := &query.GetCourses{
		Id:         request.Id,
		CategoryId: request.CategoryId,
		Filter:     request.Filter,
		PerPage:    request.PerPage,
		Page:       request.Page,
		User:       user,
	}
	res, err := h.app.Queries.GetCourses.Handle(context.TODO(), q)
	if err != nil {
		h.logger.Error(err)
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}

	err = mapstructure.Decode(res, &response)
	if err != nil {
		h.logger.Error(err)
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}

	return c.JSON(http.StatusCreated, utils.ParseResponse(http.StatusOK, "success", response))
}

func (h HttpServer) InsertCourse(c controller.Context) (err error) {
	request := new(Course)
	if err = c.ReadRequest(request); err != nil {
		h.logger.Error(err)
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}

	var cmd command.AddCourse
	err = mapstructure.Decode(request, &cmd)
	if err != nil {
		h.logger.Error(err)
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}
	err = h.app.Commands.AddCourse.Handle(context.TODO(), &cmd)
	if err != nil {
		h.logger.Error(err)
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}

	return c.JSON(http.StatusCreated, utils.ParseResponse(http.StatusCreated, "success", nil))
}

func (h HttpServer) UpdateCourse(c controller.Context) (err error) {
	request := new(Course)
	if err = c.ReadRequest(request); err != nil {
		h.logger.Error(err)
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}

	var cmd command.AddCourse
	err = mapstructure.Decode(request, &cmd)
	if err != nil {
		h.logger.Error(err)
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}
	err = h.app.Commands.ReplaceCourse.Handle(context.TODO(), &cmd)
	if err != nil {
		h.logger.Error(err)
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}

	return c.JSON(http.StatusCreated, utils.ParseResponse(http.StatusCreated, "success", nil))
}

func (h HttpServer) UploadFile(c controller.Context) (err error) {
	file, fileHeader, err := c.Request().FormFile(content)
	if err != nil {
		h.logger.Error(err)
		return c.JSON(http.StatusBadRequest, utils.ParseMessage(err))
	}
	fileName := c.Request().Form.Get(filename)

	uploadFile := &query.UploadFile{
		File:     file,
		Header:   fileHeader,
		FileName: fileName,
	}

	out, err := h.app.Queries.UploadFile.Handle(context.TODO(), uploadFile)
	if err != nil {
		h.logger.Error(err)
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}

	return c.JSON(http.StatusOK, utils.ParseResponse(http.StatusOK, "success", out))
}
