package http_utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/valyala/fasthttp"
	"mime/multipart"
	"net/http"
	"time"
)

type ServerConfig struct {
	Host    string
	Port    string
	Timeout int
}

// 服务后端向其他http服务发送请求时的http客户端对象结构体
type httpClient struct {
	Request      *fasthttp.Request
	Response     *fasthttp.Response
	ServerConfig *ServerConfig
}

func Client(config *ServerConfig) *httpClient {
	return &httpClient{
		Request:      fasthttp.AcquireRequest(),
		Response:     fasthttp.AcquireResponse(),
		ServerConfig: config,
	}
}

func (c *httpClient) Close() {
	defer c.Request.ConnectionClose()
	defer c.Response.ConnectionClose()
}

func (c *httpClient) SetRequestURI(uri string) {
	c.Request.SetRequestURI(uri)
}

func (c *httpClient) request() ([]byte, error) {
	c.Request.SetHost(c.ServerConfig.Host + ":" + c.ServerConfig.Port)
	c.Request.Header.SetContentType("application/json")
	if err := fasthttp.DoTimeout(c.Request, c.Response, time.Duration(c.ServerConfig.Timeout)*time.Second); err != nil {
		return nil, err
	} else {
		if len(c.Response.Body()) == 0 {
			return nil, errors.New("no response content return")
		}
		if c.Response.StatusCode() != http.StatusOK {
			return nil, errors.New(fmt.Sprintf("reponse status code:%d", c.Response.StatusCode()))
		}
		result := ResponseJSON{}
		if err = json.Unmarshal(c.Response.Body(), &result); err != nil {
			return nil, err
		}
		if result.Code >= http.StatusBadRequest {
			return nil, errors.New(fmt.Sprintf("bad request with message:%s,code:%d,data:%v", result.Message, result.Code, result.Object))
		}
		if body, err := json.Marshal(result.Object); err != nil {
			return nil, err
		} else {
			return body, nil
		}
	}
}

func (c *httpClient) Get() ([]byte, error) {
	c.Request.Header.SetMethod("GET")
	return c.request()
}

func (c *httpClient) Post(body []byte) ([]byte, error) {
	c.Request.Header.SetMethod("POST")
	c.Request.SetBody(body)
	return c.request()
}

type FormContent struct {
	Type    int // 1=file 2=form
	FiName  string
	Content []byte
}

func (c *httpClient) PostForm(fm map[string]FormContent) ([]byte, error) {
	var (
		body = new(bytes.Buffer)
	)
	writer := multipart.NewWriter(body)
	if len(fm) == 0 {
		return nil, errors.New("form arg map is empty")
	}
	for k, v := range fm {
		if v.Type == 1 {
			filePart, err := writer.CreateFormFile(k, v.FiName)
			if err != nil {
				return nil, err
			}
			_, _ = filePart.Write(v.Content)
		} else if v.Type == 2 {
			formPart, err := writer.CreateFormField(k)
			if err != nil {
				return nil, err
			}
			_, _ = formPart.Write(v.Content)
		}
	}
	_ = writer.Close()
	bd := body.Bytes()
	if bd == nil {
		return nil, errors.New("body is empty")
	}
	c.Request.Header.SetContentType(writer.FormDataContentType())
	c.Request.Header.SetMethod("POST")
	c.Request.SetHost(c.ServerConfig.Host + ":" + c.ServerConfig.Port)
	c.Request.SetBody(bd)
	if err := fasthttp.DoTimeout(c.Request, c.Response, time.Duration(c.ServerConfig.Timeout)*time.Second); err != nil {
		return nil, err
	} else {
		if len(c.Response.Body()) == 0 {
			return nil, errors.New("no response content return")
		}
		if c.Response.StatusCode() != http.StatusOK {
			return nil, errors.New(fmt.Sprintf("reponse status code:%d", c.Response.StatusCode()))
		}
		result := ResponseJSON{}
		if err = json.Unmarshal(c.Response.Body(), &result); err != nil {
			return nil, err
		}
		if result.Code >= http.StatusBadRequest {
			return nil, errors.New(fmt.Sprintf("bad request with message:%s,code:%d,data:%v", result.Message, result.Code, result.Object))
		}
		if body, err := json.Marshal(result.Object); err != nil {
			return nil, err
		} else {
			return body, nil
		}
	}
}
