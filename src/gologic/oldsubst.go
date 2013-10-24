// http://gradworks.umi.com/3380156.pdf
package gologic

func osubst_name(s *SubsT) V {
        return s.name
}

func osubst_thing(s *SubsT) interface {} {
        return s.thing
}

func osubst_more(s *SubsT) *SubsT {
        if s != nil {
		return s.more
        } else {
                return s
        }
}

func osubst_find (v V, s *SubsT) (*SubsT, bool) {
        if s == nil {
                return nil, false
        } else {
                if v == osubst_name(s) {
                        return s, true
                } else {
                        return osubst_find(v, osubst_more(s))
                }
        }
}


func (s *SubsT) val_at (v V) (interface {}, bool) {
	s, ok := osubst_find(v,s)
	if ok {
		return osubst_thing(s), true
	} else {
		return nil, false
	}
}

func (s *SubsT) with (v V, t interface{}) substitution_map {
	x := &SubsT{v,t,s}
	return x
}

func (s *SubsT) count () int {
        if s == nil {
                return 0
        } else {
                return 1+osubst_more(s).count()
        }
}

func (s *SubsT) fold (f func(interface{},V,interface{}) (interface{},bool), init interface{}) (interface{},bool) {
	r := init
	for x := s; x != nil;x = s.more {
		r1, ok := f(r,x.name,x.thing)
		if !ok {
			return r1,ok
		} 
	}
	return r, true
}


type empty_subst_value struct {}

func (s empty_subst_value) val_at (v V) (interface {}, bool) {
	return nil,false
}

func (s empty_subst_value) with (v V, t interface{}) substitution_map {
	return &SubsT{v,t,nil}
}

func (s empty_subst_value) count () int {
	return 0
}

func (s empty_subst_value) fold (f func(interface{},V,interface{}) (interface{},bool), init interface{}) (interface{},bool) {
	return init,true
}


func new_subst () substitution_map {
	return empty_subst_value{}
}

