# pigCache
Distributed Cache System by Golang

# 相关名词
- lru
- cache
- group
- peers

**lru -> cache -> group -> peers**



在底层通过 lru 实现缓存存储的具体逻辑(hash表+双链表实现的一个lru算法)。
通过 cache 模块对lru模块做了一层并发处理的封装，通过线程锁实现并发处理 `mu.Lock()` and `mu.UnLock()`。
通过 cache 模块，我们解决了数据在分布式读写时候，可能会遇到的读写冲突的问题，之后我们在 cache 模块的上层做了一层**资源**的概念 - `group`, 一个group结构体代表的是一个节点（更或者说是一个namespace的一环真实节点），可以把这个group当成数据库里的表或者是库。每个group里都会有一个环形状真实节点，每个真实节点下又有多个虚拟节点，客户端是通过访问真实节点的api，间接访问虚拟节点中存储的数据。

** LRU **

LRU（Least Recently Used）是一种常见的页面置换算法，在计算中，所有的文件操作都要放在内存中进行，然而计算机内存大小是固定的，所以我们不可能把所有的文件都加载到内存，因此我们需要制定一种策略对加入到内存中的文件进项选择。

精髓：删除最近最久未使用的数据，将新输入置为队首
