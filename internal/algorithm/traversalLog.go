package algorithm

import (
	"runtime"
	"sync"
	"time"
)

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

//public wrappers
func BFSSearchReport(root *Node, selector string, topN int) TraversalReport {
	return bfsSearchReportParallel(root, selector, topN, runtime.NumCPU())
}

func DFSSearchReport(root *Node, selector string, topN int) TraversalReport {
	return dfsSearchReportParallel(root, selector, topN, runtime.NumCPU())
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
			if topN > 0 && len(report.Matches) >= topN {
				break
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

// bfsSearchReportParallel — BFS level-by-level dengan match check diparalelkan.
// Urutan log, visited, dan matches persis identik dengan versi sequential.
func bfsSearchReportParallel(root *Node, selector string, topN int, maxProcess int) TraversalReport {
	if maxProcess <= 1 {
		return bfsSearchReport(root, selector, topN)
	}

	report := TraversalReport{}
	if root == nil {
		return report
	}

	start := time.Now()

	type queueItem struct {
		node  *Node
		depth int
	}

	level := []queueItem{{node: root, depth: 0}}
	step := 0
	sem := make(chan struct{}, maxProcess)
	stopped := false

	for len(level) > 0 && !stopped {
		matches := make([]bool, len(level))
		var wg sync.WaitGroup

		for i, it := range level {
			if it.node == nil {
				continue
			}
			// includeNode murni (read-only) — aman dirunning paralel
			if !acquireSem(sem) {
				matches[i] = includeNode(it.node, selector)
				continue
			}
			wg.Add(1)
			go func(idx int, n *Node) {
				defer wg.Done()
				defer releaseSem(sem)
				matches[idx] = includeNode(n, selector)
			}(i, it.node)
		}
		wg.Wait()

		// merge secara sequential → mempertahankan urutan BFS untuk log/step/matches
		next := make([]queueItem, 0)
		for i, it := range level {
			if it.node == nil {
				continue
			}
			step++
			report.VisitedNodeCount++
			report.VisitedOrder = append(report.VisitedOrder, it.node)
			report.Log = append(report.Log, TraversalLogEntry{
				Step:    step,
				Node:    it.node,
				Depth:   it.depth,
				Matched: matches[i],
			})
			if matches[i] {
				report.Matches = append(report.Matches, it.node)
				if topN > 0 && len(report.Matches) >= topN {
					stopped = true
					break
				}
			}
			if it.depth > report.TraversalMaxDepth {
				report.TraversalMaxDepth = it.depth
			}
			for _, c := range it.node.Children {
				if c != nil {
					next = append(next, queueItem{node: c, depth: it.depth + 1})
				}
			}
		}
		level = next
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

// dfsSubReport — hasil traversal subtree (step belum diassign).
type dfsSubReport struct {
	visited []*Node
	log     []TraversalLogEntry
	matches []*Node
	maxDep  int
}

// dfsSearchReportParallel — divide & conquer. Setiap subtree dikerjakan di goroutine,
// lalu hasilnya digabung dalam urutan children (preserving pre-order DFS).
// Step number diassign sequential di akhir agar identik dengan DFS reguler.
func dfsSearchReportParallel(root *Node, selector string, topN int, maxProcess int) TraversalReport {
	if maxProcess <= 1 {
		return dfsSearchReport(root, selector, topN)
	}
	report := TraversalReport{}
	if root == nil {
		return report
	}

	start := time.Now()
	sem := make(chan struct{}, maxProcess)

	sub := dfsSubtreeParallel(root, selector, 0, sem)

	step := 0
	for i, entry := range sub.log {
		step++
		entry.Step = step
		report.Log = append(report.Log, entry)
		report.VisitedOrder = append(report.VisitedOrder, sub.visited[i])
		report.VisitedNodeCount++
		if entry.Matched {
			report.Matches = append(report.Matches, sub.visited[i])
			if topN > 0 && len(report.Matches) >= topN {
				break
			}
		}
		if entry.Depth > report.TraversalMaxDepth {
			report.TraversalMaxDepth = entry.Depth
		}
	}

	report.TreeMaxDepth = getMaxDepth(root)
	report.TimeTaken = time.Since(start).Milliseconds()
	return report
}

func dfsSubtreeParallel(n *Node, selector string, depth int, sem chan struct{}) dfsSubReport {
	var r dfsSubReport
	if n == nil {
		return r
	}

	matched := includeNode(n, selector)
	r.visited = append(r.visited, n)
	r.log = append(r.log, TraversalLogEntry{
		Node:    n,
		Depth:   depth,
		Matched: matched,
	})
	if matched {
		r.matches = append(r.matches, n)
	}
	r.maxDep = depth

	if !hasChildren(n) {
		return r
	}

	childReps := make([]dfsSubReport, len(n.Children))
	var wg sync.WaitGroup
	for i, c := range n.Children {
		if c == nil {
			continue
		}
		if !acquireSem(sem) {
			childReps[i] = dfsSubtreeParallel(c, selector, depth+1, sem)
			continue
		}
		wg.Add(1)
		go func(idx int, child *Node) {
			defer wg.Done()
			defer releaseSem(sem)
			childReps[idx] = dfsSubtreeParallel(child, selector, depth+1, sem)
		}(i, c)
	}
	wg.Wait()

	for _, cr := range childReps {
		r.visited = append(r.visited, cr.visited...)
		r.log = append(r.log, cr.log...)
		r.matches = append(r.matches, cr.matches...)
		if cr.maxDep > r.maxDep {
			r.maxDep = cr.maxDep
		}
	}
	return r
}
