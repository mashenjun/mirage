package service

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"testing"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/aliyun-sts-go-sdk/sts"
	"github.com/oklog/ulid/v2"
)

var (
	b64 = ""
	ak  = ""
	sk  = ""
	arn = ""
)

func TestService_FetchAndBase64(t *testing.T) {
	ossCli, err := oss.New("http://oss-cn-hongkong.aliyuncs.com", ak, sk)
	if err != nil {
		t.Fatal(err)
	}
	bucket, err := ossCli.Bucket("mirage-test")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(bucket)
	r, err := bucket.GetObject("image/autistic-child6.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(base64.StdEncoding.EncodeToString(buf))
}

func TestService_Upload(t *testing.T) {
	ossCli, err := oss.New("http://oss-cn-hongkong.aliyuncs.com", ak, sk)
	if err != nil {
		t.Fatal(err)
	}
	bucket, err := ossCli.Bucket("mirage-test")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(bucket)
	b, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		t.Fatal(err)
	}
	uuid, err := ulid.New(ulid.Now(), nil)
	if err != nil {
		t.Fatal(err)
	}
	err = bucket.PutObject(uuid.String(), bytes.NewBuffer(b), oss.ContentType("image/jpg"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(uuid.String())
}

func TestService_UploadSignature(t *testing.T) {
	client := sts.NewClient(ak, sk, arn, "123")
	resp, err := client.AssumeRole(uint(3600))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}
