package backend

import (
	"fmt"
	"strings"
	"golang.org/x/net/html" //ini dipake buat tokenization using html.Tokenizer ajah
)

//convert raw html jadi DOM tree. html.Tokenizer bakalan ngescan html mentahnya dan ngubah jadi stream of tokens. dari situ kita bisa bangun node tree kita dengan stack :D
func parseHTML(rawhtml string) *Node {
	if strings.TrimSpace(rawHTML) == ""{
		return nil
	}

	tokenizer := html.NewTokenizer(strings.NewReader(rawHTML))
	//virtual root asal buat container sementara dan pijakan awal dulu aja, dan biar lebih seragam :v
	virRoot := NewNode("#virtualroot")
	//stack buat nyimpen nodes yang lagi diproses, it also helps tracks ancestry
	stack := []*Node{virRoot}
	for{
		tt := tokenizer.Next()
		switch tt{
		case html.ErrorToken:
			goto done
		case html.StartTagToken:
			name, hasAttribute := tokenizer.TagName()
			tag := strings.ToLower(string(name))
			node := NewNode(tag)
			node.Parent = stack[len(stack)-1]
			stack[len(stack)-1].Children = append(stack[len(stack)-1].Children, node)
			//parse semua atribut
			for hasAttribute{
				var key, value []byte
				key, value, hasAttribute = tokenizer.TagAttr()
				k := strings.ToLower(string(key))
				v := string(value)
				switch k{
				case "id":
					node.ID = v
				case "class":
					node.Classes = strings.Fields(v)
				default:
					node.Attributes[k] = v
				}
			}
			//hanya push ke stack kalau bukan self-closing tag (gapunya children soalnya)
			if !selfClosingTags[tag]{
				stack = append(stack, node)
			}
		case html.SelfClosingTagToken:
			name, hasAttribute := tokenizer.TagName()
			tag := strings.ToLower(string(name))
			node := NewNode(tag)
			node.Parent = stack[len(stack)-1]
			stack[len(stack)-1].Children = append(stack[len(stack)-1].Children, node)
			//parse semua atribut
			for hasAttribute{
				var key, value []byte
				key, value, hasAttribute = tokenizer.TagAttr()
				k := strings.ToLower(string(key))
				v := string(value)
				switch k{
				case "id":
					node.ID = v
				case "class":
					node.Classes = strings.Fields(v)
				default:
					node.Attributes[k] = v
				}
			}
		case html.EndTagToken:
			name, _ := tokenizer.TagName()
			tag := strings.ToLower(string(name))
			//pop stack sampe ketemu open tag yang matching, buat handle kasus tag yang ga ditutup atau nested dengan tag yang sama
			for len(stack) > 1{
				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if top.Tag == tag{
					break
				}
			}
		case html.TextToken:
			//attach text ke current top node sebagai innertext
			text := string.TrimSpace(string(tokenizer.Text()))
			if text != "" && len(stack) > 0{
				top := stack[len(stack)-1]
				if top.InnerText == ""{
					top.InnerText = text
				}else{
					top.InnerText += " " + text
				}
			}
		case html.CommentToken: 
			//boleh diskip
			continue
		case html.DoctypeToken:
			continue
		}
	}
	done:
	//kembaliin anak-anak virtual root, kalau punya satu anak aja kembaliin itu sebagai root
	if len(virRoot.Children) == 1{
		root := virRoot.Children[0]
		root.Parent = nil
		return root
	}
	if len(virRoot.Children) == 0{
		return nil
	}
	return virRoot
}

//ini buat basic structural checking pada DOM tree, bakalan return error string kalau ada issues, empty string kalau aman
func validateHTML(root *Node) error{
	if root == nil{
		return fmt.Errorf("html structure error, tree kosong (html kosong atau invalid)")
	}
	//cek parent-child consistency, pastiin setiap anak node punya parent yang bener
	var check func(n *Node) error
	check = func(n *Node) error{
		for _, child := range n.Children{
			if child.Parent != n{
				return fmt.Errorf("html structure error, node <%s> punya parent yang salah", child.Tag)
			}
			if err := check(child); err != nil{
				return err
			}
		}
		return nil
	}
	return check(root)
}

//lists of all html tags yang child-free
var selfClosingTags = map[string]bool{
	"area": true,
	"base": true,
	"br": true,
	"col": true,
	"embed": true,
	"hr": true,
	"img": true,
	"input": true,
	"link": true,
	"meta": true,
	"param": true,
	"source": true,
	"track": true,
	"wbr": true,
}