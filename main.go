package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"pigCache"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *pigCache.Group {
	return pigCache.NewGroup("scores", 2<<10, pigCache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

func startCacheServer(addr string, addrs []string, pig *pigCache.Group) {
	peers := pigCache.NewHTTPPool(addr) // 设置节点的路由
	peers.Set(addrs...)                 // 添加传入的节点，设置分布式节点
	pig.RegisterPeers(peers)
	log.Println("pigCache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func startAPIServer(apiAddr string, pig *pigCache.Group) {
	http.Handle("/get", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := pig.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Connection", "keep-alive")
			w.Write(view.ByteSlice())
		}))

	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func main() {

	var port int
	var api bool
	// 存储节点服务
	flag.IntVar(&port, "port", 8001, "pigCache server port")

	// 客户端服务
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	// api服务绑定了一个本地的geecache服务，这个服务miss的时候，才会去寻找其他结点
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	gee := createGroup()
	if api {
		go startAPIServer(apiAddr, gee)
	}
	startCacheServer(addrMap[port], []string(addrs), gee)
}
