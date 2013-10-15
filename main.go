package main
import "fmt"
import l "gologic"


func main() {
        v1, v2 := l.Fresh2()
        g := func (a, b int) l.Goal {
                return l.Unify(a,b)
        }
        c := l.Run(v2, l.Or(l.And(l.Unify(v1,3),
                                  l.Unify(v2,v1),
                                  l.Unify(v1,l.Fresh()),
		                  l.Unify(10,l.Fresh()),
                                  g(1,1)),
                            l.Unify(v2,5),
                            l.Unify(v2,8)))

        fmt.Println("# Results 1")

        for {
                i := <- c
                if i != nil {
                        fmt.Println(i)
                } else {
                        break
                }
        }



}
