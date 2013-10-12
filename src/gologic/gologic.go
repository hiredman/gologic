package gologic
//import "fmt"
import "strconv"
import "reflect"
import "container/list"

func is_struct (x interface{}) bool {
  //      fmt.Println("is_struct")
        v := reflect.ValueOf(x)
        k := v.Kind()
    //    fmt.Println(k == reflect.Struct)
        return k == reflect.Struct

        }

func type_name (x interface {}) string {
        v := reflect.ValueOf(x)
        t := v.Type()
        return t.PkgPath()+"/"+t.Name()
}

func field_count (x interface{}) int {
        v := reflect.ValueOf(x)
        return v.NumField()
}

func field_by_index (x interface{}, i int) interface {} {
        v := reflect.ValueOf(x)
        return v.Field(i).Interface()
}

func lvar(n string) V {
        var foo = new(LVarT)
        foo.name=n
        return foo
}

func exts_no_check (n V, v interface {}, s *SubsT) *SubsT {
        if n == nil {
                panic("foo")
        }
        var news = new(SubsT)
        news.name=n
        news.thing=v
        news.more=s
        return news
}

func lookup (thing interface{}, s *SubsT) LookupResult {
        var lr LookupResult

        v, isvar := thing.(V)

        if !isvar {
                lr.Var = false
                lr.Term = true
                lr.t = thing
                return lr
        } else {
                if s == nil {
                        lr.Var = true
                        lr.Term = false
                        lr.v = v
                        return lr
                } else if s.name.name == v.name {
                        lr.Var = false
                        lr.Term = true
                        lr.t = s.thing
                        return lr
                } else {
                        return lookup(thing,s.more)
                }

        }

}

func subst_find (v V, s *SubsT) (*SubsT, bool) {
        // fmt.Println("==subst_find==")
        // fmt.Println(v)
        // fmt.Println(s)
        if s == nil {
                return nil, false
        } else {
                // fmt.Println("A")
                // fmt.Println(v.name)
                // fmt.Println("B")
                // fmt.Println(s)
                // fmt.Println("C")
                // fmt.Println(s.name)
                if v.name == s.name.name {
                        return s, true
                } else {
                        return subst_find(v, s.more)
                }
        }
}

func walk (n interface {}, s *SubsT) LookupResult {
        // fmt.Println("==walk==")
	// fmt.Println(n)
        var lr LookupResult
        v, visvar := n.(V)
        // fmt.Println("visvar")
        // fmt.Println(visvar)
        // fmt.Println(v)
        if !visvar  || v == nil {
                lr.Term = true
                lr.Var = false
                lr.t = n
                return lr
        } else {
                // fmt.Println("yoyo")
                // fmt.Println(v)
                // fmt.Println(s)
                subs, subsfound := subst_find(v, s)
                if subsfound {
                        return walk(subs.thing, s)
                } else {
                        lr.Var = true
                        lr.Term = false
                        lr.v = v
                        return lr
                }
        }
}

func occurs_check (x V, v interface{}, s *SubsT) bool {
        thing := walk(v, s)
        if (thing.Var) {
                return thing.v.name == x.name
        } else {
                if is_struct(x) {
                        for i := 0; i < field_count(x); i++ {
                                nv, nvisvar := field_by_index(x,i).(V)
                                if nvisvar {
                                        if occurs_check(nv, v, s) {
                                                return true
                                        }
                                }
                        }
                        return false
                } else {
                        return false
                }
        }

}

func ext_s (x V, v interface{}, s *SubsT) (*SubsT, bool) {
        if x == nil {
                panic("foo")
        }
        if occurs_check(x,v,s) {
                return nil,false
        } else {
                return exts_no_check(x,v,s), true
        }
}



func unify (u interface{}, v interface{}, s *SubsT) (*SubsT, bool) {
        // fmt.Println("==unify==")
        u1 := walk(u,s)
        v1 := walk(v,s)

        // fmt.Println("u")
        // fmt.Println(u1.t)
        // fmt.Println("v")
        // fmt.Println(v1.v)

        if u1.Term && v1.Term && !is_struct(u1.t) && !is_struct(v1.t) {
                //fmt.Println("A")
                return s, u1.t == v1.t
        } else if u1.Var {
                // fmt.Println("B")
                if v1.Var {
                        // fmt.Println("B.1")
                        return exts_no_check(u1.v, v1.v, s), true
                } else {
                        // fmt.Println("B.2")
                        return ext_s(u1.v, v1.t, s)
                }
        } else if v1.Var {
                // fmt.Println("C")
                return ext_s(v1.v,u1.t,s)
        } else {
                if is_struct(u1.t) &&
                        is_struct(v1.t) &&
                        (type_name(u1.t) == type_name(v1.t)) &&
                        (field_count(v1.t) == field_count(u1.t)) {
//			fmt.Println("Here")
                        ns := s
                        for i := 0 ; i < field_count(v1.t); i++  {
                                n, ok := unify(field_by_index(u1.t,i),field_by_index(v1.t,i),ns)
                                if !ok {
                                        return ns, false
                                }
				ns = n
                        }
                        return ns, true
                } else {
                        return s, false
                }
        }
}

func unify_no_check (u, v, s *SubsT) (*SubsT, bool) {
        u1 := walk(u,s)
        v1 := walk(v,s)
        if u1 == v1 {
                return s,true
        } else if u1.Var {
                return exts_no_check(u1.v, v1.v, s), true
        } else if v1.Var {
                return ext_s(v1.v,u1.t,s)
        } else {
                return s, false
        }
}

func walk_star (v LookupResult, s *SubsT) LookupResult {
        // fmt.Println("==walk_star==")
	// fmt.Println(v)
	// fmt.Println(s)
        if v.Var {
                return walk(v.v,s)
        } else {
		if is_struct(v.t) {
//			fmt.Println("found struct")
			var lr LookupResult
			lr.Var = false
			lr.Term = true
			lr.t = 5
			return lr
		} else {
			return walk(v.t,s)
		}
        }
}

func length (s *SubsT) int32 {
        return 5
}

func reify_name (x int32) string {
        return "_."
}

func reify_s (v_ LookupResult, s *SubsT) *SubsT {
        // fmt.Println("==reify_s==")
        var v LookupResult
        if v_.Var {
                if v_.v == nil {
                        panic("foo")
                }
                v=walk(v.v,s)
        } else {
                v=walk(v.t,s)
        }
        if v.Var {
                if v.v == nil {
                        panic("foo")
                }
                // fmt.Println("reify here")
                // fmt.Println(v.v)
                s1, ok := ext_s(v.v, reify_name(length(s)), s)
                if ok {
                        return s1
                } else {
                        panic("whoops")
                }
        } else {
                return s
        }
}

func reify (v_ interface{}, s *SubsT) interface{} {
        // fmt.Println("==reify==")
        // fmt.Println(v_)
        var lr LookupResult
        va, vaisvar := v_.(V)
        if vaisvar {
                lr.Var = true
                lr.Term = false
                lr.v = va
        } else {
                lr.Var = false
                lr.Term = true
                lr.t = v_
        }
	
        // fmt.Println(lr)
	// fmt.Println("before first ws")
        v := walk_star(lr,s)
	// fmt.Println("after first ws")
        // fmt.Println("v")
        // fmt.Println(v.Var)
	// fmt.Println(v.v)
        lr2 := walk_star(v, reify_s(v,nil))
        return lr2.t
}

func mzero () *Stream {
        return nil
}

func unit (a *SubsT) *Stream {
        var x = new(Stream)
        x.first = a
        x.rest = func () *Stream {
                return mzero()
        }
        return x
}

func choice (a *SubsT, s func () *Stream) *Stream {
        var x = new(Stream)
        x.first = a
        x.rest = s
        return x
}

func Unify (u interface{}, v interface{}) Goal {
        return func (s S) R {
                s1, unify_success := unify(u,v,s)
		// fmt.Println("unify_success")
		// fmt.Println(unify_success)
		// fmt.Println(s1)
                if unify_success {
                        return unit(s1)
                } else {
                        return mzero()
                }
        }
}

func stream_concat(s1 *Stream, s2 func () *Stream) *Stream {
        if s1 == mzero() {
                return s2()
        } else {
                return choice(s1.first, func () *Stream {
                        return stream_concat(s1.rest(), s2)
                })
        }
}

func and_composer (g1s *Stream, g2 Goal) *Stream {
        if g1s == mzero() {
                return mzero()
        } else {
                return stream_concat(g2(g1s.first), func () *Stream {
                        a := g1s.rest()
                        if a == mzero() {
                                return mzero()
                        } else {
                                return and_composer(a, g2)
                        }
                })
        }
}

func and_base (g1, g2 Goal) Goal {
        return func (s S) R {
                g1s := g1(s)
                return and_composer(g1s, g2)

        }
}

func And (gs ...Goal) Goal {
        var g Goal = gs[0]
        for _,e := range gs[1:] {
                g = and_base(g,e)
        }
        return g
}

func stream_interleave (s1 *Stream, s2 *Stream) *Stream {
        if s1 == mzero() {
                return s2
        } else {
                return choice(s1.first, func () *Stream {
                        return stream_interleave(s2,s1.rest())
                })
        }

}

func or_base (g1, g2 Goal) Goal {
        return func (s S) R {
                g1s := g1(s)
                g2s := g2(s)
                return stream_interleave(g1s,g2s)
        }
}

func Fail () Goal {
        return func (s S) R {
                return mzero()
        }
}

func Or (gs ...Goal) Goal {
        var g Goal = gs[0]
        for _,e := range gs[1:] {
                g = or_base(g,e)
        }
        return g
}

func reify_as_list (v V, s *Stream, c chan interface{}) {
        for {
                if s == mzero() {
                        break
                } else {
                        c <- reify(v, s.first)
                        s=s.rest()
                }
        }
}

func Run (v V, g Goal) chan interface{} {
        c := make(chan interface{})
        go func () {
                reify_as_list(v, g(nil), c)
                close(c)
        }()
        return c
}

var c chan int

func init () {
        c = make (chan int)
        go func () {
                for i := 0; true; i++ {
                        c <- i
                }
        }()
}

func Fresh() V {
        var i int
        i = <- c
        return lvar("var"+strconv.Itoa(i))
}

func Fresh2() (V,V) {
        return Fresh(), Fresh()
}

func Fresh3() (V,V,V) {
        return Fresh(), Fresh(), Fresh()
}

func Fresh4() (V,V,V,V) {
        return Fresh(), Fresh(), Fresh(), Fresh()
}

func Fresh5() (V,V,V,V,V) {
        return Fresh(), Fresh(), Fresh(), Fresh(), Fresh()
}

func Fresh6() (V,V,V,V,V,V) {
        return Fresh(), Fresh(), Fresh(), Fresh(), Fresh(), Fresh()
}

func Db () DB {
        var x DB
        x.l = new(list.List)
        return x
}

func (d DB) Assert (entity interface{}, attribute interface{}, value interface{}) {
        d.l.PushBack(db_record{Entity:entity,Attribute:attribute,Value:value})
}

func (d DB) Find (entity interface{}, attribute interface{}, value interface{}) Goal {
        r := db_record{Entity:entity,Attribute:attribute,Value:value}
        g := Fail()
        for e := d.l.Front(); e != nil; e = e.Next() {
                g = Or(g,Unify(r,e.Value))
        }
        return g
}
