package trib

type KeyValue struct {
	Key   string
	Value string
}

type Pattern struct {
	Prefix string
	Suffix string
}

type List struct {
	L []string
}

func KV(k, v string) *KeyValue { return &KeyValue{k, v} }

type Storage interface {
	// Return an auto-incrementing clock, the returned value
	// will be no smaller than atLeast, and it will be
	// strictly larger than the value returned last time,
	// unless it was math.MaxUint64.
	Clock(atLeast uint64, ret *uint64) error

	// Key-value pair interfaces
	// Default value for all keys is empty string
	Get(key string, value *string) error
	Set(kv *KeyValue, succ *bool) error
	Keys(p *Pattern, list *List) error

	// Key-list interfaces.
	// Default value for all lists is an empty list.
	// After the call, list.L should never by nil.
	ListGet(key string, list *List) error
	ListAppend(kv *KeyValue, succ *bool) error
	ListRemove(kv *KeyValue, n *int) error
	ListKeys(p *Pattern, list *List) error
}
