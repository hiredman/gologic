//http://brownbuffalo.sourceforge.net/BreakingNewsClues.html
package main
import "fmt"
import l "gologic"

type Report struct {
        Name,Location,Story interface{}
}

type News struct {
        C1,C2,C3,C4 interface{}
}

func v () l.V {
        return l.Fresh()
}

func membero(p, t interface{}) l.Goal {
        return l.Or(
                l.Unify(News{  p,v(),v(),v()}, t),
                l.Unify(News{v(),  p,v(),v()}, t),
                l.Unify(News{v(),v(),  p,v()}, t),
		l.Unify(News{v(),v(),v(),  p}, t))
}

func newso (q l.V) l.Goal {
	baby_location := v()
	jimmy_location := v()
	lois_story := v()
	pc_story := v()
	sm_story := v()
	corey_location := v()
	whale_location := v()
        return l.And(
		l.Unify(News{Report{"Corey",corey_location,v()},Report{"Jimmy",jimmy_location,v()},Report{"Lois",v(),lois_story},Report{"Perry",v(),v()}},q),

		// 1
		l.Neq(baby_location,"South Amboy"),
		l.Neq(baby_location,"New Hope"),

		// 2
		l.Neq(jimmy_location,"Port Charles"),

		// 3
		l.Or(l.And(l.Unify(lois_story,"blimp launching"),l.Unify(pc_story,"skyscraper dedication")),
                     l.And(l.Unify(lois_story,"skyscraper dedication"),l.Unify(pc_story,"blimp launching"))),

		// 5
		l.Or(l.Unify(corey_location,"Bayonne"),
                     l.Unify(whale_location,"Bayonne"),
		     l.And(l.Unify(corey_location,"Bayonne"),
        	           l.Unify(whale_location,"Bayonne"))),

		//
		membero(Report{v(),"Bayonne",v()},q),
		membero(Report{v(),"New Hope",v()},q),
		membero(Report{v(),"Port Charles",pc_story},q),
		membero(Report{v(),"South Amboy",sm_story},q),
				
		membero(Report{v(),baby_location,"30 pound baby"},q),
		membero(Report{v(),v(),"blimp launching"},q),
		membero(Report{v(),v(),"skyscraper dedication"},q),
		membero(Report{v(),whale_location,"beached whale"},q),
		
		// 4
		l.Neq(sm_story, "beached whale"),
		l.Neq(sm_story,"skyscraper dedication"),
		
                )
}

func main() {
        q := l.Fresh()

        c := l.Run(q,newso(q))

        for n := 0 ; n < 10 ; n++ {
                i := <- c
                if i != nil {
                        fmt.Println(i)
                } else {
                        break
                }
        }

}
