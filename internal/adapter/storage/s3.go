package storage

import (
	"backend-layout/internal/config"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/gabriel-vasile/mimetype"
)

// Bucket represents an S3 storage client
type Bucket struct {
	client *s3.Client
	config config.AWSConfig
}

// NewS3Client creates a new S3 client instance
func NewS3Client(conf config.AWSConfig) (Uploader, error) {
	if conf.Region == "" || conf.AccessKeyID == "" || conf.SecretAccessKey == "" {
		return nil, fmt.Errorf("invalid AWS configuration: missing required fields")
	}

	s3Config, err := awsConfig.LoadDefaultConfig(context.Background(),
		awsConfig.WithRegion(conf.Region),
		awsConfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				conf.AccessKeyID,
				conf.SecretAccessKey,
				"",
			),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &Bucket{
		client: s3.NewFromConfig(s3Config),
		config: conf,
	}, nil
}

// UploadFile uploads a file to S3 and returns its URL
func (b *Bucket) UploadFile(ctx context.Context, file *multipart.FileHeader) (string, error) {
	if file == nil {
		return "", fmt.Errorf("file cannot be nil")
	}

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Generate unique key for the file
	key := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)

	// Detect content type

	contentType, err := mimetype.DetectReader(src)
	if err != nil {
		return "", fmt.Errorf("failed to detect content type: %w", err)
	}

	// Reset file pointer after reading content type
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("failed to reset file pointer: %w", err)
	}

	input := &s3.PutObjectInput{
		Bucket:      aws.String(b.config.Bucket),
		Key:         aws.String(key),
		Body:        src,
		ContentType: aws.String(contentType.String()),
	}

	if _, err = b.client.PutObject(ctx, input); err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s",
		b.config.Bucket,
		b.config.Region,
		key,
	), nil
}

// detectContentType determines the content type of the file
func detectContentType(src io.ReadSeeker) (string, error) {
	buffer := make([]byte, 512)
	n, err := src.Read(buffer)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to read file buffer: %w", err)
	}
	return http.DetectContentType(buffer[:n]), nil
}
