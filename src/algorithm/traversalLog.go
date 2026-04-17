package algorithm

import "time"

// TraversalLogEntry menyimpan history traversal per node
type TraversalLogEntry struct {
	Step    int
	Node    *Node
	Depth   int
	Matched bool
}

// TraversalReport menyimpan hasil pencarian + metrik traversal
type TraversalReport struct {
	Matches           []*Node
	VisitedOrder      []*Node
	Log               []TraversalLogEntry
	VisitedNodeCount  int
	TraversalMaxDepth int
	TreeMaxDepth      int
	TimeTaken         int64 // dalam miliseconds (ms)
}

// mengembalikan depth node dari root (root depth = 0)
func getDepth(node *Node) int {
	if node == nil {
		return 0
	}

	depth := 0
	for cur := node; cur.Parent != nil; cur = cur.Parent {
		depth++
	}

	return depth
}

// mengembalikan depth max subtree dari node
func getMaxDepth(node *Node) int {
	if node == nil {
		return 0
	}

	return maxDepthFrom(node, 0)
}

// helper getMaxDepth func
func maxDepthFrom(node *Node, currentDepth int) int {
	maxDepth := currentDepth
	for _, child := range node.Children {
		if child == nil {
			continue
		}

		childMax := maxDepthFrom(child, currentDepth+1)

		if childMax > maxDepth {
			maxDepth = childMax
		}
	}
	return maxDepth
}

// melakukan BFS selector search dan return log + metrik, topN <= 0 (mengambil semua kemunculan)
func bfsSearchReport(root *Node, selector string, topN int) TraversalReport {
	report := TraversalReport{}
	if root == nil {
		return report
	}

	start := time.Now()

	// queue menyimpan node yang belum diproses dengan depth-nya
	type queueItem struct {
		node  *Node
		depth int
	}

	queue := []queueItem{{node: root, depth: 0}}
	step := 0

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		if curr.node == nil {
			continue
		}

		// step untuk urutan log traversal
		step++
		report.VisitedNodeCount++
		report.VisitedOrder = append(report.VisitedOrder, curr.node)

		matched := includeNode(curr.node, selector)
		report.Log = append(report.Log, TraversalLogEntry{
			Step:    step,
			Node:    curr.node,
			Depth:   curr.depth,
			Matched: matched,
		})

		if matched {
			report.Matches = append(report.Matches, curr.node)
			if topN > 0 && len(report.Matches) >= topN { //				break
			}
		}

		// simpan depth traversal terjauh 
		if curr.depth > report.TraversalMaxDepth {
			report.TraversalMaxDepth = curr.depth
		}

		// insert child ke belakang queue (proses level by level)
		for _, child := range curr.node.Children {
			if child != nil {
				queue = append(queue, queueItem{node: child, depth: curr.depth + 1})
			}
		}
	}

	report.TreeMaxDepth = getMaxDepth(root)
	report.TimeTaken = time.Since(start).Milliseconds()
	return report
}

// melakukan DFS selector search dan mengembalikan log + metrik, topN <= 0 (mengambil semua kemunculan)
func dfsSearchReport(root *Node, selector string, topN int) TraversalReport {
	report := TraversalReport{}
	if root == nil {
		return report
	}

	start := time.Now()
	step := 0
	
	dfsSearchReportVisit(root, selector, topN, &report, &step, 0)
	report.TreeMaxDepth = getMaxDepth(root)
	report.TimeTaken = time.Since(start).Milliseconds()
	return report
}

// helper for dfs report search func
func dfsSearchReportVisit(cur *Node, selector string, topN int, report *TraversalReport, step *int, depth int) bool {
	if cur == nil {
		return false
	}

	*step++
	report.VisitedNodeCount++
	report.VisitedOrder = append(report.VisitedOrder, cur)

	matched := includeNode(cur, selector)
	report.Log = append(report.Log, TraversalLogEntry{
		Step:    *step,
		Node:    cur,
		Depth:   depth,
		Matched: matched,
	})

	if matched {
		report.Matches = append(report.Matches, cur)
		if topN > 0 && len(report.Matches) >= topN {
			return true
		}
	}

	// simpan depth maksimum
	if depth > report.TraversalMaxDepth {
		report.TraversalMaxDepth = depth
	}

	// explore hingga paling dalam, baru pindah ke cabang selanjutnya
	for _, child := range cur.Children {
		if dfsSearchReportVisit(child, selector, topN, report, step, depth+1) {
			return true
		}
	}

	return false
}
