package model

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var EmptyAdvConfig = &AdvConfig{
	Candidates: make([]AdvCandidate, 0),
}

type AdvCandidate struct {
	AdCode    string `json:"ad_code"`
	AdID      string `json:"ad_id"`
	Height    int64  `json:"height"`
	Width     int64  `json:"width"`
	AdChannel string `json:"ad_channel"`
	ImageURL  string `json:"image_url"`
	CoolDown  int64  `json:"cool_down"`
	CountDown int64  `json:"count_down"`
	Location  string `json:"location"` // 跳转地址
	Action    int64  `json:"action"` // 0 广告; 1 原生; 2 webView
}

type AdvConfig struct {
	Candidates []AdvCandidate `json:"candidates"`
}

type AdvertiseDao struct {
	rdsCli *redis.Client
}

type IAdvertiseDao interface {
	Find(ctx context.Context, adCode string) (*AdvConfig, error)
}

func NewAdvertiseDao(rdsCli *redis.Client) (*AdvertiseDao, error) {
	return &AdvertiseDao{rdsCli: rdsCli}, nil
}

func (dao *AdvertiseDao) Find(ctx context.Context, adCode string) (*AdvConfig, error) {
	key := fmt.Sprintf("adv:%s", adCode)
	bs, err := dao.rdsCli.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	adv := new(AdvConfig)
	adv.Candidates = make([]AdvCandidate, 0)
	if err := json.Unmarshal(bs, adv); err != nil {
		return nil, err
	}
	return adv, nil
}
