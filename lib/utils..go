package lib

import "strconv"

func FormatCurrency(amount int) string {
	str := strconv.Itoa(amount)
	n := len(str)
	if n <= 3 {
		return str
	}

	result := ""
	for i, digit := range str {
		if i > 0 && (n-i)%3 == 0 {
			result += "."
		}
		result += string(digit)
	}

	return result
}
