package timewheel

import "time"

var tw = NewTimeWheel(3600, time.Second)

func init() {
	tw.Start()
}

func At(key string, ttl time.Time, fn func()) {
	tw.AddTask(key, time.Until(ttl), fn)
}
