package algorithm

import "sync"

type bfsLevelResult struct {
	matchedNode *Node
	children    []*Node
}

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

// implement BFS paralel per level dengan urutan hasil tetap sama seperti bfsSearch
func (node *Node) bfsSearchParallel(selector string, maxProcess int) []*Node {
	if maxProcess <= 1 {
		return node.bfsSearch(selector)
	}

	var result []*Node
	if node == nil {
		return result
	}

	// Semaphore untuk membatasi jumlah goroutine aktif (semacam thread)
	sem := make(chan struct{}, maxProcess)
	currentLevel := []*Node{node}

	for len(currentLevel) > 0 {
		// Simpan hasil per level
		levelResults := make([]bfsLevelResult, len(currentLevel))
		var wg sync.WaitGroup

		for i, curr := range currentLevel {
			if curr == nil {
				continue
			}

			// Kalau true, ada slot kosong (run child di goroutine)
			// Kalau false, semua slot penuh (lanjut ke proses sinkron)
			if !acquireSem(sem) {
				levelResults[i] = bfsCollectNodeData(curr, selector)
				continue
			}

			wg.Add(1)
			go func(idx int, curNode *Node) {
				defer wg.Done()
				defer releaseSem(sem)
				levelResults[idx] = bfsCollectNodeData(curNode, selector)
			}(i, curr)
		}

		wg.Wait()

		nextLevel := make([]*Node, 0)
		for _, lr := range levelResults {
			if lr.matchedNode != nil {
				result = append(result, lr.matchedNode)
			}

			// Gabung child
			if len(lr.children) > 0 {
				nextLevel = append(nextLevel, lr.children...)
			}
		}

		// lanjut tingkat berikut
		currentLevel = nextLevel
	}

	return result
}

func bfsCollectNodeData(node *Node, selector string) bfsLevelResult {
	res := bfsLevelResult{}

	if node == nil {
		return res
	}

	if includeNode(node, selector) {
		res.matchedNode = node
	}

	if hasChildren(node) {
		for _, child := range node.Children {
			if child != nil {
				res.children = append(res.children, child)
			}
		}
	}

	return res
}
