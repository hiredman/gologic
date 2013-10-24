package redblack
import "testing"
import "unicode/utf8"
import "container/list"

type foo struct {
	i int
}

func (f foo) Key() int {
	return f.i
}

func (f foo) Merge(e Element) Element {
	return f
}

func TestSet(t *testing.T) {
	var x *Rbnode
	for i := 0; i < 1000; i++ {
		x = Insert(x,foo{i})
	}
	for i := 0; i < 1000; i++ {
		foo, ok := Locate(x,i)
		if ok {
			if foo.Key() != i { t.Fatal("var") }
		} else {
			t.Fatal("foo")
		}
	}
}

func hash_string (str string) int {
	n := utf8.RuneCountInString(str)
	v := 0
	for i := 0; i < n; i++ {
		b := int(byte(str[i]))
		v += b * 31^((n-1)-i)
	}
	return v
}

type bar struct {
	hash int
	l *list.List
}

func (b bar) Key() int {
	return b.hash
}

func (b bar) Merge(e Element) Element {
	x, ok := e.(bar)
	if !ok {
		panic("oh no")
	}
	l := list.New()
	for node := x.l.Front(); node != nil; node = node.Next() {
		l.PushFront(node.Value)
	}
	for node := b.l.Front(); node != nil; node = node.Next() {
		l.PushFront(node.Value)
	}
	return bar{b.hash,l}

}

type pair struct {
	k string
	v interface{}
}

func kv_pair(k string, v interface{}) Element {
	l := list.New()
	l.PushFront(pair{k,v})
	return bar{hash_string(k),l}
	
}

func TestMap(t *testing.T) {
	var x *Rbnode
	x = Insert(x,kv_pair("foo","bar"))
	x = Insert(x,kv_pair("hello","world"))
	e, _ := Locate(x,hash_string("foo"))
	b := e.(bar)
	if b.l.Front().Value.(pair).v != "bar" {t.Fail()}
	if b.l.Len() != 1 {t.Fail()}
}
