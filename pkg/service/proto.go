package service

import "errors"

type GetAdvertiseParam struct {
	AdCode string `json:"ad_code"`
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
	Image string `json:"image"`
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