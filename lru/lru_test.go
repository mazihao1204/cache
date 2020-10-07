package lru

import "testing"

type simpleStruct struct {
	int
	string
}

type complexStruct struct {
	int
	simpleStruct
}

type String string

func (d String) Len() int{
	return len(d)
}

var getTests = []struct{
	name string
	keyToAdd interface{}
	keyToGet interface{}
	expectedOk bool
}{
	{"string_hit","myKey","myKey",true},
	{"string_miss", "myKey", "nonsense", false},
	{"simple_struct_hit", simpleStruct{1, "two"}, simpleStruct{1, "two"}, true},
	{"simeple_struct_miss", simpleStruct{1, "two"}, simpleStruct{0, "noway"}, false},
	{"complex_struct_hit", complexStruct{1, simpleStruct{2, "three"}},
		complexStruct{1, simpleStruct{2, "three"}}, true},
}

func TestCache_Get(t *testing.T) {
	lru := New(int(0),nil)
	lru.Add("key1",String("1234"))
	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("cache hit key1=1234 failed")
	}
	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}