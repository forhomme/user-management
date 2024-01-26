package query

import (
	"context"
	"fmt"
	domain2 "github.com/forhomme/app-base/domain"
	"github.com/forhomme/app-base/usecase/logger"
	"github.com/forhomme/app-base/usecase/storage"
	"github.com/mitchellh/mapstructure"
	"mime/multipart"
	"path/filepath"
	"user-management/app/common/decorator"
	domain "user-management/app/domain/upload_file"
	"user-management/config"
)

type UploadFile struct {
	FileName  string
	File      multipart.File
	Header    *multipart.FileHeader
	Requester string
}

type UploadFileResponse struct {
	Filepath string
	FileUrl  string
}

type UploadFileHandler decorator.QueryHandler[*UploadFile, *UploadFileResponse]

type uploadFileHandler struct {
	cfg     *config.Config
	storage storage.Storage
	logger  logger.Logger
}

func NewUploadFileHandler(cfg *config.Config, storage storage.Storage, logger logger.Logger) decorator.QueryHandler[*UploadFile, *UploadFileResponse] {
	return decorator.ApplyQueryDecorators[*UploadFile, *UploadFileResponse](
		uploadFileHandler{
			cfg:     cfg,
			storage: storage,
			logger:  logger,
		},
		logger,
	)
}

func (u uploadFileHandler) Handle(ctx context.Context, in *UploadFile) (*UploadFileResponse, error) {
	var file *domain.UploadFile
	err := mapstructure.Decode(in, &file)
	if err != nil {
		u.logger.Error(err)
		return nil, err
	}

	if err = file.Validate(); err != nil {
		u.logger.Error(err)
		return nil, err
	}

	filename := fmt.Sprintf("%s/%s/%s%s", u.cfg.BucketName, in.Requester, file.FileName, filepath.Ext(in.Header.Filename))
	if err = u.storage.PutMultipartObject(ctx, &domain2.PutObjectRequest{
		Key:        filename,
		File:       file.File,
		FileHeader: file.Header,
	}); err != nil {
		u.logger.Error(err)
		return nil, err
	}

	url, err := u.storage.PreSignedGetObject(ctx, &domain2.GetObjectRequest{
		Filename: filename,
		Duration: 1,
	})
	if err != nil {
		u.logger.Error(err)
		return nil, err
	}

	return &UploadFileResponse{
		Filepath: filename,
		FileUrl:  url.String(),
	}, nil
}
