package iface

import (
	"github.com/justcy/ygo/ygo/ytimer"
	"time"
)

type ITimeWheel interface {
	Start()
	Stop()
	Remove(task *ytimer.Task) error
	AddAfter(delay time.Duration, async bool, callback func(task *ytimer.Task), args []interface{}) *ytimer.Task
	AddCron(delay time.Duration, async bool, callback func(task *ytimer.Task), args []interface{}) *ytimer.Task
	AddAt(delay time.Duration, async bool, callback func(task *ytimer.Task), args []interface{}) *ytimer.Task
}
