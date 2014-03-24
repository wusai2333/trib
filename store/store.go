package store

import (
	"bytes"

	"trib"
)

type Storage struct {
	kvs map[string]*bytes.Buffer
}

var _ trib.Storage = new(Storage)

func NewStorage() *Storage {
	return &Storage{
		make(map[string]*bytes.Buffer),
	}
}

func (self *Storage) Get(key string, value *string) error {
	buf := self.kvs[key]
	if buf == nil {
		*value = ""
	} else {
		*value = buf.String()
	}

	return nil
}

func (self *Storage) Set(kv *trib.KeyValue, succ *bool) error {
	buf := self.kvs[kv.Key]
	if buf == nil {
		if kv.Value != "" {
			buf = new(bytes.Buffer)
			buf.WriteString(kv.Value)
			self.kvs[kv.Key] = buf
		}
	} else {
		if kv.Value != "" {
			buf.Reset()
			buf.WriteString(kv.Value)
		} else {
			delete(self.kvs, kv.Key)
		}
	}

	*succ = true
	return nil
}

func (self *Storage) Append(kv *trib.KeyValue, succ *bool) error {
	if kv.Value != "" {
		buf := self.kvs[kv.Key]
		if buf == nil {
			buf = new(bytes.Buffer)
			buf.WriteString(kv.Value)
			self.kvs[kv.Key] = buf
		} else {
			buf.WriteString(kv.Value)
		}
	}

	*succ = true
	return nil
}
