package algorithm

import (
	"fmt"
	"strings"
)

type simpleKind int

const (
	simpleTag simpleKind = iota
	simpleClass
	simpleID
	simpleUniversal
)

type simpleSelector struct {
	kind  simpleKind
	value string
}

type compoundSelector struct {
	parts []simpleSelector
}

type combinator rune

const (
	combDescendant combinator = ' '
	combChild      combinator = '>'
	combAdjacent   combinator = '+'
	combGeneral    combinator = '~'
)

type selectorStep struct {
	compound compoundSelector
	next     combinator
}

type parsedSelector struct {
	steps []selectorStep
}

//divide string selector jd urutan compound+combinator.
func parseSelector(s string) (parsedSelector, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return parsedSelector{}, fmt.Errorf("empty selector")
	}

	var steps []selectorStep
	var buf strings.Builder

	flush := func(next combinator) error {
		if buf.Len()==0 {
			return fmt.Errorf("dangling combinator")
		}
		compound, err := parseCompound(buf.String())
		if err != nil {
			return err
		}
		steps = append(steps, selectorStep{compound: compound, next: next})
		buf.Reset()
		return nil
	}

	i := 0
	for i < len(s) {
		c := s[i]
		if c == ' ' || c == '\t' {
			j := i
			for j<len(s) && (s[j] == ' ' || s[j] == '\t') {
				j++
			}
			if j>=len(s) {
				break
			}
			nc := s[j]
			if nc =='>'|| nc == '+' || nc == '~' {
				if err := flush(combinator(nc)); err != nil {
					return parsedSelector{}, err
				}
				i = j + 1
				for i < len(s) && (s[i] == ' ' || s[i] == '\t') {
					i++
				}
				continue
			}
			// descendant
			if err := flush(combDescendant); err != nil {
				return parsedSelector{}, err
			}
			i = j
			continue
		}

		// combinator tanpa leading space
		if c == '>' || c == '+' || c == '~' {
			if err := flush(combinator(c)); err != nil {
				return parsedSelector{}, err
			}
			i++
			for i<len(s) && (s[i] == ' ' || s[i] == '\t') {
				i++
			}
			continue
		}

		buf.WriteByte(c)
		i++
	}

	if buf.Len() == 0 {
		return parsedSelector{}, fmt.Errorf("trailing combinator")
	}
	compound, err := parseCompound(buf.String())
	if err != nil {
		return parsedSelector{}, err
	}
	steps = append(steps, selectorStep{compound: compound, next: 0})

	return parsedSelector{steps: steps}, nil
}

func parseCompound(s string) (compoundSelector, error) {
	var cs compoundSelector
	i := 0
	for i < len(s) {
		c := s[i]
		switch {
		case c == '.':
			j := i + 1
			for j < len(s) && !isCompoundTerminator(s[j]) {
				j++
			}
			if j == i+1 {
				return cs, fmt.Errorf("empty class name")
			}
			cs.parts = append(cs.parts, simpleSelector{kind: simpleClass, value: s[i+1 : j]})
			i = j
		case c == '#':
			j := i + 1
			for j < len(s) && !isCompoundTerminator(s[j]) {
				j++
			}
			if j == i+1 {
				return cs, fmt.Errorf("empty id name")
			}
			cs.parts = append(cs.parts, simpleSelector{kind: simpleID, value: s[i+1 : j]})
			i = j
		case c == '*':
			cs.parts = append(cs.parts, simpleSelector{kind: simpleUniversal})
			i++
		default:
			j := i
			for j < len(s) && !isCompoundTerminator(s[j]) {
				j++
			}
			cs.parts = append(cs.parts, simpleSelector{kind: simpleTag, value: s[i:j]})
			i = j
		}
	}
	if len(cs.parts) == 0 {
		return cs, fmt.Errorf("empty compound")
	}
	return cs, nil
}

func isCompoundTerminator(c byte) bool {
	return c == '.' || c == '#' || c == '*' || c == ' ' || c == '\t' ||
		c == '>' || c == '+' || c == '~'
}


//buat matching
func matchSimple(n *Node, s simpleSelector) bool {
	switch s.kind {
	case simpleUniversal:
		return true
	case simpleTag:
		return n.Elmt.Data == s.value
	case simpleClass:
		for _, a := range n.Elmt.Attributes {
			if a.Name != "class" {
				continue
			}
			for _, cls := range strings.Fields(a.Value) {
				if cls == s.value {
					return true
				}
			}
		}
		return false
	case simpleID:
		for _, a := range n.Elmt.Attributes {
			if a.Name == "id" && strings.TrimSpace(a.Value) == s.value {
				return true
			}
		}
		return false
	}
	return false
}

func matchCompound(n *Node, c compoundSelector) bool {
	if n == nil || n.Elmt.IsText {
		return false
	}
	for _, p := range c.parts {
		if !matchSimple(n, p) {
			return false
		}
	}
	return true
}

//buat cek apakah node n cocok dengan seluruh selector chain
func matchSelector(n *Node, sel parsedSelector) bool {
	if n == nil || len(sel.steps) == 0 {
		return false
	}
	last := len(sel.steps) - 1
	if !matchCompound(n, sel.steps[last].compound) {
		return false
	}
	return matchChainBackward(n, sel, last)
}

func matchChainBackward(n *Node, sel parsedSelector, idx int) bool {
	if idx == 0 {
		return true
	}
	comb := sel.steps[idx-1].next
	prev := sel.steps[idx-1].compound

	switch comb {
	case combChild:
		p := n.Parent
		if p == nil || !matchCompound(p, prev) {
			return false
		}
		return matchChainBackward(p, sel, idx-1)

	case combDescendant:
		for p := n.Parent; p != nil; p = p.Parent {
			if matchCompound(p, prev) && matchChainBackward(p, sel, idx-1) {
				return true
			}
		}
		return false

	case combAdjacent:
		if n.Parent == nil {
			return false
		}
		var immediatePrev *Node
		for _, s := range n.Parent.Children {
			if s == n {
				break
			}
			immediatePrev = s
		}
		if immediatePrev == nil || !matchCompound(immediatePrev, prev) {
			return false
		}
		return matchChainBackward(immediatePrev, sel, idx-1)

	case combGeneral:
		if n.Parent == nil {
			return false
		}
		for _, s := range n.Parent.Children {
			if s == n {
				break
			}
			if matchCompound(s, prev) && matchChainBackward(s, sel, idx-1) {
				return true
			}
		}
		return false
	}
	return false
}
