package luhn

import (
	"strconv"
)

func Check(orderID string) bool {
	var total = 0
	var slice []int

	if _, err := strconv.Atoi(orderID); err != nil {
		return false
	}
	for _, digit := range orderID {
		slice = append(slice, int(digit-'0'))
	}

	for pos, digit := range slice {
		if pos%2 == len(orderID)%2 {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		total += digit
	}

	return total%10 == 0
}
