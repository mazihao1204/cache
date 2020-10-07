package lru

import "container/list"

type Cache struct {
	//允许使用的最大内存
	maxBytes int
	//当前已使用的内存
	nbytes int
	ll *list.List
	cache map[interface{}]*list.Element
	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func New(maxBytes int,onEvicted func(string,Value)) *Cache {
	return &Cache{
		maxBytes: maxBytes,
		ll: list.New(),
		cache: make(map[interface{}]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Add(key string,value Value){
	if c.cache == nil{
		c.cache = make(map[interface{}]*list.Element)
		c.ll = list.New()
	}
	if ee,ok := c.cache[key];ok{
		c.ll.PushFront(ee)
		kv := ee.Value.(*entry)
		c.nbytes += int(value.Len()) - int(kv.value.Len())
		kv.value = value
	}else {
		ele := c.ll.PushFront(&entry{key,value})
		c.cache[key] = ele
		c.nbytes += int(len(key)) + int(value.Len())
	}
	if c.maxBytes != 0 && c.maxBytes < c.nbytes{
		c.RemoveOldest()
	}
}

func (c *Cache) Get(key string) (value Value,ok bool){
	if ele,ok := c.cache[key];ok{
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value,true
	}
	return
}

func (c *Cache) Remove(key string){
	if c.cache == nil{
		return
	}
	if ele,ok := c.cache[key];ok{
		c.removeElement(ele)
	}
}

func (c *Cache) RemoveOldest(){
	if c.cache == nil{
		return
	}
	ele := c.ll.Back()
	if ele != nil{
		c.removeElement(ele)
	}
}

func (c *Cache) removeElement(e *list.Element){
	c.ll.Remove(e)
	kv := e.Value.(*entry)
	delete(c.cache,kv.key)
	if c.OnEvicted != nil{
		c.OnEvicted(kv.key,kv.value)
	}
}

func (c *Cache) Len() int{
	if c.cache == nil{
		return 0
	}
	return c.ll.Len()
}
