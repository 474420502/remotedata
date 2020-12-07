package remotedata

import (
	"github.com/474420502/gcurl"
	"github.com/474420502/requests"
)

// UpdateMethod 更新方法. 可以是自定义
type UpdateMethod func(param interface{}) interface{}

// GCurlMethod 这个是由本人gcurl实现的方法. 支持curlbash直接请求
var GCurlMethod UpdateMethod = func(param interface{}) interface{} {
	c := param.(string)
	ses := gcurl.Parse(c)
	resp, err := ses.Temporary().Execute()
	if err == nil {
		return resp.Content()
	}
	return err
}

// HTTPGetMethod Http get请求.
var HTTPGetMethod UpdateMethod = func(param interface{}) interface{} {
	hurl := param.(string)
	resp, err := requests.NewSession().Get(hurl).Execute()
	if err == nil {
		return resp.Content()
	}
	return err
}
