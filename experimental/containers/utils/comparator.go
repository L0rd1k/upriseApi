package utils

import "fmt"

type Comparator func(val_1, val_2 interface{}) int

// String comparison
func Comparator_String(val_1, val_2 interface{}) int {
	fmt.Println("------------")
	str1 := val_1.(string)
	str2 := val_2.(string)

	min := len(str2)

	if len(str1) < len(str2) {
		min = len(str1)
	}

	difference := 0

	for i := 0; i < min && difference == 0; i++ {
		difference = int(str1[i]) - int(str2[i])
		fmt.Println(difference, "=", string(str1[i]), int(str1[i]), "-", string(str2[i]), int(str2[i]))
	}

	if difference == 0 {
		difference = len(str1) - len(str2)
	} else if difference < 0 {
		return -1
	} else if difference > 0 {
		return 1
	}
	return 0
}
