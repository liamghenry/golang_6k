package concurrentmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNextPowerOfTwo(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want int
	}{
		{
			name: "case1",
			n:    1,
			want: 16,
		},
		{
			name: "case2",
			n:    16,
			want: 16,
		},
		{
			name: "case3",
			n:    17,
			want: 32,
		},
		{
			name: "case4",
			n:    31,
			want: 32,
		},
		{
			name: "case5",
			n:    32,
			want: 32,
		},
		{
			name: "case6",
			n:    33,
			want: 64,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nextPowerOfTwo(tt.n); got != tt.want {
				t.Errorf("nextPowerOfTwo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetAndThenGet(t *testing.T) {
	cm := NewConcurrentMap(16)
	cm.Set("key", "value")
	v, ok := cm.Get("key")
	assert.Equal(t, true, ok)
	assert.Equal(t, "value", v)
}

func TestSetMulti(t *testing.T) {
	cm := NewConcurrentMap(16)

	cm.SetMulti([]string{"key1", "key2"}, []any{"value1", "value2"})

	v1, ok1 := cm.Get("key1")
	assert.Equal(t, true, ok1)
	assert.Equal(t, "value1", v1)

	v2, ok2 := cm.Get("key2")
	assert.Equal(t, true, ok2)
	assert.Equal(t, "value2", v2)
}
