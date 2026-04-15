package backend

import (
	"strings"
)

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
    Relation string // "child", "descendant", "adjacent", "sibling"
    Prev *SelectorNode
}

func ParseSelector(input string) *SelectorNode {
	// TOKENIZE
	var isAttribute bool
	var token strings.Builder
	var tokens []string

	for i := 0; i < len(input); i++ {
		char := input[i]

		if isAttribute {
			if char == ']' {
				token.WriteByte(char)
				tokens = append(tokens, token.String())
				token.Reset()
				isAttribute = false
			} else {
				token.WriteByte(char)
			}
			continue
		}

		if char == '[' {
			if token.Len() > 0 {
				tokens = append(tokens, token.String())
				token.Reset()
			}
			isAttribute = true
			token.WriteByte(char)
		} else if char == ' ' {
			if token.Len() > 0 {
				tokens = append(tokens, token.String())
				token.Reset()
			}
			if len(tokens) > 0 && tokens[len(tokens)-1] != " " {
				tokens = append(tokens, " ")
			}
		} else if char == '.' || char == '#' {
			if token.Len() > 0 {
				tokens = append(tokens, token.String())
				token.Reset()
			}
			token.WriteByte(char)
		} else if char == '>' || char == '+' || char == '~' || char == '*' {
			if token.Len() > 0 {
				tokens = append(tokens, token.String())
				token.Reset()
			}
			tokens = append(tokens, string(char))
		} else {
			token.WriteByte(char)
		}
	}

	if token.Len() > 0 {
		tokens = append(tokens, token.String())
	}

	// CONVERT KE SELECTORNODE
	var head *SelectorNode = &SelectorNode{}
	current := head

	for _, t := range tokens {
		if t == " " || t == ">" || t == "+" || t == "~" {
			if t == " " {
				current.Relation = "descendant"
			} else if t == ">" {
				current.Relation = "child"
			} else if t == "+" {
				current.Relation = "adjacent"
			} else if t == "~" {
				current.Relation = "sibling"
			}

			newNode := &SelectorNode{}
			current.Prev = newNode
			current = newNode

		} else if t[0] == '.' {
			current.Classes = append(current.Classes, t[1:]) 
		} else if t[0] == '#' {
			current.ID = t[1:]
		} else if t[0] == '[' {
			attr := AttributeSelector{}
			var builder strings.Builder
			var isVal bool

			for i := 1; i < len(t)-1; i++ {
				char := t[i]

				if isVal {
					if char == '"' {
						attr.Value = builder.String()
						builder.Reset()
						isVal = false
					} else {
						builder.WriteByte(char)
					}
					continue
				}

				if i+1 < len(t)-1 && strings.ContainsRune("~|^$*", rune(char)) && t[i+1] == '=' {
					attr.Name = builder.String()
					builder.Reset()
					attr.Operator = string(char) + "="
					i++ 
				} else if char == '=' {
					attr.Name = builder.String()
					builder.Reset()
					attr.Operator = "="
				} else if char == '"' {
					isVal = true
				} else if char == 'i' || char == 'I' {
					attr.CaseInsensitive = true
				} else if char == ' ' {
					continue
				} else {
					builder.WriteByte(char)
				}
			}

			if attr.Name == "" {
				attr.Name = builder.String()
			}
			current.Attributes = append(current.Attributes, attr)

		} else {
			current.Tag = t
		}
	}
	return head
}

