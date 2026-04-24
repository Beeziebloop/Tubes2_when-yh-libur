package backend
import "strings"

type Node struct{
	Tag string //buat nyimpen tag HTML seperti div, p, span, ect.
	ID string //buat nyimpen id (id="...")
	Classes []string //buat nyimpen kelas-kelas (class="...")
	Attributes map[string]string //buat atribut-atribut lain
	Children []*Node //karena node anaknya bisa banyak
	Parent *Node
	InsideText string //isi dalam tag
}

//ctor
func NewNode(tag string) *Node{
	return &Node{
		Tag: strings.ToLower(tag),
		ID: "",
		Classes: []string{},
		Attributes: make(map[string]string),
		Children: []*Node{},
		Parent: nil,
		InsideText: "",
	}
}

//buat ngecek apakah node punya kelas-kelas tertentu
func (n *Node) HasClasses(class string) bool{
	for _, c := range n.Classes{
		if c == class{
			return true
		}
	}
	return false
}

//ngembaliin nilai atribut dari map, atau string kosong kalo gaada
func (n *Node) GetAttribute(key string) (string, bool) {
    val, ok := n.Attributes[key]
    return val, ok
}

//return semua sanak saudara yang datang sebelum node sekarang, dipake buat adjacent dan general sibling selector
func (n *Node) SiblingsBeforeNow() []*Node{
	if n.Parent == nil{
		return []*Node{}
	}
	result := []*Node{}
	for _, sibling := range n.Parent.Children{
		if sibling == n{
			break
		}
		result = append(result, sibling)
	}
	return result
}

//ngasih tau siapa exactly saudara (singular) yang datang sebelum node sekarang
func (n *Node) ImmediateSiblingsBeforeNow() *Node{
	siblings := n.SiblingsBeforeNow()
	if len(siblings) == 0{
		return nil
	}
	return siblings[len(siblings)-1]
}

//buat ngitung kedalaman node
func (n *Node) Depth() int{
	depth := 0
	current := n.Parent
	for current != nil{
		depth++
		current = current.Parent
	}
	return depth
}