package gologic

type node struct {
	r *Rbnode
	size int
}

func (p subs_pair) Key() int {
	return p.v.id
}

func (p subs_pair) Merge(e Element) Element {
	return p
}

func (n node) val_at (v V) (interface {}, bool) {
	if n.r != nil {
		x, found := Locate(n.r,v.id)
		if found {
			a,ok := x.(subs_pair)
			if ok {
				return a.t, true
			} else {
				panic("oh no")
			}
		} else {
			return nil, false
		}
	} else {
		return nil, false
	}
}

func (n node) with (v V, t interface{}) substitution_map {
	return node{Insert(n.r,subs_pair{v,t}),n.size+1}
}

func (n node) count () int {
	return n.size
}

func (n node) fold (f func(interface{},V,interface{}) (interface{},bool), init interface{}) (interface{},bool) {
	return Fold(init,func(x interface{}, e Element) (interface{}, bool){
		p, ok := e.(subs_pair)
		if !ok {panic("uh oh")}
		return f(x,p.v,p.t)
	},n.r)
}

func (s empty_subst_value) with (v V, t interface{}) substitution_map {
	//return &SubsT{v,t,nil}
	return node{Node(subs_pair{v,t}),1}
}

