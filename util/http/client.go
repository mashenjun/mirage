package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/mashenjun/mirage/util"
)

// default timeout 1000ms
var DefaultClient = New(TimeoutOption(1000 * time.Millisecond))

// var DefaultClient = New()

type Client struct {
	*http.Client
}

func New(opts ...func(client *Client)) *Client {
	cli := Client{Client: &http.Client{Transport: http.DefaultTransport}}
	for _, opt := range opts {
		opt(&cli)
	}
	return &cli
}

func (r Client) DoRequestWithForm(
	ctx context.Context, method, uri string, data map[string][]string) (resp *http.Response, err error) {

	msg := url.Values(data).Encode()
	if method == "GET" || method == "HEAD" || method == "DELETE" {
		if strings.ContainsRune(uri, '?') {
			uri += "&"
		} else {
			uri += "?"
		}
		return r.DoRequest(ctx, method, uri+msg)
	}
	return r.DoRequestWith(
		ctx, method, uri, "application/x-www-form-urlencoded", strings.NewReader(msg), int64(len(msg)))
}

func (r Client) DoRequestWithFormHeader(
	ctx context.Context, method, uri string, data map[string][]string, h map[string]string) (resp *http.Response, err error) {

	msg := url.Values(data).Encode()
	if method == "GET" || method == "HEAD" || method == "DELETE" {
		if strings.ContainsRune(uri, '?') {
			uri += "&"
		} else {
			uri += "?"
		}
		return r.DoRequestWithHeader(ctx, method, uri+msg, "application/x-www-form-urlencoded", nil, 0, h)
	}
	return r.DoRequestWithHeader(
		ctx, method, uri, "application/x-www-form-urlencoded", strings.NewReader(msg), int64(len(msg)), h)
}

func (r Client) DoRequestWithJson(
	ctx context.Context, method, uri string, data interface{}) (resp *http.Response, err error) {

	msg, err := json.Marshal(data)
	if err != nil {
		return
	}
	return r.DoRequestWith(
		ctx, method, uri, "application/json", bytes.NewReader(msg), int64(len(msg)))
}

func (r Client) DoRequestWithJsonHeader(
	ctx context.Context, method, uri string, data interface{}, h map[string]string) (resp *http.Response, err error) {

	msg, err := json.Marshal(data)
	if err != nil {
		return
	}
	return r.DoRequestWithHeader(
		ctx, method, uri, "application/json", bytes.NewReader(msg), int64(len(msg)), h)
}

func (r Client) DoRequest(ctx context.Context, method, uri string) (resp *http.Response, err error) {

	req, err := http.NewRequest(method, uri, nil)
	if err != nil {
		return
	}
	return r.Do(ctx, req)
}

func (r Client) DoRequestWithHeader(
	ctx context.Context, method, uri string,
	bodyType string, body io.Reader, bodyLength int64, m map[string]string) (resp *http.Response, err error) {

	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", bodyType)
	for key, value := range m {
		req.Header.Set(key, value)
	}
	req.ContentLength = bodyLength
	return r.Do(ctx, req)
}

func (r Client) DoRequestWith(
	ctx context.Context, method, uri string,
	bodyType string, body io.Reader, bodyLength int64) (resp *http.Response, err error) {

	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", bodyType)
	req.ContentLength = bodyLength
	return r.Do(ctx, req)
}

func (r Client) Do(ctx context.Context, req *http.Request) (resp *http.Response, err error) {
	if ctx == nil {
		ctx = context.Background()
	}

	resp, err = r.Client.Do(req)
	if err != nil {
		recordRpcNetError(ctx, err, req)
	}

	return
}

func (r Client) CallWithForm(
	ctx context.Context, ret interface{}, method, url1 string, param map[string][]string) (err error) {

	resp, err := r.DoRequestWithForm(ctx, method, url1, param)
	if err != nil {
		return err
	}
	return r.callRet(ctx, ret, resp)
}

func (r Client) CallWithFormHeader(
	ctx context.Context, ret interface{}, method, url1 string, param map[string][]string, h map[string]string) (err error) {

	resp, err := r.DoRequestWithFormHeader(ctx, method, url1, param, h)
	if err != nil {
		return err
	}
	return r.callRet(ctx, ret, resp)
}

func (r Client) CallWithJson(
	ctx context.Context, ret interface{}, method, url1 string, param interface{}) (err error) {

	resp, err := r.DoRequestWithJson(ctx, method, url1, param)
	if err != nil {
		return err
	}
	return r.callRet(ctx, ret, resp)
}

func (r Client) CallWithJsonHeader(
	ctx context.Context, ret interface{}, method, url1 string, param interface{}, h map[string]string) (err error) {

	resp, err := r.DoRequestWithJsonHeader(ctx, method, url1, param, h)
	if err != nil {
		return err
	}
	return r.callRet(ctx, ret, resp)
}

func (r Client) callRet(ctx context.Context, ret interface{}, resp *http.Response) (err error) {

	defer func() {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close() // must close http body
	}()

	if ret == nil {
		return
	}

	err = json.NewDecoder(resp.Body).Decode(ret)
	if err != nil {
		recordBizError(ctx, err, nil, resp)
		return
	}

	// handle biz code error
	err = handleBizCodeError(ctx, ret, resp)

	return
}

func handleBizCodeError(ctx context.Context, ret interface{}, resp *http.Response) error {
	var err = fmt.Errorf("rpc biz error: %+v", ret)
	if ret != nil {
		retValue := reflect.Indirect(reflect.ValueOf(ret))
		if retValue.Kind() == reflect.Struct {
			codeVal := retValue.FieldByName("Code")
			if codeVal.IsValid() {
				if codeVal.Int() != 0 {
					recordBizError(ctx, err, ret, resp)
					return err
				}
			}
		} else {
			recordBizError(ctx, err, ret, resp)
			return err
		}
	}
	return nil
}

func recordBizError(ctx context.Context, err error, ret interface{}, resp *http.Response) {
	var (
		status string
		req    *http.Request
	)
	if resp != nil {
		status = strconv.Itoa(resp.StatusCode)
		req = resp.Request
	}
	metrics := getRequestCommonFields(ctx, req)


	util.RPCError.WithLabelValues(getRpcName(metrics.RemoteServiceName, metrics.UrlPath), status).Inc()
}

func recordRpcNetError(ctx context.Context, err error, req *http.Request) {
	metrics := getRequestCommonFields(ctx, req)

	// net error code default 599??
	util.RPCError.WithLabelValues(getRpcName(metrics.RemoteServiceName, metrics.UrlPath), "599").Inc()
}

type requestCommonMetrics struct {
	RemoteServiceName string
	Method            string
	UrlPath           string
	Query             string
	Body              string
}

func getRequestCommonFields(ctx context.Context, req *http.Request) (ret *requestCommonMetrics) {
	if req != nil {
		var body string
		if req.GetBody != nil {
			reader, err := req.GetBody()
			if err == nil {
				buf, _ := ioutil.ReadAll(reader)
				_ = reader.Close()
				body = string(buf)
			}
		}
		ret = &requestCommonMetrics{
			Method:            req.Method,
			UrlPath:           req.URL.Path,
			Query:             req.URL.RawQuery,
			Body:              body,
		}
	} else {
		ret = &requestCommonMetrics{}
	}
	return
}

func getRpcName(serviceName, urlPath string) string {
	return serviceName + urlPath
}
