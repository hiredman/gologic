// http://gradworks.umi.com/3380156.pdf
package gologic
import "strconv"
import "reflect"
import "fmt"

var c chan int

var U Unique = Unique{}

func is_struct (x interface{}) bool {
        v := reflect.ValueOf(x)
        k := v.Kind()
        return k == reflect.Struct
        }

func type_of (x interface {}) reflect.Type {
        v := reflect.ValueOf(x)
        return v.Type()
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

func field_by_index (x interface{}, i int) (foo interface {}) {
	defer func () {
		r := recover()
		if r != nil {
			foo = U
		}
	}()
        v := reflect.ValueOf(x)
        foo = v.Field(i).Interface()
	return foo
}

func set_field (x reflect.Value, i int, y interface {}) (success bool) {
	defer func () {
		r := recover()
		if r != nil {
			success = false
		}
	}()
        x.Elem().Field(i).Set(reflect.ValueOf(y))
	success = true
	return success
}

func s_of(p S) substitution_map {
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

func constraints_of(p S) *Constraints {
        if p != nil {
                return p.constraint_store
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
	c := constraints_of(s)

        if a == nil {
                return &Package{s:new_subst().with(n,v),c:b,constraint_store:c}
        } else {
                news := a.with(n,v)
                return &Package{s:news,c:b,constraint_store:c}
        }
}

// func subst_name(s S) V {
//         return s.s.name
// }

// func subst_thing(s S) interface {} {
//         return s.s.thing
// }

// func subst_more(s S) S {
//         if s_of(s) != nil {
//                 a := s_of(s)
//                 b := c_of(s)
//                 if a != nil {
//                         return &Package{s:a.more,c:b}
//                 } else {
//                         return &Package{s:nil,c:b}
//                 }
//         } else {
//                 return s
//         }
// }

func empty_subst(s S) bool {
        return s_of(s) == nil
}

func AVar(v V) LookupResult {
	var lr LookupResult
        lr.Var = true
        lr.Term = false
        lr.v = v
        return lr
}

func ATerm(t interface{}) LookupResult {
	var lr LookupResult
        lr.Var = false
        lr.Term = true
        lr.t = t
        return lr
}

// func lookup (thing interface{}, s S) LookupResult {
//         v, isvar := thing.(V)
//         if !isvar {
// 		return ATerm(thing)
//         } else {
//                 if empty_subst(s) {
// 			return AVar(v)
// 		} else {
// 			thing, found := s_of(s).val_at(v)
// 			if found {
// 				return ATerm(thing)
// 			} else {
// 				return 
// 			}
//                 } else if subst_name(s) == v {
// 			return ATerm(subst_thing(s))
//                 } else {
//                         return lookup(thing,subst_more(s))
//                 }
//         }
// }

func subst_find(v V, s S) (interface{}, bool) {
	x := s_of(s)
	if x != nil {
		item, ok := x.val_at(v)
		if ok {
			return item, true
		} else {
			return nil, false
		}
	} else {
		return nil,false
	}
}

func walk (n interface {}, s S) LookupResult {
        v, visvar := n.(V)
        if !visvar {
		return ATerm(n)
        } else {
                thing, subsfound := subst_find(v, s)
                if subsfound {
                        return walk(thing, s)
                } else {
			return AVar(v)
                }
        }
}

func occurs_check (x V, v interface{}, s S) bool {
        thing := walk(v, s)
        if (thing.Var) {
                return thing.v == x
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

func unify_no_constraints (u interface{}, v interface{}, s S) (S, bool) {
        u1 := walk(u,s)
        v1 := walk(v,s)
        if u1.Term && v1.Term && !is_struct(u1.t) && !is_struct(v1.t) {
                return s, u1.t == v1.t
        } else if (u1.Term || v1.Term) && (u1.t == U || v1.t == U) {
                return s,false
        } else if u1.Var && v1.Var {
                return exts_no_check(u1.v, v1.v, s), true
        } else if u1.Var {
                return ext_s(u1.v, v1.t, s)
        } else if v1.Var {
                return ext_s(v1.v,u1.t,s)
        } else if is_struct(u1.t) &&
                is_struct(v1.t) &&
                (type_of(u1.t) == type_of(v1.t)) &&
                (field_count(v1.t) == field_count(u1.t)) {
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

func apply_constraints(s S) (S, bool) {
	var newc *Constraints
	ns := s
	var r ConstraintResult
	for c := constraints_of(s); c != nil; c = c.rest {
                if true {
                        ns, r = c.first.F(ns)
                        if r == No {
                                return ns, false
                        } else if r == Maybe {
                                newc = &Constraints{c.first,newc}
                        } else {
				
                        }
                } else {
                        newc = &Constraints{c.first,newc}
                }
        }
        return make_a(s_of(ns),c_of(ns),newc), true
}

func unify (u interface{}, v interface{}, s S) (S, bool) {
        ns, ok := unify_no_constraints(u,v,s)
        if !ok {
                return ns,ok
        } else {
		nns, success := apply_constraints(ns)
		if success {
			return nns, success
		} else {
			return ns, false
		}

        }
}

// func unify_no_check (u, v, s S) (S, bool) {
//         u1 := walk(u,s)
//         v1 := walk(v,s)
//         if u1 == v1 {
//                 return s,true
//         } else if u1.Var {
//                 return exts_no_check(u1.v, v1.v, s), true
//         } else if v1.Var {
//                 return ext_s(v1.v,u1.t,s)
//         } else {
//                 return s, false
//         }
// }

func walk_star (v LookupResult, s S) LookupResult {
        if v.Var {
                x := walk(v.v,s)
                if x.Var {
                        return x
                } else {
                        return walk_star(x, s)
                }
        } else {
                if is_struct(v.t) {
                        x := zero_of(v.t)
                        var lr LookupResult
                        lr.Var = false
                        lr.Term = true
                        for i := 0 ; i < field_count(v.t); i++  {
                                a := field_by_index(v.t,i)
                                var b LookupResult
                                vt, vtisvar := a.(V)
                                if vtisvar {
					b = AVar(vt)
                                } else {
					b = ATerm(a)
                                }
                                c := walk_star(b,s)
				good_set := false
                                if c.Var {
                                        good_set = set_field(x,i,c.v)
                                } else {
                                        good_set = set_field(x,i,c.t)
                                }
				if !good_set {
					return v
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
		return s_of(s).count()
	}
}

func reify_name (x int) Symbol {
        return Symbol{"_."+strconv.Itoa(x)}
}

func reify_s (v_ LookupResult, s S) S {
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
                        ns := s
                        for i := 0; i < field_count(v.t); i++ {
                                x := field_by_index(v.t,i)
                                var t LookupResult
                                d, disvar := x.(V)
                                if disvar {
					t = AVar(d)
                                } else {
					t = ATerm(x)
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
        var lr LookupResult
        va, vaisvar := v_.(V)
        if vaisvar {
		lr = AVar(va)
        } else {
		lr = ATerm(v_)
        }
        v := walk_star(lr,s)
        x := reify_s(v,nil)
        lr2 := walk_star(v, make_a(s_of(x),nil,nil))
        return lr2.t
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

// Run takes a logic variable to solve for, and a goal
func Run (v V, g Goal) chan interface{} {
        c := make(chan interface{})
        go func () {
                reify_as_list(v, g(nil), c)
                close(c)
        }()
        return c
}



func init () {
        c = make(chan int, 10)
        go func () {
                for i := 0; true; i++ {
                        c <- i
                }
        }()
}

// Fresh returns a new logic variable
func Fresh() V {
        foo := new(LVarT)
        foo.id = <- c
        return foo
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

func cons_c (c substitution_map, cs *SubsTNode) *SubsTNode {
        return &SubsTNode{e:c,r:cs}
}

func make_a (s substitution_map, c *SubsTNode, con *Constraints) S {
        return &Package{s:s,c:c,constraint_store:con}
}

// func prefix_s(s substitution_map, ss substitution_map) *SubsT {
//         if s == ss {
//                 return nil
//         } else {
//                 //return &SubsT{name:s.name,thing:s.thing,more:prefix_s(s.more, ss)}
// 		return 
// 		return &SubsT{name:s.name,thing:s.thing,more:prefix_s(s.more, ss)}
//         }
// }

// func neq_verify(s *SubsT, a S, unify_success bool) R {
//         if !unify_success {
//                 return unit(a)
//         } else if s_of(a) == s {
//                 return mzero()
//         } else {
//                 c := prefix_s(s,s_of(a))
//                 b := make_a(s_of(a), cons_c(c, c_of(a)), constraints_of(a))
//                 return unit(b)
//         }
// }


// // Neq returns a goal that suceeds when u and v do not unify
// func Neq (u interface{}, v interface{}) Goal {
//         return func (s S) R {
//                 s1, unify_success := unify(u,v,s)
//                 return neq_verify(s_of(s1),s,unify_success)
//         }
// }

func unify_star(p substitution_map, s S) (S, bool){
        if nil == p {
                return s, true
        } else {
		i, ok := p.fold(func (i interface{}, v V, t interface{}) (interface{}, bool) {
			s1, ok := i.(S)
			if !ok {panic("oh no")}
			s1, unify_success := unify(v,t,s1)
			if unify_success {
				return s1,true
			} else {
				return nil, false
			}
		}, s)
		if ok {
			s1, ok2 := i.(S)
			if !ok2 {panic("oh no")}
			return s1, true
		} else {
			return nil, false
		}
        }
}

// func verify_c(c *SubsTNode, cs *SubsTNode, s S) (*SubsTNode, bool) {
//         if c == nil {
//                 return cs, true
//         } else {
//                 s1, unify_success := unify_star(c.e,s)
//                 if unify_success {
//                         if s == s1 {
//                                 return nil,false
//                         } else {
//                                 cc := prefix_s(s_of(s1),s_of(s))
//                                 return verify_c(c.r, &SubsTNode{e:cc,r:cs}, s)
//                         }
//                 } else {
//                         return verify_c(c.r, cs, s)
//                 }
//         }
// }

func unify_verify(s S, a S, unify_success bool) R {
        if !unify_success {
                return mzero()
        } else if s_of(a) == s_of(s) {
                return unit(a)
        } else  {
                // c, verified := verify_c(c_of(a), nil, s)
                // if verified {
                //         return unit(make_a(s_of(s), c, constraints_of(s)))
                // } else {
                //         return mzero()
                // }
		return unit(s)
        }
}

// Unify returns a goal that succeeds when u and v unify
func Unify (u interface{}, v interface{}) Goal {
        return func (s S) R {
                s1, unify_success := unify(u,v,s)
                return unify_verify(s1,s,unify_success)
        }
}

func (v LVarT) String () string {
        return "<lvar "+string(v.id)+">"
}

func (v Symbol) String () string {
        return v.name
}

// helper for constructing recursive goals
func Call(constructor interface{}, args ...interface{}) Goal {
        foo := make([]reflect.Value, len(args))
        for i,e := range args {
                foo[i] = reflect.ValueOf(e)
        }
        fun := reflect.ValueOf(constructor)
        return func (s S) R {
                r := fun.Call(foo)
                x := r[0].Interface()
                g,ok := x.(Goal)
                if ok {
                        return g(s)
                } else {
                        panic("whoops")
                }

        }
}

func Project(a interface{}, s S) interface{} {
	v, vok := a.(V)
	if vok {
		lr := walk_star(AVar(v),s)
		if lr.Var {
			return lr.v
		} else {
			return lr.t
		}
	} else {
		return a
	}
}

func IsSymbol(s interface{}) bool {
        _, ok := s.(Symbol)
        return ok
}

func AddC (c Constraint) Goal {
        return func (s S) R {
		new_s := make_a(s_of(s),c_of(s),&Constraints{c,constraints_of(s)})
		ns, ok := apply_constraints(new_s)
		if ok {
			return unit(ns)
		} else {
			return mzero()
		}
        }
}

func Unifi(a,b interface{}, s S) (S, bool) {
	return unify_no_constraints(a,b,s)
}

func StructMemberoConstructor5 (f func (interface{},interface{},interface{},interface{},interface{}) interface{}) func (interface{}, interface{}) Goal {
	return func (p, t interface{}) Goal {
		return Or(
			Unify(f(      p,Fresh(),Fresh(),Fresh(),Fresh()), t),
			Unify(f(Fresh(),      p,Fresh(),Fresh(),Fresh()), t),
			Unify(f(Fresh(),Fresh(),      p,Fresh(),Fresh()), t),
			Unify(f(Fresh(),Fresh(),Fresh(),      p,Fresh()), t),
			Unify(f(Fresh(),Fresh(),Fresh(),Fresh(),      p), t))
	}
}

func StructMemberoConstructor4 (f func (interface{},interface{},interface{},interface{}) interface{}) func (interface{}, interface{}) Goal {
	return func (p, t interface{}) Goal {
		return Or(
			Unify(f(      p,Fresh(),Fresh(),Fresh()), t),
			Unify(f(Fresh(),      p,Fresh(),Fresh()), t),
			Unify(f(Fresh(),Fresh(),      p,Fresh()), t),
			Unify(f(Fresh(),Fresh(),Fresh(),      p), t))
	}
}

func StructMemberoConstructor3 (f func (interface{},interface{},interface{}) interface{}) func (interface{}, interface{}) Goal {
	return func (p, t interface{}) Goal {
		return Or(
			Unify(f(      p,Fresh(),Fresh()), t),
			Unify(f(Fresh(),      p,Fresh()), t),
			Unify(f(Fresh(),Fresh(),      p), t))
	}
}

func PrintChannel(n int, c chan interface{}) {
	for i := 0; i < n; i++ {
		e := <- c
		if e != nil {
			fmt.Println(e)
		}
	}
}
