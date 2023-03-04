package utils

func Merge[T any](acc map[string]T, maps ...map[string]T) map[string]T {
	if acc == nil {
		acc = map[string]T{}
	}
	for _, rule := range maps {
		for k, v := range rule {
			acc[k] = v
		}
	}
	return acc
}
