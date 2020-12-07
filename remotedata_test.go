package remotedata

import (
	"testing"
	"time"

	"github.com/tidwall/gjson"
)

func TestCase(t *testing.T) {
	var completedCount = 0
	var create = func() *RemoteData {
		data := Default()
		cbash := `curl 'https://www.xe.com/zh-CN/api/stats.php?fromCurrency=VND&toCurrency=CNY' \
	-H 'Accept: application/json, text/plain, */*' \
	-H 'Referer: https://www.xe.com/zh-CN/currencyconverter/convert/?Amount=1&From=VND&To=CNY' \
	-H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.66 Safari/537.36' \
	--compressed`
		data.AddParam(cbash)
		data.SetOnUpdateCompleted(func(content interface{}) (value interface{}, ok bool) {
			r := gjson.ParseBytes(content.([]byte))
			average := r.Get("payload.Last_30_Days.average")
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
