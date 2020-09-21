package faceai

const (
	accessTokenPath = "/oauth/2.0/token"
	editAttrPath    = "/rest/2.0/face/v1/editattr"
	detectPath      = "/rest/2.0/face/v3/detect"
	styleTransPath  = "/rest/2.0/image-process/v1/style_trans"
	selfieAnimePath = "/rest/2.0/image-process/v1/selfie_anime"
	mergeFacePath   = "/rest/2.0/face/v1/merge"
	bodySegPath     = "/rest/2.0/image-classify/v1/body_seg"

	accessTokenKey = "accesstoken"
)

const (
	ActionTypeToKid    = "TO_KID"
	ActionTypeToOld    = "TO_OLD"
	ActionTypeToFemale = "TO_FEMALE"
	ActionTypeToMale   = "TO_MALE"
)
