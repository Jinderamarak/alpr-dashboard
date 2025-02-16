package templates

import "time"

func Sequence(n int) []int {
	arr := make([]int, n)
	for i := range arr {
		arr[i] = i + 1
	}
	return arr
}

func FormatDateTime(t time.Time) string {
	return t.Format("02/01/2006 15:04")
}
