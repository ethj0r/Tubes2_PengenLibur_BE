package algorithm

// implement DFS secara rekursif
func (node *Node) dfsSearch(selector string) []*Node {
	var result []*Node

	// empty node
	if node == nil {
		return result
	}

	// insert node ke hasil sesuai filter
	if node.includeNode(selector) {
		result = append(result, node)
	}

	if node.hasChildren() {
		// explore hingga paling dalam, baru pindah ke cabang selanjutnya
		for _, child := range node.Children {
			if child != nil {
				childResult := child.dfsSearch(selector)
				result = append(result, childResult...)
			}
		}
	}

	return result
}
