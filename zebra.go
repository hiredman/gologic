package main
import "fmt"
//import "reflect"
import l "gologic"

type H struct {
	N interface {}
	P interface {}
	D interface {}
	C interface {}
}

type B struct {
	H1 interface {}
	H2 interface {}
	H3 interface {}
	H4 interface {}
	H5 interface {}
}


type cons struct {
	CAR interface {}
	CDR interface {}
}

func li (foo ...interface{}) interface {} {
	var x interface {} = nil
	for i:=len(foo) ; i>-1 ; i++ {
		x = &cons{foo[i],x}
	}
	return x
}

func main() {
	hs := l.Fresh()
	c := l.Run(hs, l.And(l.Unify(hs,B{l.A, l.A, H{l.A, l.A, "milk", l.A}, l.A, l.A}),
	                     l.Unify(hs,B{H{"norwegian", l.A, l.A, l.A}, l.A, l.A, l.A, l.A})))
	for {
		i := <- c
		if i != nil {
			fmt.Println(i)
		} else {
			break
		}
	}

   // type t struct {
   //      N int
   //  }
   //  var n = t{42}
   //  fmt.Println(n.N)
   //  reflect.ValueOf(&n).Elem().FieldByName("N").SetInt(7)
   //  fmt.Println(n.N)
}
