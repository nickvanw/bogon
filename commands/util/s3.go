package util

import (
	"errors"
	"fmt"
	"io"

	"github.com/minio/minio-go"
	"github.com/nickvanw/bogon/commands/config"
)

var (
	// ErrInvalidEnvironment is returned when the required env variables are not present
	ErrInvalidEnvironment = errors.New("required keys were not in the environment")
)

// UploadWithEnv takes a bucket, object name, content type and reader and uploads
// using the s3 endpoint, access key and secret key from the environment
// It does not close the reader, that is the responsibility of the caller
func UploadWithEnv(bucket, name, cType string, r io.Reader) (string, error) {
	s3Endpoint, eOk := config.Get("S3_ENDPOINT")
	s3ID, idOk := config.Get("S3_ACCESS_KEY")
	s3Key, keyOk := config.Get("S3_SECRET_KEY")
	if !eOk || !idOk || !keyOk {
		return "", ErrInvalidEnvironment
	}

	client, err := minio.New(s3Endpoint, s3ID, s3Key, false)
	if err != nil {
		return "", err
	}

	if _, err := client.PutObject(bucket, name, r, cType); err != nil {
		return "", err
	}
	// We assume the URL can be accessed over HTTP here.
	return fmt.Sprintf("http://%s/%s/%s", s3Endpoint, bucket, name), nil
}
