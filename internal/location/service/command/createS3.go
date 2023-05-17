package command

import (
	"github.com/theandrew168/dripfile/internal/location"
	locationStorage "github.com/theandrew168/dripfile/internal/location/storage"
)

type CreateS3 struct {
	ID string

	Endpoint        string
	Bucket          string
	AccessKeyID     string
	SecretAccessKey string
}

type CreateS3Handler struct {
	locationStorage locationStorage.Storage
}

func NewCreateS3Handler(locationStorage locationStorage.Storage) *CreateS3Handler {
	h := CreateS3Handler{
		locationStorage: locationStorage,
	}
	return &h
}

func (h *CreateS3Handler) Handle(cmd CreateS3) error {
	l, err := location.NewS3(
		cmd.ID,
		cmd.Endpoint,
		cmd.Bucket,
		cmd.AccessKeyID,
		cmd.SecretAccessKey,
	)
	if err != nil {
		return err
	}

	return h.locationStorage.Create(l)
}
