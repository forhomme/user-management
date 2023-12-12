package _interface

import (
	"context"
	"encoding/json"
	"github.com/forhomme/app-base/errs"
	"github.com/forhomme/app-base/usecase/controller"
	"github.com/forhomme/app-base/usecase/logger"
	"github.com/mitchellh/mapstructure"
	"net/http"
	_interface "user-management/app/audit_trail/interface"
	"user-management/app/courses/usecase"
	"user-management/app/courses/usecase/command"
	"user-management/app/courses/usecase/query"
	"user-management/config"
	"user-management/utils"
)

type HttpServer struct {
	cfg    *config.Config
	logger logger.Logger
	app    usecase.Application
	audit  _interface.AuditTrailService
}

func NewHttpServer(cfg *config.Config, log logger.Logger, app usecase.Application, audit _interface.AuditTrailService) HttpServer {
	return HttpServer{cfg: cfg, logger: log, app: app, audit: audit}
}

func (h HttpServer) InsertCategory(c controller.Context) (err error) {
	request := new(Category)
	defer func() {
		requestString, _ := json.Marshal(request)
		h.audit.InsertAuditTrail(context.TODO(), &_interface.AuditTrail{
			UserId:   c.GetAuthUser().UserId,
			Menu:     "Insert Category",
			Path:     c.Path(),
			Request:  string(requestString),
			Response: "",
		})
	}()
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
	defer func() {
		responseString, _ := json.Marshal(response)
		h.audit.InsertAuditTrail(context.TODO(), &_interface.AuditTrail{
			UserId:   c.GetAuthUser().UserId,
			Menu:     "Get All Category",
			Path:     c.Path(),
			Request:  "",
			Response: string(responseString),
		})
	}()
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
	defer func() {
		requestString, _ := json.Marshal(request)
		responseString, _ := json.Marshal(response)
		h.audit.InsertAuditTrail(context.TODO(), &_interface.AuditTrail{
			UserId:   c.GetAuthUser().UserId,
			Menu:     "Get All Category",
			Path:     c.Path(),
			Request:  string(requestString),
			Response: string(responseString),
		})
	}()
	if err = c.ReadRequest(request); err != nil {
		h.logger.Error(err)
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}
	q := &query.GetCourses{
		Id:         request.Id,
		CategoryId: request.CategoryId,
		Filter:     request.Filter,
		PerPage:    request.PerPage,
		Page:       request.Page,
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
	defer func() {
		requestString, _ := json.Marshal(request)
		h.audit.InsertAuditTrail(context.TODO(), &_interface.AuditTrail{
			UserId:   c.GetAuthUser().UserId,
			Menu:     "Insert Course",
			Path:     c.Path(),
			Request:  string(requestString),
			Response: "",
		})
	}()
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
	defer func() {
		requestString, _ := json.Marshal(request)
		h.audit.InsertAuditTrail(context.TODO(), &_interface.AuditTrail{
			UserId:   c.GetAuthUser().UserId,
			Menu:     "Update Course",
			Path:     c.Path(),
			Request:  string(requestString),
			Response: "",
		})
	}()
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
	err = h.app.Commands.UpdateCourse.Handle(context.TODO(), &cmd)
	if err != nil {
		h.logger.Error(err)
		return c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}

	return c.JSON(http.StatusCreated, utils.ParseResponse(http.StatusCreated, "success", nil))
}
