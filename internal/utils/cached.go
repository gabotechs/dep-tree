package utils

func Cached[I comparable, O any](f func(I) O) func(I) O {
	cache := make(map[I]O)
	return func(x I) O {
		if _, ok := cache[x]; !ok {
			cache[x] = f(x)
		}
		value, _ := cache[x]
		return value
	}
}
