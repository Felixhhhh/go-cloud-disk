package test

import (
	"bytes"
	"context"
	"fmt"
	_ "fmt"
	"go-cloud-disk/core/define"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/tencentyun/cos-go-sdk-v5"
)

func TestFileUploadByFilepath(t *testing.T) {
	u, _ := url.Parse(define.CosBucket)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.TencentSecretID,
			SecretKey: define.TencentSecretKey,
		},
	})

	key := "go-cloud-disk/exampleobject.jpg"

	_, _, err := client.Object.Upload(
		context.Background(), key, "./img/9.png", nil,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFileUploadByReader(t *testing.T) {
	u, _ := url.Parse(define.CosBucket)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.TencentSecretID,
			SecretKey: define.TencentSecretKey,
		},
	})

	key := "go-cloud-disk/exampleobject2.jpg"

	f, err := os.ReadFile("./img/9.png")
	if err != nil {
		return
	}
	_, err = client.Object.Put(
		context.Background(), key, bytes.NewReader(f), nil,
	)
	if err != nil {
		t.Fatal(err)
	}
}

//分片上传初始化
func TestInitPartUpload(t *testing.T) {
	u, _ := url.Parse(define.CosBucket)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.TencentSecretID,
			SecretKey: define.TencentSecretKey,
		},
	})

	key := "go-cloud-disk/exampleTest.jpg"
	v, _, err := client.Object.InitiateMultipartUpload(context.Background(), key, nil)
	if err != nil {
		t.Fatal(err)
	}
	UploadID := v.UploadID // 16710039863102499e6a451e168f085f3cb5f254a2ec1e289054be89613e7f08e16195795b
	fmt.Println(UploadID)
}

//分片上传
func TestPartUpload(t *testing.T) {
	u, _ := url.Parse(define.CosBucket)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.TencentSecretID,
			SecretKey: define.TencentSecretKey,
		},
	})

	key := "go-cloud-disk/exampleTest.jpg"
	UploadID := "16710039863102499e6a451e168f085f3cb5f254a2ec1e289054be89613e7f08e16195795b"
	f, err := os.ReadFile("0.chunk") // md5 : a8a41a0935a4e32c59850b6defa074c4
	if err != nil {
		t.Fatal(err)
	}
	// opt可选
	resp, err := client.Object.UploadPart(
		context.Background(), key, UploadID, 1, bytes.NewReader(f), nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	PartETag := resp.Header.Get("ETag")
	fmt.Println(PartETag)
}

//分片上传完成
func TestPartUploadComplete(t *testing.T) {
	u, _ := url.Parse(define.CosBucket)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.TencentSecretID,
			SecretKey: define.TencentSecretKey,
		},
	})

	key := "go-cloud-disk/exampleTest.jpg"
	UploadID := "16710039863102499e6a451e168f085f3cb5f254a2ec1e289054be89613e7f08e16195795b"

	opt := &cos.CompleteMultipartUploadOptions{}
	opt.Parts = append(opt.Parts, cos.Object{
		PartNumber: 1, ETag: "a8a41a0935a4e32c59850b6defa074c4"},
	)
	_, _, err := client.Object.CompleteMultipartUpload(
		context.Background(), key, UploadID, opt,
	)
	if err != nil {
		t.Fatal(err)
	}
}
