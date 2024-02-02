package query

import (
	"context"
	"fmt"
	domain2 "github.com/forhomme/app-base/domain"
	"github.com/forhomme/app-base/infrastructure/baselogger"
	"github.com/forhomme/app-base/infrastructure/telemetry"
	"github.com/forhomme/app-base/usecase/storage"
	"path/filepath"
	"time"
	"user-management/app/common/decorator"
	domain "user-management/app/domain/upload_file"
	"user-management/config"
)

type UploadFileHandler decorator.QueryHandler[*domain.UploadFile, *domain.UploadFileResponse]

type uploadFileHandler struct {
	cfg     *config.Config
	storage storage.Storage
}

func NewUploadFileHandler(cfg *config.Config, storage storage.Storage, logger *baselogger.Logger, tracer *telemetry.OtelSdk) decorator.QueryHandler[*domain.UploadFile, *domain.UploadFileResponse] {
	return decorator.ApplyQueryDecorators[*domain.UploadFile, *domain.UploadFileResponse](
		uploadFileHandler{
			cfg:     cfg,
			storage: storage,
		},
		logger,
		tracer,
	)
}

func (u uploadFileHandler) Handle(ctx context.Context, in *domain.UploadFile) (out *domain.UploadFileResponse, err error) {
	if err = in.Validate(); err != nil {
		return nil, err
	}

	filename := fmt.Sprintf("%s/%s/%s%s", u.cfg.BucketName, in.Requester, in.FileName, filepath.Ext(in.Header.Filename))
	if err = u.storage.PutMultipartObject(ctx, &domain2.PutObjectRequest{
		Key:        filename,
		File:       in.File,
		FileHeader: in.Header,
	}); err != nil {
		return nil, err
	}

	url, err := u.storage.PreSignedGetObject(ctx, &domain2.GetObjectRequest{
		Filename: filename,
		Duration: 30 * time.Minute,
	})
	if err != nil {
		return nil, err
	}

	return &domain.UploadFileResponse{
		Filepath: filename,
		FileUrl:  url.String(),
	}, nil
}
