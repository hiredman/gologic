// http://brownbuffalo.sourceforge.net/ArchFriendsClues.html
package main
import "fmt"
import l "gologic"

type Purchase struct {
        Shop interface {}
	Shoe interface {}
}

type Trip struct {
        Stop1 interface {}
        Stop2 interface {}
        Stop3 interface {}
        Stop4 interface {}
}

func v () l.V {
	return l.Fresh()
}

// func membero(p, t interface{}) l.Goal {
// 	return l.Or(
// 		l.Unify(Trip{  p,v(),v(),v()}, t),
// 		l.Unify(Trip{v(),  p,v(),v()}, t),
// 		l.Unify(Trip{v(),v(),  p,v()}, t),
// 		l.Unify(Trip{v(),v(),v(),  p}, t))
		
// }

func aftero(p1, p2, t interface{}) l.Goal {
	return l.Or(
		l.Unify(Trip{ p1, p2,v(),v()}, t),
		l.Unify(Trip{v(), p1, p2,v()}, t),
		l.Unify(Trip{v(),v(), p1, p2}, t))		
}

func t (a,b,c,d interface{}) interface{} {
	return Trip{a,b,c,d}
}

func archo (q l.V) l.Goal {
	a,b := l.Fresh2()
	membero := l.StructMemberoConstructor4(t)
        return l.And(
		membero(Purchase{"The Foot Farm", v()}, q),
		membero(Purchase{"Heels in a Handcart", v()}, q),
		membero(Purchase{"The Shoe Palace", v()}, q),
		membero(Purchase{"Tootsies", v()}, q),
		membero(Purchase{v(),"ecru espadrilles"}, q),
		membero(Purchase{v(),"fuchsia flats"}, q),
		membero(Purchase{v(),"purple pumps"}, q),
		membero(Purchase{v(),"suede sandals"}, q),
		// 1
		membero(Purchase{"Heels in a Handcart", "fuchsia flats"}, q),
		// 3
		l.Unify(Trip{v(),Purchase{"The Foot Farm", v()},v(),v()}, q),
		// 4
		aftero(Purchase{"The Shoe Palace",v()}, b, q),
		aftero(b, Purchase{v(),"suede sandals"}, q),
		// 2
		aftero(Purchase{v(),"purple pumps"},Purchase{a,v()},q),
		l.Neq(a,"Tootsies"),
		)

}

func main() {
        q := l.Fresh()

        c := l.Run(q,archo(q))
	
        for {
                i := <- c
                if i != nil {
                        fmt.Println(i)
                } else {
                        break
                }
        }

}
