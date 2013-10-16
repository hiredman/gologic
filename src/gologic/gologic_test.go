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

        // match.com, for triples!
        c3 := Run(who, And(db.Find("Bob", Likes, something),
                           db.Find(who, Likes, something),
                           Or(db.Find(something, Is, "cheap"),
                              db.Find(something, Is, "delicious")),
                           db.Find(who, Sex, "female")))

        var foo [2]interface{}

        for ii := 0 ; true; ii++ {
                i := <- c3
                if i != nil {
                        foo[ii] = i
                } else {
                        break
                }
        }

	Assert(t,foo[0] == "Jill",foo[0])
	Assert(t,foo[1] == "Alice",foo[1])

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

func ancestoro (db DB, a,b interface{}) Goal {
        c := Fresh()
        return Or(db.Find(a,"parent",b),
                  And(db.Find(a,"parent",c),
                      Call(ancestoro,db,c,b)))
}

func TestGenealogy (t *testing.T) {
        db := Db()
        db.Assert("bill","parent","mary")
        db.Assert("mary","parent","john")
	q := Fresh()
	c := Run(q,ancestoro(db,"bill",q))
	Assert(t,<- c == "mary","not mary")
	Assert(t,<- c == "john","not mary")
}
