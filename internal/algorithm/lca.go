package algorithm

import "fmt"

type LCASolver struct {
	nodes []*Node
	idByNode map[*Node]int
	nodeByID map[string]int
	depth []int
	ancestor [][]int
	logN int
}

func NewLCASolver(root *Node, ids []string) (*LCASolver, error) {
	if root == nil {
		return nil, fmt.Errorf("root is nil")
	}

	s := &LCASolver{
		idByNode: map[*Node]int{},
		nodeByID: map[string]int{},
	}

	var walk func(n *Node, d int)
	walk = func(n *Node, d int) {
		if n == nil {
			return
		}
		idx := len(s.nodes)
		s.nodes = append(s.nodes, n)
		s.idByNode[n] = idx
		s.depth = append(s.depth, d)
		for _, c := range n.Children {
			walk(c, d+1)
		}
	}
	walk(root, 0)

	n := len(s.nodes)
	if len(ids) != n {
		return nil, fmt.Errorf("id slice length %d != node count %d", len(ids), n)
	}
	for i, id := range ids {
		s.nodeByID[id] = i
	}

	logN := 1
	for (1 << logN) < n {
		logN++
	}
	s.logN = logN

	s.ancestor = make([][]int, logN+1)
	for k := 0; k <= logN; k++ {
		s.ancestor[k] = make([]int, n)
		for i := range s.ancestor[k] {
			s.ancestor[k][i] = -1
		}
	}

	for i, node := range s.nodes {
		if node.Parent != nil {
			if pi, ok := s.idByNode[node.Parent]; ok {
				s.ancestor[0][i] = pi
			}
		}
	}

	for k := 1; k <= logN; k++ {
		for v := 0; v < n; v++ {
			mid := s.ancestor[k-1][v]
			if mid == -1 {
				s.ancestor[k][v] = -1
				continue
			}
			s.ancestor[k][v] = s.ancestor[k-1][mid]
		}
	}

	return s, nil
}

func (s *LCASolver) LCAByID(idA, idB string) (string, error) {
	a, okA := s.nodeByID[idA]
	b, okB := s.nodeByID[idB]
	if !okA {
		return "", fmt.Errorf("unknown node id: %s", idA)
	}
	if !okB {
		return "", fmt.Errorf("unknown node id: %s", idB)
	}

	lca := s.lcaIdx(a, b)
	if lca == -1 {
		return "", fmt.Errorf("no common ancestor")
	}

	for id, idx := range s.nodeByID {
		if idx == lca {
			return id, nil
		}
	}
	return "", fmt.Errorf("lca id not found")
}

func (s *LCASolver) lcaIdx(u, v int) int {
	if s.depth[u] < s.depth[v] {
		u, v = v, u
	}
	diff := s.depth[u] - s.depth[v]
	for k := 0; k <= s.logN && diff > 0; k++ {
		if diff&1 == 1 {
			u = s.ancestor[k][u]
			if u == -1 {
				return -1
			}
		}
		diff >>= 1
	}
	if u == v {
		return u
	}
	for k := s.logN; k >= 0; k-- {
		au := s.ancestor[k][u]
		av := s.ancestor[k][v]
		if au != av {
			u, v = au, av
			if u == -1 || v == -1 {
				return -1
			}
		}
	}
	return s.ancestor[0][u]
}
