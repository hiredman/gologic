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
        name_ V
        thing_ interface {}
        more_ *SubsT
	c *ConsT
}

type LookupResult struct {
        Var bool
        v V
        Term bool
        t interface{}
}

type S *SubsT

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
