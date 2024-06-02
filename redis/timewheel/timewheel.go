package timewheel

import (
	"container/list"
	"time"
)

type location struct {
	slot int
	e    *list.Element
}
type task struct {
	circle int
	job    func()
	key    string
}

type TimeWheel struct {
	slotNum int
	// list element is task
	slots    []*list.List
	interval time.Duration
	timer    map[string]*location
	ticker   *time.Ticker

	currentPos int
}

func NewTimeWheel(slotNum int, interval time.Duration) *TimeWheel {
	tw := &TimeWheel{
		slotNum:    slotNum,
		slots:      make([]*list.List, slotNum),
		interval:   interval,
		timer:      make(map[string]*location),
		currentPos: 0,
	}
	for i := 0; i < slotNum; i++ {
		tw.slots[i] = list.New()
	}
	return tw
}

func (tw *TimeWheel) Start() {
	ticker := time.NewTicker(tw.interval)
	tw.ticker = ticker
	go func() {
		for range ticker.C {
			tw.tickHandler()
		}
	}()
}

func (tw *TimeWheel) tickHandler() {
	l := tw.slots[tw.currentPos]
	tw.scanAndRun(l)
	tw.currentPos = (tw.currentPos + 1) % tw.slotNum
}

func (tw *TimeWheel) scanAndRun(l *list.List) {
	for e := l.Front(); e != nil; {
		task := e.Value.(*task)
		if task.circle > 0 {
			task.circle--
			e = e.Next()
			continue
		}
		go task.job()
		next := e.Next()
		l.Remove(e)
		delete(tw.timer, task.key)
		e = next
	}
}

func (tw *TimeWheel) Stop() {
	tw.ticker.Stop()
}

func (tw *TimeWheel) AddTask(key string, delay time.Duration, job func()) {
	circle, pos := tw.getCircleAndPos(delay)
	e := tw.slots[pos].PushBack(&task{
		circle: circle,
		job:    job,
		key:    key,
	})
	tw.timer[key] = &location{
		slot: pos,
		e:    e,
	}
}

func (tw *TimeWheel) getCircleAndPos(delay time.Duration) (circle int, pos int) {
	afterPos := tw.currentPos + int(delay)/int(tw.interval)
	circle = afterPos / tw.slotNum
	pos = afterPos % tw.slotNum
	return
}

func (tw *TimeWheel) RemoveTask(key string) {
	if loc, ok := tw.timer[key]; ok {
		tw.slots[loc.slot].Remove(loc.e)
		delete(tw.timer, key)
	}
}
