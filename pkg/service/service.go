package service

import (
	apiErr "github.com/mashenjun/mirage/errors"

	"context"

	"github.com/go-redis/redis/v8"
)

type Service struct {
	rdsCli *redis.Client
}

type Option func(service *Service) error

func New(rdsCli *redis.Client, opts ...Option) (*Service, error) {
	srv := &Service{
		rdsCli:    rdsCli,
	}
	for _, opt := range opts {
		if err := opt(srv); err != nil {
			return nil, err
		}
	}
	return srv, nil
}

func (srv *Service) GetAdvertise(ctx context.Context, param GetAdvertiseParam) (*GetAdvertiseData, error) {
	return nil, apiErr.ErrUnimplemented()
}

func (srv *Service) GetAccessCode(ctx context.Context, param GetAccessCodeParam) (*GetAccessCodeData, error) {
	return nil, apiErr.ErrUnimplemented()
}
