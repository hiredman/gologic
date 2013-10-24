package gologic
import "testing"
import "fmt"

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

func (b B) String() string {
	return fmt.Sprintf("Block:\n  %s\n  %s\n  %s\n  %s\n  %s",b.H1,b.H2,b.H3,b.H4,b.H5)
}

// func v() V {
//         return Fresh()
// }


func righto (a,b,block interface{}) Goal {
        return Or(
                Unify(B{  a,  b,v(),v(),v()}, block),
                Unify(B{v(),  a,  b,v(),v()}, block),
                Unify(B{v(),v(),  a,  b,v()}, block),
                Unify(B{v(),v(),v(),  a,  b}, block))
}

func nexto (a,b,block interface{}) Goal {
        return Or(
                righto(a,b,block),
                righto(b,a,block))
}

func firsto (block, a interface {}) Goal {
        return Unify(block,B{a,v(),v(),v(),v()})
}

func membero (a, block interface{}) Goal {
        return Or(
                Unify(B{  a,v(),v(),v(),v()}, block),
                Unify(B{v(),  a,v(),v(),v()}, block),
                Unify(B{v(),v(),  a,v(),v()}, block),
                Unify(B{v(),v(),v(),  a,v()}, block),
                Unify(B{v(),v(),v(),v(),  a}, block))

}

func zerbao (q V) Goal {
        return And(
                Unify(B{v(),v(),H{v(),v(),"milk",v(),v()}, v(), v()}, q),
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
                nexto(H{v(),v(),v(),"fox",v()},H{v(),"chesterfields",v(),v(),v()},q),
                membero(H{v(),v(),v(),"zebra",v()},q),
		membero(H{v(),v(),"water",v(),v()},q),)

}

func BenchmarkZerbra(b *testing.B) {
	for  i := 0; i < b.N; i++ {
		q := Fresh()
		
		c := Run(q,zerbao(q))
		
		<- c
	}
}
