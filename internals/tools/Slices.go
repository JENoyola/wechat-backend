package tools

// FilterSliceValues filters values that "t" contains on "s"
func FilterSliceValues[T comparable](s []T, t []T) []T {

	targetMap := make(map[T]bool)
	for _, target := range t {
		targetMap[target] = true
	}

	var res []T
	for _, value := range s {
		if !targetMap[value] {
			res = append(res, value)
		}
	}
	return res
}

// AddSliceValues adds v to s
func AddSliceValues[T comparable](s []T, v ...T) []T {
	return append(s, v...)
}
