package service

import (
	"errors"
)

type GetAdvertiseParam struct {
	AdCode string `json:"ad_code"`
	Extra  string `json:"extra"`
}

func (param *GetAdvertiseParam) validate() error {
	if len(param.AdCode) == 0 {
		return errors.New("param is invalid")
	}
	return nil
}

type GetTemplateParam struct {
	Type string `json:"type"`
}

func (param *GetTemplateParam) validate() error {
	if len(param.Type) == 0 {
		return errors.New("param is invalid")
	}
	return nil
}

type GetAccessCodeParam struct {
}

type GetAccessCodeData struct {
}

type EditAttrParam struct {
	Image      string `json:"image"`
	ActionType string `json:"action_type"`
}

func (param *EditAttrParam) validate() error {
	if len(param.Image) == 0 {
		return errors.New("param.image is invalid")
	}
	if len(param.ActionType) == 0 {
		return errors.New("param.action_type is invalid")
	}
	return nil
}

type EditAttrData struct {
	Image string `json:"image"`
}

type DetectFaceParam struct {
	Image string `json:"image"` // image is the url point to oss url
}

func (param *DetectFaceParam) validate() error {
	if len(param.Image) == 0 {
		return errors.New("param.image is invalid")
	}
	return nil
}

type DetectFaceData struct {
	Age float64 `json:"age"`
}

type StyleTransParam struct {
	Image  string `json:"image"`
	Option string `json:"option"`
}

func (param *StyleTransParam) validate() error {
	if len(param.Image) == 0 {
		return errors.New("param.image is invalid")
	}
	if len(param.Option) == 0 {
		return errors.New("param.option is invalid")
	}
	return nil
}

type StyleTransData struct {
	Image string `json:"image"`
}

type SelfieAnimeParam struct {
	Image string `json:"image"`
}

func (param *SelfieAnimeParam) validate() error {
	if len(param.Image) == 0 {
		return errors.New("param.image is invalid")
	}
	return nil
}

type SelfieAnimeData struct {
	Image string `json:"image"`
}

type UploadSignatureData struct {
	EndPoint        string `json:"end_point"`
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	BucketName      string `json:"bucket_name"`
	Expiration      int64  `json:"expiration"`
	SecurityToken   string `json:"security_token"`
	Path            string `json:"path"`
}

type MergeFaceParam struct {
	TemplateImage string `json:"template_image"`
	TargetImage string `json:"target_image"`
}

func (param *MergeFaceParam) validate() error {
	if len(param.TemplateImage) == 0 {
		return errors.New("param.template_image is invalid")
	}
	if len(param.TargetImage) == 0 {
		return errors.New("param.target_image is invalid")
	}
	return nil
}

type MergeFaceData struct {
	Image string `json:"image"`
}

type BodySegParam struct {
	Image string `json:"image"`
}

func (param *BodySegParam) validate() error {
	if len(param.Image) == 0 {
		return errors.New("param.image is invalid")
	}
	return nil
}

type BodySegData struct {
	Image string `json:"image"`
}
