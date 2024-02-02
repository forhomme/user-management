package _interface

import (
	"github.com/forhomme/app-base/errs"
	"github.com/forhomme/app-base/infrastructure/baselogger"
	"github.com/forhomme/app-base/infrastructure/telemetry"
	"github.com/forhomme/app-base/usecase/controller"
	"net/http"
	"user-management/app/common/utils"
	"user-management/app/domain/course"
	domain "user-management/app/domain/upload_file"
	"user-management/app/domain/user"
	"user-management/app/usecase"
	"user-management/config"
)

type HttpServer struct {
	cfg    *config.Config
	logger *baselogger.Logger
	tracer *telemetry.OtelSdk
	app    usecase.Application
}

func NewHttpServer(cfg *config.Config, log *baselogger.Logger, tracer *telemetry.OtelSdk, app usecase.Application) HttpServer {
	return HttpServer{cfg: cfg, logger: log, tracer: tracer, app: app}
}

func (h HttpServer) SignUp(c controller.Context) (err error) {
	ctx, span := h.tracer.Tracer.Start(c.Request().Context(), "http.signup")
	defer span.End()

	request := new(user.SignUp)
	response := new(user.Token)
	if err = c.ReadRequest(request); err != nil {
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}

	response, err = h.app.Queries.SignUp.Handle(ctx, request)
	if err != nil {
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}
	return c.JSON(http.StatusCreated, utils.ParseResponse(http.StatusOK, "success", response))
}

func (h HttpServer) Login(c controller.Context) (err error) {
	ctx, span := h.tracer.Tracer.Start(c.Request().Context(), "http.login")
	defer span.End()

	request := new(user.Login)
	response := new(user.Token)
	if err = c.ReadRequest(request); err != nil {
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}

	response, err = h.app.Queries.Login.Handle(ctx, request)
	if err != nil {
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}
	return c.JSON(http.StatusOK, utils.ParseResponse(http.StatusOK, "success", response))
}

func (h HttpServer) RefreshToken(c controller.Context) (err error) {
	ctx, span := h.tracer.Tracer.Start(c.Request().Context(), "http.refresh_token")
	defer span.End()

	request := new(user.RefreshToken)
	response := new(user.Token)
	if err = c.ReadRequest(request); err != nil {
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}

	response, err = h.app.Queries.RefreshToken.Handle(ctx, request)
	if err != nil {
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}
	return c.JSON(http.StatusOK, utils.ParseResponse(http.StatusOK, "success", response))
}

func (h HttpServer) ChangePassword(c controller.Context) (err error) {
	ctx, span := h.tracer.Tracer.Start(c.Request().Context(), "http.change_password")
	defer span.End()

	request := new(user.ChangePassword)
	if err = c.ReadRequest(request); err != nil {
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}
	request.UserId = c.GetAuthUser().UserId

	err = h.app.Commands.ChangePassword.Handle(ctx, request)
	if err != nil {
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}
	return c.JSON(http.StatusOK, utils.ParseResponse(http.StatusOK, "success", nil))
}

func (h HttpServer) InsertCategory(c controller.Context) (err error) {
	ctx, span := h.tracer.Tracer.Start(c.Request().Context(), "http.insert_category")
	defer span.End()

	request := new(course.Category)
	if err = c.ReadRequest(request); err != nil {
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}

	err = h.app.Commands.AddCategory.Handle(ctx, request)
	if err != nil {
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}

	return c.JSON(http.StatusCreated, utils.ParseResponse(http.StatusCreated, "success", nil))
}

func (h HttpServer) GetAllCategory(c controller.Context) (err error) {
	ctx, span := h.tracer.Tracer.Start(c.Request().Context(), "http.get_category")
	defer span.End()

	categories, err := h.app.Queries.GetAllCategories.Handle(ctx, &course.Category{})
	if err != nil {
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}

	return c.JSON(http.StatusOK, utils.ParseResponse(http.StatusOK, "success", categories))
}

func (h HttpServer) GetCourses(c controller.Context) (err error) {
	ctx, span := h.tracer.Tracer.Start(c.Request().Context(), "http.get_courses")
	defer span.End()

	request := new(course.FilterCourse)
	if err = c.ReadRequest(request); err != nil {
		h.logger.Error(err)
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}

	user, err := newUserFromContext(c)
	if err != nil {
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}
	request.User = user

	res, err := h.app.Queries.GetCourses.Handle(ctx, request)
	if err != nil {
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}

	return c.JSON(http.StatusCreated, utils.ParseResponse(http.StatusOK, "success", res))
}

func (h HttpServer) InsertCourse(c controller.Context) (err error) {
	ctx, span := h.tracer.Tracer.Start(c.Request().Context(), "http.insert_course")
	defer span.End()

	request := new(course.CoursePath)
	if err = c.ReadRequest(request); err != nil {
		h.logger.Error(err)
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}

	err = h.app.Commands.AddCourse.Handle(ctx, request)
	if err != nil {
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}

	return c.JSON(http.StatusCreated, utils.ParseResponse(http.StatusCreated, "success", nil))
}

func (h HttpServer) UpdateCourse(c controller.Context) (err error) {
	ctx, span := h.tracer.Tracer.Start(c.Request().Context(), "http.update_course")
	defer span.End()

	request := new(course.CoursePath)
	if err = c.ReadRequest(request); err != nil {
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}

	err = h.app.Commands.ReplaceCourse.Handle(ctx, request)
	if err != nil {
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}

	return c.JSON(http.StatusCreated, utils.ParseResponse(http.StatusCreated, "success", nil))
}

func (h HttpServer) UploadFile(c controller.Context) (err error) {
	ctx, span := h.tracer.Tracer.Start(c.Request().Context(), "http.get_category")
	defer span.End()

	file, fileHeader, err := c.Request().FormFile(content)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ParseMessage(err))
	}
	fileName := c.Request().Form.Get(filename)

	uploadFile := &domain.UploadFile{
		File:      file,
		Header:    fileHeader,
		FileName:  fileName,
		Requester: c.GetAuthUser().UserId,
	}

	out, err := h.app.Queries.UploadFile.Handle(ctx, uploadFile)
	if err != nil {
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}

	return c.JSON(http.StatusOK, utils.ParseResponse(http.StatusOK, "success", out))
}
