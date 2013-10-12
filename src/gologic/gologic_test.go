package gologic
import "testing"

type rel int

// keywords if you squit hard
const (
	Likes rel =  iota
	Is
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

	who, something := Fresh2()

        c3 := Run(who, And(db.Find("Bob", Likes, something),
                           db.Find(who, Likes, something),
                           Or(db.Find(something, Is, "cheap"), 
                              db.Find(something, Is, "delicious"))))

	var foo [6]interface{}

        for ii := 0 ; true; ii++ {
                i := <- c3
                if i != nil {
			foo[ii] = i
                } else {
                        break
                }
        }

	if foo[0] != "Bob" { t.Fatal(foo[0])}
	if foo[1] != "Tom" { t.Fatal(foo[1])}
	if foo[2] != "Jill" { t.Fatal(foo[2])}
	if foo[3] != "Bob" { t.Fatal(foo[3])}
	if foo[4] != "Alice" { t.Fatal(foo[4])}
	if foo[5] != "Tom" { t.Fatal(foo[5])}
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
