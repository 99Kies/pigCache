package lru

import (
	"fmt"
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}

func TestAddAndGet(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("key", String("1"))
	lru.Add("key", String("111"))

	val, ok := lru.Get("key")
	if ok {
		fmt.Println(val)
	} else {
		fmt.Println("get error")
	}

	lru.RemoveOldest()
	val, ok = lru.Get("key")
	if ok {
		fmt.Println(val)
	} else {
		fmt.Println("get error")
	}
}
