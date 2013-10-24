package gologic

type LVarT struct {
        id int
}

type V *LVarT

type SubsT struct {
        name V
        thing interface {}
        more *SubsT
}

type subs_pair struct {
	v V
	t interface {}
}

type substitution_map interface {
	val_at (V) (interface{}, bool)
	with (V, interface{}) substitution_map
	count () int
	fold(func (interface{},V,interface{}) (interface{},bool), interface{}) (interface{}, bool)
}

type LookupResult struct {
        Var bool
        v V
        Term bool
        t interface{}
}

type SubsTNode struct {
	e substitution_map
	r *SubsTNode
}

type ConstraintResult int

const (
	No ConstraintResult = iota
	Yes
	Maybe
)

type Vars struct {
	first V
	rest *Vars
}

type Constraint struct {
//	vars *Vars
	F func(S) (S, ConstraintResult)
	
}

type Constraints struct {
	first Constraint
	rest *Constraints
}

type Package struct {
	s substitution_map
	c *SubsTNode
	constraint_store *Constraints
}


type S *Package

type Stream struct {
        first S
        rest func () *Stream
}

type R *Stream

type Goal func(S) R

type db_record struct {
        Entity interface{}
        Attribute interface{}
        Value interface{}
}

type DBCons struct {
        car db_record
        cdr *DBCons
}

type DB struct {
        c chan func (*DBCons) *DBCons
}

type DBValue struct {
        d *DBCons
}

type Symbol struct {
	name string
}

type Unique struct {}
