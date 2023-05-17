package service

import (
	"github.com/theandrew168/dripfile/internal/location/service/command"
	"github.com/theandrew168/dripfile/internal/location/service/query"
	locationStorage "github.com/theandrew168/dripfile/internal/location/storage"
)

type Command struct {
	CreateS3 *command.CreateS3Handler
}

type Query struct {
	Read *query.ReadHandler
}

type Service struct {
	Command Command
	Query   Query
}

func New(locationStorage locationStorage.Storage) *Service {
	s := Service{
		Command: Command{
			CreateS3: command.NewCreateS3Handler(locationStorage),
		},
		Query: Query{
			Read: query.NewReadHandler(locationStorage),
		},
	}
	return &s
}
