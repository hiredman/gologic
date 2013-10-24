//http://brownbuffalo.sourceforge.net/FourIslandsClues.html
package main
import "fmt"
import l "gologic"

type Island struct {
        Location,Name,Export,Attraction interface{}
}

type Nation struct {
        C1,C2,C3,C4 interface{}
}

func v () l.V {
        return l.Fresh()
}

func n (a,b,c,d interface{}) interface{} {
        return Nation{a,b,c,d}
}

const (
        North_of int = iota
        East_of
)

var db l.DB = l.Db()

func north_south_bridge(a,b interface{}) l.Goal {
        return l.Or(
                db.Deref().Find(a,North_of,b),
                db.Deref().Find(b,North_of,a))
}

func west_east_bridge(a,b interface{}) l.Goal {
        return l.Or(
                db.Deref().Find(a,East_of,b),
                db.Deref().Find(b,East_of,a))
}

func disconnected(a,b interface{}) l.Goal {
	x,y := l.Fresh2()
	return l.And(
		north_south_bridge(a,x),
		l.Neq(x,b),
		west_east_bridge(a,y),
		l.Neq(y,b))
}

func north_of(a,b interface{}) l.Goal {
	return db.Deref().Find(a,North_of,b)
}

func east_of(a,b interface{}) l.Goal {
	return db.Deref().Find(a,East_of,b)
}

func islando (q l.V) l.Goal {
        db.Assert("A",North_of,"C")
        db.Assert("B",North_of,"D")
        db.Assert("D",East_of,"C")
        db.Assert("B",East_of,"A")
        
	membero := l.StructMemberoConstructor4(n)
        
        alabaster_island,durian_island,banana_island := l.Fresh3()
        pwana_island,quero_island,skern_island,rayou_island := l.Fresh4()
        hotel_island,koala_island,jai_island,skating_island := l.Fresh4()
	
        return l.And(
                l.Unify(Nation{Island{"A",v(),v(),v()},Island{"B",v(),v(),v()},Island{"C",v(),v(),v()},Island{"D",v(),v(),v()}},q),
                // 1
		north_of(pwana_island,koala_island),
                //2
		east_of(quero_island,alabaster_island),
                //3
		east_of(hotel_island,durian_island),
                //4
		north_south_bridge(skern_island,jai_island),
                //5
		west_east_bridge(rayou_island,banana_island),
                //6
		disconnected(skating_island,jai_island),
		//
                membero(Island{pwana_island,"Pwana",v(),v()},q),
                membero(Island{quero_island,"Quero",v(),v()},q),
                membero(Island{rayou_island,"Rayou",v(),v()},q),
                membero(Island{skern_island,"Skern",v(),v()},q),
		//
                membero(Island{alabaster_island,v(),"alabaster",v()},q),
                membero(Island{banana_island,v(),"bananas",v()},q),
                membero(Island{v(),v(),"coconuts",v()},q),
                membero(Island{durian_island,v(),"durian fruit",v()},q),
		//
                membero(Island{hotel_island,v(),v(),"hotel"},q),
                membero(Island{skating_island,v(),v(),"ice skating rink"},q),
                membero(Island{jai_island,v(),v(),"jai alai stadium"},q),
                membero(Island{koala_island,v(),v(),"koala preserve"},q),
                )
}

func main() {
        q := l.Fresh()

        c := l.Run(q,islando(q))

        for n := 0 ; n < 10 ; n++ {
                i := <- c
                if i != nil {
                        fmt.Println(i)
                } else {
                        break
                }
        }

}
