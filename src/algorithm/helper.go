package algorithm

import "strings"

// mengecek apakah node sesuai dengan selector yang diinginkan
func (node *Node) isMatch(selector string) bool {
	// node / selector kosong
	if node == nil || selector == "" {
		return false
	}

	// tag selector (no prefix)
	if selector[0] != '.' && selector[0] != '#' {
		return node.Elmt.Data == selector
	}

	// class dan id selector
	for _, atr := range node.Elmt.Attributes {
		// class selector (bisa ada spasi, misal "p1 p2 p3")
		if atr.Name == "class" && selector[0] == '.' {
			classes := strings.Split(atr.Value, " ") // split spasi

			for _, cls := range classes {
				if "." + strings.TrimSpace(cls) == selector { // remove spasi yang masuk
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