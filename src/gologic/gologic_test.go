package gologic
import "testing"

type rel int

// keywords if you squint hard
const (
        Likes rel =  iota
        Is
        Sex
)

func Assert (t *testing.T, exp bool, msg interface{}) {
        if !exp {t.Fatal(msg)}
}

func TestReasoningOverDB (t *testing.T) {
        db := Db()
        db.Assert("Bob",   Likes, "Pizza")
        db.Assert("Bob",   Likes, "Sushi")
        db.Assert("Jill",  Likes, "Pizza")
        db.Assert("Tom",   Likes, "Pizza")
        db.Assert("Tom",   Likes, "Sushi")
        db.Assert("Alice", Likes, "Sushi")
        db.Assert("Pizza", Is,    "cheap")
        db.Assert("Sushi", Is,    "delicious")
        db.Assert("Bob",   Sex,    "male")
        db.Assert("Tom",   Sex,    "male")
        db.Assert("Jill",  Sex,    "female")
        db.Assert("Alice", Sex,    "female")

        who, something := Fresh2()

        d := db.Deref()

        // match.com, for triples!
        c3 := Run(who, And(d.Find("Bob", Likes, something),
                d.Find(who, Likes, something),
                Or(d.Find(something, Is, "cheap"),
                d.Find(something, Is, "delicious")),
                d.Find(who, Sex, "female")))

        var foo [2]interface{}

        for ii := 0 ; true; ii++ {
                i := <- c3
                if i != nil {
                        foo[ii] = i
                } else {
                        break
                }
        }

        Assert(t,foo[0] == "Alice",foo[0])
        Assert(t,foo[1] == "Jill",foo[1])

}

func TestReasoningOverStructs (t *testing.T) {
        type person struct {
                Name interface{}
                Last string
        }
        p := person{Name:"Bob",Last:"Villa"}
        vp := Fresh()
        c2 := Run(vp, Unify(p,person{Name:vp,Last:"Villa"}))
        Assert(t, "Bob" == <- c2, "not bob")
}

func TestGoals (t *testing.T) {
        a,b,c:=Fresh3()
        ch := Run(a,Or(And(Unify(a,b),Unify(b,c),Unify(c,1)),
                And(Unify(a,b),Unify(b,c),Unify(c,3))))
        Assert(t,1 == <- ch,"not a 1")
        Assert(t,3 == <- ch,"not a 3")
}

func TestInequality (t *testing.T) {
        a,b,c:=Fresh3()
        ch := Run(a,Or(And(Unify(a,b),Unify(b,c),Unify(c,1)),
                And(Unify(a,b),Unify(b,c),Unify(c,3),Neq(c,3))))
        Assert(t, 1 == <- ch, "not an 1")
        Assert(t, nil == <- ch, "not closed")
}

func TestReifyStructs (t *testing.T) {
        type X struct {
                A interface {}
        }
        v := Fresh()
        c := Run(v,Unify(v,X{"Hello"}))
        z := <- c
        d, ok := z.(X)
        Assert(t,ok,"not an X")
        Assert(t,d.A == "Hello","not Hello")
}

func TestReifyNestedStructs (t *testing.T) {
        type X struct {
                A interface {}
        }
        v,l := Fresh2()
        c := Run(l,And(Unify(v,X{l}), Unify(l,X{"Hello"})))
        z := <- c
        d, ok := z.(X)
        Assert(t,ok,"not an X")
        Assert(t,d.A == "Hello","not Hello")
}

func ancestoro (db DBValue, a,b interface{}) Goal {
        c := Fresh()
        return Or(db.Find(a,"parent",b),
                And(db.Find(a,"parent",c),
                // Call delays calling the goal constructor till goal execution time
                Call(ancestoro,db,c,b)))
}

func TestGenealogy (t *testing.T) {
        db := Db()
        db.Assert("bill","parent","mary")
        db.Assert("mary","parent","john")
        q := Fresh()
        v := db.Deref()
        c := Run(q,ancestoro(v,"bill",q))
        Assert(t,<- c == "mary","not mary")
        Assert(t,<- c == "john","not mary")
}

func TestP4AI (t *testing.T) {
        const (
                parent rel =  iota
		is
        )

        db := Db()
        db.Assert("bob",parent,"pam")
	db.Assert("bob",parent,"tom")
	db.Assert("liz",parent,"tom")
	db.Assert("ann",parent,"bob")
	db.Assert("pat",parent,"bob")
	db.Assert("jim",parent,"pat")
	
	d := db.Deref()
	q := Fresh()
	c := Run(q,d.Find("pat",parent,"bob"))
	if nil == <-c {t.Fail()}

	c = Run(q,d.Find("pat",parent,"liz"))
	if nil != <-c {t.Fail()}

	c = Run(q,d.Find("liz",parent,q))
	if "tom" != <-c {t.Fail()}

	c = Run(q,d.Find(q,parent,"bob"))
	if "pat" != <-c {t.Fail()}
	if "ann" != <-c {t.Fail()}

	x := Fresh()
	db.Assert(x,is,"man")
	db.Assert(x,is,"fallible")

	d = db.Deref()
	
	c = Run(q,And(d.Find("socrates",is,"man"),
		      d.Find("socrates",is,"fallible")))
	if nil == <-c {t.Fail()}
}

func increasing_or_equal(a,b interface{}) Goal {
	return AddC(Constraint{
		func (s S) (S, ConstraintResult) {
			o := Project(a,s)
			x, xok := o.(int)
			if xok {
				o = Project(b,s)
				y, yok := o.(int)
				if yok {
					if y >= x {
						return s,Yes
					} else {
						return s,No
					}
				} else {
					return s,Maybe
				}
			} else {
				return s,Maybe
			}
	}})
}

func difference(a,b,c interface{}) Goal {
        return AddC(Constraint{
                func (s S) (S, ConstraintResult) {
                        xo := Project(a,s)
                        x, xok := xo.(int)
                        yo := Project(b,s)
                        y, yok := yo.(int)
                        zo := Project(c,s)
                        z, zok := zo.(int)

                        if xok && yok && zok && x - y == z {
                                return s,Yes
                        } else if xok && yok && zok && x - y != z {
                                return s,No
                        } else if xok && yok  {
                                ns,success := Unifi(c,x-y,s)
                                if success {
                                        return ns, Yes
                                } else {
                                        return ns, No
                                }
                        } else if xok && zok  {
                                ns,success := Unifi(b,x-z,s)
                                if success {
                                        return ns, Yes
                                } else {
                                        return ns, No
                                }
                        } else if yok && zok {
                                ns,success := Unifi(a,y+z,s)
                                if success {
                                        return ns, Yes
                                } else {
                                        return ns, No
                                }
                        } else {
                                return s, Maybe
                        }


        }})
}


func TestConstraints(t *testing.T) {
	
	q := Fresh()
	c := Run(q,And(
		increasing_or_equal(q,5),
		Unify(q,6)))

	if <- c != nil { t.Fatal("foo")}

	
	q = Fresh()
	c = Run(q,Or(difference(1,1,q),difference(2,1,q)))
	
	n := <- c
	if n != 0 { t.Fatal(n)}

}
