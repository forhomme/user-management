package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"path/filepath"
	s "user-management/app/audit_trail/service"

	in "user-management/app/audit_trail/interface"
	i "user-management/app/courses/interface"
	"user-management/app/courses/service"
	"user-management/app/user/infrastructure"
	"user-management/app/user/interface"
	"user-management/app/user/usecase"
	"user-management/config"
	"user-management/constants"

	"github.com/forhomme/app-base/cmd"
	"github.com/forhomme/app-base/infrastructure/baselogger"
	"github.com/forhomme/app-base/infrastructure/router"
	"github.com/forhomme/app-base/usecase/database"
	"github.com/forhomme/app-base/usecase/http_handler"
	"github.com/forhomme/app-base/usecase/logger"
	"github.com/forhomme/app-base/usecase/storage"
	"github.com/spf13/viper"
)

const (
	StdOut = "stdout"
	StdErr = "stderr"
)

func main() {
	appLogger := baselogger.NewBaseLogger(StdOut)
	// Execute command
	err := cmd.Execute(appLogger, newInjectEndpoints(appLogger))
	if err != nil {
		appLogger.Fatal(fmt.Errorf("Error in main.cmd.Execute: %w", err))
	}
}

func newInjectEndpoints(appLogger logger.Logger) cmd.InjectEndpoints {
	return func(
		viper *viper.Viper,
		route *router.Router,
		logger logger.Logger,
		httpHandler http_handler.HttpHandler,
		sqlHandler map[string]database.SqlHandler,
		storageHandler map[string]storage.Storage,
	) error {

		cfg, err := config.LoadLocalConfig(viper)
		if err != nil {
			appLogger.Fataf(err, "LoadConfig: %v", err)
		}

		// Use the SetServerAPIOptions() method to set the Stable API version to 1
		serverAPI := options.ServerAPI(options.ServerAPIVersion1)
		opts := options.Client().ApplyURI(cfg.DatabaseUri).SetServerAPIOptions(serverAPI)

		// Create a new client and connect to the server
		client, err := mongo.Connect(context.TODO(), opts)
		if err != nil {
			panic(err)
		}

		// user service
		db := infrastructure.NewDatabase(cfg, logger, sqlHandler[constants.Write])
		uc := usecase.NewCoreUsecase(cfg, logger, db, storageHandler[constants.Upload])
		userCtrl := _interface.NewHttpTransport(cfg, logger, uc)

		// audit service
		auditApp := s.NewAuditApplication(cfg, logger, client)
		auditCtrl := in.NewAuditTrailService(cfg, logger, auditApp)

		// course service
		courseApp := service.NewApplication(cfg, logger, client)
		courseCtrl := i.NewHttpServer(cfg, logger, courseApp, auditCtrl)

		userManagement := route.API.Group("/users")
		authUserManagement := route.API.Group("/users")
		authUserManagement.Use(route.AuthUserMiddleware)

		courseManagement := route.API.Group("/course")
		courseManagement.Use(route.AuthUserMiddleware)

		// the endpoint
		userManagement.POST("/signup", router.C(userCtrl.SignUp))
		userManagement.POST("/login", router.C(userCtrl.Login))
		userManagement.POST("/refresh-token", router.C(userCtrl.RefreshToken))
		authUserManagement.POST("/change-password", router.C(userCtrl.ChangePassword))
		authUserManagement.POST("/menu", router.C(userCtrl.GetMenu))

		// course endpoint
		courseManagement.GET("/category", router.C(courseCtrl.GetAllCategory))
		courseManagement.POST("/category/insert", router.C(courseCtrl.InsertCategory))
		courseManagement.POST("/insert", router.C(courseCtrl.InsertCourse))
		courseManagement.POST("/get", router.C(courseCtrl.GetCourses))
		courseManagement.PUT("/update", router.C(courseCtrl.UpdateCourse))

		err = InitDB(sqlHandler[constants.Write], constants.DBInitFile)
		if err != nil {
			appLogger.Fataf(err, "Database init: %w", err)
		}

		return nil
	}
}

// InitDB Method for initializing Database.
func InitDB(sqlHandler database.SqlHandler, dbfilepath string) error {

	query, err := os.ReadFile(filepath.Clean(dbfilepath))
	if err != nil {
		return err
	}
	if err = sqlHandler.MultiExec(string(query)); err != nil {
		return err
	}
	return nil
}
