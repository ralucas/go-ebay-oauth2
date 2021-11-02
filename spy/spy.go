package spy

import (
	"runtime"
	"strings"
	"sync"
)

type Spy struct {
	calls []*Call
	mu    sync.RWMutex
}

type Call struct {
	name string
	args []interface{}
}

func (s *Spy) Called(args ...interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	pc, _, _, _ := runtime.Caller(1)
	frame, _ := runtime.CallersFrames([]uintptr{pc}).Next()
	spl := strings.Split(frame.Function, ".")
	funcName := spl[len(spl)-1]

	c := &Call{funcName, args}
	s.calls = append(s.calls, c)
}

func (s *Spy) Calls(callName string) []*Call {
	calls := make([]*Call, 0)

	for _, call := range s.calls {
		if call.name == callName {
			calls = append(calls, call)
		}
	}

	return calls
}

func (s *Spy) CallCount(callName string) int {
	count := 0

	for _, call := range s.calls {
		if call.name == callName {
			count++
		}
	}

	return count
}

func (s *Spy) Reset() {
	s.calls = make([]*Call, 0)
}

func (c *Call) Arguments() []interface{} {
	return c.args
}
