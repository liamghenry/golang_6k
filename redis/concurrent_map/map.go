package concurrentmap

import "sync"

// ConcurrentMap is a thread-safe map, divided int multi slot, every slot is a map with lock
type ConcurrentMap struct {
	slotNum int
	slots   []slot
}

type slot struct {
	sync.RWMutex
	m map[string]interface{}
}

func nextPowerOfTwo(n int) int {
	if n < 16 {
		return 16
	}
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n++
	return n
}

// NewConcurrentMap create a new ConcurrentMap
func NewConcurrentMap(slotNum int) *ConcurrentMap {
	// NOTE: power(2) -1 使得所有 bit 都是 1，后续用来 & fnv32 可以避免存在 0 位导致得到的 slotIndex 分布不均
	slotNum = nextPowerOfTwo(slotNum)
	slots := make([]slot, slotNum)
	for i := 0; i < slotNum; i++ {
		slots[i].m = make(map[string]interface{})
	}
	return &ConcurrentMap{
		slotNum: slotNum,
		slots:   slots,
	}
}

// Get get value by key, first get slot by key, then get value by key in slot, the slot number is hash(key) % slotNum
func (m *ConcurrentMap) Get(key string) (interface{}, bool) {
	slot := &m.slots[m.getSlotIndex(key)]
	slot.RLock()
	defer slot.RUnlock()
	v, ok := slot.m[key]
	return v, ok
}

// Set set value by key, return 0 if key not exist, 1 if key exist and update value
func (m *ConcurrentMap) Set(key string, value interface{}) int {
	slot := &m.slots[m.getSlotIndex(key)]
	slot.Lock()
	defer slot.Unlock()
	if _, ok := slot.m[key]; ok {
		slot.m[key] = value
		return 1
	}
	slot.m[key] = value
	return 0
}

// getSlotIndex get slot index by key, its fnv32 & len(slotNum) -1
func (m *ConcurrentMap) getSlotIndex(key string) uint32 {
	return fnv32(key) & uint32(m.slotNum - 1)
}

const (
	offset32 = 2166136261
	prime32  = 16777619
)

func fnv32(key string) uint32 {
	hash := uint32(offset32)
	for i := 0; i < len(key); i++ {
		hash ^= uint32(key[i])
		hash *= prime32
	}
	return hash
}