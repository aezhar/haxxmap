package haxmap

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type atomicUInt32 struct {
	// Signal to go vet not to copy this type
	_ sync.Locker
	v uint32
}

type atomicPointer[T any] struct {
	// Signal to go vet not to copy this type
	_   sync.Locker
	ptr unsafe.Pointer
}

type atomicUintptr struct {
	// Signal to go vet not to copy this type
	_   sync.Locker
	ptr uintptr
}

func (u *atomicUInt32) Load() uint32            { return atomic.LoadUint32(&u.v) }
func (u *atomicUInt32) Store(v uint32)          { atomic.StoreUint32(&u.v, v) }
func (u *atomicUInt32) Add(delta uint32) uint32 { return atomic.AddUint32(&u.v, delta) }
func (u *atomicUInt32) Swap(v uint32) uint32    { return atomic.SwapUint32(&u.v, v) }
func (u *atomicUInt32) CompareAndSwap(old, new uint32) bool {
	return atomic.CompareAndSwapUint32(&u.v, old, new)
}

func (p *atomicPointer[T]) Load() *T     { return (*T)(atomic.LoadPointer(&p.ptr)) }
func (p *atomicPointer[T]) Store(v *T)   { atomic.StorePointer(&p.ptr, unsafe.Pointer(v)) }
func (p *atomicPointer[T]) Swap(v *T) *T { return (*T)(atomic.SwapPointer(&p.ptr, unsafe.Pointer(v))) }
func (p *atomicPointer[T]) CompareAndSwap(old, new *T) bool {
	return atomic.CompareAndSwapPointer(&p.ptr, unsafe.Pointer(old), unsafe.Pointer(new))
}

func (u *atomicUintptr) Load() uintptr             { return atomic.LoadUintptr(&u.ptr) }
func (u *atomicUintptr) Store(v uintptr)           { atomic.StoreUintptr(&u.ptr, v) }
func (u *atomicUintptr) Add(delta uintptr) uintptr { return atomic.AddUintptr(&u.ptr, delta) }
func (u *atomicUintptr) Swap(v uintptr) uintptr    { return atomic.SwapUintptr(&u.ptr, v) }
func (u *atomicUintptr) CompareAndSwap(old, new uintptr) bool {
	return atomic.CompareAndSwapUintptr(&u.ptr, old, new)
}
