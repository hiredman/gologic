package gologic

type lvart struct {
        id int
}

type V *lvart

type subs_pair struct {
	v V
	t interface {}
}

type substitution_map interface {
	val_at (V) (interface{}, bool)
	with (V, interface{}) substitution_map
	count () int
}

type empty_subst_value struct {}

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

type ConstraintResult uint

const (
	No ConstraintResult = iota
	Yes
	Maybe
)


type Constraint struct {
	F func(S) (S, ConstraintResult)
	
}

type Constraints struct {
	first Constraint
	rest *Constraints
}

type the_package struct {
	s substitution_map
	c *SubsTNode
	constraint_store *Constraints
}
type S *the_package

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
