package ytimer

import (
	"fmt"
	"testing"
	"time"
)

func SayHello(t *Task) {
	fmt.Printf("%s 执行次数 %d hello %s,%s\n", time.Now(),t.circleTimes, t.args[0], t.args[1])
}
func TestNewTimeWheel(t *testing.T) {
	tw, _ := NewTimeWheel(1*time.Second, 10, TickSafeMode())
	tw.AddAfter(1*time.Second, modeIsAsync, SayHello, []interface{}{"zhangsan", "122"})
	task2 := tw.AddCron(2*time.Second, modeIsAsync, SayHello, []interface{}{"lisi", "1111"})
	tw.AddAt(time.Duration(time.Now().UnixNano()+int64(3*time.Second)), modeIsAsync, SayHello, []interface{}{"wangwu", "233445"})
	tw.Start()
	time.Sleep(10 * time.Second)
	//tw.Remove(task)
	tw.Remove(task2)
	tw.Stop()
	time.Sleep(20 * time.Second)
}
