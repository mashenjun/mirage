package faceai

import (
	"errors"
	"net/url"

	"github.com/gorilla/schema"
)

var formEncoder *schema.Encoder

func init() {
	formEncoder = schema.NewEncoder()
}

type BaseResp struct {
	Code      int    `json:"error_code"`
	Message   string `json:"error_msg"`
	LogID     int64  `json:"log_id"`
	Timestamp int64  `json:"timestamp"`
}

type GetAccessTokenResp struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        int64  `json:"expires_in"`
	Scope            string `json:"scope"`
	SessionKey       string `json:"session_key"`
	AccessToken      string `json:"access_token"`
	SessionSecret    string `json:"session_secret"`
}

type EditAttrParam struct {
	Image          string `json:"image"`
	ImageType      string `json:"image_type"`
	ActionType     string `json:"action_type"`
	QualityControl string `json:"quality_control,omitempty"`
	FaceLocation   string `json:"face_location,omitempty"`
}

type EditAttData struct {
	Image string `json:"image"`
}

type EditAttrResp struct {
	BaseResp
	Result EditAttData `json:"result"`
}

type DetectParam struct {
	Image     string `json:"image"`
	ImageType string `json:"image_type"`
	FaceField string `json:"face_field,omitempty"`
}

type FaceInfo struct {
	FaceToken string  `json:"face_token"`
	Age       float64 `json:"age"`
}

type DetectResp struct {
	BaseResp
	FaceNum  int64      `json:"face_num"`
	FaceList []FaceInfo `json:"face_list"`
}

type StyleTransParam struct {
	Image  string `json:"image" schema:"image"`
	Option string `json:"option" schema:"option"`
}

func (param *StyleTransParam) ToForm() (url.Values, error) {
	if formEncoder == nil {
		return nil, errors.New("formEncoder is nil")
	}
	reqForm := make(url.Values)
	if err := formEncoder.Encode(param, reqForm); err != nil {
		return nil, err
	}
	return reqForm, nil
}

type StyleTransResp struct {
	BaseResp
	Image string `json:"image"`
}

type SelfieAnimeParam struct {
	Image string `json:"image" schema:"image"`
	Type  string `json:"type" schema:"type"`
}

func (param *SelfieAnimeParam) ToForm() (url.Values, error) {
	if formEncoder == nil {
		return nil, errors.New("formEncoder is nil")
	}
	reqForm := make(url.Values)
	if err := formEncoder.Encode(param, reqForm); err != nil {
		return nil, err
	}
	return reqForm, nil
}

type SelfieAnimeResp struct {
	BaseResp
	Image string `json:"image"`
}

