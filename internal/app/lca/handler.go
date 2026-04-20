package lca

import (
	"fmt"
	"net/http"

	"backend/internal/algorithm"

	"github.com/gin-gonic/gin"
)

type AttributeDTO struct {
	Name string `json:"name"`
	Value string `json:"value"`
}

type NodeDTO struct {
	ID string `json:"id"`
	Tag string `json:"tag"`
	IsText bool `json:"isText"`
	Text string `json:"text,omitempty"`
	Attributes []AttributeDTO `json:"attributes"`
	Children []*NodeDTO `json:"children"`
}

type Request struct {
	Tree *NodeDTO `json:"tree" binding:"required"`
	NodeA string `json:"nodeA" binding:"required"`
	NodeB string `json:"nodeB" binding:"required"`
}

type Response struct {
	LCANodeID string `json:"lcaNodeId"`
}

type Handler struct{}

func NewHandler() *Handler { return &Handler{} }

func (h *Handler) Compute(c *gin.Context) {
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	root, ids, err := buildTree(req.Tree)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid tree: %v", err)})
		return
	}

	solver, err := algorithm.NewLCASolver(root, ids)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lcaID, err := solver.LCAByID(req.NodeA, req.NodeB)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, Response{LCANodeID: lcaID})
}

func buildTree(dto *NodeDTO) (*algorithm.Node, []string, error) {
	if dto == nil {
		return nil, nil, fmt.Errorf("nil tree")
	}

	var ids []string
	var walk func(d *NodeDTO) *algorithm.Node
	walk = func(d *NodeDTO) *algorithm.Node {
		if d == nil {
			return nil
		}
		attrs := make([]algorithm.Attribute, 0, len(d.Attributes))
		for _, a := range d.Attributes {
			attrs = append(attrs, algorithm.Attribute{Name: a.Name, Value: a.Value})
		}
		n := algorithm.NewNode(algorithm.Element{
			Data: d.Tag,
			IsText: d.IsText,
			Attributes: attrs,
		})
		ids = append(ids, d.ID)
		for _, c := range d.Children {
			if child := walk(c); child != nil {
				n.AddChild(child)
			}
		}
		return n
	}

	root := walk(dto)
	if root == nil {
		return nil, nil, fmt.Errorf("empty tree")
	}
	return root, ids, nil
}
