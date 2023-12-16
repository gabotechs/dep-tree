package utils

import "strconv"

func ItoAArr(arr []int) string {
	result := ""
	for _, el := range arr {
		result += strconv.Itoa(el)
	}
	return result
}

func ItoAArr2(arr [2]int) string {
	result := ""
	for _, el := range arr {
		result += strconv.Itoa(el)
	}
	return result
}
