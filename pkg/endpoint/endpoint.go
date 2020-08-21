package endpoint

import (
	"github.com/gin-gonic/gin"

	apiErr "github.com/mashenjun/mirage/errors"
	"github.com/mashenjun/mirage/pkg/service"
	"github.com/mashenjun/mirage/util"
)

type Endpoint struct {
	srv *service.Service
}

func (ep *Endpoint) MountOn(router *gin.Engine) {
	router.POST("/advertise", ep.GetAdvertise)
	router.POST("/access_code", ep.GetAccessCode)
}

func New(srv *service.Service) (*Endpoint, error) {
	return &Endpoint{
		srv: srv,
	}, nil
}

func (ep *Endpoint) GetAdvertise(ctx *gin.Context) {
	param := service.GetAdvertiseParam{}
	if err := ctx.ShouldBindJSON(&param); err != nil {
		util.EncodeError(ctx, apiErr.ErrInvalidParameter(err.Error()))
		return
	}
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
