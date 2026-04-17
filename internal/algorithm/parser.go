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
// Yang dipanggil untuk ngelakuin full

// https://pkg.go.dev/golang.org/x/net/html

func TokenToNode(token *net_html.Token) *Node {
	isText := false
	attributes := []Attribute{}
	if token.Type == net_html.TextToken {
		isText = false
	}

	if !isText {
		for i := 0; i < (len(token.Attr)); i++ {
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

	// Iterate through each token

	var currentParent *Node = nil
	fmt.Println("Here!")
	for {
		tok := tokenizer.Next()
		token := tokenizer.Token()

		var currentNode *Node = nil

		switch tok {
		case net_html.ErrorToken: // aman
			if tokenizer.Err() == io.EOF {
				return currentParent, nil
			}
			return nil, tokenizer.Err()
		case net_html.TextToken: // aman
			// Proses Text
			if strings.TrimSpace(token.Data) == "" {
				continue
			}
			// add ke parent
			if currentParent != nil {
				currentParent.AddChild(TokenToNode(&token))
			}
		case net_html.SelfClosingTagToken:
			// add ke parent
			if currentParent != nil {
				currentParent.AddChild(TokenToNode(&token))
			}
		case net_html.StartTagToken:
			// kasus khusus kalo <meta> ni somehow di-classify jadi start :/

			// buat node
			currentNode = TokenToNode(&token)
			if currentParent != nil {
				currentParent.AddChild(currentNode)
			}
			// pindah ke dalam bagian si child yang baru
			if token.Data != "meta" {
				currentParent = currentNode
			}
		case net_html.EndTagToken:
			// balik naik ke atas
			if currentParent != nil && currentParent.Parent != nil {
				currentParent = currentParent.Parent
			}
		default:
			continue
		}
	}
}
