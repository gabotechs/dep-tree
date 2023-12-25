package utils

import (
	"errors"
	"strings"

	"github.com/elliotchance/orderedmap/v2"
)

type CallStack struct {
	m *orderedmap.OrderedMap[string, bool]
}

func NewCallStack() *CallStack {
	return &CallStack{
		m: orderedmap.NewOrderedMap[string, bool](),
	}
}

// Push only errors if a cycle is detected.
func (cs *CallStack) Push(entry string) error {
	if _, ok := cs.m.Get(entry); ok {
		msg := "cycle detected:\n"
		for _, el := range cs.m.Keys() {
			msg += el + "\n"
		}
		msg += entry
		return errors.New(msg)
	} else {
		cs.m.Set(entry, true)
		return nil
	}
}

func (cs *CallStack) Pop() {
	last := cs.m.Back()
	if last != nil {
		cs.m.Delete(last.Key)
	}
}

func (cs *CallStack) Back() (string, bool) {
	back := cs.m.Back()
	if back == nil {
		return "", false
	} else {
		return back.Key, true
	}
}

func (cs *CallStack) Stack() []string {
	return cs.m.Keys()
}

func (cs *CallStack) Hash() string {
	return strings.Join(cs.m.Keys(), "-")
}
