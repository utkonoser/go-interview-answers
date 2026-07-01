//go:build solution

package racecounter

import "sync/atomic"

// Counter — потокобезопасный счётчик.
// fix: data race на c.n++ — несколько горутин читают/пишут без синхронизации.
// atomic проще mutex для одного int; альтернатива — sync.Mutex.
type Counter struct {
	n atomic.Int64
}

func (c *Counter) Inc() {
	c.n.Add(1)
}

func (c *Counter) Value() int {
	return int(c.n.Load())
}
