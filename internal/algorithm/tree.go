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
	Data       string
	IsText     bool
	Attributes []Attribute
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
	if node == nil {
		return
	}
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
	fmt.Printf("Current Depth: %d, Data: %s \n", depth, root.Elmt.Data)
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
		Data:       "html",
		IsText:     false,
		Attributes: nil,
	}

	headElement := Element{
		Data:       "head",
		IsText:     false,
		Attributes: nil,
	}

	titleElement := Element{
		Data:       "title",
		IsText:     false,
		Attributes: nil,
	}

	textInsideTitle := Element{
		Data:       "Hello",
		IsText:     true,
		Attributes: nil,
	}

	bodyElement := Element{
		Data:       "body",
		IsText:     false,
		Attributes: []Attribute{kelas1Attribute},
	}

	textInsideBody := Element{
		Data:       "Hi ges !",
		IsText:     true,
		Attributes: nil,
	}

	h1Element := Element{
		Data:       "h1",
		IsText:     false,
		Attributes: []Attribute{kelas2Attribute, kelas3Attribute},
	}

	textInsideh1 := Element{
		IsText:     true,
		Data:       "Hello, World!",
		Attributes: nil,
	}

	pElement := Element{
		Data:       "p",
		IsText:     false,
		Attributes: []Attribute{kelas4Attribute, idUnikAttribute},
	}

	textInsideP := Element{
		IsText:     true,
		Data:       "Aku ganteng!",
		Attributes: nil,
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
