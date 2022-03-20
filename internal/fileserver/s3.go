package fileserver

import (
	"context"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Info struct {
	Endpoint        string `json:"endpoint"`
	BucketName      string `json:"bucket_name"`
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
}

type s3Conn struct {
	info   S3Info
	client *minio.Client
}

func NewS3(info S3Info) (FileServer, error) {
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
		return nil, ErrInvalidEndpoint
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
	ctx := context.Background()
	obj, err := c.client.GetObject(
		ctx,
		c.info.BucketName,
		path,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (c *s3Conn) Write(path string, r io.Reader) error {
	ctx := context.Background()
	_, err := c.client.PutObject(
		ctx,
		c.info.BucketName,
		path,
		r,
		-1,
		minio.PutObjectOptions{},
	)
	if err != nil {
		return err
	}

	return nil
}

func normalize(err error) error {
	// check for net.Error first (invalid / unreachable endpoint)
	if _, ok := err.(net.Error); ok {
		return ErrInvalidEndpoint
	}

	// check for invalid credentials and / or bucket
	s3Err := minio.ToErrorResponse(err)
	if s3Err.Code == "InvalidAccessKeyId" {
		return ErrInvalidCredentials
	}
	if s3Err.Code == "AccessDenied" {
		return ErrInvalidCredentials
	}
	if s3Err.Code == "SignatureDoesNotMatch" {
		return ErrInvalidCredentials
	}
	if s3Err.Code == "NoSuchBucket" {
		return ErrInvalidBucket
	}

	// else bubble
	fmt.Printf("unhandled S3 error code: %s\n", s3Err.Code)
	return err
}
