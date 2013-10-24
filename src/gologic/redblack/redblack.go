package redblack
import "container/list"

type color int

const (
	red color = iota
	black
)

type Element interface {
	Key() int
	Merge(Element) Element
}

type Rbnode struct {
	color color
	element Element
	left *Rbnode
	right *Rbnode
}

func Node(e Element) *Rbnode {
	return &Rbnode{red,e,nil,nil}
}

func case1 (color color, a *Rbnode, y Element, b *Rbnode) bool {
	if color == black && a != nil && a.color == red {
		a_ := a.left
		return a_.color == red
	} else {
		return false
	}
}

func case2 (color color, a *Rbnode, y Element, b *Rbnode) bool {
	if color == black && a != nil && red == a.color {
		a_ := a.right
		return red == a_.color
	} else {
		return false
	}
}

func case3 (color color, a *Rbnode, y Element, b *Rbnode) bool {
	if color == black && b != nil && red == b.color {
		b_ := b.left
		return red == b_.color
	} else {
		return false
	}
}

func case4 (color color, a *Rbnode, y Element, b *Rbnode) bool {
	if color == black && b != nil && red == b.color {
		b_ := b.right
		return red == b_.color
	} else {
		return false
	}
}

func make_black (s *Rbnode) *Rbnode {
	if s.color == red {
		return s
	} else {
		return &Rbnode{black, s.element, s.left, s.right}
	}
}

func tree(a *Rbnode,x Element,b *Rbnode,y Element,c *Rbnode, z Element,d *Rbnode) *Rbnode {
	return &Rbnode{red,y,&Rbnode{black,x,a,b},&Rbnode{black,z,c,d}}
}

func balance (x Element, color color, a *Rbnode, y Element, b *Rbnode) *Rbnode {
	if case1(color, a, y, b) {
		z := y
		d := b
		y := a.element
		a_left := a.left
		c := a.right
		x := a_left.element
		a := a_left.left
		b := a_left.right
		return tree(a,x,b,y,c,z,d)
	} else if case2(color,a,y,b) {
		z := y
		d := b
		x := a.element
		a := a.left
		a_right := a.right
		y := a_right.element
		b := a_right.left
		c := a_right.right
		return tree(a,x,b,y,c,z,d)
	} else if case3(color,a,y,b) {
		z := b.element
		b_left := b.left
		d := b.right
		y := b_left.element
		b := b_left.left
		c := b_left.right
		return tree(a,x,b,y,c,z,d)
	} else if case4(color,a,y,b) {
		y := b.element
		b := b.left
		b_right := b.right
		z := b_right.element
		c := b_right.left
		d := b_right.right
		return tree(a,x,b,y,c,z,d)
	} else {
		return &Rbnode{color,y,a,b}
	}
}

func ins (tree *Rbnode, x Element) *Rbnode {
	if tree == nil {
		return Node(x);
	} else {
		color := tree.color
		y := tree.element
		a := tree.left
		b := tree.right
		y_key := y.Key()
		x_key := x.Key()
		if y_key > x_key {
			return balance(x, color, ins(a,x), y, b)
		} else if x_key == y_key {
			return &Rbnode{color,y.Merge(x),a,b}
		} else {
			return balance(x,color,a,y,ins(b,x))
		}
	}
}

func Insert (t *Rbnode, e Element) *Rbnode {
	return make_black(ins(t,e))
}

func Locate (tree *Rbnode, n int) (Element, bool) {
	for  {
		if tree == nil {
			return nil, false
		} else {
			x := tree.element
			l := tree.left
			r := tree.right
			if x.Key() == n {
				return tree.element, true
			} else if x.Key() > n {
				tree = l
 			} else {
				tree = r
			}
		}
	}
}

type ducer func (interface{}, Element) (interface{}, bool)

func Fold (init interface{}, f ducer, tree *Rbnode) (interface{}, bool) {
	c := make(chan Element)
	go func () {
		stack := list.New()
		stack.PushFront(tree)
		for {
			if stack.Len() > 0 {
				a := stack.Front().Value
				x, ok := a.(Rbnode)
				if !ok {panic("oh no")}
				c <- x.element
				if x.left != nil {
					stack.PushFront(x.left)
				}
				if x.right != nil {
					stack.PushFront(x.right)
				}

			} else {
				break
			}
		}
	}()
	r := init
	for item := <- c; item != nil; item = <- c {
		a,cont := f(r,item)
		if !cont {
			close(c)
			return a,cont
		}
		r = a
	}
	return r,true
} 
