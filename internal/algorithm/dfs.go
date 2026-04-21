package algorithm

import "sync"

// implement DFS secara rekursif
func (node *Node) dfsSearch(selector string) []*Node {
	var result []*Node

	// empty node
	if node == nil {
		return result
	}

	// insert node ke hasil sesuai filter
	if includeNode(node, selector) {
		result = append(result, node)
	}

	if hasChildren(node) {
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

// implement DFS paralel (concurrent)
func (node *Node) dfsSearchParallel(selector string, maxProcess int) []*Node {
	if maxProcess <= 1 {
		return node.dfsSearch(selector)
	}

	// Semaphore untuk membatasi jumlah goroutine aktif (semacam thread)
	sem := make(chan struct{}, maxProcess)
	return dfsSearchConcurrent(node, selector, sem)
}

func dfsSearchConcurrent(node *Node, selector string, sem chan struct{}) []*Node {
	var result []*Node

	if node == nil {
		return result
	}

	if includeNode(node, selector) {
		result = append(result, node)
	}

	if !hasChildren(node) {
		return result
	}

	childResults := make([][]*Node, len(node.Children))
	var wg sync.WaitGroup

	for i, child := range node.Children {
		if child == nil {
			continue
		}

		// Kalau true, ada slot kosong (run child di goroutine)
		// Kalau false, semua slot penuh (lanjut ke proses sinkron)
		if !acquireSem(sem) {
			childResults[i] = dfsSearchConcurrent(child, selector, sem)
			continue
		}

		wg.Add(1)
		go func(idx int, childNode *Node) {
			defer wg.Done()

			// Semaphore dilepaskan supaya slot process bisa dipakai subtree lain
			defer releaseSem(sem)
			childResults[idx] = dfsSearchConcurrent(childNode, selector, sem)
		}(i, child)
	}

	wg.Wait()

	// Gabungkan child
	for _, childResult := range childResults {
		result = append(result, childResult...)
	}

	return result
}

func acquireSem(sem chan struct{}) bool {
	// Kalau channel buffer penuh, langsung return false (tidak menunggu)
	select {
	case sem <- struct{}{}:
		return true
	default:
		return false
	}
}

func releaseSem(sem chan struct{}) {
	// release 1 semaphore yang sebelumnya diambil oleh acquireSem
	<-sem
}
