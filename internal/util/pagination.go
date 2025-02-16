package util

func NumberOfPages(items int64, page int64) int {
	if items < 0 || page <= 0 {
		return 0
	}

	n := items / page
	if items > page*n {
		n++
	}
	return max(1, int(n))
}
