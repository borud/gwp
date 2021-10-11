package gwp

import "sync/atomic"

// AtomicBool is just a helper type that gives us an atomic bool. The null
// value is 0, which means false.
type AtomicBool int32

// IsTrue is true if it is true :-).
func (b *AtomicBool) IsTrue() bool {
	return atomic.LoadInt32((*int32)(b)) != 0
}

// SetTrue sets the value to true.
func (b *AtomicBool) SetTrue() {
	atomic.StoreInt32((*int32)(b), 1)
}

// SetFalse sets the value to false.
func (b *AtomicBool) SetFalse() {
	atomic.StoreInt32((*int32)(b), 0)
}
