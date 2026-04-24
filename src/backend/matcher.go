package main

import "strings"

func MakeMatchFunc(selector string) MatchFunc {
    parsed := ParseSelector(selector)
    return func(node *Node, _ string) bool {
        return MatchSelector(node, parsed)
    }
}

func MatchSelector(dom *Node, selector *SelectorNode) bool {
	return matchSelectorRecursive(dom, selector)
}

func matchSelectorRecursive(node *Node, sel *SelectorNode) bool {
	if node == nil {
		return false
	}

	// cek node sekarang apakah cocok
	if !matchNode(node, sel) {
		return false
	}

	// kalo paling kiri udah berarti selesai
	if sel.Prev == nil {
		return true
	}

	// lanjut sesuai relation
	switch sel.Relation {
	case "child":
		return matchSelectorRecursive(node.Parent, sel.Prev)

	case "descendant":
		cur := node.Parent
		for cur != nil {
			if matchSelectorRecursive(cur, sel.Prev) {
				return true
			}
			cur = cur.Parent
		}
		return false

	case "adjacent_sibling":
		return matchSelectorRecursive(node.ImmediateSiblingsBeforeNow(), sel.Prev)

	case "general_sibling":
		siblings := node.SiblingsBeforeNow()
		for i := len(siblings) - 1; i >= 0; i-- {
			if matchSelectorRecursive(siblings[i], sel.Prev) {
				return true
			}
		}
		return false
	}

	return false
}

func matchNode(dom *Node, selector *SelectorNode) bool {
	if dom == nil {
		return false
	}

	// tag
	if selector.Tag != "" && selector.Tag != "*" && selector.Tag != dom.Tag {
		return false
	}

	// id
	if selector.ID != "" && selector.ID != dom.ID {
		return false
	}

	// class
	for _, c := range selector.Classes {
		if !dom.HasClasses(c) {
			return false
		}
	}

	// attribute
	if !matchAttributes(dom, selector) {
		return false
	}

	return true
}

func matchAttributes(dom *Node, selector *SelectorNode) bool {
	for _, attr := range selector.Attributes {
		var val string
		var found bool

		// ambil dari dom
		if attr.Name == "id" {
			val = dom.ID
			found = dom.ID != ""
		} else if attr.Name == "class" {
			val = strings.Join(dom.Classes, " ")
			found = len(dom.Classes) > 0
		} else {
			val, found = dom.GetAttribute(attr.Name)
		}

		if !found {
			return false
		}

		a := val
		b := attr.Value

        // cek case insensitive
		if attr.CaseInsensitive {
			a = strings.ToLower(a)
			b = strings.ToLower(b)
		}

        // cek logic operator
		switch attr.Operator {
		case "":
			continue

		case "=":
			if a != b {
				return false
			}

		case "~=":
			words := strings.Fields(a)
			ok := false
			for _, w := range words {
				if w == b {
					ok = true
					break
				}
			}
			if !ok {
				return false
			}

		case "|=":
			if a != b && !strings.HasPrefix(a, b+"-") {
				return false
			}

		case "^=":
			if b == "" || !strings.HasPrefix(a, b) {
				return false
			}

		case "$=":
			if b == "" || !strings.HasSuffix(a, b) {
				return false
			}

		case "*=":
			if b == "" || !strings.Contains(a, b) {
				return false
			}
		}
	}
	return true
}