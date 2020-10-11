package geecache

import (
	"fmt"
	"github.com/my/repo/geecache/singlefight"
	"log"
	"sync"
)

type Group struct {
	name string
	getter Getter
	mainCache cache
	peers PeerPicker
	loader *singlefight.Group
}

type Getter interface {
	Get(key string) ([]byte,error)
}

type GetterFunc func(key string) ([]byte,error)

func (f GetterFunc) Get(key string) ([]byte,error){
	return f(key)
}

var (
	mu sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string,cacheBytes int,getter Getter) *Group{
	if getter == nil{
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name: name,
		getter: getter,
		mainCache: cache{cacheBytes: cacheBytes},
		loader: &singlefight.Group{},
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group{
	mu.RLock() //只读锁
	g := groups[name]
	mu.RUnlock()
	return g
}

func (g *Group) Get(key string) (ByteView,error){
	if key == ""{
		return ByteView{},fmt.Errorf("key is required")
	}

	if v,ok := g.mainCache.get(key);ok{
		log.Fatalln("[GeeCache] hit")
		return v,nil
	}
	return g.load(key)
}

func (g *Group) getLocally(key string) (ByteView,error){
	bytes,err := g.getter.Get(key)
	if err != nil{
		return ByteView{},err
	}
	value := ByteView{b:cloneBytes(bytes)}
	g.populateCache(key,value)
	return value,nil
}

func (g *Group) populateCache(key string,value ByteView){
	g.mainCache.add(key,value)
}

func (g *Group) RegisterPeers(peers PeerPicker){
	if g.peers != nil{
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

func (g *Group) load(key string) (value ByteView,err error){
	//if g.peers != nil{
	//	if peer,ok := g.peers.PickPeer(key);ok{
	//
	//	}
	//}
	return g.getLocally(key)
}

func (g *Group) getFromPeer(peer PeerGetter,key string) (ByteView,err error){
	return nil,nil

}