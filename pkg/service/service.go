package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/aliyun-sts-go-sdk/sts"
	"github.com/go-redis/redis/v8"
	"github.com/oklog/ulid/v2"

	apiErr "github.com/mashenjun/mirage/errors"
	"github.com/mashenjun/mirage/log"
	"github.com/mashenjun/mirage/model"
	"github.com/mashenjun/mirage/third_party/faceai"
)

type Service struct {
	advDao      model.IAdvertiseDao
	templateDao model.ITemplateImageDao
	globalDao   model.IGlobalConfigDao
	faceAICli   *faceai.Client

	ossCli *oss.Client

	ossPublicEndpoint       string
	ossPublicBucketEndpoint string
	ossBucketName           string
	pathPrefix              string //需要以/结尾

	ramAK string
	ramSK string
	arn   string
}

type Option func(service *Service) error

func OSSOption(bucketName string, publicEndpoint string, publicBucketEndpoint string, pathPrefix string) Option {
	return func(service *Service) error {
		service.ossBucketName = bucketName
		service.ossPublicEndpoint = publicEndpoint
		service.ossPublicBucketEndpoint = publicBucketEndpoint
		service.pathPrefix = pathPrefix
		return nil
	}
}

func STSOption(ak string, sk string, arn string) Option {
	return func(service *Service) error {
		service.ramAK = ak
		service.ramSK = sk
		service.arn = arn
		return nil
	}
}

func New(advDao model.IAdvertiseDao, templateDao model.ITemplateImageDao, globalConfigDao model.IGlobalConfigDao,
	faceAICli *faceai.Client, ossCli *oss.Client, opts ...Option) (*Service, error) {
	srv := &Service{
		advDao:      advDao,
		templateDao: templateDao,
		globalDao:   globalConfigDao,
		faceAICli:   faceAICli,
		ossCli:      ossCli,
	}
	for _, opt := range opts {
		if err := opt(srv); err != nil {
			return nil, err
		}
	}
	return srv, nil
}

// GetAdvertise 从redis中获取广告配置
func (srv *Service) GetAdvertise(ctx context.Context, param GetAdvertiseParam) (*model.AdvConfig, error) {
	adv, err := srv.advDao.Find(ctx, param.AdCode)
	if err != nil {
		log.Errorf("advDao.Find err:%+v", err)
		return nil, err
	}
	return adv, nil
}

// GetTemplate 从redis中获取素材模板
func (srv *Service) GetTemplate(ctx context.Context, param GetTemplateParam) (*model.TemplateImageConfig, error) {
	template, err := srv.templateDao.Find(ctx, param.Type)
	if err != nil && err != redis.Nil {
		log.Errorf("templateDao.Find err:%+v", err)
		return nil, err
	}
	return template, nil
}

// GetVersionUpdate 从redis版本配置
func (srv *Service) GetVersionUpdate(ctx context.Context) ([]byte, error) {
	bs, err := srv.globalDao.Find(ctx, "version_update")
	if err != nil && err != redis.Nil {
		log.Errorf("globalDao.Find err:%+v", err)
		return nil, err
	} else if err == redis.Nil {
		return nil, nil
	}
	return bs, nil
}

// GetAccessCode returns access code from baidu ai
func (srv *Service) GetAccessCode(ctx context.Context, param GetAccessCodeParam) (*GetAccessCodeData, error) {
	return nil, apiErr.ErrUnimplemented()
}

// UploadSignature 返回给客户端上传oss所需要的配置
func (srv *Service) UploadSignature(ctx context.Context) (*UploadSignatureData, error) {
	sessionName := fmt.Sprintf("mirage@%d", time.Now().UnixNano())
	client := sts.NewClient(srv.ramAK, srv.ramSK, srv.arn, sessionName)
	resp, err := client.AssumeRole(uint(3600))
	if err != nil {
		log.Errorf("AssumeRole err:%+v", err)
		return nil, err
	}
	data := &UploadSignatureData{
		EndPoint:        srv.ossPublicEndpoint,
		AccessKeyId:     resp.Credentials.AccessKeyId,
		AccessKeySecret: resp.Credentials.AccessKeySecret,
		BucketName:      srv.ossBucketName,
		Expiration:      resp.Credentials.Expiration.Unix(),
		SecurityToken:   resp.Credentials.SecurityToken,
		Path:            srv.pathPrefix,
	}
	return data, err
}

func (srv *Service) DetectFace(ctx context.Context, param DetectFaceParam) (*DetectFaceData, error) {
	if err := param.validate(); err != nil {
		return nil, apiErr.ErrInvalidParameter(err.Error())
	}
	_, info, err := srv.faceAICli.Detect(ctx, faceai.DetectParam{
		Image:     param.Image,
		ImageType: "URL",
		FaceField: "age",
	})
	if err != nil {
		log.Errorf("faceAICli.Detect err:%+v", err)
		return nil, err
	}
	if len(info) == 0 {
		return nil, apiErr.ErrFaceNotFound()
	}
	data := new(DetectFaceData)
	data.Age = info[0].Age
	return data, nil
}

func (srv *Service) EditAttr(ctx context.Context, param EditAttrParam) (*EditAttrData, error) {
	if err := param.validate(); err != nil {
		return nil, apiErr.ErrInvalidParameter(err.Error())
	}
	ret, err := srv.faceAICli.EditAttr(ctx, faceai.EditAttrParam{
		Image:      param.GetImageReSize(),
		ImageType:  "URL",
		ActionType: param.ActionType,
	})
	if err != nil {
		log.Errorf("faceAICli.EditAttr err:%+v", err)
		return nil, err
	}
	imageURL, err := srv.upload(ctx, ret.Image)
	if err != nil {
		log.Errorf("upload err:%+v", err)
		return nil, err
	}
	data := new(EditAttrData)
	data.Image = imageURL
	return data, nil
}

func (srv *Service) StyleTrans(ctx context.Context, param StyleTransParam) (*StyleTransData, error) {
	if err := param.validate(); err != nil {
		return nil, apiErr.ErrInvalidParameter(err.Error())
	}
	u, err := url.Parse(param.Image)
	if err != nil {
		log.Errorf("parse err:%+v", err)
		return nil, apiErr.ErrInvalidParameter(err.Error())
	}
	key := strings.TrimPrefix(u.Path, "/")
	b64, err := srv.fetchAndBas64Encode(ctx, key)
	if err != nil {
		log.Errorf("fetchAndBas64Encode with key:%+v err:%+v", key, err)
		return nil, err
	}
	p := faceai.StyleTransParam{
		Image:  b64,
		Option: param.Option,
	}
	img, err := srv.faceAICli.StyleTrans(ctx, p)
	if err != nil {
		log.Errorf("faceAICli.StyleTrans err:%+v", err)
		return nil, err
	}
	imageURL, err := srv.upload(ctx, img)
	if err != nil {
		log.Errorf("upload err:%+v", err)
		return nil, err
	}
	data := new(StyleTransData)
	data.Image = imageURL
	return data, nil
}

func (srv *Service) SelieAnime(ctx context.Context, param SelfieAnimeParam) (*SelfieAnimeData, error) {
	if err := param.validate(); err != nil {
		return nil, apiErr.ErrInvalidParameter(err.Error())
	}
	u, err := url.Parse(param.Image)
	if err != nil {
		return nil, apiErr.ErrInvalidParameter(err.Error())
	}
	key := strings.TrimPrefix(u.Path, "/")
	b64, err := srv.fetchAndBas64Encode(ctx, key)
	if err != nil {
		log.Errorf("fetchAndBas64Encode with key %+v err:%+v", key, err)
		return nil, err
	}
	log.Debugf("encode finish with length %+v", len(b64))
	img, err := srv.faceAICli.SelfieAnime(ctx, faceai.SelfieAnimeParam{
		Image: b64,
		Type:  "anime",
	})
	if err != nil {
		log.Errorf("faceAICli.SelfieAnime err:%+v", err)
		return nil, err
	}
	imageURL, err := srv.upload(ctx, img)
	if err != nil {
		log.Errorf("upload err:%+v", err)
		return nil, err
	}
	data := new(SelfieAnimeData)
	data.Image = imageURL
	return data, nil
}

func (srv *Service) MergeFace(ctx context.Context, param MergeFaceParam) (*MergeFaceData, error) {
	if err := param.validate(); err != nil {
		return nil, apiErr.ErrInvalidParameter(err.Error())
	}
	mergeParam := faceai.MergeFaceParam{
		ImageTemplate: faceai.MergeFaceImage{
			Image:     param.TemplateImage,
			ImageType: "URL",
		},
		ImageTarget: faceai.MergeFaceImage{
			Image:     param.TargetImage,
			ImageType: "URL",
		},
	}
	img, err := srv.faceAICli.MergeFace(ctx, mergeParam)
	if err != nil {
		log.Errorf("faceAICli.MergeFace err:%+v", err)
		return nil, err
	}
	imageURL, err := srv.upload(ctx, img)
	if err != nil {
		log.Errorf("upload err:%+v", err)
		return nil, err
	}
	data := new(MergeFaceData)
	data.Image = imageURL
	return data, nil
}

func (srv *Service) BodySeg(ctx context.Context, param BodySegParam) (*BodySegData, error) {
	if err := param.validate(); err != nil {
		return nil, apiErr.ErrInvalidParameter(err.Error())
	}
	u, err := url.Parse(param.Image)
	if err != nil {
		return nil, apiErr.ErrInvalidParameter(err.Error())
	}
	key := strings.TrimPrefix(u.Path, "/")
	b64, err := srv.fetchAndBas64Encode(ctx, key)
	if err != nil {
		log.Errorf("fetchAndBas64Encode with key %+v err:%+v", key, err)
		return nil, err
	}
	log.Debugf("encode finish with length %+v", len(b64))
	ret, err := srv.faceAICli.BodySeg(ctx, faceai.BodySegParam{
		Image: b64,
		Type:  "foreground",
	})
	if err != nil {
		log.Errorf("faceAICli.BodySeg err:%+v", err)
		return nil, err
	}
	if len(ret.Foreground) == 0 {
		return nil, apiErr.ErrBodyNotFound()
	}
	imageURL, err := srv.upload(ctx, ret.Foreground)
	if err != nil {
		log.Errorf("upload err:%+v", err)
		return nil, err
	}
	data := new(BodySegData)
	data.Image = imageURL
	return data, nil
}

func (srv *Service) fetchAndBas64Encode(ctx context.Context, key string) (string, error) {
	bucket, err := srv.ossCli.Bucket(srv.ossBucketName)
	if err != nil {
		log.Errorf("Bucket err:%+v", err)
		return "", err
	}
	r, err := bucket.GetObject(key)
	if err != nil {
		log.Errorf("GetObject err:%+v", err)
		return "", err
	}
	defer r.Close()
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf), nil
}

func (srv *Service) upload(ctx context.Context, b64 string) (string, error) {
	bucket, err := srv.ossCli.Bucket(srv.ossBucketName)
	if err != nil {
		return "", err
	}
	b, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", err
	}
	uuid, err := ulid.New(ulid.Now(), nil)
	if err != nil {
		return "", err
	}
	key := fmt.Sprintf("%s%s.jpg", srv.pathPrefix, uuid)
	err = bucket.PutObject(key, bytes.NewBuffer(b), oss.ContentType("image/jpg"))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", srv.ossPublicBucketEndpoint, key), nil
}
