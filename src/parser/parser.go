package parser

import (
	"backend/src/algorithm"
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

func TokenToNode(token *net_html.Token) *algorithm.Node {
	isText := false
	attributes := []algorithm.Attribute{}
	if token.Type == net_html.TextToken {
		isText = false
	}

	if !isText {
		for i := 0; i < (len(token.Attr)); i++ {
			attributes = append(attributes, algorithm.Attribute{
				Name:  token.Attr[i].Key,
				Value: token.Attr[i].Val,
			})

		}
	}

	newElmt := algorithm.Element{
		Data:       token.Data,
		IsText:     isText,
		Attributes: attributes,
	}

	return algorithm.NewNode(newElmt)
}

func Parse(src []byte) (*algorithm.Node, error) {
	source := bytes.NewReader(src)
	tokenizer := net_html.NewTokenizer(source)

	// Iterate through each token

	var currentParent *algorithm.Node = nil
	fmt.Println("Here!")
	for {
		tok := tokenizer.Next()
		token := tokenizer.Token()

		var currentNode *algorithm.Node = nil

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
			// buat node
			currentNode = TokenToNode(&token)
			if currentParent != nil {
				currentParent.AddChild(currentNode)
			}
			// pindah ke dalam bagian si child yang baru
			currentParent = currentNode
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
