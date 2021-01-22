package remotedata

import (
	"io/ioutil"
	"log"
	"testing"
	"time"

	"github.com/tidwall/gjson"
)

func init() {
	log.Println(`
	----- testing needs: 
	docker run -p 80:80 kennethreitz/httpbin
	echo "127.0.0.1   httpbin.org" >> /etc/hosts`)
}

var httpbinGetCurl = `curl 'http://httpbin.org/get' \
-H 'authority: www.httpbin.org' \
-H 'cache-control: max-age=0' \
-H 'upgrade-insecure-requests: 1' \
-H 'user-agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.66 Safari/537.36' \
-H 'accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9' \
-H 'sec-fetch-site: none' \
-H 'sec-fetch-mode: navigate' \
-H 'sec-fetch-user: ?1' \
-H 'sec-fetch-dest: document' \
-H 'accept-language: zh-CN,zh;q=0.9,ja;q=0.8' \
--compressed`

func TestCaseIntervalUpdate(t *testing.T) {
	var completedCount = 0
	var create = func() *RemoteData {
		data := Default()
		cbash := httpbinGetCurl
		data.AddParam(cbash)
		data.SetOnUpdateCompleted(func(content interface{}) (value interface{}, ok bool) {
			r := gjson.ParseBytes(content.([]byte))
			average := r.Get("headers.Host")
			if average.Exists() {
				completedCount++
				return average.Float(), true
			}
			return nil, false
		})

		data.SetInterval(time.Second * 4)
		return data
	}
	data := create()
	// data.SetDisableInterval(true)
	// ４秒更新间隔, 应该２次更新
	for i := 0; i < 5; i++ {
		// t.Error(data.Value())
		data.Value()
		time.Sleep(time.Second)
	}
	if completedCount < 2 {
		t.Error("completedCount != 2, completedCount = ", completedCount)
	}

	completedCount = 0
	data = create()
	data.SetDisableInterval(true) // 不按时间间隔更新
	for i := 0; i < 5; i++ {
		// t.Error(data.Value())
		data.Value()
		time.Sleep(time.Second)
	}

	if completedCount != 1 {
		t.Error("completedCount != 1, completedCount = ", completedCount)
	}

	data.Update() // 主动更新
	if completedCount != 2 {
		t.Error("completedCount != 2, completedCount = ", completedCount)
	}

	// t.Error(data.Value())
	// t.Error(data.Value())
	// t.Error(data.Value())
}

func TestHttpGet(t *testing.T) {
	data := New(MethodHTTPGet)
	data.AddParam("http://httpbin.org/get")
	if data.Value() == nil { // 默认更新一次
		t.Error("error get")
	}

	if data.Value() == nil { // 默认更新一次
		t.Error("error get")
	}
}

func TestOnError(t *testing.T) {
	data := New(MethodGcurl)
	data.SetUpdateMethod(MethodHTTPGet)
	data.AddParam("http://httpbin.org:1/get")
	data.SetOnError(func(err error) {
		if err == nil {
			t.Error(err)
		}
	})
	data.Value()

}

func TestNoParam(t *testing.T) {
	data := New(func(param interface{}) interface{} {
		return 1
	})

	if data.Value().(int) != 1 {
		t.Error("TestNoParam error")
	}
}

func TestReadFile(t *testing.T) {
	data := New(func(param interface{}) interface{} {
		data, err := ioutil.ReadFile("test.json")
		if err != nil {
			return err
		}
		return data
	})

	data.SetInterval(time.Second * 1) // 每一次更新一次. 可以做配置热更新

	if string(data.Value().([]byte)) != `{ "a": 1, "b": 2 }` {
		t.Error("TestReadFile error")
	}

}

func TestReadFile1(t *testing.T) {
	data := New(MethodReadFile)
	data.AddParam("test.json")
	data.SetInterval(time.Second * 1) // 每一次更新一次. 可以做配置热更新

	if string(data.Value().([]byte)) != `{ "a": 1, "b": 2 }` {
		t.Error("TestReadFile error")
	}
}
