package utils

func Merge[T any](rules ...map[string]T) map[string]T {
	acc := make(map[string]T)
	for _, rule := range rules {
		for k, v := range rule {
			acc[k] = v
		}
	}
	return acc
}
