// Package store provides a simple in-memory key value store.
package store

import (
	"container/list"
	"fmt"
	"math"
	"strings"
	"sync"

	"trib"
)

type strList []string

type Storage struct {
	clock uint64

	kvs   map[string]string
	lists map[string]*list.List
	lock  sync.Mutex
}

var _ trib.Storage = new(Storage)

func NewStorageId(id int) *Storage {
	return &Storage{
		kvs:   make(map[string]string),
		lists: make(map[string]*list.List),
	}
}

func NewStorage() *Storage {
	return NewStorageId(0)
}

func (self *Storage) Clock(atLeast uint64, ret *uint64) error {
	self.lock.Lock()
	defer self.lock.Unlock()

	self.clock++
	if self.clock < atLeast {
		self.clock = atLeast
	}
	if self.clock == math.MaxUint64 {
		return fmt.Errorf("clock overflow")
	}

	*ret = self.clock

	return nil
}

func (self *Storage) Get(key string, value *string) error {
	self.lock.Lock()
	defer self.lock.Unlock()

	*value = self.kvs[key]
	return nil
}

func (self *Storage) Set(kv *trib.KeyValue, succ *bool) error {
	self.lock.Lock()
	defer self.lock.Unlock()

	if kv.Value != "" {
		self.kvs[kv.Key] = kv.Value
	} else {
		delete(self.kvs, kv.Key)
	}

	*succ = true
	return nil
}

func matching(p *trib.Pattern, s string) bool {
	if !strings.HasPrefix(s, p.Prefix) {
		return false
	}
	if !strings.HasSuffix(s, p.Suffix) {
		return false
	}
	return true
}

func (self *Storage) Keys(p *trib.Pattern, r *trib.List) error {
	self.lock.Lock()
	defer self.lock.Unlock()

	ret := make([]string, 0, len(self.kvs))

	for k := range self.kvs {
		if matching(p, k) {
			ret = append(ret, k)
		}
	}

	r.L = ret
	return nil
}

func (self *Storage) ListKeys(p *trib.Pattern, r *trib.List) error {
	self.lock.Lock()
	defer self.lock.Unlock()

	ret := make([]string, 0, len(self.lists))
	for k := range self.lists {
		if matching(p, k) {
			ret = append(ret, k)
		}
	}

	r.L = ret
	return nil
}

func (self *Storage) ListGet(key string, ret *trib.List) error {
	self.lock.Lock()
	defer self.lock.Unlock()

	if lst, found := self.lists[key]; !found {
		ret.L = []string{}
	} else {
		ret.L = make([]string, 0, lst.Len())
		for i := lst.Front(); i != nil; i = i.Next() {
			ret.L = append(ret.L, i.Value.(string))
		}
	}

	return nil
}

func (self *Storage) ListAppend(kv *trib.KeyValue, succ *bool) error {
	self.lock.Lock()
	defer self.lock.Unlock()

	lst, found := self.lists[kv.Key]
	if !found {
		lst = list.New()
		self.lists[kv.Key] = lst
	}

	lst.PushBack(kv.Value)

	*succ = true
	return nil
}

func (self *Storage) ListRemove(kv *trib.KeyValue, n *int) error {
	self.lock.Lock()
	defer self.lock.Unlock()

	*n = 0

	lst, found := self.lists[kv.Key]
	if !found {
		return nil
	}

	i := lst.Front()
	for i != nil {
		if i.Value.(string) == kv.Value {
			hold := i
			i = i.Next()
			lst.Remove(hold)
			*n++
			continue
		}

		i = i.Next()
	}

	return nil
}
