# use curl or other method update remote data
# 更新远程的数据. 用于缓存. 按时间间隔更新. 可以是数据库数据. 页面数据. 接口数据. 减少负载, 增加多个地址. 稳定更新.


## http get url data. 
## http get url 方式更新. 按时间间隔. 或者主动更新
```go 
func TestHttpGet(t *testing.T) {
	data := New(MethodHTTPGet)
	data.AddParam("http://httpbin.org/get")
	if data.Value() == nil { // 默认更新一次 only update onces
		t.Error("error get")
	}
	// data.Value() == response (curl http://httpbin.org/get)
	if data.Value() == nil { // 默认更新一次 only update onces
		t.Error("error get")
	}
}
```

## 设置自定义更新方法. 可以是多种方式.
```go 
func TestNoParam(t *testing.T) {
	data := New(func(param interface{}) interface{} {
		//DO: you get data method
		return 1
	})

	if data.Value().(int) != 1 {
		t.Error("TestNoParam error")
	}
}
```

## 时间间隔更新文件. file grpc sqldriver http都可以
```go 
func TestReadFile(t *testing.T) {
	data := New(func(param interface{}) interface{} {
		data, err := ioutil.ReadFile("test.json")
		if err != nil {
			return err
		}
		return data
	})

	data.SetInterval(time.Second * 1) // 每一秒更新一次. 可以做配置热更新.最好的还是用fsnotify

	if string(data.Value().([]byte)) != `{ "a": 1, "b": 2 }` {
		t.Error("TestReadFile error")
	}
}
```

## 主动更新文件. file grpc sqldriver http都可以
```go 
func TestReadFile(t *testing.T) {
	data := New(func(param interface{}) interface{} {
		data, err := ioutil.ReadFile("test.json")
		if err != nil {
			return err
		}
		return data
	})

	data.SetDisableInterval(true) // 每一次Value()更新一次. 可以做配置热更新.最好的还是用fsnotify

	if string(data.Value().([]byte)) != `{ "a": 1, "b": 2 }` {
		t.Error("TestReadFile error")
	}
}
```


