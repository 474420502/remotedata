package remotedata

import (
	"log"
	"sync"
	"time"
)

// RemoteData 远程数据
type RemoteData struct {
	current int
	params  []interface{}

	updateWithInterval *updateInterval

	updateContent     interface{}
	onUpdateCompleted func(content interface{}) (value interface{}, ok bool)
	onError           func(err error)

	updateMethod UpdateMethod

	valuelock sync.Mutex
	value     interface{}
}

type updateInterval struct {
	Update          time.Time
	Interval        time.Duration
	UpdateCondition func(update time.Time) bool
}

// DefaultUpdateComplete 默认完成更新后的处理事件
var DefaultUpdateComplete = func(content interface{}) (value interface{}, ok bool) {
	return content, true
}

// New remotedata 必须由New创建
func New(updateMethod UpdateMethod) *RemoteData {
	rd := &RemoteData{}

	rd.updateMethod = updateMethod

	rd.onUpdateCompleted = DefaultUpdateComplete
	rd.onError = func(err error) {
		log.Println("default error handler:", err)
	}
	return rd
}

// Default 默认使用gcurl更新方法
func Default() *RemoteData {
	rd := New(MethodGcurl)
	return rd
}

// SetDisableInterval 设置不允许时间间隔更新. 默认true.(不以时间间隔更新)
func (rd *RemoteData) SetDisableInterval() {
	rd.valuelock.Lock()
	defer rd.valuelock.Unlock()
	// rd.isUpdateWithInterval = !is
	rd.updateWithInterval = nil
}

// SetInterval 设置时间间隔 同时SetDisableInterval(). nil为不更新.
func (rd *RemoteData) SetInterval(dur time.Duration) {
	rd.valuelock.Lock()
	defer rd.valuelock.Unlock()
	rd.updateWithInterval = &updateInterval{
		Interval: dur,
		UpdateCondition: func(update time.Time) bool {
			now := time.Now()
			if now.Sub(update) >= rd.updateWithInterval.Interval {
				return true
			}
			return false
		},
	}
}

// SetIntervalCondition 设置时间间隔 同时SetDisableInterval(). 不更新.
func (rd *RemoteData) SetIntervalCondition(cond func(update time.Time) bool) {
	rd.valuelock.Lock()
	defer rd.valuelock.Unlock()
	rd.updateWithInterval = &updateInterval{
		UpdateCondition: cond,
	}
}

// SetOnUpdateCompleted 设置更新后处理远程获取内容的事件
func (rd *RemoteData) SetOnUpdateCompleted(event func(content interface{}) (value interface{}, ok bool)) {
	rd.valuelock.Lock()
	defer rd.valuelock.Unlock()

	rd.onUpdateCompleted = event
}

// SetUpdateMethod 更新数据的方式
func (rd *RemoteData) SetUpdateMethod(method UpdateMethod) {
	rd.valuelock.Lock()
	defer rd.valuelock.Unlock()

	rd.updateMethod = method
}

// AddParam 添加UpdateDataMethod的参数. 多个AddParam会像循环链表一样调用. 例子: 多个地址负载
func (rd *RemoteData) AddParam(c interface{}) {
	rd.valuelock.Lock()
	defer rd.valuelock.Unlock()

	rd.params = append(rd.params, c)
}

// SetOnError 错误处理
func (rd *RemoteData) SetOnError(onError func(err error)) {
	rd.valuelock.Lock()
	defer rd.valuelock.Unlock()
	rd.onError = onError
}

// Value 获取值. 如果不设置SetInterval. 默认只更新一次. 可以使用Update做主动更新
func (rd *RemoteData) Value() interface{} {
	rd.valuelock.Lock()
	defer rd.valuelock.Unlock()
	rd.checkUpdate()
	return rd.value
}

// Update 主动更新
func (rd *RemoteData) Update() {
	rd.valuelock.Lock()
	defer rd.valuelock.Unlock()
	rd.remoteUpdate()

	if rd.updateWithInterval != nil {
		rd.updateWithInterval.Update = time.Now()
	}
}

func (rd *RemoteData) remoteUpdate() {

	if rd.updateMethod == nil {
		panic("UpdateMethod is nil. please Set this.")
	}

	var param interface{}

	if len(rd.params) != 0 {
		param = rd.params[rd.current]
		rd.current++
		if rd.current >= len(rd.params) {
			rd.current = 0
		}
	}

	data := rd.updateMethod(param)
	switch content := data.(type) {
	case nil:
	case error:
		if rd.onError != nil {
			rd.onError(content)
		}
	default:
		if value, ok := rd.onUpdateCompleted(content); ok {
			rd.value = value
			rd.updateContent = content
			return
		}
	}
}

func (rd *RemoteData) checkUpdate() { //

	if rd.value == nil {
		rd.remoteUpdate()
		if rd.updateWithInterval != nil {
			rd.updateWithInterval.Update = time.Now()
		}
	}

	if rd.updateWithInterval != nil {
		if rd.updateWithInterval.UpdateCondition(rd.updateWithInterval.Update) {
			rd.remoteUpdate()
			rd.updateWithInterval.Update = time.Now()
		}
	}
}
