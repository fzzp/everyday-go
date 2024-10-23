package lru

import "container/list"

type Cache struct {
	maxBytes int                      // 允许使用的最大内存
	nbytes   int64                    //  已使用的内存
	ll       *list.List               // Go 语言标准库实现的双向链表
	cache    map[string]*list.Element // 键是字符串，值是双向链表中对应节点的指针

	// 删除元素时回调
	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func New() {}
