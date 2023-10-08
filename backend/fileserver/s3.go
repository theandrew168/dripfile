package fileserver

import (
	"context"
	"errors"
	"io"
	"net"
	"path/filepath"
	"slices"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// ensure FileServer interface is satisfied
var _ FileServer = (*S3FileServer)(nil)

var (
	ErrInvalidEndpoint    = errors.New("s3: invalid endpoint")
	ErrInvalidCredentials = errors.New("s3: invalid credentials")
	ErrInvalidBucket      = errors.New("s3: invalid bucket")
)

type S3Info struct {
	Endpoint        string
	Bucket          string
	AccessKeyID     string
	SecretAccessKey string
}

type S3FileServer struct {
	info   S3Info
	client *minio.Client
}

func NewS3(info S3Info) (*S3FileServer, error) {
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

	fs := S3FileServer{
		info:   info,
		client: client,
	}

	return &fs, nil
}

func (fs *S3FileServer) Ping() error {
	ctx := context.Background()
	buckets, err := fs.client.ListBuckets(ctx)
	if err != nil {
		return checkError(err)
	}

	found := slices.ContainsFunc(buckets, func(bucket minio.BucketInfo) bool {
		return bucket.Name == fs.info.Bucket
	})

	if !found {
		return ErrInvalidBucket
	}

	return nil
}

func (fs *S3FileServer) Search(pattern string) ([]FileInfo, error) {
	ctx := context.Background()
	objects := fs.client.ListObjects(
		ctx,
		fs.info.Bucket,
		minio.ListObjectsOptions{},
	)

	var files []FileInfo
	for object := range objects {
		err := object.Err
		if err != nil {
			return nil, checkError(err)
		}

		ok, _ := filepath.Match(pattern, object.Key)
		if !ok {
			continue
		}

		file := FileInfo{
			Name: object.Key,
			Size: int(object.Size),
		}
		files = append(files, file)
	}

	return files, nil
}

func (fs *S3FileServer) Read(name string) (io.Reader, error) {
	ctx := context.Background()
	obj, err := fs.client.GetObject(
		ctx,
		fs.info.Bucket,
		name,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, checkError(err)
	}

	return obj, nil
}

func (fs *S3FileServer) Write(file FileInfo, r io.Reader) error {
	ctx := context.Background()
	_, err := fs.client.PutObject(
		ctx,
		fs.info.Bucket,
		file.Name,
		r,
		int64(file.Size),
		minio.PutObjectOptions{},
	)
	if err != nil {
		return checkError(err)
	}

	return nil
}

func checkError(err error) error {
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
	return err
}
