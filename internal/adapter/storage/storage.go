package storage

import (
	"context"
	"mime/multipart"
)

type Uploader interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader) (string, error)
}
