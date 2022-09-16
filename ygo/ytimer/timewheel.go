package ytimer

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)
const (
	ModeIsCircle  = true
	ModeNotCircle = false

	ModeIsAsync  = true
	ModeNotAsync = false
)

type TimeWheel struct {
	randomId int64 //自增任务ID

	tick      time.Duration
	ticker    *time.Ticker
	tickQueue chan time.Time

	bucketsNum    int
	buckets       []map[taskId]*Task // key: added item, value: *Task
	bucketIndexes map[taskId]int     // key: added item, value: bucket position

	currentIndex int

	onceStart sync.Once

	addChan    chan *Task
	removeChan chan *Task
	stopChan   chan struct{}

	exited   bool
	syncPool bool

	//互斥锁（继承RWMutex的 RWLock,UnLock 等方法）
	sync.RWMutex
}
// NewTimeWheel create new time wheel
func NewTimeWheel(tick time.Duration, bucketsNum int, options ...optionCall) (*TimeWheel, error) {
	if tick.Seconds() < 0.1 {
		return nil, errors.New("invalid params, must tick >= 100 ms")
	}
	if bucketsNum <= 0 {
		return nil, errors.New("invalid params, must bucketsNum > 0")
	}

	tw := &TimeWheel{
		// tick
		tick:      tick,
		tickQueue: make(chan time.Time, 10),

		// store
		bucketsNum:    bucketsNum,
		bucketIndexes: make(map[taskId]int, 1024*100),
		buckets:       make([]map[taskId]*Task, bucketsNum),
		currentIndex:  0,

		// signal
		addChan:    make(chan *Task, 1024*5),
		removeChan: make(chan *Task, 1024*2),
		stopChan:   make(chan struct{}),
	}

	for i := 0; i < bucketsNum; i++ {
		tw.buckets[i] = make(map[taskId]*Task, 16)
	}

	for _, op := range options {
		op(tw)
	}

	return tw, nil
}
func (tw *TimeWheel) Start() {
	tw.onceStart.Do(
		func() {
			tw.ticker = time.NewTicker(tw.tick)
			go tw.run()
			go tw.tickGenerator()
		}, )
}
func (tw *TimeWheel) tickGenerator() {
	if tw.tickQueue != nil {
		return
	}

	for !tw.exited {
		select {
		case <-tw.ticker.C:
			select {
			case tw.tickQueue <- time.Now():
			default:
				panic("raise long time blocking")
			}
		}
	}
}

func (tw *TimeWheel) Stop() {
	tw.stopChan <- struct{}{}
}

func (tw *TimeWheel) AddAfter(delay time.Duration, async bool, callback func(task *Task), args []interface{}) *Task{
	return tw.createTask(delay, ModeNotCircle, async, callback, args)
}

func (tw *TimeWheel) AddCron(delay time.Duration, async bool, callback func(task *Task), args []interface{}) *Task{
	return tw.createTask(delay, ModeIsCircle, async, callback, args)
}

func (tw *TimeWheel) AddAt(unixNano time.Duration , async bool, callback func(task *Task), args []interface{}) *Task{
	delay := int64(unixNano) - time.Now().UnixNano()
	if delay <= 0 {
		errors.New("invalid unixNano, must unixNano >= now ")
	}
	return tw.createTask(time.Duration(delay), ModeNotCircle, async, callback, args)
}

func (tw *TimeWheel) addTask(task *Task,circleMode bool) {
	if task.callback == nil {
		errors.New("task callback is nil")
	}
	index,round := tw.getIndexAndCircle(task.delay)
	if round > 0 && circleMode {
		task.round = round - 1
	} else {
		task.round = round
	}
	tw.bucketIndexes[task.id] = index
	tw.buckets[index][task.id] = task
}
func (tw *TimeWheel) Remove(task *Task) error {
	tw.removeChan <- task
	return nil
}
func (tw *TimeWheel) removeTask(task *Task) {
	tw.collectTask(task)
}

func (tw *TimeWheel) createTask(delay time.Duration, circle, async bool, callback func(task *Task), args []interface{}) *Task {
	if delay <= 0 {
		delay = tw.tick
	}
	var task *Task
	if tw.syncPool {
		task = defaultTaskPool.get()
		task.callback = callback
		task.Args = args
	}else{
		task = NewTask(callback, args)
	}
	task.id = tw.genUniqueId()
	task.circle = circle
	task.async = async
	task.delay = delay
	tw.addChan <- task
	return task
}
func (tw *TimeWheel) run() {
	queue := tw.ticker.C
	if tw.tickQueue == nil {
		queue = tw.tickQueue
	}
	for {
		select {
		case <-queue:
			tw.handleTick()
		case task := <-tw.addChan:
			tw.addTask(task, ModeNotCircle)
		case task := <-tw.removeChan:
			tw.removeTask(task)
		case <-tw.stopChan:
			tw.exited = true
			tw.ticker.Stop()
			return
		}
	}
}
type optionCall func(*TimeWheel) error

func TickSafeMode() optionCall {
	return func(tw *TimeWheel) error {
		tw.tickQueue = make(chan time.Time, 10)
		return nil
	}
}

func SetSyncPool(state bool) optionCall {
	return func(tw *TimeWheel) error {
		tw.syncPool = state
		return nil
	}
}
func (tw *TimeWheel) collectTask(task *Task) {
	index := tw.bucketIndexes[task.id]
	delete(tw.bucketIndexes, task.id)
	delete(tw.buckets[index], task.id)

	if tw.syncPool {
		defaultTaskPool.put(task)
	}
}
func (tw *TimeWheel) handleTick() {
	bucket := tw.buckets[tw.currentIndex]

	for k, task := range bucket {
		if task.stop {
			tw.collectTask(task)
			continue
		}
		if bucket[k].round > 0 {
			bucket[k].round--
			continue
		}
		if task.callback != nil{
			task.circleTimes ++
			task.Run()
		}

		// circle
		if task.circle == ModeIsCircle {
			//tw.collectTask(task)
			tw.addTask(task, ModeIsCircle)
			continue
		}
		// gc
		tw.collectTask(task)
	}

	if tw.currentIndex == tw.bucketsNum-1 {
		tw.currentIndex = 0
		return
	}

	tw.currentIndex++
}
// get the task position
func (tw *TimeWheel) getIndexAndCircle(d time.Duration) (index int, circle int) {
	delaySeconds := int(d.Seconds())
	intervalSeconds := int(tw.tick.Seconds())
	circle = int(delaySeconds / intervalSeconds / tw.bucketsNum)
	index = int(tw.currentIndex+delaySeconds/intervalSeconds) % tw.bucketsNum
	return
}
func (tw *TimeWheel) genUniqueId() taskId {
	id := atomic.AddInt64(&tw.randomId, 1)
	return taskId(id)
}