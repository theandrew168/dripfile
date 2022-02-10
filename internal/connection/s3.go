package connection

import (
	"context"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/theandrew168/dripfile/internal/core"
)

type s3Conn struct {
	info   core.S3Info
	client *minio.Client
}

func NewS3(info core.S3Info) (core.Connection, error) {
	creds := credentials.NewStaticV4(
		info.AccessKeyID,
		info.SecretAccessKey,
		"",
	)

	// disable HTTPS requirement for local development / testing
	secure := true
	if strings.Contains(info.Endpoint, "localhost") {
		secure = false
	}
	if strings.Contains(info.Endpoint, "127.0.0.1") {
		secure = false
	}

	client, err := minio.New(info.Endpoint, &minio.Options{
		Creds:  creds,
		Secure: secure,
	})
	if err != nil {
		return nil, core.ErrInvalidEndpoint
	}

	ctx := context.Background()
	if _, err = client.ListBuckets(ctx); err != nil {
		return nil, normalize(err)
	}

	conn := s3Conn{
		info:   info,
		client: client,
	}

	return &conn, nil
}

func (c *s3Conn) List() ([]string, error) {
	ctx := context.Background()
	objects := c.client.ListObjects(
		ctx,
		c.info.BucketName,
		minio.ListObjectsOptions{},
	)

	var files []string
	for object := range objects {
		err := object.Err
		if err != nil {
			return nil, normalize(err)
		}

		files = append(files, object.Key)
	}

	return files, nil
}

func (c *s3Conn) Read(path string) (io.Reader, error) {
	return nil, nil
}

func (c *s3Conn) Write(path string, r io.Reader) error {
	return nil
}

func normalize(err error) error {
	// check for net.Error first (invalid / unreachable endpoint)
	if _, ok := err.(net.Error); ok {
		return core.ErrInvalidEndpoint
	}

	// check for invalid credentials and / or bucket
	s3Err := minio.ToErrorResponse(err)
	if s3Err.Code == "AccessDenied" {
		return core.ErrInvalidCredentials
	}
	if s3Err.Code == "SignatureDoesNotMatch" {
		return core.ErrInvalidCredentials
	}
	if s3Err.Code == "NoSuchBucket" {
		return core.ErrInvalidBucket
	}

	// else bubble
	fmt.Printf("unhandled S3 error code: %s\n", s3Err.Code)
	return err
}
