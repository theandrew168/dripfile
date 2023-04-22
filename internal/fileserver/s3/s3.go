package s3

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/theandrew168/dripfile/internal/fileserver"
)

var (
	ErrInvalidEndpoint    = errors.New("s3: invalid endpoint")
	ErrInvalidCredentials = errors.New("s3: invalid credentials")
	ErrInvalidBucket      = errors.New("s3: invalid bucket")
)

type Info struct {
	Endpoint        string `json:"endpoint"`
	Bucket          string `json:"bucket"`
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
}

func NewInfo(endpoint, bucket, accessKeyID, secretAccessKey string) Info {
	info := Info{
		Endpoint:        endpoint,
		Bucket:          bucket,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
	}
	return info
}

func NewInfoFromJSON(data []byte) (Info, error) {
	var info Info
	err := json.Unmarshal(data, &info)
	if err != nil {
		return Info{}, nil
	}

	return info, nil
}

func (info Info) ToJSON() ([]byte, error) {
	data, err := json.Marshal(info)
	if err != nil {
		return nil, err
	}

	return data, nil
}

type FileServer struct {
	info   Info
	client *minio.Client
}

func New(info Info) (*FileServer, error) {
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

	fs := FileServer{
		info:   info,
		client: client,
	}

	return &fs, nil
}

func (fs *FileServer) Ping() error {
	ctx := context.Background()
	buckets, err := fs.client.ListBuckets(ctx)
	if err != nil {
		return normalize(err)
	}

	found := false
	for _, bucket := range buckets {
		if bucket.Name == fs.info.Bucket {
			found = true
			break
		}
	}

	if !found {
		return ErrInvalidBucket
	}

	return nil
}

func (fs *FileServer) Search(pattern string) ([]fileserver.FileInfo, error) {
	ctx := context.Background()
	objects := fs.client.ListObjects(
		ctx,
		fs.info.Bucket,
		minio.ListObjectsOptions{},
	)

	var files []fileserver.FileInfo
	for object := range objects {
		err := object.Err
		if err != nil {
			return nil, normalize(err)
		}

		ok, _ := filepath.Match(pattern, object.Key)
		if !ok {
			continue
		}

		file := fileserver.NewFileInfo(object.Key, object.Size)
		files = append(files, file)
	}

	return files, nil
}

func (fs *FileServer) Read(file fileserver.FileInfo) (io.Reader, error) {
	ctx := context.Background()
	obj, err := fs.client.GetObject(
		ctx,
		fs.info.Bucket,
		file.Name,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (fs *FileServer) Write(file fileserver.FileInfo, r io.Reader) error {
	ctx := context.Background()
	_, err := fs.client.PutObject(
		ctx,
		fs.info.Bucket,
		file.Name,
		r,
		file.Size,
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
