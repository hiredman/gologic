// http://gradworks.umi.com/3380156.pdf
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

func zero_of (x interface {}) reflect.Value {
        v := reflect.ValueOf(x)
        t := v.Type()
	c := reflect.New(t)
        return c
}

func field_count (x interface{}) int {
        v := reflect.ValueOf(x)
        return v.NumField()
}

func field_by_index (x interface{}, i int) interface {} {
        v := reflect.ValueOf(x)
        return v.Field(i).Interface()
}

func set_field (x reflect.Value, i int, y interface {}) {
        x.Elem().Field(i).Set(reflect.ValueOf(y))
}

func lvar(n string) V {
        var foo = new(LVarT)
        foo.name=n
        return foo
}

func s_of(p S) *SubsT {
        if p != nil {
                return p.s
        } else {
                return nil
        }
}

func c_of(p S) *SubsTNode {
        if p != nil {
                return p.c
        } else {
                return nil
        }
}

func exts_no_check (n V, v interface {}, s S) S {
        if n == nil {
                panic("foo")
        }

        a := s_of(s)
        b := c_of(s)

        if a == nil {
                return &Package{s:&SubsT{name:n,thing:v,more:nil},c:b}
        } else {
                news := &SubsT{name:n,thing:v,more:a}
                return &Package{s:news,c:s.c}
        }
}

func subst_name(s S) V {
        return s.s.name
}

func subst_thing(s S) interface {} {
        return s.s.thing
}

func subst_more(s S) S {
        if s_of(s) != nil {
                a := s_of(s)
                b := c_of(s)
                if a != nil {
                        return &Package{s:a.more,c:b}
                } else {
                        return &Package{s:nil,c:b}
                }
        } else {
                return s
        }
}

func empty_subst(s S) bool {
        return s_of(s) == nil
}

func lookup (thing interface{}, s S) LookupResult {
        var lr LookupResult

        v, isvar := thing.(V)

        if !isvar {
                lr.Var = false
                lr.Term = true
                lr.t = thing
                return lr
        } else {
                if empty_subst(s) {
                        lr.Var = true
                        lr.Term = false
                        lr.v = v
                        return lr
                } else if subst_name(s).name == v.name {
                        lr.Var = false
                        lr.Term = true
                        lr.t = subst_thing(s)
                        return lr
                } else {
                        return lookup(thing,subst_more(s))
                }

        }

}

func subst_find (v V, s S) (S, bool) {
        // fmt.Println("==subst_find==")
        // fmt.Println(v)
        // fmt.Println(s)
        if empty_subst(s) {
                return nil, false
        } else {
                // fmt.Println("A")
                // fmt.Println(v.name)
                // fmt.Println("B")
                // fmt.Println(s)
                // fmt.Println("C")
                // fmt.Println(s.name)
                if v.name == subst_name(s).name {
                        return s, true
                } else {
                        return subst_find(v, subst_more(s))
                }
        }
}

func walk (n interface {}, s S) LookupResult {
        var lr LookupResult
        v, visvar := n.(V)
        if !visvar {
                lr.Term = true
                lr.Var = false
                lr.t = n
                return lr
        } else {
                subs, subsfound := subst_find(v, s)
                if subsfound {
                        return walk(subst_thing(subs), s)
                } else {
                        lr.Var = true
                        lr.Term = false
                        lr.v = v
                        return lr
                }
        }
}

func occurs_check (x V, v interface{}, s S) bool {
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

func ext_s (x V, v interface{}, s S) (S, bool) {
        if x == nil {
                panic("foo")
        }
        if occurs_check(x,v,s) {
                return nil,false
        } else {
                return exts_no_check(x,v,s), true
        }
}



func unify (u interface{}, v interface{}, s S) (S, bool) {
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
        // } else if u1.Term && u1.t == A {
        //         return s, true
        // } else if v1.Term && v1.t == A {
        //         return s, true
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
                        //                      fmt.Println("Here")
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

func unify_no_check (u, v, s S) (S, bool) {
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

func walk_star (v LookupResult, s S) LookupResult {
        // fmt.Println("==walk_star==")
        // fmt.Println(v)
        // fmt.Println(s)
        if v.Var {
                x := walk(v.v,s)
		if x.Var {
			return x
		} else {
			return walk_star(x, s)
		}
        } else {
                if is_struct(v.t) {
			//fmt.Println("walking struct")
			x := zero_of(v.t)
                        var lr LookupResult
                        lr.Var = false
                        lr.Term = true
                        for i := 0 ; i < field_count(v.t); i++  {
                                a := field_by_index(v.t,i)
                                var b LookupResult
                                vt, vtisvar := a.(V)
                                if vtisvar {
                                        b.Var = true
                                        b.Term = false
                                        b.v = vt
                                } else {
                                        b.Var = false
                                        b.Term = true
                                        b.t = a
                                }
                                c := walk_star(b,s)
                                if c.Var {
                                        set_field(x,i,c.v)
                                } else {
                                        set_field(x,i,c.t)
                                }
                        }
			lr.t = x.Elem().Interface()
                        return lr
                } else {
                        return walk(v.t,s)
                }
        }
}

func length (s S) int {
	if s_of(s) == nil {
		return 0
	} else {
		return 1+length(subst_more(s))
	}
}

func reify_name (x int) string {
        return "_."+strconv.Itoa(x)
}

func reify_s (v_ LookupResult, s S) S {
        //fmt.Println("==reify_s==")
	v := walk_star(v_,s)
        if v.Var {
                if v.v == nil {
                        panic("foo")
                }
                s1, ok := ext_s(v.v, reify_name(length(s)), s)
                if ok {
                        return s1
                } else {
                        panic("whoops")
                }
        } else {
		if is_struct(v.t) {
			//fmt.Println("reify struct")
			ns := s
			for i := 0; i < field_count(v.t); i++ {
				x := field_by_index(v.t,i)
				var t LookupResult
				d, disvar := x.(V)
				if disvar {
					t.Var = true
					t.Term = false
					t.v = d
				} else {
					t.Var = false
					t.Term = true
					t.t = x
				}
				ns = reify_s(t,ns)
			}
			return ns
		} else {
			return s
		}
        }
}

func reify (v_ interface{}, s S) interface{} {
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
        // fmt.Println(v.t)
	// fmt.Println("before refiy_s")
        x := reify_s(v,nil)
	// fmt.Println("after refiy_s")
        lr2 := walk_star(v, x)
        return lr2.t
}

func mzero () *Stream {
        return nil
}

func unit (a S) *Stream {
        var x = new(Stream)
        x.first = a
        x.rest = func () *Stream {
                return mzero()
        }
        return x
}

func choice (a S, s func () *Stream) *Stream {
        var x = new(Stream)
        x.first = a
        x.rest = s
        return x
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

func cons_c (c *SubsT, cs *SubsTNode) *SubsTNode {
        return &SubsTNode{e:c,r:cs}
}

func make_a (s *SubsT, c *SubsTNode) S {
        return &Package{s:s,c:c}
}

func prefix_s(s *SubsT, ss *SubsT) *SubsT {
        if s == ss {
                return nil
        } else {
                return &SubsT{name:s.name,thing:s.thing,more:prefix_s(s.more, ss)}
        }
}

func neq_verify(s *SubsT, a S, unify_success bool) R {
        if !unify_success {
                return unit(a)
        } else if s_of(a) == s {
                return mzero()
        } else {
                c := prefix_s(s,s_of(a))
                b := make_a(s_of(a), cons_c(c, c_of(a)))
                return unit(b)
        }
}


func Neq (u interface{}, v interface{}) Goal {
        return func (s S) R {
                s1, unify_success := unify(u,v,s)
                return neq_verify(s_of(s1),s,unify_success)
        }
}

func unify_star(p *SubsT, s S) (S, bool){
        if nil == p {
                return s, true
        } else {
                s1, unify_success := unify(p.name,p.thing,s)
                if unify_success {
                        return unify_star(p.more,s1)
                } else {
                        return nil, false
                }
        }
}

func verify_c(c *SubsTNode, cs *SubsTNode, s S) (*SubsTNode, bool) {
        if c == nil {
                return cs, true
        } else {
                s1, unify_success := unify_star(c.e,s)
                if unify_success {
                        if s == s1 {
                                return nil,false
                        } else {
                                cc := prefix_s(s_of(s1),s_of(s))
                                return verify_c(c.r, &SubsTNode{e:cc,r:cs}, s)
                        }
                } else {
                        return verify_c(c.r, cs, s)
                }
        }
}

func unify_verify(s S, a S, unify_success bool) R {
        if !unify_success {
                return mzero()
        } else if s_of(a) == s_of(s) {
                return unit(a)
        } else  {
                c, verified := verify_c(c_of(a), nil, s)
                if verified {
                        return unit(make_a(s_of(s), c))
                } else {
                        return mzero()
                }
        }
}

func Unify (u interface{}, v interface{}) Goal {
        return func (s S) R {
                s1, unify_success := unify(u,v,s)
		// fmt.Println("unify_success")
		// fmt.Println(unify_success)
                return unify_verify(s1,s,unify_success)
        }
}

func (v LVarT) String () string {
	return "<lvar "+v.name+">"
}
