package main

import "fmt"

var pl = fmt.Println

func main() {
	var dayOfWeek int = 3

	switch dayOfWeek {
	case 1:
		pl("Today is Monday.")
	case 2:
		pl("Today is Tuesday.")
	case 3:
		pl("Today is Wednesday.")
	case 4:
		pl("Today is Thursday.")
	case 5:
		pl("Today is Friday.")
	case 6:
		pl("Today is Saturday.")
	case 7:
		pl("Today is Sunday.")
	default:
		pl("Invalid day of the week.")
	}
}
