package backend

type AttributeSelector struct {
    Name string
    Operator string
    Value string
	CaseInsensitive bool
}

type SelectorNode struct {
    Tag string
    ID string
    Classes []string
	Attributes []AttributeSelector
    Relation string // "child", "descendant", "adjacent_sibling", "general_sibling"
    Prev *SelectorNode
}
