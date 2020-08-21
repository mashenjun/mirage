package util

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"

	apiErr "github.com/mashenjun/mirage/errors"
)

type BaseResp struct {
	apiErr.ErrorInfo
	Data interface{} `json:"data"`
}

func EncodeError(gCtx *gin.Context, err error) {
	var (
		code      int
		errDetail apiErr.ErrorInfo
	)

	if errInfo, ok := err.(apiErr.ErrorInfo); ok {
		code = errInfo.StatusCode()
		errDetail = errInfo
	} else {
		errInfo := apiErr.ErrInternal(err.Error())
		code = errInfo.StatusCode()
		errDetail = errInfo
	}

	serverErrors.WithLabelValues(gCtx.Request.Method, gCtx.Request.URL.Path,
		fmt.Sprintf("%d", errDetail.Code), errDetail.Message).Inc()

	gCtx.AbortWithStatusJSON(code, errDetail)
}

func EncodeResp(gCtx *gin.Context, data interface{}) {
	resp := BaseResp{}
	resp.Message = "ok"
	resp.Data = struct{}{}
	if data != nil {
		resp.Data = data
	}
	gCtx.JSON(http.StatusOK, resp)
}

func IsToday(t time.Time) bool {
	return time.Now().Format("2006-01-02") == t.Local().Format("2006-01-02")
}

func CalSeconds(str string) (int64, error) {
	t, err := time.ParseInLocation("15:04:05", str, time.Local)
	if err != nil {
		return 0, err
	}
	return int64(t.Hour()*3600 + t.Minute()*60 + t.Second()), nil
}

func CalDate(base string, offset int) (string, error) {
	baseTime, err := time.ParseInLocation("2006-01-02", base, time.Local)
	if err != nil {
		return "", err
	}
	newTime := baseTime.AddDate(0, 0, offset)
	return newTime.Format("2006-01-02"), nil
}

var (
	serverErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "rest",
			Subsystem: "server",
			Name:      "errors_total",
			Help:      "Total number of errors that happen during RESTful processing on server side.",
		},
		[]string{"method", "uri", "err_code", "err_message"},
	)

	daysInCN = [...]string{
		"星期日",
		"星期一",
		"星期二",
		"星期三",
		"星期四",
		"星期五",
		"星期六",
	}
)

func init() {
	prometheus.MustRegister(serverErrors)
}

func ConvertToWeekCN(date time.Time) string {
	d := date.Weekday()
	if time.Sunday <= d && d <= time.Saturday {
		return daysInCN[date.Weekday()]
	}
	return d.String()
}
