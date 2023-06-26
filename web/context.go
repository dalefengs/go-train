package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

type Context struct {
	Req        *http.Request
	Resp       http.ResponseWriter
	PathParams map[string]string

	// 缓存路由参数， Form 有自带缓存
	QueryValues url.Values
}

// ResponseJSON 响应 JSON 数据
func (c *Context) ResponseJSON(status int, val any) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	c.Resp.WriteHeader(status)
	c.Resp.Header().Set("Content-Type", "application/json")
	c.Resp.Header().Set("Content-Length", strconv.Itoa(len(data)))
	_, err = c.Resp.Write(data)
	return err
}

func (c *Context) ResponseJSONOK(val any) error {
	return c.ResponseJSON(200, val)
}

func (c *Context) ResponseJSONError(val any) error {
	return c.ResponseJSON(500, val)
}

// SetCookie 设置 Cookie
func (c *Context) SetCookie(ck *http.Cookie) {
	http.SetCookie(c.Resp, ck)
}

// BindJSON 绑定 json
func (c *Context) BindJSON(val any) error {
	if val == nil {
		return errors.New("web: 输入不能为空")
	}
	if c.Req.Body == nil {
		return errors.New("web: body 为空")
	}
	decoder := json.NewDecoder(c.Req.Body)
	return decoder.Decode(val)
}

// FormValue 获取 Form 表单数据
func (c *Context) FormValue(key string) StringValue {
	err := c.Req.ParseForm()
	if err != nil {
		return NewStringValue("", err)
	}
	return NewStringValue(c.Req.FormValue(key), nil)
}

// QueryValue 获取路由参数
func (c *Context) QueryValue(key string) StringValue {
	if c.QueryValues != nil {
		c.QueryValues = c.Req.URL.Query()
	}
	vals, ok := c.QueryValues[key]
	if !ok {
		return NewStringValue("", errors.New("web：key 不存在"))
	}
	return NewStringValue(vals[0], nil)
}

// PathValue 获取路由路径参数
func (c *Context) PathValue(key string) StringValue {
	val, ok := c.PathParams[key]
	if !ok {
		return NewStringValue("", errors.New("web：key 不存在"))
	}
	return NewStringValue(val, nil)
}

// StringValue 用于 String 字符串转换
type StringValue struct {
	val string
	err error
}

func NewStringValue(val string, err error) StringValue {
	return StringValue{val: val, err: err}
}

func (s *StringValue) String() (string, error) {
	if s.err != nil {
		return " ", s.err
	}
	return s.val, nil
}

func (s *StringValue) AsInt64() (int64, error) {
	if s.err != nil {
		return 0, s.err
	}
	return strconv.ParseInt(s.val, 10, 64)
}

func (s *StringValue) AsInt() (int, error) {
	if s.err != nil {
		return 0, s.err
	}
	return strconv.Atoi(s.val)
}

func (s *StringValue) AsFloat64() (float64, error) {
	if s.err != nil {
		return 0, s.err
	}
	return strconv.ParseFloat(s.val, 64)
}
