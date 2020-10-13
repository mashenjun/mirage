package endpoint

import (
	"net/http"

	"github.com/gin-gonic/gin"

	apiErr "github.com/mashenjun/mirage/errors"
	"github.com/mashenjun/mirage/log"
	"github.com/mashenjun/mirage/pkg/service"
	"github.com/mashenjun/mirage/util"
)

type Endpoint struct {
	srv *service.Service
}

// @title Mirage Backend API
// @version 1.0
// @description This is mirage backend server.
// @contact.name Mirage Backend Support
// @contact.url https://github.com/mashenjun/mirage/issues
// @BasePath /api/v1
// MountOn prepare the router
func (ep *Endpoint) MountOn(router *gin.Engine) {
	router.GET("/ping", ep.Ping)
	apiV1 := router.Group("/api/v1")
	apiV1.GET("/config/advertise", ep.GetAdvertise)
	apiV1.GET("/config/template", ep.GetTemplates)
	apiV1.GET("/config/version_update", ep.GetVersionUpdate)
	apiV1.POST("/face/detect", ep.DetectFace)
	apiV1.POST("/face/edit_attr", ep.EditAttr)
	apiV1.POST("/face/merge", ep.MergeFace)
	apiV1.POST("/image_process/style_trans", ep.StyleTrans)
	apiV1.POST("/image_process/selie_anime", ep.SelieAnime)
	apiV1.POST("/image_process/body_seg", ep.BodySeg)
	apiV1.POST("/upload_signature", ep.UploadSignature)
}

func New(srv *service.Service) (*Endpoint, error) {
	return &Endpoint{
		srv: srv,
	}, nil
}

// @Tags 静态配置API
// @Summary 获取广告配置
// @Description 获取广告配置
// @Produce json
// @Param type query string true "广告code"
// @Param extra query string false "额外信息"
// @Success 200 {object} util.BaseResp{data=model.AdvConfig} "广告配置"
// @Failure 500 {object} errors.ErrorInfo "服务异常"
// @Router /config/advertise [get]
func (ep *Endpoint) GetAdvertise(ctx *gin.Context) {
	param := service.GetAdvertiseParam{}
	param.AdCode = ctx.Query("type")
	param.Extra = ctx.Query("extra")
	data, err := ep.srv.GetAdvertise(ctx, param)
	if err != nil {
		util.EncodeError(ctx, err)
		return
	}
	util.EncodeResp(ctx, data)
}

// @Tags 静态配置API
// @Summary 获取素材模板配置
// @Description 获取素材模板配置
// @Produce json
// @Param type query string true "素材模板名称"
// @Success 200 {object} util.BaseResp{data=model.TemplateImageConfig} "素材模板配置"
// @Failure 500 {object} errors.ErrorInfo "服务异常"
// @Router /config/template [get]
func (ep *Endpoint) GetTemplates(ctx *gin.Context) {
	param := service.GetTemplateParam{}
	param.Type = ctx.Query("type")
	data, err := ep.srv.GetTemplate(ctx, param)
	if err != nil {
		util.EncodeError(ctx, err)
		return
	}
	util.EncodeResp(ctx, data)
}

// @Tags 静态配置API
// @Summary 获取版本配置
// @Description 获取版本配置
// @Produce json
// @Success 200 {object} util.BaseResp{data=[]byte} "版本配置"
// @Failure 500 {object} errors.ErrorInfo "服务异常"
// @Router /config/template [get]
func (ep *Endpoint) GetVersionUpdate(ctx *gin.Context) {
	data, err := ep.srv.GetVersionUpdate(ctx)
	if err != nil {
		util.EncodeError(ctx, err)
		return
	}
	util.EncodeBytes(ctx, data)
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

// @Tags 人脸处理API
// @Summary 检测人脸年龄
// @Description 检测人脸年龄
// @Accept  json
// @Produce json
// @Param body body service.DetectFaceParam true "图片地址"
// @Success 200 {object} util.BaseResp{data=service.DetectFaceData} "检测结果，如果没有检测出人脸，code不为0"
// @Failure 400 {object} errors.ErrorInfo "请求参数不正确"
// @Failure 500 {object} errors.ErrorInfo "服务异常"
// @Router /face/detect [post]
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

// @Tags 人脸处理API
// @Summary 人脸变老，变年轻，变性别
// @Description 人脸变老，变年轻，变性别
// @Accept  json
// @Produce json
// @Param body body service.EditAttrParam true "json parameter"
// @Success 200 {object} util.BaseResp{data=service.EditAttrData} "处理结果"
// @Failure 400 {object} errors.ErrorInfo "请求参数不正确"
// @Failure 500 {object} errors.ErrorInfo "服务异常"
// @Router /face/edit_attr [post]
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
	log.Info("SelieAnime")
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

// @Tags 人脸处理API
// @Summary 人脸融合
// @Description 人脸融合
// @Accept  json
// @Produce json
// @Success 200 {object} util.BaseResp{data=service.MergeFaceData} "处理结果"
// @Failure 400 {object} errors.ErrorInfo "请求参数不正确"
// @Failure 500 {object} errors.ErrorInfo "服务异常"
// @Router /face/merge [post]
func (ep *Endpoint) MergeFace(ctx *gin.Context) {
	param := service.MergeFaceParam{}
	if err := ctx.ShouldBindJSON(&param); err != nil {
		util.EncodeError(ctx, apiErr.ErrInvalidParameter(err.Error()))
		return
	}
	data, err := ep.srv.MergeFace(ctx, param)
	if err != nil {
		util.EncodeError(ctx, err)
		return
	}
	util.EncodeResp(ctx, data)
}

func (ep *Endpoint) BodySeg(ctx *gin.Context) {
	param := service.BodySegParam{}
	if err := ctx.ShouldBindJSON(&param); err != nil {
		util.EncodeError(ctx, apiErr.ErrInvalidParameter(err.Error()))
		return
	}
	data, err := ep.srv.BodySeg(ctx, param)
	if err != nil {
		util.EncodeError(ctx, err)
		return
	}
	util.EncodeResp(ctx, data)
}
