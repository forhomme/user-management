package domain

import (
	"github.com/pkg/errors"
	"mime/multipart"
)

type UploadFile struct {
	FileName  string `json:"file_name" validate:"required"`
	Requester string
	File      multipart.File
	Header    *multipart.FileHeader
}

var EmptyFileError = errors.New("file uploaded is empty")

func (u *UploadFile) IsEmpty() bool {
	return u.File == nil && u.Header == nil
}

func (u *UploadFile) Validate() error {
	if u.FileName == "" || u.IsEmpty() {
		return EmptyFileError
	}
	return nil
}

type UploadFileResponse struct {
	Filepath string
	FileUrl  string
}
