package model

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type GlobalConfigDao struct {
	rdsCli *redis.Client
}

type IGlobalConfigDao interface {
	Find(ctx context.Context, key string) ([]byte, error)
}

func NewGlobalConfigDao(rdsCli *redis.Client) (*GlobalConfigDao, error) {
	return &GlobalConfigDao{rdsCli: rdsCli}, nil
}

func (dao *GlobalConfigDao) Find(ctx context.Context, key string) ([]byte, error) {
	bs, err := dao.rdsCli.Get(ctx, fmt.Sprintf("config:%s", key)).Bytes()
	if err != nil {
		return nil, err
	}
	return bs, nil
}
