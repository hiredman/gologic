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

func membero(p, t interface{}) l.Goal {
        return l.Or(
                l.Unify(Children{  p,v(),v(),v(),v()}, t),
                l.Unify(Children{v(),  p,v(),v(),v()}, t),
                l.Unify(Children{v(),v(),  p,v(),v()}, t),
                l.Unify(Children{v(),v(),v(),  p,v()}, t),
                l.Unify(Children{v(),v(),v(),v(),  p}, t))

}

func difference(a,b,c interface{}) l.Goal {
        return l.AddC(l.Constraint{
                func (s l.S) (l.S, l.ConstraintResult) {
                        xo := l.Project(a,s)
                        x, xok := xo.(int)
                        yo := l.Project(b,s)
                        y, yok := yo.(int)
                        zo := l.Project(c,s)
                        z, zok := zo.(int)

                        if xok && yok && zok && x - y == z {
                                return s,l.Yes
                        } else if xok && yok && zok && x - y != z {
                                return s,l.No
                        } else if xok && yok  {
                                ns,success := l.Unifi(c,x-y,s)
                                if success {
                                        return ns, l.Yes
                                } else {
                                        return ns, l.No
                                }
                        } else if xok && zok  {
                                ns,success := l.Unifi(b,x-z,s)
                                if success {
                                        return ns, l.Yes
                                } else {
                                        return ns, l.No
                                }
                        } else if yok && zok {
                                ns,success := l.Unifi(a,y+z,s)
                                if success {
                                        return ns, l.Yes
                                } else {
                                        return ns, l.No
                                }
                        } else {
                                return s, l.Maybe
                        }


        }})
}

func mult(a,b,c interface{}) l.Goal {
        return l.AddC(l.Constraint{
                func (s l.S) (l.S, l.ConstraintResult) {
			//fmt.Println("mult")
                        xo := l.Project(a,s)
                        x, xok := xo.(int)
                        yo := l.Project(b,s)
                        y, yok := yo.(int)
                        zo := l.Project(c,s)
                        z, zok := zo.(int)
			
                        if xok && yok && zok && x * y == z {
				//fmt.Println("a")
                                return s,l.Yes
                        } else if xok && yok && zok && x * y != z {
				//fmt.Println("b")
                                return s,l.No
                        } else if xok && yok  {
				//fmt.Println("c")
                                ns,success := l.Unifi(c,x*y,s)
                                if success {
                                        return ns, l.Yes
                                } else {
                                        return ns, l.No
                                }
                        } else if xok && zok  {
                                ns,success := l.Unifi(b,z/x,s)
                                if success {
                                        return ns, l.Yes
                                } else {
                                        return ns, l.No
                                }
                        } else if yok && zok {
                                ns,success := l.Unifi(a,z/y,s)
                                if success {
                                        return ns, l.Yes
                                } else {
                                        return ns, l.No
                                }
                        } else {
                                return s, l.Maybe
                        }


        }})
}


func babbyo (q l.V) l.Goal {
	keiths_last_name,libbys_last_name,margos_last_name,noras_last_name,ottos_last_name := l.Fresh5()
	keiths_age,noras_age,margos_age,ottos_age,libbys_age := l.Fresh5()
	iveys_age,fells_age,halls_age := l.Fresh3()
        return l.And(
                l.Unify(Children{Child{"Keith",keiths_last_name,keiths_age},
		                 Child{"Libby",libbys_last_name,libbys_age},
		                 Child{"Margo",margos_last_name,margos_age},
                                 Child{"Nora",noras_last_name,noras_age},
                                 Child{"Otto",ottos_last_name,ottos_age}},q),
		
		//2
		difference(keiths_age,1,iveys_age),
		difference(iveys_age,1,noras_age),

		//3
		difference(fells_age,3,margos_age),

		//4
		mult(halls_age,2,ottos_age),

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
