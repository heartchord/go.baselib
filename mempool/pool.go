package mempool

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/heartchord/goblazer/racedetect"
)

// Pool implements a generic object pool, based on sync.Pool(/Go/src/sync/pool.go).
//     Improvement: 1. add statistics function
//                  2. create local pool in initialization, unlike pinSlow in sync.Pool.
type Pool struct {
	localPool    unsafe.Pointer // localPool set, the number is specified by runtime.GOMAXPROCS(0).
	localPoolNum uintptr        // the number of localPool set. it equals to runtime.GOMAXPROCS(0).
	localPoolCap int            // the capacity of each localPool. len(localPool.shared) + 1(localPool.private) <= localPoolCap
	newFunc      NewFunc        // the new function to create a new object when none is in all localPools.
}

// localPool holds all appendixes  belonged to Per-P-Pool .
type localPool struct {
	private    interface{}   // can be used only by the respective P
	shared     []interface{} // can be used by any P
	sync.Mutex               // protects shared when accessing localPool.shared
	pad        [128]byte     // prevents false sharing
}

// NewPool creates a new PCachedPool.
func NewPool(localPoolCap int, newFunc NewFunc) *Pool {
	if localPoolCap <= 0 {
		panic(fmt.Sprintf("[NewPool error] expected localPoolCap > 0, got localPoolCap = %d", localPoolCap))
	}

	procs := runtime.GOMAXPROCS(0)
	local := make([]localPool, procs)

	return &Pool{
		localPool:    unsafe.Pointer(&local[0]),
		localPoolNum: uintptr(procs),
		localPoolCap: localPoolCap,
		newFunc:      newFunc,
	}
}

// Reset the pool to initialized state.
func (p *Pool) Reset() {
	procs := runtime.GOMAXPROCS(0)
	local := make([]localPool, procs)
	atomic.StoreUintptr(&p.localPoolNum, uintptr(procs))
	atomic.StorePointer(&p.localPool, unsafe.Pointer(&local[0]))
}

// Get obtains an object from the pool. The steps are P.private => P.shared => Other P.shared => new
func (p *Pool) Get() interface{} {
	if racedetect.Enabled {
		return p.getByNewFunc()
	}

	l := p.pin()
	if l == nil { // can't find P's localPool, GOMAXPROCS(0) may be chaged, so just call newFunc to create a new object.
		syncRuntimeProcUnpin()
		return p.getByNewFunc()
	}

	// under pin() operation, juse access l.private.
	x := l.private
	l.private = nil
	syncRuntimeProcUnpin()
	if x != nil {
		return x
	}

	// try to obtain one object from l.shared.
	l.Lock()
	last := len(l.shared) - 1
	if last >= 0 { // obtain the last object
		x = l.shared[last]
		l.shared = l.shared[:last]
	}
	l.Unlock()

	if x != nil {
		return x
	}

	// try to steal from other P' shared.
	return p.getSlow()
}

// Put returns an object to the pool. The steps are P.private => P.shared => Other P.shared => drop
func (p *Pool) Put(x interface{}) {
	if x == nil {
		return
	}

	if racedetect.Enabled {
		return
	}

	l := p.pin()
	if l == nil { // can't find P's localPool, GOMAXPROCS(0) may be chaged, so just call newFunc to create a new object.
		syncRuntimeProcUnpin()
		return
	}

	// under pin() operation, juse access l.private.
	if l.private == nil {
		l.private = x
		x = nil
	}
	syncRuntimeProcUnpin()

	if x == nil {
		return
	}

	// try to returns one object to l.shared.
	l.Lock()
	if len(l.shared)+1 < p.localPoolCap {
		l.shared = append(l.shared, x)
		x = nil
	}
	l.Unlock()
}

// getSlow trys to steal one object from other procs or call newFunc to create a new one.
func (p *Pool) getSlow() (x interface{}) {
	num := atomic.LoadUintptr(&p.localPoolNum)
	local := atomic.LoadPointer(&p.localPool)

	pid := syncRuntimeProcPin()
	syncRuntimeProcUnpin()

	for i := 0; i < int(num); i++ {
		l := getLocalPoolIdx(local, (pid+i+1)%int(num))
		l.Lock()
		last := len(l.shared) - 1
		if last >= 0 {
			x = l.shared[last]
			l.shared = l.shared[:last]
			l.Unlock()
			break
		}
		l.Unlock()
	}

	if x == nil {
		x = p.getByNewFunc()
		return
	}

	return
}

func (p *Pool) getByNewFunc() (x interface{}) {
	if p.newFunc != nil {
		x = p.newFunc()
	}
	return
}

func (p *Pool) pin() *localPool {
	n := atomic.LoadUintptr(&p.localPoolNum)
	l := atomic.LoadPointer(&p.localPool)

	pid := syncRuntimeProcPin()
	if uintptr(pid) < n {
		return getLocalPoolIdx(l, pid)
	}

	// GOMAXPROCS(n) changes to a larger value. This will happen when Pool.Reset() is invoked.
	// When local pool can't be found, the Get() will call newFunc to create a new object, then Put() will juse throw the object away.
	return nil
}

func getLocalPoolIdx(local unsafe.Pointer, i int) *localPool {
	return &(*[1000000]localPool)(local)[i]
}

//go:linkname syncRuntimeProcPin sync.runtime_procPin
func syncRuntimeProcPin() int

//go:linkname syncRuntimeProcUnpin sync.runtime_procUnpin
func syncRuntimeProcUnpin()
