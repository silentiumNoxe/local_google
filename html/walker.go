package html

type Walker struct {
	root    *Node
	current *Node
}

func NewWalker(root *Node) *Walker {
	return &Walker{root: root}
}

func (w *Walker) Next() *Node {
	var x = w.current

	if x == nil {
		x = w.root
		w.current = x
		return x
	}

	if x.FirstChild != nil {
		w.current = x.FirstChild
		return x.FirstChild
	}

	if x.Next != nil {
		w.current = x.Next
		return x.Next
	}

	x = x.Parent
	if x == nil {
		return nil
	}

	k := x
	y := k.Next
	for {
		if y != nil {
			break
		}

		k = k.Parent
		if k == nil {
			break
		}

		y = k.Next
	}

	w.current = y
	return y
}
