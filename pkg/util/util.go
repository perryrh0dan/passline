package util

func ArrayContains(l []string, i string) bool {
	for _, s := range l {
		if s == i {
			return true
		}
	}

	return false
}
