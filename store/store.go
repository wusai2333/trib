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

func (self *Storage) Clock(_ int, ret *uint64) error {
	self.lock.Lock()
	defer self.lock.Unlock()

	if self.clock == math.MaxUint64 {
		return fmt.Errorf("clock overflow")
	}

	*ret = self.clock
	self.clock++

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

func (self *Storage) Keys(p *trib.Pattern, r *trib.List) error {
	ret := make([]string, 0, len(self.kvs))

	for k := range self.kvs {
		if !strings.HasPrefix(k, p.Prefix) {
			continue
		}
		if !strings.HasSuffix(k, p.Suffix) {
			continue
		}

		ret = append(ret, k)
	}

	r.L = ret
	return nil
}

func (self *Storage) List(key string, ret *trib.List) error {
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
