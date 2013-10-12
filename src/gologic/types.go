package gologic
import "container/list"

type LVarT struct {
        name string
}

type V *LVarT

type SubsT struct {
        name V
        thing interface {}
        more *SubsT
}

type LookupResult struct {
        Var bool
        v V
        Term bool
        t interface{}
}

type SubsTNode struct {
	e *SubsT
	r *SubsTNode
}

type Package struct {
	s *SubsT
	c *SubsTNode
}

type S *Package

type Stream struct {
        first S
        rest func () *Stream
}

type R *Stream

type Goal func(S) R

type DB struct {
        l *list.List
}

type db_record struct {
        Entity interface{}
        Attribute interface{}
        Value interface{}
}
