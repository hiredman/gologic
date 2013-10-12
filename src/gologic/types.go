package gologic
import "container/list"

type LVarT struct {
        name string
}

type V *LVarT

type ConsT struct {
	name V
	thing interface {}
	more *ConsT
}

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

type Package struct {
	s *SubsT
	c *ConsT
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
