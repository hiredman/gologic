//https://github.com/swannodette/logic-tutorial#zebras
package main
import "fmt"
import l "gologic"

type H struct {
        Nat interface {}
	Smo interface {}
        Dri interface {}
        Pet interface {}
        Col interface {}
}

type B struct {
        H1 interface {}
        H2 interface {}
        H3 interface {}
        H4 interface {}
        H5 interface {}
}

func v() l.V {
        return l.Fresh()
}


func righto (a,b,block interface{}) l.Goal {
        return l.Or(
                l.Unify(B{  a,  b,v(),v(),v()}, block),
                l.Unify(B{v(),  a,  b,v(),v()}, block),
                l.Unify(B{v(),v(),  a,  b,v()}, block),
                l.Unify(B{v(),v(),v(),  a,  b}, block))
}

func nexto (a,b,block interface{}) l.Goal {
        return l.Or(
                righto(a,b,block),
                righto(b,a,block))
}

func firsto (block, a interface {}) l.Goal {
        return l.Unify(block,B{a,v(),v(),v(),v()})
}

func membero (a, block interface{}) l.Goal {
        return l.Or(
                l.Unify(B{  a,v(),v(),v(),v()}, block),
                l.Unify(B{v(),  a,v(),v(),v()}, block),
                l.Unify(B{v(),v(),  a,v(),v()}, block),
                l.Unify(B{v(),v(),v(),  a,v()}, block),
                l.Unify(B{v(),v(),v(),v(),  a}, block))

}

func zerbao (q l.V) l.Goal {
        return l.And(
                l.Unify(B{v(),v(),H{v(),v(),"milk",v(),v()}, v(), v()}, q),
                firsto(q, H{"norwegian",v(),v(),v(),v()}),
                nexto(H{"norwegian",v(),v(),v(),v()}, H{v(),v(),v(),v(),"blue"}, q),
                righto(H{v(),v(),v(),v(),"ivory"},H{v(),v(),v(),v(),"green"},q),
                membero(H{"englishmen",v(),v(),v(),"red"},q),
                membero(H{v(),"kools",v(),v(),"yello"},q),
                membero(H{"spaniard",v(),v(),"dog",v()},q),
                membero(H{v(),v(),"coffee",v(),"green"},q),
                membero(H{"ukranian",v(),"tea",v(),v()},q),
                membero(H{v(),"lucky strikes","oj",v(),v()},q),
                membero(H{"japanese","parliaments",v(),v(),v()},q),
                membero(H{v(),"oldgolds",v(),"snails",v()},q),
                nexto(H{v(),v(),v(),"horse",v()},H{v(),"kools",v(),v(),v()},q),
                nexto(H{v(),v(),v(),"fox",v()},H{v(),"chesterfields",v(),v(),v()},q))

}

func main() {
        q := l.Fresh()

        c := l.Run(q,zerbao(q))
	
        for {
                i := <- c
                if i != nil {
                        fmt.Println(i)
                } else {
                        break
                }
        }

}
