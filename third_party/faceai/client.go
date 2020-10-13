package faceai

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-redis/redis/v8"

	httprpc "github.com/mashenjun/mirage/util/http"
)

type Client struct {
	client       *httprpc.Client
	cacheBackend *redis.Client
	cacheEnable  bool

	endpoint string
	ak       string
	sk       string
}

type Option func(client *Client) error

func CacheOption(cache *redis.Client) Option {
	return func(client *Client) error {
		if cache != nil {
			client.cacheEnable = true
			client.cacheBackend = cache
		}
		return nil
	}
}

func TimeoutOption(timeout time.Duration) Option {
	opt := httprpc.TimeoutOption(timeout)
	return func(client *Client) error {
		opt(client.client)
		return nil
	}
}

func MaxIdleConnOption(connCount int) Option {
	opt := httprpc.MaxIdleConnOption(connCount)
	return func(client *Client) error {
		opt(client.client)
		return nil
	}
}

func New(endpoint, ak, sk string, opts ...Option) (*Client, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	if u.Scheme == "" {
		endpoint = "http://" + endpoint
	}
	c := &Client{
		client:   httprpc.New(),
		endpoint: endpoint,
		ak:       ak,
		sk:       sk,
	}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (cli *Client) getAccessToken(ctx context.Context) (string, error) {
	if cli.cacheEnable {
		if token, err := cli.getAccessTokenFromCache(ctx); err == nil && len(token) > 0 {
			return token, nil
		}
	}
	u := fmt.Sprintf("%s%s?grant_type=client_credentials&client_id=%s&client_secret=%s",
		cli.endpoint, accessTokenPath, cli.ak, cli.sk)
	ret := new(GetAccessTokenResp)
	if err := cli.client.CallWithJson(ctx, ret, http.MethodGet, u, nil); err != nil {
		return "", err
	}
	if len(ret.Error) != 0 {
		return "", fmt.Errorf("%s with %s", ret.Error, ret.ErrorDescription)
	}
	if cli.cacheEnable {
		cli.setAccessTokenToCache(ctx, ret.AccessToken)
	}
	return ret.AccessToken, nil
}

func (cli *Client) getAccessTokenFromCache(ctx context.Context) (string, error) {
	return cli.cacheBackend.Get(ctx, accessTokenKey).Result()
}

func (cli *Client) setAccessTokenToCache(ctx context.Context, accessToken string) error {
	if _, err := cli.cacheBackend.Set(ctx, accessTokenKey, accessToken, 2592000*time.Second).Result(); err != nil {
		return err
	}
	return nil
}

// EditAttr 调用人脸属性编辑接口
func (cli *Client) EditAttr(ctx context.Context, param EditAttrParam) (*EditAttData, error) {
	token, err := cli.getAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("%s%s?access_token=%s", cli.endpoint, editAttrPath, token)

	ret := new(EditAttrResp)
	if err := cli.client.CallWithJson(ctx, ret, http.MethodPost, u, param); err != nil {
		return nil, err
	}
	if ret.Code != 0 {
		return nil, &ret.BaseResp
	}
	return &ret.Result, nil
}

// Detect 调用人脸检测接口
func (cli *Client) Detect(ctx context.Context, param DetectParam) (int64, []FaceInfo, error) {
	token, err := cli.getAccessToken(ctx)
	if err != nil {
		return 0, nil, err
	}
	u := fmt.Sprintf("%s%s?access_token=%s", cli.endpoint, detectPath, token)

	ret := new(DetectResp)
	if err := cli.client.CallWithJson(ctx, ret, http.MethodPost, u, param); err != nil {
		return 0, nil, err
	}
	if ret.Code != 0 {
		return 0, nil, &ret.BaseResp
	}
	return ret.Result.FaceNum, ret.Result.FaceList, nil
}

// StyleTrans 调用图片风格转换接口
func (cli *Client) StyleTrans(ctx context.Context, param StyleTransParam) (string, error) {
	token, err := cli.getAccessToken(ctx)
	if err != nil {
		return "", err
	}
	u := fmt.Sprintf("%s%s?access_token=%s", cli.endpoint, styleTransPath, token)

	ret := new(StyleTransResp)
	form, err := param.ToForm()
	if err != nil {
		return "", err
	}
	if err := cli.client.CallWithForm(ctx, ret, http.MethodPost, u, form); err != nil {
		return "", err
	}
	if ret.Code != 0 {
		return "", &ret.BaseResp
	}
	return ret.Image, nil
}

// StyleTrans 调用人像动漫化接口
func (cli *Client) SelfieAnime(ctx context.Context, param SelfieAnimeParam) (string, error) {
	token, err := cli.getAccessToken(ctx)
	if err != nil {
		return "", err
	}
	u := fmt.Sprintf("%s%s?access_token=%s", cli.endpoint, selfieAnimePath, token)

	ret := new(SelfieAnimeResp)
	form, err := param.ToForm()
	if err != nil {
		return "", err
	}
	if err := cli.client.CallWithForm(ctx, ret, http.MethodPost, u, form); err != nil {
		return "", err
	}
	if ret.Code != 0 {
		return "", &ret.BaseResp
	}
	return ret.Image, nil
}

// MergeFace 调用人脸融合接口
func (cli *Client) MergeFace(ctx context.Context, param MergeFaceParam) (string, error) {
	token, err := cli.getAccessToken(ctx)
	if err != nil {
		return "", err
	}
	u := fmt.Sprintf("%s%s?access_token=%s", cli.endpoint, mergeFacePath, token)

	ret := new(MergeFaceResp)
	if err := cli.client.CallWithJson(ctx, ret, http.MethodPost, u, param); err != nil {
		return "", err
	}
	if ret.Code != 0 {
		return "", &ret.BaseResp
	}
	return ret.Result.MergeImage, nil
}

// BodySeg 调用人像分割接口
func (cli *Client) BodySeg(ctx context.Context, param BodySegParam) (*BodySegResp, error) {
	token, err := cli.getAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("%s%s?access_token=%s", cli.endpoint, bodySegPath, token)
	form, err := param.ToForm()
	if err != nil {
		return nil, err
	}
	ret := new(BodySegResp)
	if err := cli.client.CallWithForm(ctx, ret, http.MethodPost, u, form); err != nil {
		return nil, err
	}
	if ret.Code != 0 {
		return nil, &ret.BaseResp
	}
	return ret, nil
}

