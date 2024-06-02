package timewheel

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetCircleAndPos(t *testing.T) {
	interval := time.Millisecond * 10
	tw := NewTimeWheel(3, interval)
	circle, pos := tw.getCircleAndPos(0)
	assert.Equal(t, 0, circle)
	assert.Equal(t, 0, pos)

	tw.currentPos = 1
	circle, pos = tw.getCircleAndPos(4 * interval)
	assert.Equal(t, 1, circle)
	assert.Equal(t, 2, pos)
}

func TestTimeWheel(t *testing.T) {
	task1Result := 0
	task2Result := 0
	task3Result := 0
	interval := time.Millisecond * 10

	tw := NewTimeWheel(3, interval)
	tw.AddTask("task1", 2*interval, func() {
		task1Result = 1
	})
	tw.AddTask("task2", 4*interval, func() {
		task2Result = 2
	})
	tw.AddTask("task3", 3*interval, func() {
		task3Result = 3
	})
	go tw.Start()
	tw.RemoveTask("task3")

	time.Sleep(3*interval + time.Millisecond*5)
	assert.Equal(t, 1, task1Result)

	time.Sleep(2 * interval)
	assert.Equal(t, 2, task2Result)

	// decause task3 is removed, so task3Result should be 0
	assert.Equal(t, 0, task3Result)
}
