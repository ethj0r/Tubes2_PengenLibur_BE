package algorithm

import "fmt"

// root node :
// depth = 0
// parent : nil

type Attribute struct {
	Name  string
	Value string
}

type Element struct {
	TagName     string
	TextContent string
	Attributes  []Attribute
}

type Node struct {
	Elmt     Element
	Children []*Node
	Parent   *Node
}

// Making a new node
func NewNode(element Element) *Node {
	return &Node{
		Elmt:     element,
		Children: nil,
		Parent:   nil,
	}
}

// Adding new child
func (node *Node) AddChild(childNode *Node) {
	childNode.Parent = node // set Parent to current node
	node.Children = append(node.Children, childNode)
}

// Getting parent of current node
func (node *Node) GetParentNode() *Node {
	return node.Parent
}

// Getting children of current node
func (node *Node) GetChildren() []*Node {
	return node.Children
}

// Get siblings of the current node
func (node *Node) GetSiblings() []*Node {
	if node.Parent == nil {
		return nil
	}
	return node.Parent.Children
}

// Print out tree
func (root *Node) PrintTree(depth int) {
	if root == nil {
		return
	}

	// proses sekarang
	for i := 0; i < depth; i++ {
		fmt.Print(" ")

	}
	fmt.Printf("Current Depth: %d, Current Tag: %s, Current Text: %s \n", depth, root.Elmt.TagName, root.Elmt.TextContent)
	for i := 0; i < len(root.Elmt.Attributes); i++ {
		for i := 0; i < depth; i++ {
			fmt.Print(" ")

		}
		fmt.Printf("Attribut ke-%d : Name : %s with Value : %s\n", +1, root.Elmt.Attributes[i].Name, root.Elmt.Attributes[i].Value)

	}
	// proses child
	for _, child := range root.Children {
		child.PrintTree(depth + 1)
	}
}

// Build DOM Tree Example Hardcoded
func BuildDomTreeExampleHardcoded() *Node {
	/*
		input :
		<html>
			<head>
				<title>Hello</title>
			</head>
			<body class = "kelas1">
				<h1 class = "kelas2 kelas3">Hello, World!</h1>
				<p id = "idunik" class = "kelas4">Aku ganteng!</p>
				Hi ges!
				<!-- Just a comment -->
			</body>
		</html>
	*/

	kelas1Attribute := Attribute{
		Name:  "class",
		Value: "kelas1",
	}

	kelas2Attribute := Attribute{
		Name:  "class",
		Value: "kelas2",
	}

	kelas3Attribute := Attribute{
		Name:  "class",
		Value: "kelas3",
	}

	kelas4Attribute := Attribute{
		Name:  "class",
		Value: "kelas4",
	}

	idUnikAttribute := Attribute{
		Name:  "id",
		Value: "idunik",
	}

	htmlElement := Element{
		TagName:     "html",
		TextContent: "",
		Attributes:  nil,
	}

	headElement := Element{
		TagName:     "head",
		TextContent: "",
		Attributes:  nil,
	}

	titleElement := Element{
		TagName:     "title",
		TextContent: "",
		Attributes:  nil,
	}

	textInsideTitle := Element{
		TagName:     "",
		TextContent: "Hello",
		Attributes:  nil,
	}

	bodyElement := Element{
		TagName:     "body",
		TextContent: "",
		Attributes:  []Attribute{kelas1Attribute},
	}

	textInsideBody := Element{
		TagName:     "",
		TextContent: "Hi ges !",
		Attributes:  nil,
	}

	h1Element := Element{
		TagName:     "h1",
		TextContent: "",
		Attributes:  []Attribute{kelas2Attribute, kelas3Attribute},
	}

	textInsideh1 := Element{
		TagName:     "",
		TextContent: "Hello, World!",
		Attributes:  nil,
	}

	pElement := Element{
		TagName:     "p",
		TextContent: "",
		Attributes:  []Attribute{kelas4Attribute, idUnikAttribute},
	}

	textInsideP := Element{
		TagName:     "",
		TextContent: "Aku ganteng!",
		Attributes:  nil,
	}

	h1TextNode := NewNode(textInsideh1)
	h1Node := NewNode(h1Element)
	h1Node.AddChild(h1TextNode)

	pTextNode := NewNode((textInsideP))
	pNode := NewNode((pElement))
	pNode.AddChild(pTextNode)

	bodyTextNode := NewNode(textInsideBody)

	bodyNode := NewNode(bodyElement)
	bodyNode.AddChild(h1Node)
	bodyNode.AddChild(pNode)
	bodyNode.AddChild(bodyTextNode)

	textTitleNode := NewNode(textInsideTitle)
	titleNode := NewNode(titleElement)
	titleNode.AddChild(textTitleNode)

	headNode := NewNode(headElement)
	headNode.AddChild(titleNode)

	root := NewNode(htmlElement)
	root.AddChild(headNode) // add si head
	root.AddChild(bodyNode)
	return root
}
