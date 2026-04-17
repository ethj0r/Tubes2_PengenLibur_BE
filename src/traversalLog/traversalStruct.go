package traversalLog

import (
	"backend/src/algorithm"
)

// TraversalLogEntry menyimpan history traversal per node
type TraversalLogEntry struct {
	Step    int
	Node    *algorithm.Node
	Depth   int
	Matched bool
}

// TraversalReport menyimpan hasil pencarian + traversal
type TraversalReport struct {
	Matches             []*algorithm.Node
	VisitedOrder        []*algorithm.Node
	Log                 []TraversalLogEntry
	VisitedNodeCount    int
	TraversalMaxDepth   int
	TreeMaxDepth        int
	TimeTakenNanosecond int64
}

type SearchMethod string