package algorithm

import "strings"

// mengecek apakah node sesuai selector (tag, .class, #id, *)
func (node *Node) isMatch(selector string) bool {
	// empty node
	if node == nil {
		return false
	}

	// empty selector
	selector = strings.TrimSpace(selector)
	if selector == "" {
		return false
	}

	// universal
	if selector == "*" {
		return true
	}

	// tag selector (tanpa prefix)
	if selector[0] != '.' && selector[0] != '#' {
		return node.Elmt.Data == selector
	}

	// class & id selector
	for _, atr := range node.Elmt.Attributes {
		// class selector (bisa multi-class, misal "a b c")
		if atr.Name == "class" && selector[0] == '.' {
			classes := strings.Fields(atr.Value)

			for _, cls := range classes {
				if "." + cls == selector {
					return true
				}
			}
		}

		// id selector
		if atr.Name == "id" && selector[0] == '#' {
			if "#" + strings.TrimSpace(atr.Value) == selector {
				return true
			}
		}
	}

	return false
}