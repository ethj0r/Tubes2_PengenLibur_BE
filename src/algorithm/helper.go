package algorithm

import "strings"

// mengecek apakah node sesuai selector (tag, .class, #id, *)
func isMatch(node *Node, selector string) bool {
	// empty node
	if node == nil {
		return false
	}

	// empty selector
	selector = strings.TrimSpace(selector)
	if selector == "" {
		return false
	}

	// universal
	if selector == "*" {
		return true
	}

	// tag selector (tanpa prefix)
	if selector[0] != '.' && selector[0] != '#' {
		return node.Elmt.Data == selector
	}

	// class & id selector
	for _, atr := range node.Elmt.Attributes {
		// class selector (bisa multi-class, misal "a b c")
		if atr.Name == "class" && selector[0] == '.' {
			classes := strings.Fields(atr.Value)

			for _, cls := range classes {
				if "." + cls == selector {
					return true
				}
			}
		}

		// id selector
		if atr.Name == "id" && selector[0] == '#' {
			if "#" + strings.TrimSpace(atr.Value) == selector {
				return true
			}
		}
	}

	return false
}

// mengecek apakah node dapat ditelusuri lagi (bukan nil / text node)
func isSearchableNode(node *Node) bool {
	if node == nil {
		return false
	}

	if node.Elmt.IsText {
		return false
	}

	return true
}

// filter untuk include node ke result (gabungan 2 func sebelumnya)
func includeNode(node *Node, selector string) bool {
	return isSearchableNode(node) && isMatch(node, selector)
}

// mengecek apakah ada child atau tidak
func hasChildren(node *Node) bool {
	return node != nil && len(node.Children) > 0
}

// mengecek apakah node adalah child langsung dari parent (parent > child)
func isDirectChildOf(node *Node, parent *Node) bool {
	return node != nil && parent != nil && node.Parent == parent
}

// mengecek apakah node turunan dari ancestor (ancestor descendant)
func isDescendantOf(node *Node, ancestor *Node) bool {
	if node == nil || ancestor == nil {
		return false
	}

	// telusuri hingga atas
	for cur := node.Parent; cur != nil; cur = cur.Parent {
		if cur == ancestor {
			return true
		}
	}

	return false
}

// mengecek apakah node sibling tepat setelah sibling lain (sibling1 + sibling2)
func isAdjacentSiblingOf(node *Node, sibling *Node) bool {
	if node == nil || sibling == nil || node.Parent == nil || node.Parent != sibling.Parent {
		return false
	}

	siblings := node.Parent.Children
	// mengecek bersebelahan
	for i := 1; i < len(siblings); i++ {
		if siblings[i] == node && siblings[i-1] == sibling {
			return true
		}
	}

	return false
}

// mengecek apakah node adalah sibling setelah sibling lain (sibling1 ~ sibling2)
func isGeneralSiblingOf(node *Node, sibling *Node) bool {
	if node == nil || sibling == nil || node.Parent == nil || node.Parent != sibling.Parent {
		return false
	}

	siblings := node.Parent.Children
	idxNode := -1
	idxSibling := -1

	// traversal semuanya
	for i, s := range siblings {
		if s == node {
			idxNode = i
		}

		if s == sibling {
			idxSibling = i
		}
	}

	// sibling harus sebelum node
	return idxNode > idxSibling && idxSibling >= 0
}
