package algorithm

// implement BFS pake queue
func (node *Node) bfsSearch(selector string) []*Node {
	var result []*Node

	// node kosong
	if node == nil {
		return result
	}

	queue := []*Node{node} // antrian node

	for len(queue) > 0 {
		// ambil node paling depan & hapus dari queue
		curr := queue[0]
		queue = queue[1:]

		if curr == nil {
			continue
		}

		// insert node ke hasil sesuai filter
		if includeNode(curr, selector) {
			result = append(result, curr)
		}

		if hasChildren(curr) {
			// insert child ke belakang queue (proses level by level)
			for _, child := range curr.Children {
				if child != nil {
					queue = append(queue, child)
				}
			}
		}
	}

	return result
}
