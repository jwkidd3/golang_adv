package util

import "strconv"

func IntSliceToStrSlice(intSl []int) []string {
	var strSl []string
	for _, val := range intSl {
		strSl = append(strSl, strconv.Itoa(val))
	}
	return strSl
}
