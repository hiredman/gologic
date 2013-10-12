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

	if foo[0] != "Bob" { t.Fail() }
	if foo[1] != "Tom" { t.Fail() }
	if foo[2] != "Jill" { t.Fail() }
	if foo[3] != "Bob" { t.Fail() }
	if foo[4] != "Alice" { t.Fail() }
	if foo[5] != "Tom" { t.Fail() }
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
