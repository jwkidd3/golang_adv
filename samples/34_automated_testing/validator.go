package validator

import (
	"fmt"
	"regexp"
)

var (
	isValidEmail = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString
)

func ValidateEmail(input string) error {
	if !isValidEmail(input) {
		return fmt.Errorf("invalid email")
	}
	return nil
}
