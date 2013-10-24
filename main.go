package main
import "fmt"
import l "gologic"

func difference(a,b,c interface{}) l.Goal {
        return l.AddC(l.Constraint{
                func (s l.S) (l.S, l.ConstraintResult) {
			fmt.Println("====")
			fmt.Println(a)
			fmt.Println(b)
			fmt.Println(c)
			
                        xo := l.Project(a,s)
                        x, xok := xo.(int)
                        yo := l.Project(b,s)
                        y, yok := yo.(int)
                        zo := l.Project(c,s)
                        z, zok := zo.(int)

                        if xok && yok && zok && x - y == z {
				fmt.Println("a")
                                return s,l.Yes
                        } else if xok && yok && zok && x - y != z {
				fmt.Println("b")
				fmt.Println(x)
				fmt.Println(y)
				fmt.Println(z)
                                return s,l.No
                        } else if xok && yok  {
				fmt.Println("c")
                                ns,success := l.Unifi(c,x-y,s)
				fmt.Println(success)
                                if success {
                                        return ns, l.Yes
                                } else {
                                        return ns, l.No
                                }
                        } else if xok && zok  {
				fmt.Println("d")
                                ns,success := l.Unifi(b,x-z,s)
                                if success {
                                        return ns, l.Yes
                                } else {
                                        return ns, l.No
                                }
                        } else if yok && zok {
				fmt.Println("e")
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

func main() {
        // v1, v2 := l.Fresh2()
        // g := func (a, b int) l.Goal {
        //         return l.Unify(a,b)
        // }
        // c := l.Run(v2, l.Or(l.And(l.Unify(v1,3),
        //                           l.Unify(v2,v1),
        //                           l.Unify(v1,l.Fresh()),
	// 	                  l.Unify(10,l.Fresh()),
        //                           g(1,1)),
        //                     l.Unify(v2,5),
        //                     l.Unify(v2,8)))

        // fmt.Println("# Results 1")

        // for {
        //         i := <- c
        //         if i != nil {
        //                 fmt.Println(i)
        //         } else {
        //                 break
        //         }
        // }

	q,x := l.Fresh2()
	c := l.Run(q,l.And(mult(q,x,4),l.Unify(x,2)))

	fmt.Println(<- c)


}
