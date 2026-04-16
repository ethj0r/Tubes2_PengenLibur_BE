package algorithm

func (node* Node) isMatch(selector string) bool {
	// mengecek tag name / class
	if node.Elmt.Data == selector {
		return true
	}

	for _, atr := range node.Elmt.Attributes {
		// mengecek class selector
		if atr.Name == "class" && "."+atr.Value == selector {
			return true
		}

		// mengecek id selector
		if atr.Name == "id" && "#"+atr.Value == selector {
			return true
		}
	}

	return false
}