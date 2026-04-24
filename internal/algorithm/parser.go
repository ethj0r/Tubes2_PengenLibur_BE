package algorithm

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	net_html "golang.org/x/net/html"
)

// curr relevant
// https://brightdata.com/blog/web-data/parse-html-with-golang
// https://pkg.go.dev/golang.org/x/net/html

// Void elmts tokenizer emits StartTagToken, tapi secara spec
// HTML5 mereka nggak punya closing tag. so hrus handle sebagai self-closing
// supaya hierarchy tree nggak bocor.
// ref: https://html.spec.whatwg.org/multipage/syntax.html#void-elements
var voidElements = map[string]bool{
	"area": true, "base": true, "br": true, "col": true, "embed": true,
	"hr": true, "img": true, "input": true, "link": true, "meta": true,
	"param": true, "source": true, "track": true, "wbr": true,
}

func TokenToNode(token *net_html.Token) *Node {
	isText := token.Type == net_html.TextToken
	attributes := []Attribute{}
	if !isText {
		for i := 0; i < len(token.Attr); i++ {
			attributes = append(attributes, Attribute{
				Name:  token.Attr[i].Key,
				Value: token.Attr[i].Val,
			})
		}
	}

	newElmt := Element{
		Data:       token.Data,
		IsText:     isText,
		Attributes: attributes,
	}

	return NewNode(newElmt)
}

func Parse(src []byte) (*Node, error) {
	source := bytes.NewReader(src)
	tokenizer := net_html.NewTokenizer(source)

	var root *Node = nil
	var currentParent *Node = nil

	for {
		tok := tokenizer.Next()
		token := tokenizer.Token()

		switch tok {
		case net_html.ErrorToken:
			if tokenizer.Err() == io.EOF {
				if root != nil {
					return root, nil
				}
				return currentParent, nil
			}
			return nil, tokenizer.Err()
		case net_html.TextToken:
			if strings.TrimSpace(token.Data) == "" {
				continue
			}
			if currentParent != nil {
				currentParent.AddChild(TokenToNode(&token))
			}
		case net_html.SelfClosingTagToken:
			node := TokenToNode(&token)
			if currentParent != nil {
				currentParent.AddChild(node)
			} else if root == nil {
				root = node
			}
		case net_html.StartTagToken:
			node := TokenToNode(&token)
			if currentParent != nil {
				currentParent.AddChild(node)
			} else if root == nil {
				root = node
			}
			if !voidElements[token.Data] {
				currentParent = node
			} else if currentParent == nil {
				currentParent = node
			}
		case net_html.EndTagToken:
			if currentParent != nil && currentParent.Parent != nil {
				currentParent = currentParent.Parent
			}
		default:
			continue
		}
	}
}

// Check Struktur HTML (apakah valid atau tidak)
func ValidateHTML(src []byte) error {
	tokenizer := net_html.NewTokenizer(bytes.NewReader(src))
	stack := make([]string, 0, 32)

	for {
		tok := tokenizer.Next()
		token := tokenizer.Token()

		switch tok {
		case net_html.ErrorToken:
			if tokenizer.Err() == io.EOF {
				if len(stack) > 0 {
					return fmt.Errorf("unclosed tag <%s>", stack[len(stack)-1])
				}
				return nil
			}
			return tokenizer.Err()

		case net_html.StartTagToken:
			if !voidElements[token.Data] {
				stack = append(stack, token.Data)
			}

		case net_html.EndTagToken:
			if len(stack) == 0 {
				return fmt.Errorf("unexpected closing tag </%s>", token.Data)
			}

			top := stack[len(stack)-1]
			if token.Data != top {
				return fmt.Errorf("mismatched closing tag </%s>, expected </%s>", token.Data, top)
			}

			stack = stack[:len(stack)-1]
		}
	}
}
