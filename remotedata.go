package remotedata

import (
	"log"
	"sync"
	"time"

	linkedlist "github.com/474420502/focus/list/linked_list"
)

// RemoteData 远程数据
type RemoteData struct {
	currentCurl *linkedlist.CircularIterator
	targetCurl  *linkedlist.LinkedList

	update   time.Time
	interval time.Duration

	isUpdateWithInterval bool // 是否按照时间间隔更新
	isAsync              bool // 异步
	aysncData            *RemoteData

	updateContent     interface{}
	onUpdateCompleted func(content interface{}) (value interface{}, ok bool)
	onError           func(err error)

	updateMethod UpdateMethod

	valuelock sync.Mutex
	value     interface{}
}

// DefaultUpdateComplete 默认完成更新后的处理事件
var DefaultUpdateComplete = func(content interface{}) (value interface{}, ok bool) {
	return content, true
}

// New remotedata 必须由New创建
func New(updateMethod UpdateMethod) *RemoteData {
	rd := &RemoteData{}
	rd.targetCurl = linkedlist.New()
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

// SetAsync 设置不允许时间间隔更新
func (rd *RemoteData) SetAsync(is bool) {
	rd.isAsync = is
}

// SetDisableInterval 设置不允许时间间隔更新
func (rd *RemoteData) SetDisableInterval(is bool) {
	rd.isUpdateWithInterval = !is
}

// SetInterval 设置时间间隔. nil为不更新
func (rd *RemoteData) SetInterval(dur time.Duration) {
	rd.valuelock.Lock()
	defer rd.valuelock.Unlock()
	rd.interval = dur
	rd.isUpdateWithInterval = true
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

	rd.targetCurl.PushBack(c)
	if rd.currentCurl == nil {
		rd.currentCurl = rd.targetCurl.CircularIterator()
	}
}

// SetOnError 错误处理
func (rd *RemoteData) SetOnError(onError func(err error)) {
	rd.valuelock.Lock()
	defer rd.valuelock.Unlock()
	rd.onError = onError
}

// Value 获取值
func (rd *RemoteData) Value() interface{} {
	rd.valuelock.Lock()
	defer rd.valuelock.Unlock()
	if rd.isAsync {

	} else {
		rd.checkUpdate()
	}

	return rd.value
}

// Update 主动更新
func (rd *RemoteData) Update() {
	rd.valuelock.Lock()
	defer rd.valuelock.Unlock()
	rd.remoteUpdate()
	rd.update = time.Now()
}

func (rd *RemoteData) remoteUpdate() {

	if rd.updateMethod == nil {
		panic("UpdateMethod is nil. please Set this.")
	}

	var param interface{}

	if rd.currentCurl != nil {
		if rd.currentCurl.Next() {
			param = rd.currentCurl.Value()
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
	if rd.isUpdateWithInterval || rd.value == nil {
		now := time.Now()
		if now.Sub(rd.update) >= rd.interval {
			rd.remoteUpdate()
			rd.update = time.Now()
		}
	}
}
