# use curl or other method update remote data

```go 
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
	// ４秒更新间隔, 应该２次更新
	for i := 0; i < 5; i++ {
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
}

```