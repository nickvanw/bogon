package util

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
)

var (
	// ErrInvalidEnvironment is returned when the required env variables are not present
	ErrInvalidEnvironment = errors.New("required keys were not in the environment")
)

// UploadWithEnv takes a bucket, object name, content type and reader and uploads
// using the s3 endpoint, access key and secret key from the environment
// It does not close the reader, that is the responsibility of the caller
func UploadWithEnv(bucket, region, name, cType string, r io.Reader) (string, error) {
	auth, err := aws.EnvAuth()
	if err != nil {
		return "", ErrInvalidEnvironment
	}
	rgn, ok := aws.Regions[region]
	if !ok {
		return "", errors.New("invalid AWS region")
	}
	bkt := s3.New(auth, rgn).Bucket(bucket)
	var data bytes.Buffer
	if _, err := io.Copy(&data, r); err != nil {
		return "", err
	}
	if err := bkt.Put(name, data.Bytes(), cType, s3.PublicRead); err != nil {
		return "", err
	}
	return fmt.Sprintf("http://%s/%s", bucket, name), nil
}
