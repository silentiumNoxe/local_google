package html

type Data struct {
	Lang    string
	Title   string
	Links   []string
	Content []string
}

type Node struct {
	Tag     string
	Attr    map[string]string
	Content string

	FirstChild *Node
	LastChild  *Node
	Parent     *Node
	Next       *Node
	Prev       *Node
}

func (n *Node) Append(child *Node) {
	child.Parent = n
	if n.FirstChild == nil {
		n.FirstChild = child
		n.LastChild = child
		return
	}

	var last = n.LastChild
	n.LastChild = child
	last.Next = child
	child.Prev = last
}

type Token struct {
	Type    int
	Content string
}

const OpenTagToken = 1
const CloseTagToken = 2
const TextToken = 3
