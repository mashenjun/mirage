package service

import (
	apiErr "github.com/mashenjun/mirage/errors"
	"github.com/mashenjun/mirage/model"
	"github.com/mashenjun/mirage/third_party/faceai"

	"context"
)

type Service struct {
	advDao    model.IAdvertiseDao
	faceAICli *faceai.Client
}

type Option func(service *Service) error

func New(advDao model.IAdvertiseDao, faceAICli *faceai.Client, opts ...Option) (*Service, error) {
	srv := &Service{
		advDao:    advDao,
		faceAICli: faceAICli,
	}
	for _, opt := range opts {
		if err := opt(srv); err != nil {
			return nil, err
		}
	}
	return srv, nil
}

// GetAdvertise returns Advertise Data
func (srv *Service) GetAdvertise(ctx context.Context, param GetAdvertiseParam) (*model.AdvConfig, error) {
	return srv.advDao.Find(ctx, param.AdCode)
}

// GetAccessCode returns access code from baidu ai
func (srv *Service) GetAccessCode(ctx context.Context, param GetAccessCodeParam) (*GetAccessCodeData, error) {
	return nil, apiErr.ErrUnimplemented()
}

func (srv *Service) DetectFace(ctx context.Context, param DetectFaceParam) (*DetectFaceData, error) {
	if err := param.validate(); err != nil {
		return nil, apiErr.ErrInvalidParameter(err.Error())
	}
	_, info, err := srv.faceAICli.Detect(ctx, faceai.DetectParam{
		Image:     param.Image,
		ImageType: "BASE64",
		FaceField: "age",
	})
	if err != nil {
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
		Image:      param.Image,
		ImageType:  "BASE64",
		ActionType: param.ActionType,
	})
	if err != nil {
		return nil, err
	}
	data := new(EditAttrData)
	data.Image = ret.Image
	return data, nil
}

func (srv *Service) StyleTrans(ctx context.Context, param StyleTransParam) (*StyleTransData, error) {
	if err := param.validate(); err != nil {
		return nil, apiErr.ErrInvalidParameter(err.Error())
	}
	img, err := srv.faceAICli.StyleTrans(ctx, faceai.StyleTransParam{
		Image:  param.Image,
		Option: param.Option,
	})
	if err != nil {
		return nil, err
	}
	data := new(StyleTransData)
	data.Image = img
	return data, nil
}

func (srv *Service) SelieAnime(ctx context.Context, param SelfieAnimeParam) (*SelfieAnimeData, error) {
	if err := param.validate(); err != nil {
		return nil, apiErr.ErrInvalidParameter(err.Error())
	}
	img, err := srv.faceAICli.SelfieAnime(ctx, faceai.SelfieAnimeParam{
		Image: param.Image,
		Type:  "anime",
	})
	if err != nil {
		return nil, err
	}
	data := new(SelfieAnimeData)
	data.Image = img
	return data, nil
}
