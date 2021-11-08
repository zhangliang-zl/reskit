package transport

import (
	"context"
	"encoding/json"
	"github.com/zhangliang-zl/reskit/logs"
	"github.com/zhangliang-zl/reskit/transport/http"
	"net/http"
	"time"
)

type HttpClient struct {
	httpClient *client.Req
	urlPrefix  string
	logger     logs.Logger
}

func New(urlPrefix string, timeout time.Duration, logger logs.Logger) Client {
	req := client.New()
	req.SetTimeout(timeout)
	return &HttpClient{urlPrefix: urlPrefix, logger: logger, httpClient: req}
}

func (c *HttpClient) Get(ctx context.Context, uri string, param map[string]interface{}, result interface{}) ErrorInterface {
	queryParam := client.QueryParam(param)
	return c.Do(ctx, http.MethodGet, uri, result, queryParam)
}

func (c *HttpClient) PostForm(ctx context.Context, uri string, formData map[string]interface{}, result interface{}) ErrorInterface {
	postParam := client.Param(formData)
	return c.Do(ctx, http.MethodPost, uri, result, postParam)
}

func (c *HttpClient) PostJson(ctx context.Context, uri string, jsonData interface{}, result interface{}) ErrorInterface {
	jsonParam := client.BodyJSON(&jsonData)
	return c.Do(ctx, http.MethodPost, uri, result, jsonParam)
}

func (c *HttpClient) Do(ctx context.Context, method string, uri string, result interface{}, params ...interface{}) ErrorInterface {
	url := c.makeUrl(uri)
	params = append(params, ctx)
	resp, err := c.httpClient.Do(method, url, params...)
	if err == nil {
		err = c.parseResp(resp, result)
	}
	return c.wrapHttpError(ctx, err)
}

func (c *HttpClient) makeUrl(url string) string {
	return c.urlPrefix + url
}

func (c *HttpClient) parseResp(resp *client.Resp, result interface{}) error {
	var httpErr Error
	byteData, err := resp.ToBytes()
	if err != nil {
		return err
	}
	if err := json.Unmarshal(byteData, &httpErr); err != nil {
		return err
	}
	if httpErr.GetCode() != 0 {
		return httpErr
	}
	if err := json.Unmarshal(byteData, result); err != nil {
		return err
	}
	return nil
}

func (c *HttpClient) wrapHttpError(ctx context.Context, err error) ErrorInterface {
	if httpErr, ok := err.(Error); ok {
		return httpErr
	} else if err != nil {
		c.logger.Error(ctx, err.Error())
		return NewInternalError("服务器错误")
	} else {
		return nil
	}
}
