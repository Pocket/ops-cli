package util

func RemoveDuplicatesFromSlice(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

func ExcludeMainBranchFromSlice(elements []string, mainBranch string) []string {
	result := []string{}

	for v := range elements {
		// Append element to result if it is not the main branch
		if elements[v] != mainBranch {
			result = append(result, elements[v])
		}
	}

	// Return the new slice.
	return result
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
