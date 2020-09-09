package endpoint

import (
	"net/http"

	"github.com/gin-gonic/gin"

	apiErr "github.com/mashenjun/mirage/errors"
	"github.com/mashenjun/mirage/pkg/service"
	"github.com/mashenjun/mirage/util"
)

type Endpoint struct {
	srv *service.Service
}

func (ep *Endpoint) MountOn(router *gin.Engine) {
	router.GET("/ping", ep.Ping)
	apiV1 := router.Group("/api/v1")
	apiV1.GET("/config/advertise", ep.GetAdvertise)
	apiV1.POST("/face/detect", ep.DetectFace)
	apiV1.POST("/face/edit_attr", ep.EditAttr)
	apiV1.POST("/image_process/style_trans", ep.StyleTrans)
	apiV1.POST("/image_process/selie_anime", ep.SelieAnime)
	apiV1.POST("/upload_signature", ep.UploadSignature)
}

func New(srv *service.Service) (*Endpoint, error) {
	return &Endpoint{
		srv: srv,
	}, nil
}

func (ep *Endpoint) GetAdvertise(ctx *gin.Context) {
	param := service.GetAdvertiseParam{}
	param.AdCode = ctx.Param("ad_code")
	data, err := ep.srv.GetAdvertise(ctx, param)
	if err != nil {
		util.EncodeError(ctx, err)
		return
	}
	util.EncodeResp(ctx, data)
}

func (ep *Endpoint) GetAccessCode(ctx *gin.Context) {
	param := service.GetAccessCodeParam{}
	if err := ctx.ShouldBindJSON(&param); err != nil {
		util.EncodeError(ctx, apiErr.ErrInvalidParameter(err.Error()))
		return
	}
	data, err := ep.srv.GetAccessCode(ctx, param)
	if err != nil {
		util.EncodeError(ctx, err)
		return
	}
	util.EncodeResp(ctx, data)
}

func (ep *Endpoint) Ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, "ok")
}

func (ep *Endpoint) DetectFace(ctx *gin.Context) {
	param := service.DetectFaceParam{}
	if err := ctx.ShouldBindJSON(&param); err != nil {
		util.EncodeError(ctx, apiErr.ErrInvalidParameter(err.Error()))
		return
	}
	data, err := ep.srv.DetectFace(ctx, param)
	if err != nil {
		util.EncodeError(ctx, err)
		return
	}
	util.EncodeResp(ctx, data)
}

func (ep *Endpoint) EditAttr(ctx *gin.Context) {
	param := service.EditAttrParam{}
	if err := ctx.ShouldBindJSON(&param); err != nil {
		util.EncodeError(ctx, apiErr.ErrInvalidParameter(err.Error()))
		return
	}
	data, err := ep.srv.EditAttr(ctx, param)
	if err != nil {
		util.EncodeError(ctx, err)
		return
	}
	util.EncodeResp(ctx, data)
}

func (ep *Endpoint) StyleTrans(ctx *gin.Context) {
	param := service.StyleTransParam{}
	if err := ctx.ShouldBindJSON(&param); err != nil {
		util.EncodeError(ctx, apiErr.ErrInvalidParameter(err.Error()))
		return
	}
	data, err := ep.srv.StyleTrans(ctx, param)
	if err != nil {
		util.EncodeError(ctx, err)
		return
	}
	util.EncodeResp(ctx, data)
}

func (ep *Endpoint) SelieAnime(ctx *gin.Context) {
	param := service.SelfieAnimeParam{}
	if err := ctx.ShouldBindJSON(&param); err != nil {
		util.EncodeError(ctx, apiErr.ErrInvalidParameter(err.Error()))
		return
	}
	data, err := ep.srv.SelieAnime(ctx, param)
	if err != nil {
		util.EncodeError(ctx, err)
		return
	}
	util.EncodeResp(ctx, data)
}

func (ep *Endpoint) UploadSignature(ctx *gin.Context) {
	data, err := ep.srv.UploadSignature(ctx)
	if err != nil {
		util.EncodeError(ctx, err)
		return
	}
	util.EncodeResp(ctx, data)
}