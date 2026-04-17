package backend

import (
	"strings"
)

func ParseSelector(input string) *SelectorNode {
	tokens := tokenize(input)
	return convertTokensToNodes(tokens)
}

// ubah string input menjadi tokens
func tokenize(input string) []string {
	var isAttribute bool
	var token strings.Builder
	var tokens []string

	for i := 0; i < len(input); i++ {
		char := input[i]

		if isAttribute {
			token.WriteByte(char)
			if char == ']' {
				tokens = append(tokens, token.String())
				token.Reset()
				isAttribute = false
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
	return tokens
}

// ubah tokens menjadi SelectorNode
func convertTokensToNodes(tokens []string) *SelectorNode {
	head := &SelectorNode{}
	current := head

	for _, t := range tokens {
		switch {
		case t == " " || t == ">" || t == "+" || t == "~":
			current.Relation = mapRelation(t)
			newNode := &SelectorNode{}
			current.Prev = newNode
			current = newNode

		case t[0] == '.':
			current.Classes = append(current.Classes, t[1:])
		case t[0] == '#':
			current.ID = t[1:]
		case t[0] == '[':
			attr := parseAttribute(t)
			current.Attributes = append(current.Attributes, attr)
		default:
			current.Tag = t
		}
	}
	return head
}

// helper untuk parse atribut
func parseAttribute(t string) AttributeSelector {
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
	return attr
}

// helper untuk mapping combinator ke nama relasi
func mapRelation(t string) string {
	switch t {
	case " ": return "descendant"
	case ">": return "child"
	case "+": return "adjacent_sibling"
	case "~": return "general_sibling"
	}
	return ""
}