package util

import "sync/atomic"

type Lock struct {
	seed int32
}

// 判断并“加锁”
func (l *Lock) Lock() bool {
	//对比参数返回结果
	return atomic.CompareAndSwapInt32(&l.seed, 0, 1)
}

// “解锁”
func (l *Lock) Unlock() {
	//新赋值参数值
	atomic.StoreInt32(&l.seed, 0)
}
