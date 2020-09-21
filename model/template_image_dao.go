package model

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var EmptyTemplateImageConfig = &TemplateImageConfig{
	Templates: make([]TemplateImage, 0),
}

type TemplateImage struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

type TemplateImageConfig struct {
	Templates []TemplateImage `json:"templates"`
}

type TemplateImageDao struct {
	rdsCli *redis.Client
}

type ITemplateImageDao interface {
	Find(ctx context.Context, adCode string) (*TemplateImageConfig, error)
}

func NewTemplateImageDao(rdsCli *redis.Client) (*TemplateImageDao, error) {
	return &TemplateImageDao{rdsCli: rdsCli}, nil
}

func (dao *TemplateImageDao) Find(ctx context.Context, typ string) (*TemplateImageConfig, error) {
	key := fmt.Sprintf("template:%s", typ)
	bs, err := dao.rdsCli.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	cfg := new(TemplateImageConfig)
	cfg.Templates = make([]TemplateImage, 0)
	if err := json.Unmarshal(bs, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
