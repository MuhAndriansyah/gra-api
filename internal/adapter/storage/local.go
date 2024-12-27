package storage

import (
	"context"
	"mime/multipart"
)

type LocalStorage struct {
}

func (l *LocalStorage) Uploader(ctx context.Context, file *multipart.FileHeader) (string, error) {
	return "", nil
}
