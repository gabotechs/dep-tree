package utils

func Cached1In1Out[I comparable, O any](f func(I) O) func(I) O {
	cache := make(map[I]O)
	return func(x I) O {
		if _, ok := cache[x]; !ok {
			cache[x] = f(x)
		}
		value, _ := cache[x]
		return value
	}
}

type out2[O1 any, O2 any] struct {
	o1 O1
	o2 O2
}

func Cached1In2OutErr[I comparable, O1 any, O2 any](f func(I) (O1, O2, error)) func(I) (O1, O2, error) {
	cache := make(map[I]out2[O1, O2])
	return func(x I) (O1, O2, error) {
		if _, ok := cache[x]; !ok {
			o1, o2, err := f(x)
			if err != nil {
				return o1, o2, err
			}
			cache[x] = out2[O1, O2]{o1, o2}
		}
		value, _ := cache[x]
		return value.o1, value.o2, nil
	}
}
