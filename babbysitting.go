//http://brownbuffalo.sourceforge.net/BabysittingClues.html
package main
import "fmt"
import l "gologic"

type Child struct {
        First, Last, Age interface{}
}

type Children struct {
        C1,C2,C3,C4,C5 interface{}
}

func v () l.V {
        return l.Fresh()
}

func c (a,b,c,d,e interface{}) interface {} {
	return Children{a,b,c,d,e}
}

func babbyo (q l.V) l.Goal {
	keiths_last_name,libbys_last_name,margos_last_name,noras_last_name,ottos_last_name := l.Fresh5()
	keiths_age,noras_age,margos_age,ottos_age,libbys_age := l.Fresh5()
	iveys_age,fells_age,halls_age := l.Fresh3()
	membero := l.StructMemberoConstructor5(c)
        return l.And(
                l.Unify(Children{Child{"Keith",keiths_last_name,keiths_age},
		                 Child{"Libby",libbys_last_name,libbys_age},
		                 Child{"Margo",margos_last_name,margos_age},
                                 Child{"Nora",noras_last_name,noras_age},
                                 Child{"Otto",ottos_last_name,ottos_age}},q),
		
		//2
		l.Difference(keiths_age,1,iveys_age),
		l.Difference(iveys_age,1,noras_age),

		//3
		l.Difference(fells_age,3,margos_age),

		//4
		l.Mult(halls_age,2,ottos_age),

		// 1
		membero(Child{"Libby","Jule",v()},q),
				
		//
		membero(Child{v(),"Fell",fells_age},q),
		membero(Child{v(),"Gant",v()},q),
		membero(Child{v(),"Hall",halls_age},q),
		membero(Child{v(),"Ivey",iveys_age},q),
		membero(Child{v(),"Jule",v()},q),
		membero(Child{v(),v(),2},q),
		membero(Child{v(),v(),3},q),
		membero(Child{v(),v(),4},q),
		membero(Child{v(),v(),5},q),
		membero(Child{v(),v(),6},q),

                )
}

func main() {
        q := l.Fresh()

        c := l.Run(q,babbyo(q))

        for n := 0 ; n < 10 ; n++ {
                i := <- c
                if i != nil {
                        fmt.Println(i)
                } else {
                        break
                }
        }

}
