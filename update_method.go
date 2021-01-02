package remotedata

import (
	"io/ioutil"

	"github.com/474420502/gcurl"
	"github.com/474420502/requests"
)

// UpdateMethod 更新方法. 可以是自定义
type UpdateMethod func(param interface{}) interface{}

// MethodGcurl 这个是由本人gcurl实现的方法. 支持curlbash直接请求
var MethodGcurl UpdateMethod = func(param interface{}) interface{} {
	c := param.(string)
	ses := gcurl.Parse(c)
	resp, err := ses.Temporary().Execute()
	if err == nil {
		return resp.Content()
	}
	return err
}

// MethodHTTPGet Http get请求.
var MethodHTTPGet UpdateMethod = func(param interface{}) interface{} {
	hurl := param.(string)
	resp, err := requests.NewSession().Get(hurl).Execute()
	if err == nil {
		return resp.Content()
	}
	return err
}

// MethodReadFile 读文件更新方法
var MethodReadFile UpdateMethod = func(param interface{}) interface{} {
	filepath := param.(string)
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	return data
}

// Context 上下文
// type Context struct {
// 	param interface{}
// 	share interface{} //用于更新共享数据. 传递上下文, 记录...
// }

// // GetShare Get return share interface{}
// func (cxt *Context) GetShare() interface{} {
// 	return cxt.share
// }

// // SetShare Set share interface{}
// func (cxt *Context) SetShare(share interface{}) {
// 	cxt.share = share
// }

// // GetParam 获取添加参数, 循环列表每次参数都会轮循. 便于多个不同地址的数据循环更新. 减少目标的负载
// func (cxt *Context) GetParam() interface{} {
// 	return cxt.param
// }
