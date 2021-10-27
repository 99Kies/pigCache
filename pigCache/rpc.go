package pigCache

import (
	"sync"
)

type RPCGetter interface {
	Get(key string) ([]byte, error)
}

type RPCGetterFunc func(key string) ([]byte, error)

func (f RPCGetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type RPCGroup struct {
	name      string
	getter    Getter
	mainCache cache
}

var (
	rpcMu     sync.RWMutex
	rpcGroups = make(map[string]*Group)
)
