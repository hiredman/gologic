//http://brownbuffalo.sourceforge.net/RoundOfGolfClues.html
package main
import "fmt"
import l "gologic"

type Dude struct {
        First, Last, Job, Score interface{}
}

type Round struct {
        D1,D2,D3,D4 interface{}
}

func v () l.V {
        return l.Fresh()
}

func betweeno(p interface{}, low, high int) l.Goal {
        if low > high {
                return l.Fail()
        } else {
                return l.Or(
                        l.Unify(low,p),
                        l.Call(betweeno,p,low+1,high))
        }
}


func scoreo (q interface{}) l.Goal {
        a,b,c,d := l.Fresh4()
        return l.And(
                l.Unify(Round{Dude{v(),v(),v(),a}, Dude{v(),v(),v(),b}, Dude{v(),v(),v(),c}, Dude{v(),v(),v(),d}},q),
                betweeno(a,70,85),
                betweeno(b,70,85),
                betweeno(c,70,85),
                betweeno(d,70,85))
}

func golfo (q l.V) l.Goal {
	membero := l.StructMemberoConstructor4(func (a,b,c,d interface{}) interface{} {return Round{a,b,c,d}})
        bills_job := v()
        bills_score := v()
        mr_clubb_first_name := v()
        mr_clubbs_score := v()
        pro_shop_clerk_score := v()
        frank_score := v()
        caddy_score := v()
        sands_score := v()
        score1,score2,score3,score4 := l.Fresh4()
	mr_carters_first_name := v()
        return l.And(
                l.Unify(Round{Dude{"Bill",v(),v(),score1},Dude{"Jack",v(),v(),score2},Dude{"Frank",v(),v(),score3},Dude{"Paul",v(),v(),score4}},q),
                l.Neq(score1,score2),
                l.Neq(score1,score3),
                l.Neq(score1,score4),
                l.Neq(score2,score3),
                l.Neq(score2,score4),
                l.Neq(score3,score4),

                membero(Dude{"Jack", v(), v(), v()}, q),
                membero(Dude{v(), "Green", v(), v()}, q),
                membero(Dude{v(), v(), "short-order cook", v()}, q),
				
                // // 1
                membero(Dude{"Bill", v(), bills_job, bills_score}, q),
                l.Neq(bills_job,"maintenance man"),
		membero(Dude{v(), v(), "maintenance man", v()}, q),
                l.Increasing(bills_score,score2),
		l.Increasing(bills_score,score3),
		l.Increasing(bills_score,score4),

                // // 2
                membero(Dude{mr_clubb_first_name, "Clubb", v(), mr_clubbs_score}, q),
                l.Neq(mr_clubb_first_name, "Paul"),
                membero(Dude{v(), v(), "pro-shop clerk", pro_shop_clerk_score}, q),
		l.Difference(mr_clubbs_score,10,pro_shop_clerk_score),

                // //3
                membero(Dude{"Frank", v(), v(), frank_score}, q),
                membero(Dude{v(), v(), "caddy", caddy_score}, q),
                membero(Dude{v(), "Sands", v(), sands_score}, q),
                
                l.Or(l.And(l.Difference(frank_score, 7, sands_score),
                           l.Difference(caddy_score, 4, sands_score)),
                     l.And(l.Difference(frank_score, 4, sands_score),
                           l.Difference(caddy_score, 7, sands_score))),

                // // // 4
                membero(Dude{mr_carters_first_name, "Carter", v(), 78}, q),
                l.Increasing(frank_score, 78),
		l.Neq(mr_carters_first_name,"Frank"),

                // // // 5
                l.Neq(score1,81),
                l.Neq(score2,81),
                l.Neq(score3,81),
                l.Neq(score4,81),
		
                scoreo(q),
                

                )
}

func main() {
        q := l.Fresh()

        c := l.Run(q,golfo(q))

        for n := 0 ; n < 10 ; n++ {
                i := <- c
                if i != nil {
                        fmt.Println(i)
                } else {
                        break
                }
        }

}
