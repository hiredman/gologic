package gologic
import "testing"

type rel int

// keywords if you squint hard
const (
	Likes rel =  iota
	Is
	Sex
	)

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

	// match.com, for tuples!
        c3 := Run(who, And(db.Find("Bob", Likes, something),
                           db.Find(who, Likes, something),
                           Or(db.Find(something, Is, "cheap"), 
                              db.Find(something, Is, "delicious")),
		           Neq("Bob", who),
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

	if foo[0] != "Jill" { t.Fatal(foo[0])}
	if foo[1] != "Alice" { t.Fatal(foo[1])}
}

type person struct {
        Name interface{}
        Last string
}

func TestReasoningOverStructs (t *testing.T) {

	p := person{Name:"Bob",Last:"Villa"}

        vp := Fresh()

        c2 := Run(vp, Unify(p,person{Name:vp,Last:"Villa"}))

        if "Bob" != <- c2 { t.Fail()}

}

func TestingGoals (t *testing.T) {
	a,b,c:=Fresh3()
	ch := Run(a,Or(And(Unify(a,b),Unify(b,c),Unify(c,1)),
                       And(Unify(a,b),Unify(b,c),Unify(c,3))))
	if 1 == <- ch {t.Fatal("not a 1")}
	if 3 == <- ch {t.Fatal("not a 3")}
}

func TestingInequality (t *testing.T) {
	a,b,c:=Fresh3()
	ch := Run(a,Or(And(Unify(a,b),Unify(b,c),Unify(c,1)),
                       And(Unify(a,b),Unify(b,c),Unify(c,3),Neq(c,3))))
	if 1 == <- ch {t.Fatal("not a 1")}
	if nil == <- ch {t.Fatal("not closed")}
}
