package transport

import (
	"context"
	"fmt"
	"github.com/zhangliang-zl/reskit/logs"
	client "github.com/zhangliang-zl/reskit/transport/http"
	"net/http"
	"testing"
	"time"
)

var stdHttpClient Client
var timeoutHttpClient Client

func init() {
	l := logs.DefaultLogger("_caller")
	domain := "http://api.match.test.mararun.cn:8086/"
	stdHttpClient = New(domain, 300*time.Second, l)
	timeoutHttpClient = New(domain, 1*time.Nanosecond, l)
}

func TestHttpCallerCallGet(t *testing.T) {
	uri := "v1/hotCity/lists"
	param := map[string]interface{}{"name": "北京"}
	var result struct {
		Datas []struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"datas"`
		Total string `json:"total"`
	}
	err := stdHttpClient.Get(context.TODO(), uri, param, &result)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Datas) == 0 {
		t.Errorf("没找到name为%s的数据", param["name"])
	}
}

func TestHttpCallerCallPostForm(t *testing.T) {
	uri := "v1/hotCity/create"
	param := map[string]interface{}{
		"city[name]":        "testForm" + fmt.Sprint(time.Now().Unix()),
		"city[country]":     "中国大陆",
		"city[countryType]": "1",
		"city[province]":    "河北省",
	}
	var result struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}
	err := stdHttpClient.PostForm(context.TODO(), uri, param, &result)
	if err != nil {
		t.Fatal(err)
	}
	if result.Id == "0" {
		t.Fatal(err)
	}
}

func TestHttpCallerCallPostJson(t *testing.T) {
	uri := "v1/hotCity/create"
	var result struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}
	type City struct {
		Name        string `json:"name"`
		Country     string `json:"country"`
		CountryType int    `json:"countryType"`
		Province    string `json:"province"`
	}
	type HotCity struct {
		City `json:"city"`
	}
	param := HotCity{City{
		Name:        "testJson" + fmt.Sprint(time.Now().Unix()),
		Country:     "中国大陆",
		CountryType: 1,
		Province:    "河北省",
	}}
	err := stdHttpClient.PostJson(context.TODO(), uri, param, &result)
	if err != nil {
		t.Fatal(err)
	}
	if result.Id == "0" {
		t.Fatal(err)
	}
}

func TestHttpCallerDo(t *testing.T) {
	uri := "v1/hotCity/create"
	var result struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}
	type City struct {
		Name        string `json:"name"`
		Country     string `json:"country"`
		CountryType int    `json:"countryType"`
		Province    string `json:"province"`
	}
	type HotCity struct {
		City `json:"city"`
	}
	param := HotCity{City{
		Name:        "testDo" + fmt.Sprint(time.Now().Unix()),
		Country:     "中国大陆",
		CountryType: 1,
		Province:    "河北省",
	}}
	jsonParam := client.BodyJSON(&param)
	err := stdHttpClient.Do(context.TODO(), http.MethodPost, uri, &result, jsonParam)
	if err != nil {
		t.Fatal(err)
	}
	if result.Id == "0" {
		t.Fatal(err)
	}
}

func TestErrorInterface(t *testing.T) {
	uri := "v1/notExist/get"
	param := map[string]interface{}{"id": 1}
	result := make(map[string]interface{})
	err := stdHttpClient.Get(context.TODO(), uri, param, &result)
	if httpErr, ok := err.(Error); !ok {
		t.Fatal(httpErr)
	}
}

func TestTimeout(t *testing.T) {
	uri := "v1/example/get"
	param := map[string]interface{}{"id": 1}
	result := make(map[string]interface{})
	err := timeoutHttpClient.Get(context.TODO(), uri, param, &result)
	httpErr, ok := err.(Error)
	if !(ok && httpErr.Error() == "服务器错误") {
		t.Fatal(httpErr)
	}
}
