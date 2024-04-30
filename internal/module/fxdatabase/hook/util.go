package hook

func Contains(list []string, item string) bool {
	for _, listItem := range list {
		if listItem == item {
			return true
		}
	}

	return false
}
