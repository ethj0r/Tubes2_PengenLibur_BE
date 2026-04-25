package traverse

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"backend/internal/algorithm"
	"backend/internal/scraper"

	"github.com/gin-gonic/gin"
)

type Request struct {
	Source    string `json:"source" binding:"required,oneof=url html"`
	Input     string `json:"input" binding:"required"`
	Selector  string `json:"selector" binding:"required"`
	Algorithm string `json:"algorithm" binding:"required,oneof=bfs dfs"`
	Limit     int    `json:"limit"`
}

type AttributeDTO struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type NodeDTO struct {
	ID         string         `json:"id"`
	Tag        string         `json:"tag"`
	IsText     bool           `json:"isText"`
	Text       string         `json:"text,omitempty"`
	Attributes []AttributeDTO `json:"attributes"`
	Children   []*NodeDTO     `json:"children"`
}

type LogEntryDTO struct {
	Step    int    `json:"step"`
	NodeID  string `json:"nodeId"`
	Depth   int    `json:"depth"`
	Matched bool   `json:"matched"`
}

type StatsDTO struct {
	TimeTakenMs       float64 `json:"timeTakenMs"`
	VisitedNodeCount  int     `json:"visitedNodeCount"`
	TraversalMaxDepth int     `json:"traversalMaxDepth"`
	TreeMaxDepth      int     `json:"treeMaxDepth"`
}

type Response struct {
	Tree         *NodeDTO      `json:"tree"`
	Log          []LogEntryDTO `json:"log"`
	Matches      []string      `json:"matches"`
	VisitedOrder []string      `json:"visitedOrder"`
	Stats        StatsDTO      `json:"stats"`
}

type Handler struct{}

func NewHandler() *Handler { return &Handler{} }

func (h *Handler) Traverse(c *gin.Context) {
	requestStart := time.Now()

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var htmlBytes []byte
	if req.Source == "url" {
		body, status, err := scraper.Scraper(strings.TrimSpace(req.Input))
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("fetch failed: %v", err)})
			return
		}
		if status >= 400 {
			c.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("remote returned %d", status)})
			return
		}
		htmlBytes = body
	} else {
		trimmed := strings.TrimSpace(req.Input)
		if trimmed == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "html input is empty"})
			return
		}

		htmlBytes = []byte(trimmed)
		if err := algorithm.ValidateHTML(htmlBytes); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid html structure: %v", err)})
			return
		}
	}

	root, err := algorithm.Parse(htmlBytes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("parse failed: %v", err)})
		return
	}
	if root == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty document"})
		return
	}

	var report algorithm.TraversalReport
	if req.Algorithm == "bfs" {
		report = algorithm.BFSSearchReport(root, req.Selector, req.Limit)
	} else {
		report = algorithm.DFSSearchReport(root, req.Selector, req.Limit)
	}

	idByNode := map[*algorithm.Node]string{}
	tree := serializeTree(root, idByNode, req.Algorithm)

	logDTO := make([]LogEntryDTO, 0, len(report.Log))
	for _, entry := range report.Log {
		logDTO = append(logDTO, LogEntryDTO{
			Step:    entry.Step,
			NodeID:  idByNode[entry.Node],
			Depth:   entry.Depth,
			Matched: entry.Matched,
		})
	}

	matches := make([]string, 0, len(report.Matches))
	for _, n := range report.Matches {
		if id, ok := idByNode[n]; ok {
			matches = append(matches, id)
		}
	}

	visited := make([]string, 0, len(report.VisitedOrder))
	for _, n := range report.VisitedOrder {
		if id, ok := idByNode[n]; ok {
			visited = append(visited, id)
		}
	}

	totalTimeMs := float64(time.Since(requestStart).Microseconds()) / 1000
	stats := StatsDTO{
		TimeTakenMs:       totalTimeMs,
		VisitedNodeCount:  report.VisitedNodeCount,
		TraversalMaxDepth: report.TraversalMaxDepth,
		TreeMaxDepth:      report.TreeMaxDepth,
	}

	// Save log to folder
	_ = saveTraversalLog(req, stats, logDTO, matches, visited)

	c.JSON(http.StatusOK, Response{
		Tree:         tree,
		Log:          logDTO,
		Matches:      matches,
		VisitedOrder: visited,
		Stats:        stats,
	})
}

func serializeTree(root *algorithm.Node, idByNode map[*algorithm.Node]string, traversalMode string) *NodeDTO {
	assignNodeIDs(root, idByNode, traversalMode)

	var build func(n *algorithm.Node) *NodeDTO
	build = func(n *algorithm.Node) *NodeDTO {
		if n == nil {
			return nil
		}

		attrs := make([]AttributeDTO, 0, len(n.Elmt.Attributes))
		for _, a := range n.Elmt.Attributes {
			attrs = append(attrs, AttributeDTO{Name: a.Name, Value: a.Value})
		}

		dto := &NodeDTO{
			ID:         idByNode[n],
			Tag:        n.Elmt.Data,
			IsText:     n.Elmt.IsText,
			Attributes: attrs,
			Children:   make([]*NodeDTO, 0, len(n.Children)),
		}
		if n.Elmt.IsText {
			dto.Text = n.Elmt.Data
		}
		for _, child := range n.Children {
			if c := build(child); c != nil {
				dto.Children = append(dto.Children, c)
			}
		}
		return dto
	}

	return build(root)
}

func assignNodeIDs(root *algorithm.Node, idByNode map[*algorithm.Node]string, traversalMode string) {
	if root == nil {
		return
	}

	counter := 0
	nextID := func() string {
		id := fmt.Sprintf("n%d", counter)
		counter++
		return id
	}

	if traversalMode == "bfs" {
		queue := []*algorithm.Node{root}
		for len(queue) > 0 {
			curr := queue[0]
			queue = queue[1:]
			if curr == nil {
				continue
			}

			idByNode[curr] = nextID()
			for _, child := range curr.Children {
				if child != nil {
					queue = append(queue, child)
				}
			}
		}
		return
	}

	var walk func(n *algorithm.Node)
	walk = func(n *algorithm.Node) {
		if n == nil {
			return
		}

		idByNode[n] = nextID()
		for _, child := range n.Children {
			walk(child)
		}
	}

	walk(root)
}
