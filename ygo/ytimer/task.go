package ytimer

import (
	"fmt"
	log "github.com/justcy/ygo/ygo/ylog"
	"reflect"
	"time"
)

type taskId int64
type Task struct {
	delay    time.Duration
	id       taskId
	round    int
	callback func(task *Task)
	args     []interface{}

	async  bool
	stop   bool
	circle bool
	circleTimes int64   //循环次数
}

//NewTask 创建一个延迟调用任务
func NewTask(f func(v *Task), args []interface{}) *Task {
	return &Task{
		callback: f,
		args:     args,
	}
}
func (t *Task) Run() {
	defer func() {
		if err := recover(); err != nil {
			log.Error(t.String(), "Call err: ", err)
		}
	}()
	if t.async {
		go t.callback(t)
	} else {
		t.callback(t)
	}
}
// for sync.Pool
func (t *Task) Reset() {
	t.round = 0
	t.callback = nil
	t.args = nil

	t.async = false
	t.stop = false
	t.circle = false
	t.circleTimes = 0
}


func (t *Task) String() interface{} {
	return fmt.Sprintf("{Task:%s, args:%v}", reflect.TypeOf(t.callback).Name(), t.args)
}
