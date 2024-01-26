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

type in2[I1 comparable, I2 comparable] struct {
	i1 I1
	i2 I2
}

func Cached2In1OutErr[I1 comparable, I2 comparable, O1 any](f func(I1, I2) (O1, error)) func(I1, I2) (O1, error) {
	cache := make(map[in2[I1, I2]]O1)
	return func(i1 I1, i2 I2) (O1, error) {
		key := in2[I1, I2]{i1, i2}
		if _, ok := cache[key]; !ok {
			o1, err := f(i1, i2)
			if err != nil {
				return o1, err
			}
			cache[key] = o1
		}
		value, _ := cache[key]
		return value, nil
	}
}

func Cached1In1OutErr[I comparable, O1 any](f func(I) (O1, error)) func(I) (O1, error) {
	cache := make(map[I]O1)
	return func(x I) (O1, error) {
		if _, ok := cache[x]; !ok {
			o1, err := f(x)
			if err != nil {
				return o1, err
			}
			cache[x] = o1
		}
		value, _ := cache[x]
		return value, nil
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

func Cached1In2Out[I comparable, O1 any, O2 any](f func(I) (O1, O2)) func(I) (O1, O2) {
	cache := make(map[I]out2[O1, O2])
	return func(x I) (O1, O2) {
		if _, ok := cache[x]; !ok {
			o1, o2 := f(x)
			cache[x] = out2[O1, O2]{o1, o2}
		}
		value, _ := cache[x]
		return value.o1, value.o2
	}
}
